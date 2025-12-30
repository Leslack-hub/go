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
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	TestJSONFile = "test.json"
	// 优化：减少重试延迟，抢票时时间宝贵
	RetryDelay = 10 * time.Millisecond
	// 优化：增加并发 worker 数量
	NumWorkers = 50
	// 优化：每个场地发起的请求次数
	MaxExecPerField = 2
	// 预热提前时间（毫秒）
	WarmupAdvanceMs = 100
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
	UseTestData      = false
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
	// 优化：全局 HTTP 客户端，启用连接池和 Keep-Alive
	HttpClient *http.Client
	// 优化：成功计数器
	GlobalSuccessCount int64
)

// OrderRequest 用于传递下单请求信息
type OrderRequest struct {
	URL string
}

// 优化：创建高性能 HTTP 客户端
func createHTTPClient() *http.Client {
	// 自定义传输配置，优化连接池
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		// 优化：增加最大连接数
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
		// 优化：禁用压缩以减少 CPU 开销
		DisableCompression: true,
		// 优化：启用 HTTP/2
		ForceAttemptHTTP2: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		// 优化：减少握手超时
		TLSHandshakeTimeout: 3 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
}

// 优化：预热连接，提前建立 TCP 连接
func warmupConnection() {
	// 发送一个轻量级请求来预热连接
	req, err := http.NewRequest("HEAD", "https://web.xports.cn/", nil)
	if err != nil {
		return
	}
	req.Header.Set("Connection", "keep-alive")
	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Printf("预热连接失败（可忽略）: %v", err)
		return
	}
	resp.Body.Close()
	log.Println("✓ 连接预热完成")
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
	if ExecDay == "" || NetUserId == "" || Location == "" || APISecret == "" || APIVersion <= 0 {
		showUsage()
		os.Exit(1)
	}

	maxAttempts, err := strconv.Atoi(times)
	if err != nil || maxAttempts <= 0 {
		log.Println("错误: 最大执行次数必须是正整数")
		os.Exit(1)
	}

	if SuccessExitCount <= 0 {
		log.Println("错误: 成功退出次数必须是正整数")
		os.Exit(1)
	}

	switch VenueIdIndex {
	case "2":
		VenueId = "5003000103"
		FieldType = "1837"
	default:
		VenueId = "5003000101"
		FieldType = "1841"
	}

	var shanghaiLoc *time.Location
	shanghaiLoc, err = time.LoadLocation("Asia/Shanghai")
	if err == nil {
		time.Local = shanghaiLoc
	}

	// 优化：初始化高性能 HTTP 客户端
	HttpClient = createHTTPClient()

	GCtx, GCancel = context.WithCancel(context.Background())
	defer GCancel()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("收到终止信号，正在优雅退出...")
		GCancel()
	}()

	// 优化：增加 worker 数量
	WorkerChan = make(chan OrderRequest, 1000) // 增加缓冲区
	WorkerChanWg = &sync.WaitGroup{}
	for range NumWorkers {
		go orderWorker()
	}

	if UseTestData {
		if _, err = os.Stat(TestJSONFile); os.IsNotExist(err) {
			log.Printf("错误: 找不到测试数据文件 %s\n", TestJSONFile)
			os.Exit(1)
		}
		log.Println("注意: 使用测试数据模式")
	} else {
		log.Println("注意: 使用实际HTTP请求模式（原生HTTP客户端）")
	}

	// 优化：预热连接
	warmupConnection()

	if startAt != "" {
		var start time.Time
		start, err = time.ParseInLocation(time.DateTime, startAt, shanghaiLoc)
		if err != nil {
			log.Println("时间格式错误")
			return
		}
		now := time.Now()
		if !now.Before(start) {
			log.Println("指定时间已过")
			return
		}
		// 优化：提前少量时间开始，考虑网络延迟
		advanceTime := time.Duration(WarmupAdvanceMs) * time.Millisecond
		targetTime := start.Add(-advanceTime)
		sub := targetTime.Sub(now)
		log.Printf("等待 %.2f 秒后开始（提前 %dms 启动）...\n", sub.Seconds(), WarmupAdvanceMs)

		// 使用高精度定时器
		timer := time.NewTimer(sub)
		select {
		case <-timer.C:
		case <-GCtx.Done():
			timer.Stop()
			return
		}
	}

	log.Printf("🚀 开始执行，最大尝试次数: %d，并发 Worker: %d\n", maxAttempts, NumWorkers)
	log.Println("----------------------------------------")

	startTime := time.Now()

