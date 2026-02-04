package main

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

const (
	APIKey    = "e98ce2565b09ecc0"
	CenterID  = "50030001"
	TenantID  = "82"
	ChannelID = "11"

	WarmupAdvanceMs  = 500
	DNSWarmupAdvance = 2 * time.Second
	MaxIndexOffset   = 3 // 最大索引偏移量
)

var (
	execDay         string
	location        int      // v4: 改为单个索引
	netUserIds      []string // 多账号支持
	openId          string
	venueIdIndex    string
	apiSecret       string
	apiVersion      int
	venueId         string
	fieldType       string
	debugMode       bool     // debug 模式开关
	maxOrderPerUser int      // 每用户下单次数限制
	ecardNos        []string // 会员卡号列表（与 netUserIds 一一对应）

	httpClient  *http.Client
	orderCtx    context.Context
	orderCancel context.CancelFunc

	precomputedFieldListURL string
	rateLimiter             *rate.Limiter
	dnsIPs                  []string
	dnsIPIndex              uint32
	dnsIPMu                 sync.RWMutex
	userOrderCount          sync.Map // 每用户下单计数
)

func debugLog(format string, v ...interface{}) {
	if debugMode {
		log.Printf(format, v...)
	}
}

type FieldSegment struct {
	FieldSegmentID string `json:"fieldSegmentId"`
	State          string `json:"state"`
}

type Field struct {
	FieldSegmentList []*FieldSegment `json:"fieldSgementList"`
}

type APIResponse struct {
	Error     int      `json:"error"`
	Message   string   `json:"message"`
	FieldList []*Field `json:"fieldList"`
}

func createHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 60 * time.Second,
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				ip := pickDNSIP()
				if ip == "" {
					return dialer.DialContext(ctx, network, address)
				}

				host, port, err := net.SplitHostPort(address)
				if err != nil {
					return dialer.DialContext(ctx, network, address)
				}
				if host != "web.xports.cn" {
					return dialer.DialContext(ctx, network, address)
				}
				return dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
			},
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   50,
			MaxConnsPerHost:       50,
			IdleConnTimeout:       120 * time.Second,
			DisableCompression:    true,
			ForceAttemptHTTP2:     true,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: false, ServerName: "web.xports.cn"},
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 3 * time.Second,
		},
		Timeout: 3 * time.Second,
	}
}

func setRequestHeaders(req *http.Request) {
	req.Header.Set("Host", "web.xports.cn")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080712) UnifiedPCMacWechat(0xf264160c) XWEB/18056")
	req.Header.Set("xweb_xhr", "1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://servicewechat.com/wxb75b9974eac7896e/17/page-frame.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/json")
}

func md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func generateSign(apiPath string, params map[string]string, timestamp int64) string {
	allParams := map[string]string{
		"apiKey":    APIKey,
		"timestamp": strconv.FormatInt(timestamp, 10),
		"channelId": ChannelID,
		"centerId":  CenterID,
		"tenantId":  TenantID,
		"version":   strconv.Itoa(apiVersion),
	}
	for k, v := range params {
		allParams[k] = v
	}

	keys := make([]string, 0, len(allParams))
	for k := range allParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(allParams[k])
	}

	signStr := apiPath + sb.String() + apiSecret
	encoded := url.QueryEscape(signStr)
	for _, pair := range [][2]string{{"(", "%28"}, {")", "%29"}, {"'", "%27"}, {"!", "%21"}, {"~", "%7E"}} {
		encoded = strings.ReplaceAll(encoded, pair[0], pair[1])
	}
	return md5Hash(encoded)
}

func buildFieldListURL(timestamp int64) string {
	firstUserId := netUserIds[0]
	params := map[string]string{
		"netUserId":       firstUserId,
		"venueId":         venueId,
		"serviceId":       "1002",
		"day":             execDay,
		"selectByfullTag": "0",
		"fieldType":       fieldType,
		"openId":          openId,
	}

	sign := generateSign("/aisports-api/wechatAPI/venue/fieldList", params, timestamp)

	return fmt.Sprintf(
		"https://web.xports.cn/aisports-api/wechatAPI/venue/fieldList?apiKey=%s&timestamp=%d&channelId=%s&netUserId=%s&venueId=%s&serviceId=1002&day=%s&selectByfullTag=0&centerId=%s&fieldType=%s&tenantId=%s&openId=%s&version=%d&sign=%s",
		APIKey, timestamp, ChannelID, firstUserId, venueId, execDay, CenterID, fieldType, TenantID, openId, apiVersion, sign,
	)
}

