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
	// 减少重试延迟
	RetryDelay = 5 * time.Millisecond
	// 并发 worker 数量
	NumWorkers = 16
	// 预热提前时间（毫秒）
	WarmupAdvanceMs = 50
)

type FieldSegment struct {
	Price          int    `json:"price"`
	Segment        int    `json:"segment"`
	BookingStatus  string `json:"bookingStatus"`
	Step           int    `json:"step"`
	State          string `json:"state"`
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

var (
	WorkerChan       chan OrderRequest
	WorkerChanWg     *sync.WaitGroup
	GCtx             context.Context
	GCancel          context.CancelFunc
	ExecDay          string
	Location         string
	NetUserId        string
	OpenId           string
	VenueIdIndex     string
	SuccessExitCount int64
	// 全局 HTTP 客户端
	HttpClient *http.Client
	// 成功计数器
	GlobalSuccessCount int64
)

// OrderRequest 用于传递下单请求信息
type OrderRequest struct {
	URL string
}

// 创建高性能 HTTP 客户端
func createHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxIdleConns:          500,
		MaxIdleConnsPerHost:   200,
		MaxConnsPerHost:       200,
		IdleConnTimeout:       120 * time.Second,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		TLSHandshakeTimeout:   2 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}
}

// 预热连接
func warmupConnection() {
	req, _ := http.NewRequest("HEAD", "https://web.xports.cn/", nil)
	req.Header.Set("Connection", "keep-alive")
	resp, err := HttpClient.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}

func orderWorker() {
	for req := range WorkerChan {
		executeOrder(req)
		WorkerChanWg.Done()
	}
}

func executeOrder(orderReq OrderRequest) {
	select {
	case <-GCtx.Done():
		return
	default:
	}

	if atomic.LoadInt64(&GlobalSuccessCount) >= SuccessExitCount {
		return
	}

	req, err := http.NewRequestWithContext(GCtx, "GET", orderReq.URL, nil)
	if err != nil {
		return
	}

	setRequestHeaders(req)

	var resp *http.Response
	resp, err = HttpClient.Do(req)
	if err != nil {
		return
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	_ = resp.Body.Close()

	checkOrderResponse(body)
}

func setRequestHeaders(req *http.Request) {
	req.Header.Set("Host", "web.xports.cn")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080712) UnifiedPCMacWechat(0xf2641015) XWEB/16390")
	req.Header.Set("xweb_xhr", "1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://servicewechat.com/wxb75b9974eac7896e/11/page-frame.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/json")
}

func checkOrderResponse(body []byte) {
	log.Printf("下单响应: %s", string(body))

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}

	if result.Message == "ok" || result.Message == "场地预定中，请勿重复提交" {
		count := atomic.AddInt64(&GlobalSuccessCount, 1)
		log.Printf("🎉 抢票成功！(%d/%d)", count, SuccessExitCount)
		if count >= SuccessExitCount {
			GCancel()
		}
	}
}

func extractFieldSegmentIDs(locations []string, segmentList []*FieldSegment) string {
	if len(locations) == 0 {
		return ""
	}
	available := make(map[int]string)
	for i, segment := range segmentList {
		if segment.State == "0" && segment.Price == 0 && segment.FieldSegmentID != "" {
			available[i] = segment.FieldSegmentID
		}
	}
	if len(available) == 0 {
		return ""
	}

	center := 0
	if l1, err := strconv.Atoi(locations[0]); err == nil && l1 > 0 && l1 <= len(segmentList) {
		center = l1 - 1
	}
	rightStart := center + 1
	if len(locations) >= 2 {
		if l2, err := strconv.Atoi(locations[1]); err == nil && l2 > 0 && l2 <= len(segmentList) {
			rightStart = l2 - 1
		}
	}

	withinBounds := func(idx int) bool {
		return idx >= 0 && idx < len(segmentList)
	}

	// 找连续两张
	for offset := 0; offset < len(segmentList); offset++ {
		startLeft := center - offset
		if withinBounds(startLeft) && withinBounds(startLeft+1) {
			if id1, ok1 := available[startLeft]; ok1 {
				if id2, ok2 := available[startLeft+1]; ok2 {
					return id1 + "," + id2
				}
			}
		}

		startRight := rightStart + offset
		if withinBounds(startRight) && withinBounds(startRight+1) {
			if id1, ok1 := available[startRight]; ok1 {
				if id2, ok2 := available[startRight+1]; ok2 {
					return id1 + "," + id2
				}
			}
		}
	}

	// 取最多两张
	var ids []string
	seen := make(map[int]struct{})
	for step := 0; step < len(segmentList) && len(ids) < 2; step++ {
		left := center - step
		if withinBounds(left) {
			if id, ok := available[left]; ok {
				if _, exist := seen[left]; !exist {
					ids = append(ids, id)
					seen[left] = struct{}{}
					if len(ids) == 2 {
						break
					}
				}
			}
		}

		right := rightStart + step
		if withinBounds(right) {
			if id, ok := available[right]; ok {
				if _, exist := seen[right]; !exist {
					ids = append(ids, id)
					seen[right] = struct{}{}
				}
			}
		}
	}

	return strings.Join(ids, ",")
}

