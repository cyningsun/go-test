//go:build linux && arm64

package main

import (
	"log"
	"syscall"
	"time"
)

func main() {
	clientfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}

	if err := syscall.SetNonblock(clientfd, true); err != nil {
		log.Printf("set nonblock failed: %v\n", err)
		return
	}

	if err := syscall.Connect(clientfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil && err != syscall.EINPROGRESS {
		log.Printf("connect failed: %v\n", err)
		return
	}

	var tn int
	msg := []byte("hello world")
	for {
		n, err := syscall.Write(clientfd, msg[tn:])
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				time.Sleep(time.Millisecond)
				continue
			}

			log.Printf("write failed: %v\n", err)
			return
		}

		log.Printf("write %d bytes\n", n)

		tn += n
		if tn == len([]byte("hello world")) {
			break
		}
	}
}
