package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":5003")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("Listening on port 5003")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("New client connected")

	// Read alllll
	go func() {
		buf := make([]byte, 512)

		for {
			n, err := conn.Read(buf[0:])
			if err != nil {
				break
			}

			fmt.Printf(string(buf[0:n]))
		}
	}()

	// Start shell
	scanner := bufio.NewScanner(os.Stdin)
	var cmd string

	for {
		scanner.Scan()
		cmd = scanner.Text()
		if cmd != "" {
			conn.Write([]byte(fmt.Sprintf("%s\n", cmd)))
		}
	}
}
