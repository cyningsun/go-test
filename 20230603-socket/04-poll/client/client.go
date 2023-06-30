//go:build linux && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/pkg/proto"
	"github.com/cyningsun/go-test/20230603-socket/pkg/sockaddr"
	"golang.org/x/sys/unix"
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

	sa, err := sockaddr.Parse(addr)
	if err != nil {
		log.Printf("invalid ip address: %v\n", err)
		return
	}

	if err := syscall.Connect(clientfd, sa); err != nil {
		log.Printf("connect failed: %v\n", err)
		return
	}

	stdeof := false
	const (
		MAX_OPEN = 1024
	)

	args := &proto.Args{}
	ret := &proto.Result{}
	for {
		client := []unix.PollFd{
			{
				Fd:     int32(clientfd),
				Events: unix.POLLIN,
			},
		}

		if !stdeof {
			client = append(client, unix.PollFd{
				Fd:     int32(syscall.Stdin),
				Events: unix.POLLIN,
			})
		}

		client = append(client)

		_, err := unix.Poll(client, -1)
		if err != nil {
			log.Printf("select failed: %v\n", err)
			return
		}

		switch {
		case len(client) == 2 && (client[1].Revents&(unix.POLLIN|unix.POLLERR)) != 0:
			n, err := fmt.Scanf("%d %d", &args.Args1, &args.Args2)
			if err != nil {
				log.Printf("scanf failed: %v\n", err)
				return
			}

			if n == 0 {
				stdeof = true
				client[1].Fd = -1
				syscall.Shutdown(clientfd, syscall.SHUT_WR)
				continue
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
		case (client[0].Revents & (unix.POLLIN | unix.POLLERR)) != 0:
			recvbuf := make([]byte, 1024)
			size, tn := binary.Size(*ret), 0
			for rn := 0; tn < size; tn += rn {
				var err error
				rn, err = syscall.Read(clientfd, recvbuf[tn:])
				if err != nil {
					log.Printf("read failed: %v\n", err)
					return
				}

				if rn <= 0 {
					break
				}
			}

			if tn == 0 {
				if stdeof {
					return
				} else {
					fmt.Printf("server terminated\n")
					return // server terminated
				}
			}

			if err = binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, ret); err != nil {
				log.Printf("binary read failed: %v\n", err)
				return
			}

			fmt.Printf("expect: %d, actual: %d\n", args.Args1+args.Args2, ret.Sum)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
