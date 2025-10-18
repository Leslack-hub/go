package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
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
	RetryDelay   = 100 * time.Millisecond
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
	useTestData  = false
	workerChan   chan string
	workerChanWg *sync.WaitGroup
	gCtx         context.Context
	gCancel      context.CancelFunc
	execDay      string
	location     string
	netUserId    string
	venueIdIndex string
)

func main() {
	var (
		times   string
		startAt string
	)
	flag.StringVar(&execDay, "day", "", "天数格式： 20250901")
	flag.StringVar(&netUserId, "net_user_id", "", "账号")
	flag.StringVar(&times, "times", "5", "执行次数")
	flag.StringVar(&startAt, "start", "", "开始时间格式 2025-01-01 00:59:59")
	flag.StringVar(&location, "location", "", "位置（1-10）")
	flag.StringVar(&venueIdIndex, "venue_id_index", "", "场馆")
	flag.Parse()
	if execDay == "" || netUserId == "" || location == "" {
		showUsage()
		os.Exit(1)
	}

	maxAttempts, err := strconv.Atoi(times)
	if err != nil || maxAttempts <= 0 {
		log.Println("错误: 最大执行次数必须是正整数")
		os.Exit(1)
	}

	switch venueIdIndex {
	case "2":
		VenueId = "5003000103"
		FieldType = "1837"
	default:
		VenueId = "5003000101"
		FieldType = "1841"
	}

	if err = checkDependencies(); err != nil {
		log.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	var shanghaiLoc *time.Location
	shanghaiLoc, err = time.LoadLocation("Asia/Shanghai")
	if err == nil {
		time.Local = shanghaiLoc
	}

	gCtx, gCancel = context.WithCancel(context.Background())
	defer gCancel()

	workerChan = make(chan string)
	workerChanWg = &sync.WaitGroup{}
	for range 30 {
		go func() {
			for cmd := range workerChan {
				Run(cmd, 3, 1)
				workerChanWg.Done()
			}
		}()
	}

	if useTestData {
		if _, err = os.Stat(TestJSONFile); os.IsNotExist(err) {
			log.Printf("错误: 找不到测试数据文件 %s\n", TestJSONFile)
			os.Exit(1)
		}
		log.Println("注意: 使用测试数据模式")
	} else {
		log.Println("注意: 使用实际HTTP请求模式")
	}

	if startAt != "" {
		var start time.Time
		start, err = time.ParseInLocation(time.DateTime, startAt, shanghaiLoc)
		if err != nil {
			log.Println("时间格式错误")
			return
		}
		now := time.Now()
		if !now.Before(start) {
			return
		}
		sub := start.Add(500 * time.Millisecond).Sub(now)
		log.Println("sleep time:", sub.Seconds())
		time.Sleep(sub)
	}

	log.Printf("开始执行，最大尝试次数: %d\n", maxAttempts)
	log.Println("----------------------------------------")

Attempts:
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-gCtx.Done():
			log.Println("Context cancelled, stopping attempts.")
			break Attempts
		default:
		}
		log.Printf("第 %d 次尝试，正在获取场地列表...\n", attempt)

		var response APIResponse
		var data []byte

		if useTestData {
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
			data, err = fetchFieldListWithCurl()
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
			log.Println("✓ 成功获取场地列表，正在处理数据...")

			if err = processFieldList(&response); err != nil {
				log.Printf("✗ 处理场地列表失败: %v\n", err)
			}
			log.Printf("响应内容: %s\n", data)
		} else {
			log.Printf("✗ 第 %d 次尝试失败：fieldList为空\n", attempt)
			log.Printf("等待 %v 后重试...\n", RetryDelay)
			time.Sleep(RetryDelay)
		}
	}

	workerChanWg.Wait()
	close(workerChan)
	fmt.Println("----------------------------------------")
	fmt.Println("脚本执行完成")
}

func showUsage() {
	flag.Usage()
}

func checkDependencies() error {
	if _, err := exec.LookPath("curl"); err != nil {
		return fmt.Errorf("需要安装 curl 命令")
	}

	return nil
}

