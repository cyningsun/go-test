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
	"golang.org/x/sys/unix"
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

	const (
		MAX_OPEN = 1024
	)
	var i, maxi int

	client := make([]unix.PollFd, 0, MAX_OPEN)
	client = append(client, unix.PollFd{
		Fd:     int32(listenfd),
		Events: unix.POLLIN,
	})
	for i = 1; i < MAX_OPEN; i++ {
		client = append(client, unix.PollFd{})
		client[i].Fd = -1
	}
	maxi = 0

	for {

		nready, err := unix.Poll(client, -1)
		if err != nil {
			log.Printf("poll failed: %v\n", err)
			return
		}

		if (client[0].Revents & unix.POLLIN) != 0 {
			connfd, _, err := ioutil.Accept(listenfd)
			if err != nil {
				log.Printf("accept failed: %v\n", err)
				continue
			}

			log.Printf("accepted a connection")

			for i = 0; i < MAX_OPEN; i++ {
				if client[i].Fd < 0 {
					client[i].Fd = int32(connfd)
					break
				}
			}

			if i == MAX_OPEN {
				log.Printf("too many clients\n")
				ioutil.Close(connfd)
				client[i].Fd = -1
				continue
			}

			client[i].Events = unix.POLLIN

			if i > maxi {
				maxi = i
			}

			nready = nready - 1
			if nready <= 0 {
				continue
			}
		}

		for i := 1; i <= maxi; i++ {
			if client[i].Fd < 0 {
				continue
			}

			if (client[i].Revents & (unix.POLLIN | unix.POLLERR)) != 0 {
				args := &proto.Args{}
				size := binary.Size(*args)
				recvbuf := make([]byte, 1024)

				var err error
				tn, rn := 0, 0
				for tn, rn = 0, 0; tn < size && err == nil; tn += rn {
					rn, err = ioutil.Read(int(client[i].Fd), recvbuf)
					if err != nil {
						break
					}

					if rn <= 0 {
						break
					}
				}

				if tn == 0 || err == ioutil.ECONNRESET {
					ioutil.Close(int(client[i].Fd))
					client[i].Fd = -1
					log.Printf("client closed\n")
					continue
				}

				if err != nil {
					log.Printf("read failed: %v\n", err)
					ioutil.Close(int(client[i].Fd))
					client[i].Fd = -1
				}

				if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
					log.Printf("binary read failed: %v\n", err)
					ioutil.Close(int(client[i].Fd))
					client[i].Fd = -1
					continue
				}

				ret := &proto.Result{Sum: args.Args1 + args.Args2}
				buf := bytes.NewBuffer([]byte{})
				if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
					log.Printf("binary write failed: %v\n", err)
					ioutil.Close(int(client[i].Fd))
					client[i].Fd = -1
					continue
				}

				_, err = ioutil.Write(int(client[i].Fd), buf.Bytes())
				if err != nil {
					log.Printf("write failed: %v\n", err)
					ioutil.Close(int(client[i].Fd))
					client[i].Fd = -1
					continue
				}
			}

			nready = nready - 1
			if nready <= 0 {
				continue
			}
		}
	}
}