func buildNewOrderURL(fieldInfo string, timestamp int64, userId string) string {
	params := map[string]string{
		"venueId":   venueId,
		"serviceId": "1002",
		"day":       execDay,
		"fieldType": fieldType,
		"fieldInfo": fieldInfo,
		"ticket":    "",
		"randStr":   "",
		"netUserId": userId,
		"openId":    openId,
	}

	sign := generateSign("/aisports-api/wechatAPI/order/newOrder", params, timestamp)

	return fmt.Sprintf(
		"https://web.xports.cn/aisports-api/wechatAPI/order/newOrder?apiKey=%s&timestamp=%d&channelId=%s&venueId=%s&serviceId=1002&centerId=%s&day=%s&fieldType=%s&fieldInfo=%s&ticket=&randStr=&netUserId=%s&tenantId=%s&openId=%s&version=%d&sign=%s",
		APIKey, timestamp, ChannelID, venueId, CenterID, execDay, fieldType, url.QueryEscape(fieldInfo), userId, TenantID, openId, apiVersion, sign,
	)
}

func buildPayOrderBody(tradeId string, timestamp int64, userId string, userEcardNo string) string {
	params := map[string]string{
		"netUserId": userId,
		"tradeId":   tradeId,
		"payGroup":  "[]",
		"ecardNo":   userEcardNo,
		"openId":    openId,
	}

	sign := generateSign("/aisports-api/api/pay/payOrder", params, timestamp)

	return fmt.Sprintf(
		"apiKey=%s&timestamp=%d&channelId=%s&netUserId=%s&tradeId=%s&payGroup=%s&ecardNo=%s&centerId=%s&tenantId=%s&openId=%s&version=%d&sign=%s",
		APIKey, timestamp, ChannelID, userId, tradeId, url.QueryEscape("[]"), userEcardNo, CenterID, TenantID, openId, apiVersion, sign,
	)
}

func executePayOrder(ctx context.Context, tradeId string, userId string, userIndex int) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	userEcardNo := ""
	if userIndex >= 0 && userIndex < len(ecardNos) {
		userEcardNo = ecardNos[userIndex]
	}
	if userEcardNo == "" {
		debugLog("[%s] 无对应会员卡号，跳过支付", userId)
		return nil
	}

	timestamp := time.Now().UnixMilli()
	body := buildPayOrderBody(tradeId, timestamp, userId, userEcardNo)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://web.xports.cn/aisports-api/api/pay/payOrder", strings.NewReader(body))
	if err != nil {
		debugLog("[%s] 创建支付请求失败: %v", userId, err)
		return err
	}
	setRequestHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		debugLog("[%s] 支付请求失败: %v", userId, err)
		return err
	}
	respBody, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	debugLog("[%s] 支付响应: %s", userId, string(respBody))

	var result struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
	}
	if json.Unmarshal(respBody, &result) == nil {
		if result.Error == 0 {
			log.Printf("💰 账号 %s 支付成功！tradeId: %s", userId, tradeId)
		} else {
			debugLog("[%s] 支付失败: %s", userId, result.Message)
			return fmt.Errorf(result.Message)
		}
	}
	return nil
}

func warmupDNS() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, "web.xports.cn")
	if err != nil {
		debugLog("[预热] DNS 解析失败: %v", err)
		return
	}
	if len(ips) == 0 {
		debugLog("[预热] DNS 解析结果为空")
		return
	}
	updateDNSIPs(ips)
	debugLog("[预热] DNS 解析完成，IP 数量: %d", len(ips))
}

