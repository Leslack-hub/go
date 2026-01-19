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
)

const (
	APIKey    = "e98ce2565b09ecc0"
	CenterID  = "50030001"
	TenantID  = "82"
	ChannelID = "11"

	WarmupAdvanceMs = 500
)

var (
	execDay          string
	location         string
	netUserId        string
	openId           string
	venueIdIndex     string
	successExitCount int64
	apiSecret        string
	apiVersion       int

	venueId   string
	fieldType string

	httpClient *http.Client
	gCtx       context.Context
	gCancel    context.CancelFunc

	globalSuccessCount int64

	precomputedFieldListURL string
)

type FieldSegment struct {
	FieldSegmentID string `json:"fieldSegmentId"`
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
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   50,
			MaxConnsPerHost:       50,
			IdleConnTimeout:       120 * time.Second,
			DisableCompression:    true,
			ForceAttemptHTTP2:     true,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
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
	params := map[string]string{
		"netUserId":       netUserId,
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
		APIKey, timestamp, ChannelID, netUserId, venueId, execDay, CenterID, fieldType, TenantID, openId, apiVersion, sign,
	)
}

func buildNewOrderURL(fieldInfo string, timestamp int64) string {
	params := map[string]string{
		"venueId":   venueId,
		"serviceId": "1002",
		"day":       execDay,
		"fieldType": fieldType,
		"fieldInfo": fieldInfo,
		"ticket":    "",
		"randStr":   "",
		"netUserId": netUserId,
		"openId":    openId,
	}

	sign := generateSign("/aisports-api/wechatAPI/order/newOrder", params, timestamp)

	return fmt.Sprintf(
		"https://web.xports.cn/aisports-api/wechatAPI/order/newOrder?apiKey=%s&timestamp=%d&channelId=%s&venueId=%s&serviceId=1002&centerId=%s&day=%s&fieldType=%s&fieldInfo=%s&ticket=&randStr=&netUserId=%s&tenantId=%s&openId=%s&version=%d&sign=%s",
		APIKey, timestamp, ChannelID, venueId, CenterID, execDay, fieldType, url.QueryEscape(fieldInfo), netUserId, TenantID, openId, apiVersion, sign,
	)
}

func warmupConnection() {
	urlStr := buildFieldListURL(time.Now().UnixMilli())
	req, _ := http.NewRequest("GET", urlStr, nil)
	setRequestHeaders(req)
	resp, err := httpClient.Do(req)
	if err == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		log.Println("[预热] 连接预热完成")
	} else {
		log.Printf("[预热] 连接预热失败: %v", err)
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

func extractFieldSegmentIDs(locs []string, segmentList []*FieldSegment) string {
	if len(segmentList) == 0 {
		return ""
	}

	var ids []string
	for _, loc := range locs {
		idx, err := strconv.Atoi(loc)
		if err != nil || idx < 0 || idx >= len(segmentList) {
			continue
		}
		ids = append(ids, segmentList[idx].FieldSegmentID)
	}
	return strings.Join(ids, ",")
}

func executeOrder(ctx context.Context, orderURL string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	if atomic.LoadInt64(&globalSuccessCount) >= successExitCount {
		return
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
	_ = resp.Body.Close()
	log.Printf("下单响应: %s", string(body))
	var result struct {
		Message string `json:"message"`
	}

	if json.Unmarshal(body, &result) == nil {
		if result.Message == "ok" || result.Message == "场地预定中，请勿重复提交" {
			count := atomic.AddInt64(&globalSuccessCount, 1)
			log.Printf("🎉 抢票成功！(%d/%d)", count, successExitCount)
			if count >= successExitCount {
				gCancel()
			}
		}
	}
}

func processFieldList(response *APIResponse, timestamp int64) {
	locs := strings.Split(location, ",")
	var wg sync.WaitGroup
	for _, field := range response.FieldList {
		fieldInfo := extractFieldSegmentIDs(locs, field.FieldSegmentList)
		if fieldInfo == "" {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			orderURL := buildNewOrderURL(fieldInfo, timestamp)
			executeOrder(gCtx, orderURL)
		}()
	}
	wg.Wait()
}

func main() {
	var (
		times   string
		startAt string
	)

	flag.StringVar(&execDay, "day", "", "天数格式: 20250901")
	flag.StringVar(&netUserId, "net_user_id", "", "账号")
	flag.StringVar(&openId, "open_id", "", "openId")
	flag.StringVar(&apiSecret, "api_secret", "", "API密钥")
	flag.IntVar(&apiVersion, "version", 0, "签名版本")
	flag.StringVar(&times, "times", "5", "最大尝试次数")
	flag.StringVar(&startAt, "start", "", "开始时间格式 2025-01-01 00:59:59")
	flag.StringVar(&location, "location", "", "位置（0-based索引，如 0,1）")
	flag.StringVar(&venueIdIndex, "venue_id_index", "", "场馆索引")
	flag.Int64Var(&successExitCount, "ok_count", 1, "成功次数阈值")
	flag.Parse()

	if execDay == "" || netUserId == "" || location == "" || apiSecret == "" || openId == "" || apiVersion <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	maxAttempts, _ := strconv.Atoi(times)
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	if successExitCount <= 0 {
		successExitCount = 1
	}

	switch venueIdIndex {
	case "2":
		venueId, fieldType = "5003000103", "1837"
	default:
		venueId, fieldType = "5003000101", "1841"
	}

	if loc, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		time.Local = loc
	}

	httpClient = createHTTPClient()
	gCtx, gCancel = context.WithCancel(context.Background())
	defer gCancel()

	warmupConnection()

	if startAt != "" {
		start, err := time.ParseInLocation(time.DateTime, startAt, time.Local)
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
		log.Printf("[预计算] 已预生成签名 URL")

		waitDuration := start.Add(-time.Duration(WarmupAdvanceMs) * time.Millisecond).Sub(now)
		log.Printf("等待 %.2f 秒后开始...", waitDuration.Seconds())

		select {
		case <-time.After(waitDuration):
		case <-gCtx.Done():
			return
		}
	}

	log.Printf("开始执行，最大尝试次数: %d", maxAttempts)
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-gCtx.Done():
			goto End
		default:
		}

		if atomic.LoadInt64(&globalSuccessCount) >= successExitCount {
			goto End
		}

		timestamp := time.Now().UnixMilli()

		var fieldListURL string
		if attempt == 1 && precomputedFieldListURL != "" {
			fieldListURL = precomputedFieldListURL
		} else {
			fieldListURL = buildFieldListURL(timestamp)
		}

		log.Printf("[尝试 %d] 拉取场地列表...", attempt)

		response, err := fetchFieldList(gCtx, fieldListURL)
		if err != nil {
			log.Printf("[尝试 %d] 拉取失败: %v", attempt, err)
			continue
		}

		if len(response.FieldList) == 0 {
			log.Printf("[尝试 %d] 列表为空（未到开放时间）", attempt)
			continue
		}

		log.Printf("[尝试 %d] 成功获取 %d 个场地，开始下单...", attempt, len(response.FieldList))
		processFieldList(response, timestamp)

		if atomic.LoadInt64(&globalSuccessCount) >= successExitCount {
			goto End
		}
	}

End:
	log.Printf("执行完成，成功次数: %d", atomic.LoadInt64(&globalSuccessCount))
}