func extractFieldSegmentIDs(segmentList []*FieldSegment) string {
	var fieldSegmentIDs []string
	locations := strings.Split(location, ",")
	if len(locations) == 0 {
		return ""
	}
	l1, err := strconv.Atoi(locations[0])
	if err == nil && l1 <= len(segmentList) && l1 > 0 {
		l1 = l1 - 1
		fmt.Println("位置：", l1)
		if segmentList[l1].State == "0" && segmentList[l1].Price == 0 && segmentList[l1].FieldSegmentID != "" {
			fieldSegmentIDs = append(fieldSegmentIDs, segmentList[l1].FieldSegmentID)
		}
	}

	if len(locations) >= 2 {
		var l2 int
		l2, err = strconv.Atoi(locations[1])
		if err == nil && l2 <= len(segmentList) && l2 > 0 {
			l2 = l2 - 1
			fmt.Println("位置：", l2)
			if segmentList[l2].State == "0" && segmentList[l2].Price == 0 && segmentList[l2].FieldSegmentID != "" {
				fieldSegmentIDs = append(fieldSegmentIDs, segmentList[l2].FieldSegmentID)
			}
		}
	}
	if len(fieldSegmentIDs) == 0 {
		return ""
	}

	return strings.Join(fieldSegmentIDs, ",")
}

func processFieldList(response *APIResponse) error {
	fieldCount := len(response.FieldList)
	log.Printf("找到 %d 个场地\n", fieldCount)
	wg := sync.WaitGroup{}
	rand.Shuffle(fieldCount, func(i, j int) {
		response.FieldList[i], response.FieldList[j] = response.FieldList[j], response.FieldList[i]
	})
	for i, field := range response.FieldList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("处理第 %d 个场地...\n", i+1)
			fieldSegmentIDs := extractFieldSegmentIDs(field.FieldSegmentList)
			if fieldSegmentIDs != "" {
				log.Printf("  提取到的fieldSegmentIds: %s\n", fieldSegmentIDs)
				//if rand.IntN(10) < 3 {
				//	go func() {
				//		time.Sleep(1 * time.Second)
				//		workerChan <- "echo '{\"message\":\"ok\"}'"
				//	}()
				//} else {
				// 使用Go实现的签名生成器
				signatureParams, err := GenerateNewOrderSignature(execDay, fieldSegmentIDs, netUserId, "1002", VenueId)
				if err != nil {
					log.Printf("生成newOrder签名失败: %v", err)
					return
				}
				workerChanWg.Add(1)
				workerChan <- fmt.Sprintf(`curl -s "https://web.xports.cn/aisports-api/wechatAPI/order/newOrder?%s" -H 'Host: web.xports.cn' -H 'Connection: keep-alive' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080712) UnifiedPCMacWechat(0xf2641015) XWEB/16390' -H 'xweb_xhr: 1' -H 'Accept: */*' -H 'Sec-Fetch-Site: cross-site' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Dest: empty' -H 'Referer: https://servicewechat.com/wxb75b9974eac7896e/11/page-frame.html' -H 'Accept-Language: zh-CN,zh;q=0.9' -H 'Content-Type: application/json'`, signatureParams)
				//}
			} else {
				log.Println("  未找到有效的场地时段ID")
			}
		}()
	}
	wg.Wait()
	return nil
}

func fetchFieldListWithCurl() ([]byte, error) {
	signatureParams, err := GenerateFieldListSignature(execDay, netUserId, VenueId, "1002")
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	// 跨平台shell命令执行
	curlCommand := fmt.Sprintf(`curl -s "https://web.xports.cn/aisports-api/wechatAPI/venue/fieldList?%s" -H 'Host: web.xports.cn' -H 'Connection: keep-alive' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080712) UnifiedPCMacWechat(0xf2641015) XWEB/16390' -H 'xweb_xhr: 1' -H 'Accept: */*' -H 'Sec-Fetch-Site: cross-site' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Dest: empty' -H 'Referer: https://servicewechat.com/wxb75b9974eac7896e/11/page-frame.html' -H 'Accept-Language: zh-CN,zh;q=0.9' -H 'Content-Type: application/json'`, signatureParams)
	curlCmd := exec.Command("sh", "-c", curlCommand)
	var output []byte
	output, err = curlCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("curl命令执行失败: %v", err)
	}
	return output, nil
}

type Response struct {
	Message string `json:"message"`
}

type Worker struct {
	command    string
	maxExec    int64
	execCount  *int64
	ctx        context.Context
	cancel     context.CancelFunc
	cancelOnce sync.Once
}

func (w *Worker) checkJSONResponse(output []byte) {
	log.Println("exec result: ", string(output))
	var result Response
	if err := json.Unmarshal(output, &result); err != nil {
		return
	}

	if result.Message == "ok" || result.Message == "场地预定中，请勿重复提交" {
		log.Println("Success detected in JSON output, exiting program...")
		w.cancelOnce.Do(func() {
			w.cancel()
		})
	}
}

