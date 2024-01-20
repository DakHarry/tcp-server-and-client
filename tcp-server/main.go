package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type TaskMgr struct {
	processedTaskCount int
	maxConcurrency     int
	mu                 sync.Mutex
}

func (mgr *TaskMgr) addTask() {
	mgr.mu.Lock()
	mgr.processedTaskCount += 1
	mgr.mu.Unlock()
}

func main() {
	fmt.Println("Start TCP Server...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// shutdownSignal := make(chan struct{})
	serverWg := new(sync.WaitGroup)
	mgr := &TaskMgr{
		processedTaskCount: 0,
		maxConcurrency:     100,
	}

	go shutdownHandler(ln, c, serverWg)

	serverWg.Add(1)
	go handler(ctx, ln, serverWg, mgr)

	serverWg.Wait()
	fmt.Printf("All tasks are done. Total task is %v\n", mgr.processedTaskCount)
}

func shutdownHandler(ln net.Listener, c chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()

	<-c
	fmt.Println("Received interrupt signal. Shutting down...")

	// close(shutdownSignal)

	// Close the listener to stop accepting new connections
	_ = ln.Close()
}

func handler(ctx context.Context, ln net.Listener, wg *sync.WaitGroup, mgr *TaskMgr) {
	for {
		fmt.Println("[TCP] waiting..")
		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				fmt.Println("Listener closed. Exiting.")
				break
			}
		}

		fmt.Println("Dispatch worker...")
		wg.Add(1)
		go handleMessage(ctx, conn, mgr, wg)
	}
}

func handleMessage(ctx context.Context, conn net.Conn, mgr *TaskMgr, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		conn.Close()
	}()

	fmt.Println("Received message...")
	mgr.addTask()
	fmt.Printf("Task Manager:: count %v\n", mgr.processedTaskCount)
	buffer := make([]byte, 1024)
	for {
		data, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Read buffer error:", err)
			return
		}

		fmt.Printf("[Server Received] %s\n", buffer[:data])
		waitTime := rand.Intn(10)
		time.Sleep(time.Second * time.Duration(waitTime))
		conn.Write(buffer[:data])
		fmt.Printf("[Server Response] %s\n", buffer[:data])
	}
}
