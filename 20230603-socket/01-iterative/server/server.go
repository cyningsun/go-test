// go:build linux && amd64
// intro/daytimetcpsvr.c

package main

import (
	"log"
	"syscall"
	"time"

	"github.com/cyningsun/go-test/20230603-socket/pkg/ioutil"
)

func main() {
	listenfd, err := ioutil.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer ioutil.Close(listenfd)

	if err := ioutil.Bind(listenfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil {
		log.Printf("bind failed: %v\n", err)
		return
	}

	if err := ioutil.Listen(listenfd, 1024); err != nil {
		log.Printf("listen failed: %v\n", err)
		return
	}

	for {
		connfd, _, err := ioutil.Accept(listenfd)
		if err != nil {
			log.Printf("accept failed: %v\n", err)
			continue
		}

		log.Printf("Accepted a connection")

		now := time.Now()
		buf := now.Format("2006-01-02 15:04:05")

		if _, err := ioutil.Write(connfd, []byte(buf)); err != nil {
			log.Printf("write failed: %v\n", err)
			ioutil.Close(connfd)
			return
		}

		ioutil.Close(connfd)
	}
}
