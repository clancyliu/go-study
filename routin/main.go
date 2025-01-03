package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	noBufferCh := make(chan int)
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case noBufferCh <- 1:
				fmt.Println("noBufferCh")
			default:
				fmt.Println("default")
			}
			time.Sleep(time.Second)
		}
	}()
	<-ticker.C
	go func() {
		fmt.Println(<-noBufferCh)
	}()

	time.Sleep(time.Minute)

	ctx := context.Background()
	// 超时控制
	context.WithTimeout(ctx, time.Second)
	// 取消信号
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// k-v传值
	valueCtx := context.WithValue(ctx, "key", "value")

	deadline, c := context.WithDeadline(ctx, time.Now().Add(time.Minute))

}
