package main

import (
	"fmt"
	"net"
	"os"
)

func handleClient(conn net.Conn, messages chan<- string) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		data := buffer[:n]
		fmt.Printf("Received data from %s: %s", conn.RemoteAddr(), data)

		// Send the message to the channel
		messages <- fmt.Sprintf("Message from %s: %s", conn.RemoteAddr(), data)
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on 127.0.0.1:8888")

	connections := make([]net.Conn, 0)
	messages := make(chan string, 10) // Buffered channel

	for {
		// Handle new connections outside of the select block
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		connections = append(connections, conn)

		go handleClient(conn, messages)

		// Use select to handle message broadcasting
		select {
		case message := <-messages:
			for _, conn := range connections {
				_, err := conn.Write([]byte(message))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
				}
			}
		default:
			// Do nothing if there are no messages to send
		}
	}
}
