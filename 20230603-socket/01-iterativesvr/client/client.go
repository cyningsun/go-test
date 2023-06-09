// go:build linux && arm64
// intro/daytimetcpcli.c

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
	defer syscall.Close(clientfd)

	if err := syscall.Connect(clientfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil {
		log.Printf("connect failed: %v\n", err)
		return
	}

	recvbuf := make([]byte, 1024)
	var rn int

	for {
		rn, err = syscall.Read(clientfd, recvbuf)
		if err != nil {
			log.Printf("read failed: %v\n", err)
			return
		}

		if rn <= 0 {
			break
		}

		log.Printf("read %d bytes:%v\n", rn, string(recvbuf[:rn]))
	}
}
