//go:build linux && arm64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/02-connperthread/shared"
)

// IP address args from input
var addr string

func main() {
	flag.StringVar(&addr, "addr", "", "ip address")
	flag.Parse()

	if addr == "" {
		log.Fatal("invalid ip address")
	}

	clientfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer syscall.Close(clientfd)

	sa, err := shared.ToSockaddr(addr)
	if err != nil {
		log.Printf("invalid ip address: %v\n", err)
		return
	}

	if err := syscall.Connect(clientfd, sa); err != nil {
		log.Printf("connect failed: %v\n", err)
		return
	}

	for {
		args := &shared.Args{}
		n, err := fmt.Scanf("%d%d", &args.Args1, &args.Args2)
		if err != nil {
			log.Printf("scanf failed: %v\n", err)
			return
		}

		if n != 2 {
			log.Printf("invalid input")
			return
		}

		bytesBuffer := bytes.NewBuffer([]byte{})
		if err = binary.Write(bytesBuffer, binary.BigEndian, args); err != nil {
			log.Printf("binary write failed: %v\n", err)
			return
		}

		syscall.Write(clientfd, bytesBuffer.Bytes())

		recvbuf := make([]byte, 1024)
		ret := &shared.Result{}

		size := binary.Size(*ret)

		for tn, rn := 0, 0; tn < size; tn += rn {
			rn, err := syscall.Read(clientfd, recvbuf[tn:])
			if err != nil {
				log.Printf("read failed: %v\n", err)
				return
			}

			if rn <= 0 {
				break
			}
		}

		if err = binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, ret); err != nil {
			log.Printf("binary read failed: %v\n", err)
			return
		}
	}
}