Attempts:
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-GCtx.Done():
			log.Println("Context cancelled, stopping attempts.")
			break Attempts
		default:
		}

		// 检查是否已达到成功次数
		if atomic.LoadInt64(&GlobalSuccessCount) >= SuccessExitCount {
			log.Printf("✓ 已达到成功次数 %d，停止尝试\n", SuccessExitCount)
			break Attempts
		}

		log.Printf("第 %d 次尝试，正在获取场地列表...\n", attempt)

		var response APIResponse
		var data []byte

		if UseTestData {
			data, err = os.ReadFile(TestJSONFile)
			if err != nil {
				log.Printf("✗ 第 %d 次尝试失败：无法读取测试数据文件: %v\n", attempt, err)
				if attempt == maxAttempts {
					log.Printf("已达到最大尝试次数 (%d)，停止执行\n", maxAttempts)
					os.Exit(1)
				}
				time.Sleep(RetryDelay)
				continue
			}
		} else {
			data, err = fetchFieldListWithHTTP()
			fmt.Println(string(data))
			os.Exit(1)
			if err != nil {
				log.Printf("✗ 第 %d 次尝试失败：获取数据失败: %v\n", attempt, err)
				time.Sleep(RetryDelay)
				continue
			}
		}

		if err = json.Unmarshal(data, &response); err != nil {
			log.Printf("✗ 第 %d 次尝试失败：JSON解析错误: %v\n", attempt, err)
			time.Sleep(RetryDelay)
			continue
		}

		if len(response.FieldList) > 0 {
			log.Printf("✓ 成功获取场地列表（%d个场地），正在处理数据...\n", len(response.FieldList))

			if err = processFieldList(&response); err != nil {
				log.Printf("✗ 处理场地列表失败: %v\n", err)
			}
		} else {
			log.Printf("✗ 第 %d 次尝试失败：fieldList为空（error=%d, message=%s）\n",
				attempt, response.Error, response.Message)
			time.Sleep(RetryDelay)
		}
	}

	// 等待所有下单请求完成
	WorkerChanWg.Wait()
	close(WorkerChan)

	elapsed := time.Since(startTime)
	fmt.Println("----------------------------------------")
	fmt.Printf("脚本执行完成，耗时: %.2f秒，成功次数: %d\n", elapsed.Seconds(), atomic.LoadInt64(&GlobalSuccessCount))
}

// 优化：使用原生 HTTP 客户端的 worker
func orderWorker() {
	for req := range WorkerChan {
		executeOrder(req)
		WorkerChanWg.Done()
	}
}

