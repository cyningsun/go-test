//go:build linux && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

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

	evts := make([]syscall.EpollEvent, 0, MAX_OPEN)
	evts = append(evts, syscall.EpollEvent{
		Fd:     int32(listenfd),
		Events: unix.POLLIN,
	})
	for i := 1; i < MAX_OPEN; i++ {
		evts = append(evts, syscall.EpollEvent{})
		evts[i].Fd = -1
	}

	epfd, e := syscall.EpollCreate1(0)
	if e != nil {
		fmt.Println("epoll_create1 failed: ", e)
		os.Exit(1)
	}
	defer syscall.Close(epfd)

	if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, listenfd, &evts[0]); e != nil {
		fmt.Println("epoll_ctl failed: ", e)
		return
	}

	for {

		nready, err := syscall.EpollWait(epfd, evts, -1)
		if err != nil {
			log.Printf("epollwait failed: %v\n", err)
			return
		}

		for i := 0; i <= nready; i++ {
			switch {
			case evts[i].Fd == int32(listenfd):
				connfd, _, err := syscall.Accept(listenfd)
				if err != nil {
					log.Printf("accept failed: %v\n", err)
					continue
				}

				log.Printf("Accepted a connection")

				for i = 0; i < MAX_OPEN; i++ {
					if evts[i].Fd < 0 {
						evts[i].Fd = int32(connfd)
						break
					}
				}

				if i == MAX_OPEN {
					log.Printf("too many clients\n")
					evts[i].Fd = -1
					continue
				}

				evts[i].Events = unix.POLLIN

				if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connfd, &evts[i]); e != nil {
					fmt.Println("epoll_ctl failed: ", e)
					evts[i].Fd = -1
					continue
				}

			default:
				args := &proto.Args{}
				size := binary.Size(*args)
				recvbuf := make([]byte, 1024)

				var err error
				for tn, rn := 0, 0; tn < size && err == nil; tn += rn {
					rn, err = syscall.Read(int(evts[i].Fd), recvbuf)
					if err != nil {
						log.Printf("read failed: %v\n", err)
						syscall.Close(int(evts[i].Fd))
						evts[i].Fd = -1
						break
					}

					if rn <= 0 {
						break
					}
				}

				if err == syscall.ECONNRESET {
					continue
				}

				if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
					log.Printf("binary read failed: %v\n", err)
					syscall.Close(int(evts[i].Fd))
					evts[i].Fd = -1
					continue
				}

				ret := &proto.Result{Sum: args.Args1 + args.Args2}
				buf := bytes.NewBuffer([]byte{})
				if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
					log.Printf("binary write failed: %v\n", err)
					syscall.Close(int(evts[i].Fd))
					evts[i].Fd = -1
					continue
				}

				_, err = syscall.Write(int(evts[i].Fd), buf.Bytes())
				if err != nil {
					log.Printf("write failed: %v\n", err)
					syscall.Close(int(evts[i].Fd))
					evts[i].Fd = -1
					continue
				}
			}
		}
	}
}
