package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 8000})
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	conn, err := listener.AcceptTCP()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf[0:])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[0:n]))

	time.Sleep(100 * time.Second)

	m, err := conn.Write([]byte("hello world"))
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}
