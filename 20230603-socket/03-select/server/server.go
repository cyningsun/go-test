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

	var (
		rset, allset syscall.FdSet
		client       [syscall.FD_SETSIZE]int
		i, maxi      int
	)
	maxfd := listenfd
	fdset.Zero(&allset)
	fdset.Set(&allset, listenfd)

	for {

		rset = allset

		nready, err := syscall.Select(maxfd, &rset, nil, nil, nil)
		if err != nil {
			log.Printf("select failed: %v\n", err)
			return
		}

		if fdset.IsSet(&rset, listenfd) {
			connfd, _, err := syscall.Accept(listenfd)
			if err != nil {
				log.Printf("accept failed: %v\n", err)
				continue
			}

			log.Printf("Accepted a connection")

			for i := 0; i < syscall.FD_SETSIZE; i++ {
				if client[i] < 0 {
					client[i] = connfd
					break
				}
			}

			if i == syscall.FD_SETSIZE {
				log.Printf("too many clients\n")
				syscall.Close(connfd)
				return
			}

			fdset.Set(&allset, connfd)

			if connfd > maxfd {
				maxfd = connfd
			}

			if i > maxi {
				maxi = i
			}

			nready = nready - 1
			if nready <= 0 {
				continue
			}
		}

		for i := 0; i <= maxi; i++ {
			connfd := client[i]
			if connfd < 0 {
				continue
			}

			if fdset.IsSet(&rset, connfd) {
				args := &proto.Args{}
				size := binary.Size(*args)
				recvbuf := make([]byte, 1024)

				for tn, rn := 0, 0; tn < size; tn += rn {
					var err error
					rn, err = syscall.Read(connfd, recvbuf)
					if err != nil {
						log.Printf("read failed: %v\n", err)
						syscall.Close(connfd)
						fdset.Clear(&allset, connfd)
						client[i] = -1
					}

					if rn <= 0 {
						break
					}
				}

				if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
					log.Printf("binary read failed: %v\n", err)
					syscall.Close(connfd)
					fdset.Clear(&allset, connfd)
					client[i] = -1
				}

				ret := &proto.Result{Sum: args.Args1 + args.Args2}
				buf := bytes.NewBuffer([]byte{})
				if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
					log.Printf("binary write failed: %v\n", err)
					syscall.Close(connfd)
					fdset.Clear(&allset, connfd)
					client[i] = -1
				}

				_, err = syscall.Write(connfd, buf.Bytes())
				if err != nil {
					log.Printf("write failed: %v\n", err)
					syscall.Close(connfd)
					fdset.Clear(&allset, connfd)
					client[i] = -1
				}
			}

			nready = nready - 1
			if nready <= 0 {
				continue
			}
		}
	}
}
