//go:build linux && arm64

package main

import (
	"log"
	"syscall"
)

func main() {
	listenfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("create socket failed: %v\n", err)
		return
	}

	if err := syscall.Bind(listenfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}); err != nil {
		log.Printf("bind failed: %v\n", err)
		return
	}

	if err := syscall.Listen(listenfd, 1024); err != nil {
		log.Printf("listen failed: %v\n", err)
		return
	}

	r, err := newReactor(listenfd)
	if err != nil {
		log.Printf("new reactor failed: %v\n", err)
		return
	}
	defer r.Close()

	acceptor := func(listenfd int) {
		connfd, _, err := syscall.Accept(listenfd)
		if err != nil {
			log.Printf("accept failed: %v\n", err)
			return
		}

		log.Printf("Accepted a connection")

		if err := syscall.SetNonblock(connfd, true); err != nil {
			log.Printf("set nonblock failed: %v\n", err)
			return
		}

		event := &syscall.EpollEvent{
			Events: syscall.EPOLLIN,
			Fd:     int32(connfd),
		}

		if err := r.EpollCtl(connfd, syscall.EPOLL_CTL_ADD, event); err != nil {
			log.Printf("epoll ctl failed: %v\n", err)
			return
		}
	}

	handler := func(connfd int) {
		buf := make([]byte, 1024)
		n, err := syscall.Read(connfd, buf)
		if err != nil {
			log.Printf("read failed: %v\n", err)
			return
		}

		if n == 0 {
			log.Printf("connection closed by peer")
			return
		}

		log.Printf("received: %s", buf[:n])
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
	acceptor func(listenfd int)
	handler  func(connfd int)
}

func newReactor(listenfd int) (*Reactor, error) {
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &Reactor{
		listenfd: listenfd,
		epfd:     epfd,
	}, nil
}

func (r *Reactor) RegisterAcceptor(acceptor func(listenfd int)) {
	r.acceptor = acceptor
}

func (r *Reactor) RegisterHandler(handler func(connfd int)) {
	r.handler = handler
}

func (r *Reactor) EpollWait() error {
	events := make([]syscall.EpollEvent, 1024)

	n, err := syscall.EpollWait(r.epfd, events, -1)
	if err != nil {
		log.Printf("epoll wait failed: %v\n", err)
		return err
	}

	for i := 0; i < n; i++ {
		if int(events[i].Fd) == r.listenfd {
			r.acceptor(int(events[i].Fd))
		} else {
			r.handler(int(events[i].Fd))
		}
	}

	return nil
}

func (r *Reactor) EpollCtl(fd int, op int, event *syscall.EpollEvent) error {
	return syscall.EpollCtl(r.epfd, op, fd, event)
}

func (r *Reactor) Close() error {
	return syscall.Close(r.epfd)
}