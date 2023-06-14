// go:build linux && arm64
// intro/daytimetcpsvr.c

package main

import (
	"log"
	"syscall"
	"time"
)

func main() {
	listenfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer syscall.Close(listenfd)

	if err := syscall.Bind(listenfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil {
		log.Printf("bind failed: %v\n", err)
		return
	}

	if err := syscall.Listen(listenfd, 1024); err != nil {
		log.Printf("listen failed: %v\n", err)
		return
	}

	for {
		connfd, _, err := syscall.Accept(listenfd)
		if err != nil {
			log.Printf("accept failed: %v\n", err)
			continue
		}

		log.Printf("Accepted a connection")

		now := time.Now()
		buf := now.Format("2006-01-02 15:04:05")

		if _, err := syscall.Write(connfd, []byte(buf)); err != nil {
			log.Printf("write failed: %v\n", err)
			syscall.Close(connfd)
			return
		}

		syscall.Close(connfd)
	}
}
