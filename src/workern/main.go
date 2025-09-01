package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Resp struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

// readJDFile 读取jd.txt文件内容
func readJDFile() (string, error) {
	fileByte, err := os.ReadFile("jd.txt")
	if err != nil {
		return "", err
	}
	return string(fileByte), nil
}

func worker2() {
	if len(os.Args) < 2 {
		return
	}
	long, err1 := strconv.Atoi(os.Args[1])
	if err1 != nil {
		return
	}

	bash, err1 := readJDFile()
	if err1 != nil {
		log.Println(err1)
		return
	}
	ch := make(chan os.Signal, 1)
	ch2 := make(chan struct{}, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	timeout := time.NewTimer(time.Duration(long) * time.Second)
	go func() {
		<-timeout.C
		fmt.Printf("程序运行%d秒后自动退出", long)
		close(ch2)
	}()
	go func() {
		<-ch
		fmt.Println("收到退出信号，程序立即退出")
		timeout.Stop()
		close(ch2)
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			t := time.Tick((time.Duration((i+1)*100 + rand.Intn(100))) * time.Millisecond)
			// t := time.Tick(2 * time.Second)
		Exit:
			for {
				select {
				case <-ch2:
					break Exit
				case <-t:
					output, err := exec.Command("sh", "-c", bash).CombinedOutput()
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					fmt.Println(string(output))
				}
			}
		}()
	}
	wg.Wait()
}

func worker1() {
	long, err1 := strconv.Atoi(os.Args[1])
	if err1 != nil {
		return
	}

	bash, err1 := readJDFile()
	if err1 != nil {
		log.Println(err1)
		return
	}
	ch := make(chan os.Signal, 1)
	ch2 := make(chan struct{}, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	timeout := time.NewTimer(time.Duration(long) * time.Second)
	go func() {
		<-timeout.C
		fmt.Printf("程序运行%d秒后自动退出", long)
		close(ch2)
	}()
	go func() {
		<-ch
		fmt.Println("收到退出信号，程序立即退出")
		timeout.Stop()
		close(ch2)
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			t := time.Tick((time.Duration((i+1)*100 + rand.Intn(100))) * time.Millisecond)
			//t := time.Tick(2 * time.Second)
		Exit:
			for {
				select {
				case <-ch2:
					break Exit
				case <-t:
					var err error
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

					cmd := exec.CommandContext(ctx, "sh", "-c", bash)

					var stdout io.ReadCloser
					stdout, err = cmd.StdoutPipe()
					if err != nil {
						log.Println("[error]", err)
						cancel()
						continue
					}

					if err = cmd.Start(); err != nil {
						log.Println("[error]", err)
						cancel()
						continue
					}

					go func(cmdCtx context.Context, cmdCancel context.CancelFunc) {
						defer cmdCancel()
						defer func() {
							if cmd.Process != nil {
								if err2 := cmd.Process.Kill(); err2 != nil {
									log.Println("进程终止错误", err2)
								}
								_ = cmd.Wait()
							}
						}()

						scanner := bufio.NewScanner(stdout)
						for scanner.Scan() {
							bytes := scanner.Bytes()
							log.Println(string(bytes))
							if !json.Valid(bytes) {
								continue
							}

							var ret *Resp
							_ = json.Unmarshal(bytes, &ret)
							if ret.Success {
								log.Println("请求成功")
								close(ch2)
								return
							}
						}
					}(ctx, cancel)
				}
			}
		}()
	}
	wg.Wait()
}

func main() {
	if len(os.Args) < 2 {
		worker2()
	} else {
		worker1()
	}
}
