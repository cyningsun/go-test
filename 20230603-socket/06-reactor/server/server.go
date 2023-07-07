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

const (
	MAX_OPEN = 1024
)

// IP address args from input
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

	r, err := newReactor(listenfd)
	if err != nil {
		log.Printf("new reactor failed: %v\n", err)
		return
	}
	defer r.Close()

	acceptor := func(listenfd int) (int, error) {
		connfd, _, err := ioutil.Accept(listenfd)
		if err != nil {
			log.Printf("accept failed: %v\n", err)
			return -1, err
		}

		log.Printf("Accepted a connection")

		if err := ioutil.SetNonblock(connfd, true); err != nil {
			log.Printf("set nonblock failed: %v\n", err)
			return -1, err
		}

		return connfd, nil
	}

	handler := func(connfd int) error {
		args := &proto.Args{}
		size := binary.Size(*args)
		recvbuf := make([]byte, 1024)

		var err error
		var tn, rn int
		for tn, rn = 0, 0; tn < size && err == nil; tn += rn {
			rn, err = ioutil.Read(connfd, recvbuf)
			if err != nil {
				log.Printf("read failed: %v\n", err)
				return err
			}

			if rn <= 0 {
				break
			}
		}

		if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
			log.Printf("binary read failed: %v\n", err)
			return err
		}

		ret := &proto.Result{Sum: args.Args1 + args.Args2}
		buf := bytes.NewBuffer([]byte{})
		if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
			log.Printf("binary write failed: %v\n", err)
			return err
		}

		_, err = ioutil.Write(connfd, buf.Bytes())
		if err != nil {
			log.Printf("write failed: %v\n", err)
			return err
		}

		return nil
	}

	r.RegisterAcceptor(acceptor)
	r.RegisterHandler(handler)

	for {
		if err := r.EpollWait(); err != nil {
			log.Printf("epoll wait failed: %v\n", err)
			return
		}
	}
}

type Reactor struct {
	listenfd int
	epfd     int
	acceptor func(listenfd int) (int, error)
	handler  func(connfd int) error
}

func newReactor(listenfd int) (*Reactor, error) {
	epfd, err := ioutil.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	if err = ioutil.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, listenfd, &syscall.EpollEvent{
		Fd:     int32(listenfd),
		Events: syscall.EPOLLIN,
	}); err != nil {
		fmt.Println("epoll_ctl failed: ", err)
		return nil, err
	}

	return &Reactor{
		listenfd: listenfd,
		epfd:     epfd,
	}, nil
}

func (r *Reactor) RegisterAcceptor(acceptor func(listenfd int) (int, error)) {
	r.acceptor = acceptor
}

func (r *Reactor) RegisterHandler(handler func(connfd int) error) {
	r.handler = handler
}

func (r *Reactor) EpollWait() error {
	events := make([]syscall.EpollEvent, MAX_OPEN)

	nready, err := ioutil.EpollWait(r.epfd, events, -1)
	if err != nil {
		log.Printf("epoll wait failed: %v\n", err)
		return err
	}

	for i := 0; i < nready; i++ {
		if int(events[i].Fd) == r.listenfd {
			connfd, err := r.acceptor(int(events[i].Fd))
			if err != nil {
				return err
			}

			if err = ioutil.EpollCtl(r.epfd, syscall.EPOLL_CTL_ADD, connfd, &syscall.EpollEvent{
				Fd:     int32(connfd),
				Events: syscall.EPOLLIN,
			}); err != nil {
				fmt.Println("epoll_ctl failed: ", err)
				ioutil.Close(connfd)
				return err
			}
		} else {
			err := r.handler(int(events[i].Fd))
			if err != nil {
				ioutil.Close(int(events[i].Fd))
				ioutil.EpollCtl(r.epfd, syscall.EPOLL_CTL_DEL, int(events[i].Fd), nil)
				return err
			}
		}
	}

	return nil
}

func (r *Reactor) EpollCtl(fd int, op int, event *syscall.EpollEvent) error {
	return ioutil.EpollCtl(r.epfd, op, fd, event)
}

func (r *Reactor) Close() error {
	return ioutil.Close(r.epfd)
}