func processFieldList(response *APIResponse) {
	wg := sync.WaitGroup{}
	for i, field := range response.FieldList {
		wg.Add(1)
		go func(idx int, f *Field) {
			defer wg.Done()

			fieldSegmentIDs := extractFieldSegmentIDs(strings.Split(Location, ","), f.FieldSegmentList)
			if fieldSegmentIDs == "" {
				return
			}

			signatureParams, err := GenerateNewOrderSignature(ExecDay, fieldSegmentIDs, NetUserId, "1002", VenueId, OpenId, APISecret, APIVersion)
			if err != nil {
				return
			}
			orderURL := "https://web.xports.cn/aisports-api/wechatAPI/order/newOrder?" + signatureParams

			WorkerChanWg.Add(1)
			select {
			case WorkerChan <- OrderRequest{URL: orderURL}:
			case <-GCtx.Done():
				WorkerChanWg.Done()
			}
		}(i, field)
	}
	wg.Wait()
}

func fetchFieldListWithHTTP() ([]byte, error) {
	signatureParams, err := GenerateFieldListSignature(ExecDay, NetUserId, VenueId, "1002", OpenId, APISecret, APIVersion)
	if err != nil {
		return nil, err
	}

	requestURL := "https://web.xports.cn/aisports-api/wechatAPI/venue/fieldList?" + signatureParams

	req, err := http.NewRequestWithContext(GCtx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	setRequestHeaders(req)

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type Response struct {
	Message string `json:"message"`
}

// 配置常量
const (
	APIKey    = "e98ce2565b09ecc0"
	CenterID  = "50030001"
	TenantID  = "82"
	ChannelID = "11"
)

var (
	VenueId    string
	FieldType  string
	APISecret  string
	APIVersion int
)

type KeyValue struct {
	Key   string
	Value string
}

type SignatureResult struct {
	APIKey    string
	Timestamp int64
	ChannelID string
	CenterID  string
	TenantID  string
	OpenId    string
	Version   int
	Sign      string
	Params    map[string]interface{}
}

func md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func generateSignature(apiPath string, params map[string]any, apiSecret string, version int) (*SignatureResult, error) {
	return generateSignatureWithTimestamp(apiPath, params, apiSecret, version, 0)
}

func toURLParams(result *SignatureResult) string {
	var params []string

	params = append(params, "apiKey="+url.QueryEscape(result.APIKey))
	params = append(params, "timestamp="+url.QueryEscape(strconv.FormatInt(result.Timestamp, 10)))
	params = append(params, "channelId="+url.QueryEscape(result.ChannelID))

	if _, hasFieldInfo := result.Params["fieldInfo"]; hasFieldInfo {
		if venueId, ok := result.Params["venueId"]; ok {
			params = append(params, "venueId="+url.QueryEscape(fmt.Sprintf("%v", venueId)))
		}
		if serviceId, ok := result.Params["serviceId"]; ok {
			params = append(params, "serviceId="+url.QueryEscape(fmt.Sprintf("%v", serviceId)))
		}
		if result.CenterID != "" {
			params = append(params, "centerId="+url.QueryEscape(result.CenterID))
		}
		if day, ok := result.Params["day"]; ok {
			params = append(params, "day="+url.QueryEscape(fmt.Sprintf("%v", day)))
		}
		if fieldType, ok := result.Params["fieldType"]; ok {
			params = append(params, "fieldType="+url.QueryEscape(fmt.Sprintf("%v", fieldType)))
		}
		if fieldInfo, ok := result.Params["fieldInfo"]; ok {
			params = append(params, "fieldInfo="+url.QueryEscape(fmt.Sprintf("%v", fieldInfo)))
		}
		if ticket, ok := result.Params["ticket"]; ok {
			params = append(params, "ticket="+url.QueryEscape(fmt.Sprintf("%v", ticket)))
		}
		if randStr, ok := result.Params["randStr"]; ok {
			params = append(params, "randStr="+url.QueryEscape(fmt.Sprintf("%v", randStr)))
		}
		if netUserId, ok := result.Params["netUserId"]; ok {
			params = append(params, "netUserId="+url.QueryEscape(fmt.Sprintf("%v", netUserId)))
		}
		if result.TenantID != "" {
			params = append(params, "tenantId="+url.QueryEscape(result.TenantID))
		}
	} else {
		if netUserId, ok := result.Params["netUserId"]; ok {
			params = append(params, "netUserId="+url.QueryEscape(fmt.Sprintf("%v", netUserId)))
		}
		if venueId, ok := result.Params["venueId"]; ok {
			params = append(params, "venueId="+url.QueryEscape(fmt.Sprintf("%v", venueId)))
		}
		if serviceId, ok := result.Params["serviceId"]; ok {
			params = append(params, "serviceId="+url.QueryEscape(fmt.Sprintf("%v", serviceId)))
		}
		if day, ok := result.Params["day"]; ok {
			params = append(params, "day="+url.QueryEscape(fmt.Sprintf("%v", day)))
		}
		if selectByfullTag, ok := result.Params["selectByfullTag"]; ok {
			params = append(params, "selectByfullTag="+url.QueryEscape(fmt.Sprintf("%v", selectByfullTag)))
		}
		if result.CenterID != "" {
			params = append(params, "centerId="+url.QueryEscape(result.CenterID))
		}
		if fieldType, ok := result.Params["fieldType"]; ok {
			params = append(params, "fieldType="+url.QueryEscape(fmt.Sprintf("%v", fieldType)))
		}
		if result.TenantID != "" {
			params = append(params, "tenantId="+url.QueryEscape(result.TenantID))
		}
	}

	params = append(params, "openId="+url.QueryEscape(result.OpenId))
	params = append(params, "version="+strconv.Itoa(result.Version))
	params = append(params, "sign="+url.QueryEscape(result.Sign))

	return strings.Join(params, "&")
}

func GenerateFieldListSignature(day, netUserID, venueID, serviceID, openId, apiSecret string, version int) (string, error) {
	apiPath := "/aisports-api/wechatAPI/venue/fieldList"
	params := map[string]any{
		"netUserId":       netUserID,
		"venueId":         venueID,
		"serviceId":       serviceID,
		"day":             day,
		"selectByfullTag": "0",
		"fieldType":       FieldType,
		"openId":          openId,
	}

	result, err := generateSignature(apiPath, params, apiSecret, version)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

func GenerateNewOrderSignature(day, fieldInfo, netUserID, serviceID, venueID, openId, apiSecret string, version int) (string, error) {
	apiPath := "/aisports-api/wechatAPI/order/newOrder"
	params := map[string]any{
		"serviceId": serviceID,
		"day":       day,
		"fieldType": FieldType,
		"fieldInfo": fieldInfo,
		"ticket":    "",
		"randStr":   "",
		"venueId":   venueID,
		"netUserId": netUserID,
		"openId":    openId,
	}

	result, err := generateSignature(apiPath, params, apiSecret, version)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

func GenerateFieldListSignatureWithTimestamp(day, netUserID, venueID, serviceID, openId, apiSecret string, version int, timestamp int64, fieldType string) (string, error) {
	apiPath := "/aisports-api/wechatAPI/venue/fieldList"
	params := map[string]any{
		"netUserId":       netUserID,
		"venueId":         venueID,
		"serviceId":       serviceID,
		"day":             day,
		"selectByfullTag": "0",
		"fieldType":       fieldType,
		"openId":          openId,
	}

	result, err := generateSignatureWithTimestamp(apiPath, params, apiSecret, version, timestamp)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

func generateSignatureWithTimestamp(apiPath string, params map[string]any, apiSecret string, version int, customTimestamp int64) (*SignatureResult, error) {
	if apiSecret == "" || version <= 0 {
		return nil, fmt.Errorf("invalid params")
	}

	var timestamp int64
	if customTimestamp > 0 {
		timestamp = customTimestamp
	} else {
		timestamp = time.Now().UnixMilli()
	}

	result := &SignatureResult{
		APIKey:    APIKey,
		Timestamp: timestamp,
		ChannelID: ChannelID,
		Version:   version,
		Params:    make(map[string]any),
	}

	for k, v := range params {
		result.Params[k] = v
	}

	if _, exists := result.Params["centerId"]; !exists {
		result.CenterID = CenterID
	}

	result.OpenId = result.Params["openId"].(string)
	result.TenantID = TenantID

	signParams := make(map[string]any)
	signParams["apiKey"] = result.APIKey
	signParams["timestamp"] = result.Timestamp
	signParams["channelId"] = result.ChannelID
	if result.CenterID != "" {
		signParams["centerId"] = result.CenterID
	}
	if result.TenantID != "" {
		signParams["tenantId"] = result.TenantID
	}
	signParams["openId"] = result.OpenId
	signParams["version"] = result.Version

	for k, v := range result.Params {
		signParams[k] = v
	}

	var keyValues []KeyValue
	for k, v := range signParams {
		keyValues = append(keyValues, KeyValue{Key: k, Value: fmt.Sprintf("%v", v)})
	}

	sort.Slice(keyValues, func(i, j int) bool {
		return keyValues[i].Key < keyValues[j].Key
	})

	var paramStr strings.Builder
	for _, kv := range keyValues {
		paramStr.WriteString(kv.Key)
		paramStr.WriteString("=")
		paramStr.WriteString(kv.Value)
	}

	signString := apiPath + paramStr.String() + apiSecret
	encodedString := url.QueryEscape(signString)

	encodedString = strings.ReplaceAll(encodedString, "(", "%28")
	encodedString = strings.ReplaceAll(encodedString, ")", "%29")
	encodedString = strings.ReplaceAll(encodedString, "'", "%27")
	encodedString = strings.ReplaceAll(encodedString, "!", "%21")
	encodedString = strings.ReplaceAll(encodedString, "~", "%7E")

	result.Sign = md5Hash(encodedString)

	return result, nil
}

func main() {
	var (
		times   string
		startAt string
	)
	flag.StringVar(&ExecDay, "day", "", "天数格式： 20250901")
	flag.StringVar(&NetUserId, "net_user_id", "", "账号")
	flag.StringVar(&OpenId, "open_id", "", "openId")
	flag.StringVar(&APISecret, "api_secret", "", "API密钥")
	flag.IntVar(&APIVersion, "version", 0, "签名版本")
	flag.StringVar(&times, "times", "5", "执行次数")
	flag.StringVar(&startAt, "start", "", "开始时间格式 2025-01-01 00:59:59")
	flag.StringVar(&Location, "location", "", "位置（1-10）")
	flag.StringVar(&VenueIdIndex, "venue_id_index", "", "场馆")
	flag.Int64Var(&SuccessExitCount, "ok_count", 1, "收到多少次成功响应后退出")
	flag.Parse()

	if ExecDay == "" ||
		NetUserId == "" ||
		Location == "" ||
		APISecret == "" ||
		OpenId == "" ||
		APIVersion <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	maxAttempts, err := strconv.Atoi(times)
	if err != nil || maxAttempts <= 0 {
		log.Println("错误: 最大执行次数必须是正整数")
		os.Exit(1)
	}

	if SuccessExitCount <= 0 {
		SuccessExitCount = 1
	}

	switch VenueIdIndex {
	case "2":
		VenueId = "5003000103"
		FieldType = "1837"
	default:
		VenueId = "5003000101"
		FieldType = "1841"
	}

	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
	if shanghaiLoc != nil {
		time.Local = shanghaiLoc
	}

	// 初始化 HTTP 客户端
	HttpClient = createHTTPClient()

	GCtx, GCancel = context.WithCancel(context.Background())
	defer GCancel()

	// 初始化 worker
	WorkerChan = make(chan OrderRequest, 500)
	WorkerChanWg = &sync.WaitGroup{}
	for range NumWorkers {
		go orderWorker()
	}

	// 预热连接
	warmupConnection()

	if startAt != "" {
		start, err := time.ParseInLocation(time.DateTime, startAt, shanghaiLoc)
		if err != nil {
			log.Println("时间格式错误")
			return
		}
		now := time.Now()
		if !now.Before(start) {
			log.Println("指定时间已过")
			return
		}
		targetTime := start.Add(-time.Duration(WarmupAdvanceMs) * time.Millisecond)
		sub := targetTime.Sub(now)
		log.Printf("等待 %.2f 秒后开始...\n", sub.Seconds())

		timer := time.NewTimer(sub)
		select {
		case <-timer.C:
		case <-GCtx.Done():
			timer.Stop()
			return
		}
	}

	log.Printf("开始执行，最大尝试次数: %d\n", maxAttempts)
	var data []byte
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-GCtx.Done():
			goto End
		default:
		}

		if atomic.LoadInt64(&GlobalSuccessCount) >= SuccessExitCount {
			break
		}

		data, err = fetchFieldListWithHTTP()
		if err != nil {
			log.Println("[fieldList] error", err)
			time.Sleep(RetryDelay)
			continue
		}

		var response APIResponse
		if err = json.Unmarshal(data, &response); err != nil {
			time.Sleep(RetryDelay)
			continue
		}
		if len(response.FieldList) > 0 {
			processFieldList(&response)
		} else {
			time.Sleep(RetryDelay)
		}
	}

End:
	WorkerChanWg.Wait()
	close(WorkerChan)
	log.Printf("执行完成，成功次数: %d, %s\n", atomic.LoadInt64(&GlobalSuccessCount), string(data))
}
