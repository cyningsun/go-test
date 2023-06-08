//go:build linux && arm64

package main

import (
	"log"
	"syscall"
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
		defer syscall.Close(connfd)

		log.Printf("Accepted a connection")

		go func() {
			buf := make([]byte, 1024)
			ok := true

			for ok {
				n, err := syscall.Read(connfd, buf)
				if err != nil {
					log.Printf("read failed: %v\n", err)
					return
				}

				if n == 0 {
					continue
				}

				log.Printf("read %d bytes:%v\n", n, string(buf[:n]))

				syscall.Write(connfd, buf[:n]) // echo
			}
		}()
	}
}
