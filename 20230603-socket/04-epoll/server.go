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

	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Printf("create epoll failed: %v\n", err)
		return
	}

	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, listenfd, &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(listenfd),
	}); err != nil {
		log.Printf("epoll ctl failed: %v\n", err)
		return
	}

	events := make([]syscall.EpollEvent, 1024)
	for {
		n, err := syscall.EpollWait(epfd, events, -1)
		if err != nil {
			log.Printf("epoll wait failed: %v\n", err)
			return
		}

		for i := 0; i < n; i++ {
			if events[i].Fd == int32(listenfd) {
				connfd, _, err := syscall.Accept(listenfd)
				if err != nil {
					log.Printf("accept failed: %v\n", err)
					continue
				}

				log.Printf("Accepted a connection")

				if err := syscall.SetNonblock(connfd, true); err != nil {
					log.Printf("set nonblock failed: %v\n", err)
					continue
				}

				event := syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(connfd),
				}

				if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connfd, &event); err != nil {
					log.Printf("epoll ctl failed: %v\n", err)
					continue
				}
			} else {
				fd := int(events[i].Fd)
				buf := make([]byte, 1024)
				n, err := syscall.Read(fd, buf)
				if err != nil {
					log.Printf("read failed: %v\n", err)
					continue
				}

				log.Printf("recv: %s\n", buf[:n])
			}

		}
	}
}