// 优化：使用原生 HTTP 执行下单请求
func executeOrder(orderReq OrderRequest) {
	for i := 0; i < MaxExecPerField; i++ {
		select {
		case <-GCtx.Done():
			return
		default:
		}

		// 检查是否已达到成功次数
		if atomic.LoadInt64(&GlobalSuccessCount) >= SuccessExitCount {
			return
		}

		req, err := http.NewRequestWithContext(GCtx, "GET", orderReq.URL, nil)
		if err != nil {
			continue
		}

		// 设置请求头
		setRequestHeaders(req)

		resp, err := HttpClient.Do(req)
		if err != nil {
			log.Printf("下单请求失败: %v", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		// 检查响应
		checkOrderResponse(body)
	}
}

// 设置请求头
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

// 检查下单响应
func checkOrderResponse(body []byte) {
	log.Printf("下单响应: %s", string(body))

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}

	if result.Message == "ok" {
		count := atomic.AddInt64(&GlobalSuccessCount, 1)
		log.Printf("🎉 抢票成功！(%d/%d)", count, SuccessExitCount)
		if count >= SuccessExitCount {
			log.Println("✓ 已达到目标成功次数，停止后续请求")
			GCancel()
		}
	}
}

func showUsage() {
	flag.Usage()
}

func extractFieldSegmentIDs(locations []string, segmentList []*FieldSegment) string {
	if len(locations) == 0 {
		return ""
	}
	// 可用时段索引 -> ID
	available := make(map[int]string)
	for i, segment := range segmentList {
		if segment.State == "0" && segment.Price == 0 && segment.FieldSegmentID != "" {
			available[i] = segment.FieldSegmentID
		}
	}
	if len(available) == 0 {
		return ""
	}

	// 以 l1 为中心，l1 向左递减、l2 向右递增
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

	// 优先：找到最靠近中心的连续两张（先向左递减，再向右递增）
	for offset := 0; offset < len(segmentList); offset++ {
		startLeft := center - offset
		if withinBounds(startLeft) && withinBounds(startLeft+1) {
			if id1, ok1 := available[startLeft]; ok1 {
				if id2, ok2 := available[startLeft+1]; ok2 {
					return strings.Join([]string{id1, id2}, ",")
				}
			}
		}

		startRight := rightStart + offset
		if withinBounds(startRight) && withinBounds(startRight+1) {
			if id1, ok1 := available[startRight]; ok1 {
				if id2, ok2 := available[startRight+1]; ok2 {
					return strings.Join([]string{id1, id2}, ",")
				}
			}
		}
	}

	// 其次：按左右扩散顺序取最多两张
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

func processFieldList(response *APIResponse) error {
	fieldCount := len(response.FieldList)
	log.Printf("找到 %d 个场地\n", fieldCount)
	wg := sync.WaitGroup{}

	// 优化：随机打乱以分散请求
	rand.Shuffle(fieldCount, func(i, j int) {
		response.FieldList[i], response.FieldList[j] = response.FieldList[j], response.FieldList[i]
	})

	for i, field := range response.FieldList {
		wg.Add(1)
		go func(idx int, f *Field) {
			defer wg.Done()

			fieldSegmentIDs := extractFieldSegmentIDs(strings.Split(Location, ","), f.FieldSegmentList)
			if fieldSegmentIDs != "" {
				log.Printf("场地 %d: 提取到时段ID: %s\n", idx+1, fieldSegmentIDs)

				// 生成签名
				signatureParams, err := GenerateNewOrderSignature(ExecDay, fieldSegmentIDs, NetUserId, "1002", VenueId, OpenId, APISecret, APIVersion)
				if err != nil {
					log.Printf("生成newOrder签名失败: %v", err)
					return
				}
				orderURL := fmt.Sprintf("https://web.xports.cn/aisports-api/wechatAPI/order/newOrder?%s", signatureParams)
				// 发送到 worker 队列
				WorkerChanWg.Add(1)
				select {
				case WorkerChan <- OrderRequest{URL: orderURL}:
				case <-GCtx.Done():
					WorkerChanWg.Done()
				}
			} else {
				log.Printf("场地 %d: 未找到有效的时段ID\n", idx+1)
			}
		}(i, field)
	}
	wg.Wait()
	return nil
}

// 优化：使用原生 HTTP 客户端获取场地列表
func fetchFieldListWithHTTP() ([]byte, error) {
	signatureParams, err := GenerateFieldListSignature(ExecDay, NetUserId, VenueId, "1002", OpenId, APISecret, APIVersion)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	requestURL := fmt.Sprintf("https://web.xports.cn/aisports-api/wechatAPI/venue/fieldList?%s", signatureParams)

	req, err := http.NewRequestWithContext(GCtx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	setRequestHeaders(req)

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	return body, nil
}

type Response struct {
	Message string `json:"message"`
}

//func Run(command string, maxExec int64, successLimit int64, numWorkers int) {
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
//	go func() {
//		<-sigChan
//		log.Println("Received signal, shutting down gracefully...")
//		gCancel()
//	}()
//
//	if successLimit <= 0 {
//		successLimit = 1
//	}
//
//	var execCount int64
//	var successCount int64
//	worker := &Worker{
//		command:      command,
//		maxExec:      maxExec,
//		execCount:    &execCount,
//		successLimit: successLimit,
//		successCount: &successCount,
//		ctx:          gCtx,
//		cancel:       gCancel,
//	}
//
//	var wg sync.WaitGroup
//	wg.Add(numWorkers)
//
//	for i := range numWorkers {
//		go func(workerID int) {
//			defer wg.Done()
//			if err2 := worker.executeCommand(workerID); err2 != nil &&
//				!errors.Is(err2, context.Canceled) {
//				log.Printf("Worker %d error: %v", workerID, err2)
//			}
//		}(i)
//	}
//
//	wg.Wait()
//	log.Println("All workers finished")
//}

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

// KeyValue 键值对结构
type KeyValue struct {
	Key   string
	Value string
}

// SignatureOptions 签名选项
type SignatureOptions struct {
	Prefix     string
	NoCenterID bool
}

// SignatureResult 签名结果
type SignatureResult struct {
	APIKey    string `json:"apiKey"`
	Timestamp int64  `json:"timestamp"`
	ChannelID string `json:"channelId"`
	CenterID  string `json:"centerId,omitempty"`
	TenantID  string `json:"tenantId,omitempty"`
	OpenId    string `json:"openId,omitempty"`
	Version   int    `json:"version"`
	Sign      string `json:"sign"`
	// 动态参数
	Params map[string]interface{} `json:"-"`
}

// md5Hash MD5加密函数
func md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// generateSignature 根据原始JavaScript代码逆向的签名生成函数
func generateSignature(apiPath string, params map[string]any, apiSecret string, version int, options *SignatureOptions) (*SignatureResult, error) {
	return generateSignatureWithTimestamp(apiPath, params, apiSecret, version, options, 0)
}

// generateSignatureWithTimestamp 生成签名，支持自定义时间戳（用于测试）
func generateSignatureWithTimestamp(apiPath string, params map[string]any, apiSecret string, version int, options *SignatureOptions, customTimestamp int64) (*SignatureResult, error) {
	if options == nil {
		options = &SignatureOptions{}
	}

	// 获取API密钥和密钥
	apiKey := APIKey
	if apiSecret == "" {
		return nil, fmt.Errorf("apiSecret is required")
	}
	if version <= 0 {
		return nil, fmt.Errorf("version is required")
	}
	if options.Prefix != "" {
		// 这里可以根据prefix获取不同的key，当前使用默认值
	}

	// 获取时间戳（如果提供了自定义时间戳则使用，否则使用当前时间）
	var timestamp int64
	if customTimestamp > 0 {
		timestamp = customTimestamp
	} else {
		timestamp = time.Now().UnixMilli()
	}

	// 构建基础参数对象
	result := &SignatureResult{
		APIKey:    apiKey,
		Timestamp: timestamp,
		ChannelID: ChannelID,
		Version:   version,
		Params:    make(map[string]any),
	}
	// 添加传入的参数
	for k, v := range params {
		result.Params[k] = v
	}

	// 添加centerId（对应原代码逻辑）
	if !options.NoCenterID {
		if _, exists := result.Params["centerId"]; !exists {
			result.CenterID = CenterID
		}
	}

	result.OpenId = result.Params["openId"].(string)
	// 添加tenantId
	result.TenantID = TenantID

	// 构建用于签名的参数映射
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
	if result.OpenId != "" {
		signParams["openId"] = result.OpenId
	}
	signParams["version"] = result.Version

	// 添加业务参数
	for k, v := range result.Params {
		signParams[k] = v
	}

	// 转换为键值对数组
	var keyValues []KeyValue
	for k, v := range signParams {
		keyValues = append(keyValues, KeyValue{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}

	// 按key排序
	sort.Slice(keyValues, func(i, j int) bool {
		return keyValues[i].Key < keyValues[j].Key
	})

	// 拼接参数字符串
	var paramStr strings.Builder
	for _, kv := range keyValues {
		paramStr.WriteString(kv.Key)
		paramStr.WriteString("=")
		paramStr.WriteString(kv.Value)
	}

	// 生成待签名字符串并编码
	signString := apiPath + paramStr.String() + apiSecret
	encodedString := url.QueryEscape(signString)

	// 替换特殊字符（严格按照原代码逻辑）
	if strings.Contains(encodedString, "(") {
		encodedString = strings.ReplaceAll(encodedString, "(", "%28")
	}
	if strings.Contains(encodedString, ")") {
		encodedString = strings.ReplaceAll(encodedString, ")", "%29")
	}
	if strings.Contains(encodedString, "'") {
		encodedString = strings.ReplaceAll(encodedString, "'", "%27")
	}
	if strings.Contains(encodedString, "!") {
		encodedString = strings.ReplaceAll(encodedString, "!", "%21")
	}
	if strings.Contains(encodedString, "~") {
		encodedString = strings.ReplaceAll(encodedString, "~", "%7E")
	}

	// MD5加密
	result.Sign = md5Hash(encodedString)

	return result, nil
}

// toURLParams 将签名结果转换为URL参数字符串
func toURLParams(result *SignatureResult) string {
	// 按照JavaScript版本的确切顺序构建参数
	// JavaScript输出顺序：apiKey, timestamp, channelId, [业务参数], centerId, tenantId, sign
	var params []string

	// 基础参数（固定顺序）
	params = append(params, fmt.Sprintf("apiKey=%s", url.QueryEscape(result.APIKey)))
	params = append(params, fmt.Sprintf("timestamp=%s", url.QueryEscape(strconv.FormatInt(result.Timestamp, 10))))
	params = append(params, fmt.Sprintf("channelId=%s", url.QueryEscape(result.ChannelID)))

	// 业务参数（按照JavaScript中的顺序）
	// fieldList方法顺序：netUserId, venueId, serviceId, day, selectByfullTag, fieldType
	// newOrder方法顺序：serviceId, day, fieldType, fieldInfo, ticket, randStr, venueId, netUserId

	// 检查是否为newOrder方法（包含fieldInfo参数）
	if _, hasFieldInfo := result.Params["fieldInfo"]; hasFieldInfo {
		// newOrder方法的参数顺序
		if serviceId, ok := result.Params["serviceId"]; ok {
			params = append(params, fmt.Sprintf("serviceId=%s", url.QueryEscape(fmt.Sprintf("%v", serviceId))))
		}
		if day, ok := result.Params["day"]; ok {
			params = append(params, fmt.Sprintf("day=%s", url.QueryEscape(fmt.Sprintf("%v", day))))
		}
		if fieldType, ok := result.Params["fieldType"]; ok {
			params = append(params, fmt.Sprintf("fieldType=%s", url.QueryEscape(fmt.Sprintf("%v", fieldType))))
		}
		if fieldInfo, ok := result.Params["fieldInfo"]; ok {
			params = append(params, fmt.Sprintf("fieldInfo=%s", url.QueryEscape(fmt.Sprintf("%v", fieldInfo))))
		}
		if ticket, ok := result.Params["ticket"]; ok {
			params = append(params, fmt.Sprintf("ticket=%s", url.QueryEscape(fmt.Sprintf("%v", ticket))))
		}
		if randStr, ok := result.Params["randStr"]; ok {
			params = append(params, fmt.Sprintf("randStr=%s", url.QueryEscape(fmt.Sprintf("%v", randStr))))
		}
		if venueId, ok := result.Params["venueId"]; ok {
			params = append(params, fmt.Sprintf("venueId=%s", url.QueryEscape(fmt.Sprintf("%v", venueId))))
		}
		if netUserId, ok := result.Params["netUserId"]; ok {
			params = append(params, fmt.Sprintf("netUserId=%s", url.QueryEscape(fmt.Sprintf("%v", netUserId))))
		}
	} else {
		// fieldList方法的参数顺序
		if netUserId, ok := result.Params["netUserId"]; ok {
			params = append(params, fmt.Sprintf("netUserId=%s", url.QueryEscape(fmt.Sprintf("%v", netUserId))))
		}
		if venueId, ok := result.Params["venueId"]; ok {
			params = append(params, fmt.Sprintf("venueId=%s", url.QueryEscape(fmt.Sprintf("%v", venueId))))
		}
		if serviceId, ok := result.Params["serviceId"]; ok {
			params = append(params, fmt.Sprintf("serviceId=%s", url.QueryEscape(fmt.Sprintf("%v", serviceId))))
		}
		if day, ok := result.Params["day"]; ok {
			params = append(params, fmt.Sprintf("day=%s", url.QueryEscape(fmt.Sprintf("%v", day))))
		}
		if selectByfullTag, ok := result.Params["selectByfullTag"]; ok {
			params = append(params, fmt.Sprintf("selectByfullTag=%s", url.QueryEscape(fmt.Sprintf("%v", selectByfullTag))))
		}
		if result.CenterID != "" {
			params = append(params, fmt.Sprintf("centerId=%s", url.QueryEscape(result.CenterID)))
		}
		if fieldType, ok := result.Params["fieldType"]; ok {
			params = append(params, fmt.Sprintf("fieldType=%s", url.QueryEscape(fmt.Sprintf("%v", fieldType))))
		}
		if result.TenantID != "" {
			params = append(params, fmt.Sprintf("tenantId=%s", url.QueryEscape(result.TenantID)))
		}
	}

	// 对于 newOrder，添加 centerId 和 tenantId（如果还没添加）
	if _, hasFieldInfo := result.Params["fieldInfo"]; hasFieldInfo {
		if result.CenterID != "" {
			params = append(params, fmt.Sprintf("centerId=%s", url.QueryEscape(result.CenterID)))
		}
		if result.TenantID != "" {
			params = append(params, fmt.Sprintf("tenantId=%s", url.QueryEscape(result.TenantID)))
		}
	}

	if result.OpenId != "" {
		params = append(params, fmt.Sprintf("openId=%s", url.QueryEscape(result.OpenId)))
	}
	params = append(params, fmt.Sprintf("version=%d", result.Version))
	// 最后添加签名
	params = append(params, fmt.Sprintf("sign=%s", url.QueryEscape(result.Sign)))

	return strings.Join(params, "&")
}

// GenerateFieldListSignature 生成fieldList签名
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

	result, err := generateSignature(apiPath, params, apiSecret, version, nil)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

// GenerateNewOrderSignature 生成newOrder签名
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

	result, err := generateSignature(apiPath, params, apiSecret, version, nil)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

// GenerateFieldListSignatureWithTimestamp 生成fieldList签名（测试用，支持固定时间戳）
func GenerateFieldListSignatureWithTimestamp(day, netUserID, venueID, serviceID, openId, apiSecret string, version int, timestamp int64) (string, error) {
	apiPath := "/aisports-api/wechatAPI/venue/fieldList"
	params := map[string]any{
		"netUserId":       netUserID,
		"venueId":         venueID,
		"serviceId":       serviceID,
		"day":             day,
		"selectByfullTag": "0",
		"fieldType":       "1837",
		"openId":          openId,
	}

	result, err := generateSignatureWithTimestamp(apiPath, params, apiSecret, version, nil, timestamp)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

// GenerateNewOrderSignatureWithTimestamp 生成newOrder签名（测试用，支持固定时间戳）
func GenerateNewOrderSignatureWithTimestamp(day, fieldInfo, netUserID, serviceID, venueID, apiSecret string, version int, timestamp int64) (string, error) {
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
	}

	result, err := generateSignatureWithTimestamp(apiPath, params, apiSecret, version, nil, timestamp)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}
