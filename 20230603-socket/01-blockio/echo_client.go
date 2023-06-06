//go:build linux && arm64

package main

import (
	"fmt"
	"log"
	"syscall"
)

func main() {
	clientfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer syscall.Close(clientfd)

	if err := syscall.Connect(clientfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil {
		log.Printf("connect failed: %v\n", err)
		return
	}

	sendbuf := make([]byte, 1024)
	recvbuf := make([]byte, 1024)
	for {
		n, err := fmt.Scanln(&sendbuf)
		if err != nil {
			log.Printf("scanln failed: %v\n", err)
			return
		}

		wn, err := syscall.Write(clientfd, sendbuf[:n])
		if err != nil {
			log.Printf("write failed: %v\n", err)
			return
		}

		log.Printf("write %d bytes:%v\n", wn, string(sendbuf[:wn]))

		rn, err := syscall.Read(clientfd, recvbuf)
		if err != nil {
			log.Printf("read failed: %v\n", err)
			return
		}

		log.Printf("read %d bytes:%v\n", rn, string(recvbuf[:rn]))
	}
}
