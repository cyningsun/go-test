//go:build linux && arm64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/02-connperthread/shared"
)

var addr string

func main() {
	flag.StringVar(&addr, "addr", "", "ip address")
	flag.Parse()

	if addr == "" {
		log.Fatal("invalid ip address")
	}

	listenfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer syscall.Close(listenfd)

	sa, err := shared.ToSockaddr(addr)
	if err != nil {
		log.Printf("invalid ip address: %v\n", err)
		return
	}

	if err := syscall.Bind(listenfd, sa); err != nil {
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

		go func(fd int) {
			defer syscall.Close(connfd)

			args := &shared.Args{}
			size := binary.Size(*args)
			recvbuf := make([]byte, 1024)

			for tn, rn := 0, 0; tn < size; tn += rn {
				rn, err := syscall.Read(connfd, recvbuf)
				if err != nil {
					log.Printf("read failed: %v\n", err)
					return
				}

				if rn <= 0 {
					break
				}
			}

			if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
				log.Printf("binary read failed: %v\n", err)
				return
			}

			ret := &shared.Result{Sum: args.Args1 + args.Args2}
			bytesBuffer := bytes.NewBuffer([]byte{})
			if err = binary.Write(bytesBuffer, binary.BigEndian, ret); err != nil {
				log.Printf("binary write failed: %v\n", err)
				return
			}

			syscall.Write(connfd, bytesBuffer.Bytes())
		}(connfd)
	}
}
