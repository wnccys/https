package main

import (
	"net"
	"strings"
	"fmt"
	"strconv"
	"bufio"
)

const PORT = 8080;

func main() {
	listener, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(PORT), 10))
	if err != nil {
		panic(err)
	}
	// drops listene on final of fn
	defer listener.Close();

	fmt.Printf("Listening on %d...", PORT);

	for {
		// accept incoming requests
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	parts := strings.Fields(requestLine)
	if len(parts) < 3 {
		return // Invalid HTTP
	}
	
	method, path := parts[0], parts[1]

	fmt.Printf("[%s] %s\n", method, path)

	// Response construction
	// HTTP/1.1 200 OK\r\n
	// Content-Type: text/plain\r\n
	// Content-Length: 12\r\n
	// \r\n
	// Body...

	var body string
	var status string

	if method == "GET" && path == "/" {
		status = "200 OK"
		body = "Hello from Raw Go!"
	} else {
		status = "404 Not Found"
		body = "Resource not found"
	}

	response := fmt.Sprintf(
		"HTTP/1.1 %s\r\n"+
		"Content-Type: text/plain\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n"+
		"%s",
		status, len(body), body,
	)

	conn.Write([]byte(response))
}