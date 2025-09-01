package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Response struct {
	Code    int  `json:"code"`
	Success bool `json:"success"`
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
	var result Response
	if err := json.Unmarshal(output, &result); err != nil {
		return
	}

	if result.Success {
		log.Println("Success detected in JSON output, exiting program...")
		w.cancelOnce.Do(func() {
			w.cancel()
		})
	}
}

func (w *Worker) executeCommand(workerID int) error {
	interval := time.Duration((workerID+1)*100+rand.Intn(100)) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		case <-ticker.C:
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
				continue
			}

			log.Printf("Worker %d execution %d: %s", workerID, current, string(output))

			go w.checkJSONResponse(output)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program <command> <max_executions>")
		os.Exit(1)
	}

	command := os.Args[1]
	maxExec, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil || maxExec <= 0 {
		fmt.Println("Error: max_executions must be a positive number")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received signal, shutting down gracefully...")
		cancel()
	}()

	var execCount int64
	worker := &Worker{
		command:   command,
		maxExec:   maxExec,
		execCount: &execCount,
		ctx:       ctx,
		cancel:    cancel,
	}

	var wg sync.WaitGroup
	const numWorkers = 3
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
