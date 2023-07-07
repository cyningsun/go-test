//go:build linux && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/pkg/ioutil"
	"github.com/cyningsun/go-test/20230603-socket/pkg/proto"
	"github.com/cyningsun/go-test/20230603-socket/pkg/sockaddr"
)

var addr string

func main() {
	flag.StringVar(&addr, "addr", "", "ip address")
	flag.Parse()

	if addr == "" {
		log.Fatal("invalid ip address")
	}

	listenfd, err := ioutil.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer ioutil.Close(listenfd)

	sa, err := sockaddr.Parse(addr)
	if err != nil {
		log.Printf("invalid ip address: %v\n", err)
		return
	}

	ioutil.SetsockoptInt(listenfd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	if err := ioutil.Bind(listenfd, sa); err != nil {
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

		go func(fd int) {
			defer ioutil.Close(connfd)

			for {
				args := &proto.Args{}
				size := binary.Size(*args)
				recvbuf := make([]byte, 1024)

				for tn, rn := 0, 0; tn < size; tn += rn {
					var err error
					rn, err = ioutil.Read(connfd, recvbuf)
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

				ret := &proto.Result{Sum: args.Args1 + args.Args2}
				buf := bytes.NewBuffer([]byte{})
				if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
					log.Printf("binary write failed: %v\n", err)
					return
				}

				_, _ = ioutil.Write(fd, buf.Bytes())
			}
		}(connfd)
	}
}
