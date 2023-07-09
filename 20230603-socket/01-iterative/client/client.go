// go:build linux && amd64

package main

import (
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/pkg/ioutil"
)

func main() {
	clientfd, err := ioutil.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer ioutil.Close(clientfd)

	if err := ioutil.Connect(clientfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil {
		log.Printf("connect failed: %v\n", err)
		return
	}

	recvbuf := make([]byte, 1024)
	var tn int

	for {
		rn, err := ioutil.Read(clientfd, recvbuf)
		if err != nil {
			log.Printf("read failed: %v\n", err)
			return
		}

		if rn <= 0 {
			break
		}

		tn += rn
	}

	log.Printf("read %d bytes:%v\n", tn, string(recvbuf[:tn]))
}
