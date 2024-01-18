package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type TaskMgr struct {
	proccessedTaskCount int
	maxConcurrency      int
	mu                  sync.Mutex
}

func (mgr *TaskMgr) Count() {
	mgr.mu.Lock()
	mgr.proccessedTaskCount += 1
	mgr.mu.Unlock()
}

func main() {
	fmt.Println("Start TCP Server...")
	ctx, cancel := context.WithCancel(context.Background())
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer ln.Close()

	mgr := &TaskMgr{
		proccessedTaskCount: 0,
		maxConcurrency:      100,
	}

	wg := new(sync.WaitGroup)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			break
		}
		if mgr.proccessedTaskCount == 3 {
			break
		}
		wg.Add(1)
		go handleMessage(ctx, conn, mgr, wg)
	}
	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(time.Second)
	fmt.Println("Main Goroutine exited")
	wg.Wait()
	fmt.Printf("All tasks are done. Total task is %v\n", mgr.proccessedTaskCount)
}

func handleMessage(ctx context.Context, conn net.Conn, mgr *TaskMgr, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		conn.Close()
	}()

	fmt.Println("Received message...")
	mgr.Count()
	fmt.Printf("Task Manager:: count %v\n", mgr.proccessedTaskCount)
	buffer := make([]byte, 1024)
	for {
		// fmt.Println(buffer)
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

// shutdown TODO:: graceful shutdown
func shutdown() {

}

func broadcast() {}
