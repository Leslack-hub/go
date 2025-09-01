package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	TestJSONFile = "test.json"
	RetryDelay   = 1 * time.Second
	OrderDay     = "20250901"
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
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	maxAttempts, err := strconv.Atoi(os.Args[1])
	if err != nil || maxAttempts <= 0 {
		log.Println("错误: 最大执行次数必须是正整数")
		os.Exit(1)
	}

	if err = checkDependencies(); err != nil {
		log.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
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

	if len(os.Args) == 3 {
		var start time.Time
		start, err = time.ParseInLocation(time.DateTime, os.Args[2], shanghaiLoc)
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

	for attempt := 1; attempt <= maxAttempts; attempt++ {
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
				if attempt == maxAttempts {
					log.Printf("已达到最大尝试次数 (%d)，停止执行\n", maxAttempts)
					os.Exit(1)
				}
				time.Sleep(RetryDelay)
				continue
			}
		}

		if err = json.Unmarshal(data, &response); err != nil {
			log.Printf("✗ 第 %d 次尝试失败：JSON解析错误: %v\n", attempt, err)
			if attempt == maxAttempts {
				log.Printf("已达到最大尝试次数 (%d)，停止执行\n", maxAttempts)
				log.Printf("最后一次响应内容: %s\n", string(data))
				os.Exit(1)
			}
			time.Sleep(RetryDelay)
			continue
		}

		if len(response.FieldList) > 0 {
			log.Println("✓ 成功获取场地列表，正在处理数据...")

			if err = processFieldList(&response); err != nil {
				log.Printf("✗ 处理场地列表失败: %v\n", err)
				os.Exit(1)
			}
			break
		} else {
			log.Printf("✗ 第 %d 次尝试失败：fieldList为空\n", attempt)

			if attempt == maxAttempts {
				log.Printf("已达到最大尝试次数 (%d)，停止执行\n", maxAttempts)
				os.Exit(1)
			} else {
				log.Printf("等待 %v 后重试...\n", RetryDelay)
				time.Sleep(RetryDelay)
			}
		}
	}

	workerChanWg.Wait()
	fmt.Println("----------------------------------------")
	fmt.Println("脚本执行完成")
}

func showUsage() {
	log.Printf("例如: %s 10        # 使用实际HTTP请求\n", os.Args[0])
	log.Println("")
	log.Println("说明: 该脚本会尝试获取场地列表并生成预订命令")
	log.Println("选项:")
}

func checkDependencies() error {
	if _, err := exec.LookPath("curl"); err != nil {
		return fmt.Errorf("需要安装 curl 命令")
	}

	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("需要安装 Node.js")
	}

	if _, err := os.Stat("signature_generator.js"); os.IsNotExist(err) {
		return fmt.Errorf("找不到 signature_generator.js 文件")
	}

	return nil
}

func extractFieldSegmentIDs(segmentList []*FieldSegment) string {
	var fieldSegmentIDs []string

	if len(segmentList) > 3 && segmentList[3].State == "0" && segmentList[3].Price == 0 && segmentList[3].FieldSegmentID != "" {
		fieldSegmentIDs = append(fieldSegmentIDs, segmentList[3].FieldSegmentID)
	}

	if len(segmentList) > 4 && segmentList[4].State == "0" && segmentList[4].Price == 0 && segmentList[4].FieldSegmentID != "" {
		fieldSegmentIDs = append(fieldSegmentIDs, segmentList[4].FieldSegmentID)
	}

	return strings.Join(fieldSegmentIDs, ",")
}

func processFieldList(response *APIResponse) error {
	fieldCount := len(response.FieldList)
	log.Printf("找到 %d 个场地\n", fieldCount)
	wg := sync.WaitGroup{}
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
				workerChanWg.Add(1)
				workerChan <- fmt.Sprintf(`curl -s "https://web.xports.cn/aisports-api/wechatAPI/order/newOrder?$(node signature_generator.js -m newOrder --day=%s --fieldInfo=%s)" -H 'Host: web.xports.cn' -H 'Connection: keep-alive' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080712) UnifiedPCMacWechat(0xf2641015) XWEB/16390' -H 'xweb_xhr: 1' -H 'Accept: */*' -H 'Sec-Fetch-Site: cross-site' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Dest: empty' -H 'Referer: https://servicewechat.com/wxb75b9974eac7896e/11/page-frame.html' -H 'Accept-Language: zh-CN,zh;q=0.9' -H 'Content-Type: application/json'`, OrderDay, fieldSegmentIDs)
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
	curlCmd := exec.Command("sh", "-c", fmt.Sprintf(`curl -s "https://web.xports.cn/aisports-api/wechatAPI/venue/fieldList?$(node signature_generator.js -m fieldList --day=%s)" -H 'Host: web.xports.cn' -H 'Connection: keep-alive' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080712) UnifiedPCMacWechat(0xf2641015) XWEB/16390' -H 'xweb_xhr: 1' -H 'Accept: */*' -H 'Sec-Fetch-Site: cross-site' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Dest: empty' -H 'Referer: https://servicewechat.com/wxb75b9974eac7896e/11/page-frame.html' -H 'Accept-Language: zh-CN,zh;q=0.9' -H 'Content-Type: application/json'`, OrderDay))

	output, err := curlCmd.Output()
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

	if result.Message == "ok" {
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
			log.Printf("Worker %d execution %d: %s", workerID, current, string(output))
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