func (w *Worker) executeCommand(workerID int) error {
	interval := 100 * time.Millisecond

	for {
		current := atomic.AddInt64(w.execCount, 1)
		if current > w.maxExec {
			log.Printf("Execution limit (%d) reached, exiting...", w.maxExec)
			w.cancelOnce.Do(func() {
				w.cancel()
			})
			return nil
		}

		cmdCtx, cmdCancel := context.WithTimeout(w.ctx, 2*time.Second)
		cmd := exec.CommandContext(cmdCtx, "sh", "-c", w.command)
		output, err := cmd.CombinedOutput()
		cmdCancel()

		if err != nil {
			log.Printf("Worker %d execution %d error: %v", workerID, current, err)
		} else {
			go w.checkJSONResponse(output)
		}

		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		case <-time.After(interval):
		}
	}
}

func Run(command string, maxExec int64, numWorkers int) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received signal, shutting down gracefully...")
		gCancel()
	}()

	var execCount int64
	worker := &Worker{
		command:   command,
		maxExec:   maxExec,
		execCount: &execCount,
		ctx:       gCtx,
		cancel:    gCancel,
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := range numWorkers {
		go func(workerID int) {
			defer wg.Done()
			if err2 := worker.executeCommand(workerID); err2 != nil &&
				!errors.Is(err2, context.Canceled) {
				log.Printf("Worker %d error: %v", workerID, err2)
			}
		}(i)
	}

	wg.Wait()
	log.Println("All workers finished")
}

// 配置常量
const (
	APIKey    = "e98ce2565b09ecc0"
	APISecret = "b28efc98ae90a878"
	CenterID  = "50030001"
	TenantID  = "82"
	ChannelID = "11"
	//VenueId   = "5003000103"
	//FieldType = "1837"
	//VenueId   = "5003000101"
	//FieldType = "1841"
)

var VenueId string

var FieldType string

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
func generateSignature(apiPath string, params map[string]any, options *SignatureOptions) (*SignatureResult, error) {
	return generateSignatureWithTimestamp(apiPath, params, options, 0)
}

// generateSignatureWithTimestamp 生成签名，支持自定义时间戳（用于测试）
func generateSignatureWithTimestamp(apiPath string, params map[string]any, options *SignatureOptions, customTimestamp int64) (*SignatureResult, error) {
	if options == nil {
		options = &SignatureOptions{}
	}

	// 获取API密钥和密钥
	apiKey := APIKey
	apiSecret := APISecret
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
		if fieldType, ok := result.Params["fieldType"]; ok {
			params = append(params, fmt.Sprintf("fieldType=%s", url.QueryEscape(fmt.Sprintf("%v", fieldType))))
		}
	}

	// 添加centerId和tenantId
	if result.CenterID != "" {
		params = append(params, fmt.Sprintf("centerId=%s", url.QueryEscape(result.CenterID)))
	}
	if result.TenantID != "" {
		params = append(params, fmt.Sprintf("tenantId=%s", url.QueryEscape(result.TenantID)))
	}

	// 最后添加签名
	params = append(params, fmt.Sprintf("sign=%s", url.QueryEscape(result.Sign)))

	return strings.Join(params, "&")
}

// GenerateFieldListSignature 生成fieldList签名
func GenerateFieldListSignature(day, netUserID, venueID, serviceID string) (string, error) {
	apiPath := "/aisports-api/wechatAPI/venue/fieldList"
	params := map[string]any{
		"netUserId":       netUserID,
		"venueId":         venueID,
		"serviceId":       serviceID,
		"day":             day,
		"selectByfullTag": "0",
		"fieldType":       FieldType,
	}

	result, err := generateSignature(apiPath, params, nil)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

// GenerateNewOrderSignature 生成newOrder签名
func GenerateNewOrderSignature(day, fieldInfo, netUserID, serviceID, venueID string) (string, error) {
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

	result, err := generateSignature(apiPath, params, nil)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

// GenerateFieldListSignatureWithTimestamp 生成fieldList签名（测试用，支持固定时间戳）
func GenerateFieldListSignatureWithTimestamp(day, netUserID, venueID, serviceID string, timestamp int64) (string, error) {
	apiPath := "/aisports-api/wechatAPI/venue/fieldList"
	params := map[string]any{
		"netUserId":       netUserID,
		"venueId":         venueID,
		"serviceId":       serviceID,
		"day":             day,
		"selectByfullTag": "0",
		"fieldType":       FieldType,
	}

	result, err := generateSignatureWithTimestamp(apiPath, params, nil, timestamp)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}

// GenerateNewOrderSignatureWithTimestamp 生成newOrder签名（测试用，支持固定时间戳）
func GenerateNewOrderSignatureWithTimestamp(day, fieldInfo, netUserID, serviceID, venueID string, timestamp int64) (string, error) {
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

	result, err := generateSignatureWithTimestamp(apiPath, params, nil, timestamp)
	if err != nil {
		return "", err
	}

	return toURLParams(result), nil
}
