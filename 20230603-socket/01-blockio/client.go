//go:build linux && arm64

package main

import (
	"log"
	"syscall"
)

func main() {
	clientfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}

	if err := syscall.Connect(clientfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil {
		log.Printf("connect failed: %v\n", err)
		return
	}

	n, err := syscall.Write(clientfd, []byte("hello world"))
	if err != nil {
		log.Printf("write failed: %v\n", err)
		return
	}

	log.Printf("write %d bytes\n", n)
}
