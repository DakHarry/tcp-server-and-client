package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Start TCP Client...")

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Printf("Server error: %s", err.Error())
		return
	}
	defer conn.Close()

	msg := []byte("Hi Harry!")
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("Client Error: ", err)
		return
	}

	buf := make([]byte, 1024)
	resp, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Server Error: ", err)
	}
	fmt.Printf("Message: %s", buf[:resp])

}
