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

	ioutil.SetNonblock(listenfd, true)
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

	epfd, e := ioutil.EpollCreate1(0)
	if e != nil {
		fmt.Println("epoll_create1 failed: ", e)
		os.Exit(1)
	}
	defer ioutil.Close(epfd)

	if err = ioutil.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, listenfd, &syscall.EpollEvent{
		Fd:     int32(listenfd),
		Events: syscall.EPOLLIN,
	}); e != nil {
		fmt.Println("epoll_ctl failed: ", e)
		return
	}

	for {
		evts := make([]syscall.EpollEvent, MAX_OPEN)

		nready, err := ioutil.EpollWait(epfd, evts, -1)
		if err != nil {
			log.Printf("epollwait failed: %v\n", err)
			return
		}

		for i := 0; i < nready; i++ {
			switch {
			case evts[i].Fd == int32(listenfd):
				connfd, _, err := ioutil.Accept(listenfd)
				if err != nil {
					log.Printf("accept failed: %v\n", err)
					continue
				}

				if err = ioutil.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connfd, &syscall.EpollEvent{
					Fd:     int32(connfd),
					Events: syscall.EPOLLIN,
				}); e != nil {
					fmt.Println("epoll_ctl failed: ", e)
					evts[i].Fd = -1
					continue
				}

				log.Printf("Accepted a connection")

			default:
				args := &proto.Args{}
				size := binary.Size(*args)
				recvbuf := make([]byte, 1024)

				var err error
				for tn, rn := 0, 0; tn < size && err == nil; tn += rn {
					rn, err = ioutil.Read(int(evts[i].Fd), recvbuf)
					if err != nil {
						log.Printf("read failed: %v\n", err)
						ioutil.Close(int(evts[i].Fd))
						evts[i].Fd = -1
						break
					}

					if rn <= 0 {
						break
					}
				}

				if err != nil {
					continue
				}

				if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
					log.Printf("binary read failed: %v\n", err)
					ioutil.Close(int(evts[i].Fd))
					evts[i].Fd = -1
					continue
				}

				ret := &proto.Result{Sum: args.Args1 + args.Args2}
				buf := bytes.NewBuffer([]byte{})
				if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
					log.Printf("binary write failed: %v\n", err)
					ioutil.Close(int(evts[i].Fd))
					evts[i].Fd = -1
					continue
				}

				_, err = ioutil.Write(int(evts[i].Fd), buf.Bytes())
				if err != nil {
					log.Printf("write failed: %v\n", err)
					ioutil.Close(int(evts[i].Fd))
					evts[i].Fd = -1
					continue
				}
			}
		}
	}
}
