package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {
	fmt.Println("Start TCP Server...")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go handleMessage(conn)
	}

}

func handleMessage(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Received message...")
	buffer := make([]byte, 1024)
	for {
		fmt.Println(buffer)
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