func updateDNSIPs(ips []net.IPAddr) {
	next := make([]string, 0, len(ips))
	for _, item := range ips {
		if item.IP == nil {
			continue
		}
		next = append(next, item.IP.String())
	}
	if len(next) == 0 {
		return
	}
	dnsIPMu.Lock()
	dnsIPs = next
	atomic.StoreUint32(&dnsIPIndex, 0)
	dnsIPMu.Unlock()
}

func pickDNSIP() string {
	dnsIPMu.RLock()
	defer dnsIPMu.RUnlock()
	if len(dnsIPs) == 0 {
		return ""
	}
	idx := atomic.AddUint32(&dnsIPIndex, 1)
	return dnsIPs[int(idx)%len(dnsIPs)]
}

func warmupConnection() {
	urlStr := buildFieldListURL(time.Now().UnixMilli())
	req, _ := http.NewRequest("GET", urlStr, nil)
	setRequestHeaders(req)
	resp, err := httpClient.Do(req)
	if err == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		debugLog("[预热] 连接预热完成")
	} else {
		debugLog("[预热] 连接预热失败: %v", err)
	}
}

func fetchFieldList(ctx context.Context, urlStr string) (*APIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	setRequestHeaders(req)

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response APIResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v, body: %s", err, string(body))
	}

	return &response, nil
}

func findBestIndices(targetIdx int, segmentList []*FieldSegment) []int {
	if targetIdx < 0 || targetIdx >= len(segmentList) {
		return nil
	}

	n := len(segmentList)
	if targetIdx+1 < n {
		if segmentList[targetIdx].State == "0" && segmentList[targetIdx+1].State == "0" {
			debugLog("[查找] 目标位置 %d 和 %d 都可用", targetIdx, targetIdx+1)
			return []int{targetIdx, targetIdx + 1}
		}
	}

	if targetIdx-1 >= 0 {
		if segmentList[targetIdx-1].State == "0" && segmentList[targetIdx].State == "0" {
			debugLog("[查找] 位置 %d 和 %d 都可用", targetIdx-1, targetIdx)
			return []int{targetIdx - 1, targetIdx}
		}
	}

	for offset := 1; offset <= MaxIndexOffset; offset++ {
		startIdx := targetIdx - offset - 1
		if startIdx >= 0 && startIdx+1 < n {
			if segmentList[startIdx].State == "0" && segmentList[startIdx+1].State == "0" {
				debugLog("[查找] 找到两个连续位置: %d, %d (向前偏移)", startIdx, startIdx+1)
				return []int{startIdx, startIdx + 1}
			}
		}

		startIdx = targetIdx + offset
		if startIdx >= 0 && startIdx+1 < n {
			if segmentList[startIdx].State == "0" && segmentList[startIdx+1].State == "0" {
				debugLog("[查找] 找到两个连续位置: %d, %d (向后偏移)", startIdx, startIdx+1)
				return []int{startIdx, startIdx + 1}
			}
		}
	}

	debugLog("[查找] 未找到两个连续位置，开始查找单个位置...")

	if segmentList[targetIdx].State == "0" {
		debugLog("[查找] 目标位置 %d 可用", targetIdx)
		return []int{targetIdx}
	}

	for offset := 1; offset <= MaxIndexOffset; offset++ {
		idx := targetIdx - offset
		if idx >= 0 && segmentList[idx].State == "0" {
			debugLog("[查找] 找到单个位置: %d (向前偏移%d)", idx, offset)
			return []int{idx}
		}

		idx = targetIdx + offset
		if idx < n && segmentList[idx].State == "0" {
			debugLog("[查找] 找到单个位置: %d (向后偏移%d)", idx, offset)
			return []int{idx}
		}
	}

	debugLog("[查找] 在±%d范围内未找到可用位置", MaxIndexOffset)
	return nil
}

func extractFieldSegmentIDs(segmentList []*FieldSegment) string {
	if len(segmentList) == 0 {
		return ""
	}

	indices := findBestIndices(location, segmentList)
	if len(indices) == 0 {
		return ""
	}

	var ids []string
	for _, idx := range indices {
		if idx >= 0 && idx < len(segmentList) {
			ids = append(ids, segmentList[idx].FieldSegmentID)
		}
	}
	return strings.Join(ids, ",")
}

