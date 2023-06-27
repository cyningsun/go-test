//go:build linux && arm64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/pkg/fdset"
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

	listenfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}
	defer syscall.Close(listenfd)

	sa, err := sockaddr.Parse(addr)
	if err != nil {
		log.Printf("invalid ip address: %v\n", err)
		return
	}

	syscall.SetsockoptInt(listenfd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	if err := syscall.Bind(listenfd, sa); err != nil {
		log.Printf("bind failed: %v\n", err)
		return
	}

	if err := syscall.Listen(listenfd, 1024); err != nil {
		log.Printf("listen failed: %v\n", err)
		return
	}

	const (
		MAX_OPEN = 1024
	)
	var (
		client  [MAX_OPEN]unix.PollFd
		i, maxi int
	)

	client[0].Fd = int32(listenfd)
	client[0].Events = unix.POLLRDNORM
	for i = 1; i < MAX_OPEN; i++ {
		client[i].Fd = -1
	}
	maxi = 0

	for {

		nready, err := unix.Poll(client, maxi+1, unix.INFTIM)
		if err != nil {
			log.Printf("select failed: %v\n", err)
			return
		}

		if (client[0].Revents & unix.POLLRDNORM) != 0 {
			connfd, _, err := syscall.Accept(listenfd)
			if err != nil {
				log.Printf("accept failed: %v\n", err)
				continue
			}

			log.Printf("Accepted a connection")

			for i = 0; i < MAX_OPEN; i++ {
				if client[i].Fd < 0 {
					client[i].Fd = connfd
					break
				}
			}

			if i == MAX_OPEN {
				log.Printf("too many clients\n")
				syscall.Close(connfd)
				return
			}

			client[i].Events = unix.POLLRDNORM

			fdset.Set(&allset, connfd)

			if i > maxi {
				maxi = i
			}

			nready = nready - 1
			if nready <= 0 {
				continue
			}
		}

		for i := 0; i <= maxi; i++ {
			if client[i].Fd < 0 {
				continue
			}

			if (client[i].Revents & (unix.POLLRDNORM | unix.POLLERR)) != 0 {
				args := &proto.Args{}
				size := binary.Size(*args)
				recvbuf := make([]byte, 1024)

				for tn, rn := 0, 0; tn < size; tn += rn {
					var err error
					rn, err = syscall.Read(client[i].Fd, recvbuf)
					if err != nil {
						log.Printf("read failed: %v\n", err)
						syscall.Close(client[i].Fd)
						client[i].Fd = -1
					}

					if rn <= 0 {
						break
					}
				}

				if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
					log.Printf("binary read failed: %v\n", err)
					syscall.Close(client[i].Fd)
					client[i].Fd = -1
				}

				ret := &proto.Result{Sum: args.Args1 + args.Args2}
				buf := bytes.NewBuffer([]byte{})
				if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
					log.Printf("binary write failed: %v\n", err)
					syscall.Close(client[i].Fd)
					client[i].Fd = -1
				}

				_, err = syscall.Write(client[i].Fd, buf.Bytes())
				if err != nil {
					log.Printf("write failed: %v\n", err)
					syscall.Close(client[i].Fd)
					client[i].Fd = -1
				}
			}

			nready = nready - 1
			if nready <= 0 {
				continue
			}
		}
	}
}
