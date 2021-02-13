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

	// Channel for sending commands from stdin -> client
	sendChan := make(chan string)

	go readCommands(sendChan)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}

		// Could make this threaded, but no point right now...
		handleRequest(conn, sendChan)
	}
}

// readCommands continuously reads from stdin till closed and passes everything to the sendChan
func readCommands(sendChan chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		sendChan <- scanner.Text()
	}
}

func handleRequest(conn net.Conn, sendChan <-chan string) {
	fmt.Println("New client connected")

	closeChan := make(chan struct{})

	// Pipe shell output to console
	go func() {
		buf := make([]byte, 512)

		for {
			n, err := conn.Read(buf[0:])
			if err != nil {
				fmt.Println()
				fmt.Println("Client disconnected")
				close(closeChan)
				break
			}

			fmt.Printf(string(buf[0:n]))
		}
	}()

	// Pipe commands to client
	for {
		select {
		case cmd := <-sendChan:
			conn.Write([]byte(fmt.Sprintf("%s\n", cmd)))
		case <-closeChan:
			// Break out of the loop
			return
		}
	}
}