func executeOrder(ctx context.Context, orderURL string, userId string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	req, err := http.NewRequestWithContext(ctx, "GET", orderURL, nil)
	if err != nil {
		return
	}
	setRequestHeaders(req)
	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	body, _ := io.ReadAll(resp.Body)
	debugLog("[%s] 下单响应: %s", userId, string(body))
	_ = resp.Body.Close()
	var result struct {
		Message string `json:"message"`
	}

	if json.Unmarshal(body, &result) == nil {
		if result.Message == "ok" {
			debugLog("🎉 账号 %s 下单响应：%s", userId, result.Message)
		}
	}
}

func processFieldList(response *APIResponse, timestamp int64) {
	var wg sync.WaitGroup
	for _, field := range response.FieldList {
		fieldInfo := extractFieldSegmentIDs(field.FieldSegmentList)
		if fieldInfo == "" {
			continue
		}

		for _, userId := range netUserIds {
			select {
			case <-orderCtx.Done():
				return
			default:
			}
			// 检查该用户下单次数是否超限
			countVal, _ := userOrderCount.LoadOrStore(userId, new(int32))
			count := countVal.(*int32)
			if int(atomic.LoadInt32(count)) >= maxOrderPerUser {
				debugLog("[%s] 已达到下单次数限制 %d", userId, maxOrderPerUser)
				continue
			}
			atomic.AddInt32(count, 1)

			wg.Add(1)
			go func(uid string) {
				defer wg.Done()
				if err := rateLimiter.Wait(orderCtx); err != nil {
					debugLog("速率限制等待失败: %v", err)
					return
				}
				orderURL := buildNewOrderURL(fieldInfo, timestamp, uid)
				executeOrder(orderCtx, orderURL, uid)
			}(userId)
		}
	}
	wg.Wait()
}

