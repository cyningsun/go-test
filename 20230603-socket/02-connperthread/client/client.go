//go:build linux && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/pkg/ioutil"
	"github.com/cyningsun/go-test/20230603-socket/pkg/proto"
	"github.com/cyningsun/go-test/20230603-socket/pkg/sockaddr"
)

// IP address args from input
var addr string

func main() {
	flag.StringVar(&addr, "addr", "", "ip address")
	flag.Parse()

	if addr == "" {
		log.Fatal("invalid ip address")
	}

	clientfd, err := ioutil.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer ioutil.Close(clientfd)

	sa, err := sockaddr.Parse(addr)
	if err != nil {
		log.Printf("invalid ip address: %v\n", err)
		return
	}

	if err := ioutil.Connect(clientfd, sa); err != nil {
		log.Printf("connect failed: %v\n", err)
		return
	}

	for {
		args := &proto.Args{}
		n, err := fmt.Scanf("%d %d", &args.Args1, &args.Args2)
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

		ioutil.Write(clientfd, bytesBuffer.Bytes())

		recvbuf := make([]byte, 1024)
		ret := &proto.Result{}

		size := binary.Size(*ret)
		tn, rn := 0, 0
		for tn, rn = 0, 0; tn < size; tn += rn {
			var err error
			rn, err = ioutil.Read(clientfd, recvbuf[tn:])
			if err != nil {
				log.Printf("read failed: %v\n", err)
				return
			}

			if rn <= 0 {
				break
			}
		}

		if tn == 0 {
			log.Printf("server terminated\n")
			return
		}

		if err = binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, ret); err != nil {
			log.Printf("binary read failed: %v\n", err)
			return
		}

		log.Printf("expect: %d, actual: %d\n", args.Args1+args.Args2, ret.Sum)
	}
}
