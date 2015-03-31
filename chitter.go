package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// define a struct for Messages
type Message struct {
	src int
	dst int
	msg string
}

// map clientId to tcp connections
var dict = make(map[int]net.Conn)

func main() {
	// Parse command line arguments
	port := "9999"
	if len(os.Args) != 2 {
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
	clientId := 0

	// Handle connections
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		clientId++

		// Map clientId to connections
		dict[clientId] = conn

		// Handle connections in a new goroutine.
		go handleRequest(conn, clientId)
	}
}

// Send messages to other id
func sendMessage(message Message) {
	if message.dst == 0 {
		for k := range dict {
			conn := dict[k]
			conn.Write([]byte(strconv.Itoa(message.src) + ":" + message.msg))
		}
	} else {
		conn := dict[message.dst]
		conn.Write([]byte(strconv.Itoa(message.src) + ":" + message.msg))
	}
	return
}

// Handles incoming requests.
func handleRequest(conn net.Conn, clientId int) {
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {
			continue
		}

		// Parse data
		data := buf[0 : reqLen-1]
		msg := string(data)
		id, message := strings.Split(msg, ":")[0], strings.Split(msg, ":")[1]
		// Send messages
		if idnum, err := strconv.Atoi(id); err == nil {
			msg := Message{src: clientId, dst: idnum, msg: message}
			sendMessage(msg)
		} else if id == "whoami" {
			print("Who cares who the fuck you are!")
		} else if id == "all" {
			msg := Message{src: clientId, dst: 0, msg: message}
			sendMessage(msg)
		}

		// Send a response back to person contacting us.
	}
	conn.Close()
}