func main() {
	var (
		times        string
		startAt      string
		locationStr  string
		netUserIdStr string
	)

	flag.StringVar(&execDay, "day", "", "天数格式: 20250901")
	flag.StringVar(&netUserIdStr, "net_user_id", "", "账号（多账号用逗号分隔）")
	flag.StringVar(&openId, "open_id", "", "openId")
	flag.StringVar(&apiSecret, "api_secret", "", "API密钥")
	flag.IntVar(&apiVersion, "version", 0, "签名版本")
	flag.StringVar(&times, "times", "5", "最大尝试次数")
	flag.StringVar(&startAt, "start", "", "开始时间格式 2025-01-01 00:59:59")
	flag.StringVar(&locationStr, "location", "", "位置（0-based单个索引，如 5）")
	flag.StringVar(&venueIdIndex, "venue_id_index", "", "场馆索引")
	flag.IntVar(&maxOrderPerUser, "max_order", 30, "每用户下单次数限制")
	var ecardNoStr string
	flag.StringVar(&ecardNoStr, "ecard_no", "", "会员卡号（多卡号用逗号分隔，与账号顺序对应）")
	flag.BoolVar(&debugMode, "debug", true, "启用debug日志")
	flag.Parse()

	if execDay == "" ||
		netUserIdStr == "" ||
		locationStr == "" ||
		apiSecret == "" ||
		openId == "" ||
		apiVersion <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, id := range strings.Split(netUserIdStr, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			netUserIds = append(netUserIds, id)
		}
	}
	if len(netUserIds) == 0 {
		log.Fatal("至少需要一个 netUserId")
	}

	if ecardNoStr != "" {
		log.Fatalf("⚠️ 错误: 会员卡号未设置")
	}
	for _, ecard := range strings.Split(ecardNoStr, ",") {
		ecard = strings.TrimSpace(ecard)
		ecardNos = append(ecardNos, ecard)
	}
	if len(ecardNos) != len(netUserIds) {
		log.Fatalf("⚠️ 错误: 会员卡号数量(%d)与账号数量(%d)不一致", len(ecardNos), len(netUserIds))
	}

	var err error
	location, err = strconv.Atoi(locationStr)
	if err != nil {
		log.Fatalf("location 必须是一个整数: %v", err)
	}

	maxAttempts, _ := strconv.Atoi(times)
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	if maxOrderPerUser <= 0 {
		maxOrderPerUser = 3
	}

	switch venueIdIndex {
	case "2":
		venueId, fieldType = "5003000103", "1837"
	default:
		venueId, fieldType = "5003000101", "1841"
	}

	var loc *time.Location
	if loc, err = time.LoadLocation("Asia/Shanghai"); err == nil {
		time.Local = loc
	}

	httpClient = createHTTPClient()
	orderCtx, orderCancel = context.WithCancel(context.Background())
	defer orderCancel()

	rateLimiter = rate.NewLimiter(rate.Every(300*time.Millisecond), 1)
	if startAt != "" {
		var start time.Time
		start, err = time.ParseInLocation(time.DateTime, startAt, time.Local)
		if err != nil {
			log.Fatalf("时间格式错误: %v", err)
		}

		now := time.Now()
		if !now.Before(start) {
			log.Println("指定时间已过")
			return
		}

		targetTimestamp := start.UnixMilli()
		precomputedFieldListURL = buildFieldListURL(targetTimestamp)
		debugLog("[预计算] 已预生成签名 URL")

		dnsWarmupDuration := start.Add(-DNSWarmupAdvance).Sub(now)
		if dnsWarmupDuration <= 0 {
			warmupDNS()
		} else {
			log.Printf("等待 %.2f 秒后执行 DNS 预热...", dnsWarmupDuration.Seconds())
			select {
			case <-time.After(dnsWarmupDuration):
				warmupDNS()
			case <-orderCtx.Done():
				return
			}
		}

		connWarmupDuration := start.Add(-time.Duration(WarmupAdvanceMs) * time.Millisecond).Sub(time.Now())
		if connWarmupDuration > 0 {
			log.Printf("等待 %.2f 秒后开始...", connWarmupDuration.Seconds())
			select {
			case <-time.After(connWarmupDuration):
			case <-orderCtx.Done():
				return
			}
		}
		warmupConnection()
	} else {
		warmupDNS()
		warmupConnection()
	}

	debugLog("开始执行，最大尝试次数: %d, 目标位置: %d", maxAttempts, location)

	var verifyOnce sync.Once
	var verifyWg sync.WaitGroup

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-orderCtx.Done():
			goto End
		default:
		}

		timestamp := time.Now().UnixMilli()

		var fieldListURL string
		if attempt == 1 && precomputedFieldListURL != "" {
			fieldListURL = precomputedFieldListURL
		} else {
			fieldListURL = buildFieldListURL(timestamp)
		}

		debugLog("[尝试 %d] 拉取场地列表...", attempt)

		var response *APIResponse
		response, err = fetchFieldList(orderCtx, fieldListURL)
		if err != nil {
			debugLog("[尝试 %d] 拉取失败: %v", attempt, err)
			continue
		}

		if len(response.FieldList) == 0 {
			debugLog("[尝试 %d] 列表为空（未到开放时间）", attempt)
			continue
		}

		verifyOnce.Do(func() {
			verifyWg.Add(1)
			go func() {
				defer verifyWg.Done()
				verifyOrders()
			}()
		})

		debugLog("[尝试 %d] 成功获取 %d 个场地，开始下单...", attempt, len(response.FieldList))
		processFieldList(response, timestamp)
	}

End:
	log.Printf("下单流程完成，账号数: %d", len(netUserIds))
	verifyWg.Wait()
}

type OrderItem struct {
	TradeId any `json:"tradeId"`
}

type OrderPageInfo struct {
	PageNum  int          `json:"pageNum"`
	PageSize int          `json:"pageSize"`
	Total    int          `json:"total"`
	List     []*OrderItem `json:"list"`
}

type OrderResponse struct {
	Error    int            `json:"error"`
	Message  string         `json:"message"`
	PageInfo *OrderPageInfo `json:"pageInfo"`
}

