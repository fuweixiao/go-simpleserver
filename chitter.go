package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Parse command line arguments
	port := "9999"
	if len(os.Args) == 1 {
		println("Usage: port number required! \nExample: go run chitter.go 9999")
		os.Exit(0)
	} else {
		port = os.Args[1]
	}
	// Listen for incoming connections.

	l, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error()+port)
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + "localhost" + ":" + "3333")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// Parse data
	data := buf[0 : reqLen-1]
	msg := string(data)
	id := strings.Split(msg, ":")[0]
	// Send messages
	if _, err := strconv.Atoi(id); err == nil {
		print(id)
	} else if id == "whoami" {
		print("Who cares who the fuck you are!")
	} else if id == "all" {
		print("all")
	}

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}
