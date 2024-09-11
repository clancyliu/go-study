package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: 8000})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("hello world"))
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	m, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:m]))
}