func buildGetOrdersURL(timestamp int64, userId string) string {
	params := map[string]string{
		"pageNo":     "1",
		"orderState": "2",
		"netUserId":  userId,
		"openId":     openId,
	}

	sign := generateSign("/aisports-api/api/order/user/getOrders", params, timestamp)

	return fmt.Sprintf(
		"https://web.xports.cn/aisports-api/api/order/user/getOrders?apiKey=%s&timestamp=%d&channelId=%s&pageNo=1&orderState=2&netUserId=%s&centerId=%s&tenantId=%s&openId=%s&version=%d&sign=%s",
		APIKey, timestamp, ChannelID, userId, CenterID, TenantID, openId, apiVersion, sign,
	)
}

func payOrdersForUser(userId string, userIndex int, tradeIds []string) error {
	if len(ecardNos) == 0 {
		return nil
	}
	var err error
	for _, tradeId := range tradeIds {
		log.Printf("💳 开始支付订单 %s...", tradeId)
		if err = executePayOrder(context.Background(), tradeId, userId, userIndex); err == nil {
			return nil
		}
	}
	return err
}

func verifyOrders() {
	const maxRetries = 60
	const tickInterval = 1 * time.Second

	log.Println("开始验证订单...")
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()
	var cancelOnce sync.Once
	for attempt := 1; attempt <= maxRetries; attempt++ {
		debugLog("第 %d/%d 次验证订单...", attempt, maxRetries)
		if userIdx, tradeIds := findFirstUserWithOrders(); tradeIds != nil {
			userId := netUserIds[userIdx]
			log.Printf("✅ 账号 %s 订单验证成功，找到 %d 个订单", userId, len(tradeIds))
			cancelOnce.Do(func() {
				orderCancel()
			})
			if err := payOrdersForUser(userId, userIdx, tradeIds); err == nil {
				return
			}
		}

		if attempt < maxRetries {
			<-ticker.C
		} else {
			log.Printf("❌ 已达到最大重试次数 %d，所有账号均未找到订单", maxRetries)
			cancelOnce.Do(func() {
				orderCancel()
			})
		}
	}
}

func findFirstUserWithOrders() (int, []string) {
	for idx, userId := range netUserIds {
		if tradeIds := verifyOrderForUser(userId); len(tradeIds) > 0 {
			return idx, tradeIds
		}
	}
	return -1, nil
}

func verifyOrderForUser(userId string) []string {
	timestamp := time.Now().UnixMilli()
	orderURL := buildGetOrdersURL(timestamp, userId)

	req, err := http.NewRequest("GET", orderURL, nil)
	if err != nil {
		debugLog("[%s] 创建订单请求失败: %v", userId, err)
		return nil
	}
	setRequestHeaders(req)
	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		debugLog("[%s] 获取订单失败: %v", userId, err)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		debugLog("[%s] 读取订单响应失败: %v", userId, err)
		return nil
	}

	debugLog("[%s] 订单响应: %s", userId, string(body))

	var orderResp OrderResponse
	if err = json.Unmarshal(body, &orderResp); err != nil {
		debugLog("[%s] 解析订单响应失败: %v", userId, err)
		return nil
	}

	if orderResp.Error != 0 {
		debugLog("[%s] 订单接口返回错误: %s", userId, orderResp.Message)
		return nil
	}

	if orderResp.PageInfo == nil || len(orderResp.PageInfo.List) == 0 {
		debugLog("[%s] 订单列表为空", userId)
		return nil
	}

	var tradeIds []string
	for _, order := range orderResp.PageInfo.List {
		var tradeIdStr string
		switch v := order.TradeId.(type) {
		case string:
			tradeIdStr = v
		case float64:
			tradeIdStr = strconv.FormatInt(int64(v), 10)
		case int, int32, int64, uint, uint32, uint64:
			tradeIdStr = fmt.Sprintf("%v", v)
		default:
			if v != nil {
				debugLog("[%s] 未知的 TradeId 类型: %T, 值: %v", userId, v, v)
			}
			continue
		}
		if tradeIdStr != "" {
			tradeIds = append(tradeIds, tradeIdStr)
		}
	}

	log.Printf("[%s] 找到 %d 个订单", userId, len(tradeIds))
	return tradeIds
}
