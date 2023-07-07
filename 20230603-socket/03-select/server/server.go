//go:build linux && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/pkg/fdset"
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

	ioutil.SetsockoptInt(listenfd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	if err := ioutil.Bind(listenfd, sa); err != nil {
		log.Printf("bind failed: %v\n", err)
		return
	}

	if err := ioutil.Listen(listenfd, 1024); err != nil {
		log.Printf("listen failed: %v\n", err)
		return
	}

	var (
		rset, allset syscall.FdSet
		client       [syscall.FD_SETSIZE]int
		i, maxi      int
	)

	maxi = -1
	for i = 0; i < syscall.FD_SETSIZE; i++ {
		client[i] = -1
	}

	maxfd := listenfd
	fdset.Zero(&allset)
	fdset.Set(&allset, listenfd)

	for {

		rset = allset

		nready, err := ioutil.Select(maxfd+1, &rset, nil, nil, nil)
		if err != nil {
			log.Printf("select failed: %v\n", err)
			return
		}

		if fdset.IsSet(&rset, listenfd) {
			connfd, _, err := ioutil.Accept(listenfd)
			if err != nil {
				log.Printf("accept failed: %v\n", err)
				continue
			}

			log.Printf("Accepted a connection")

			for i = 0; i < syscall.FD_SETSIZE; i++ {
				if client[i] < 0 {
					client[i] = connfd
					break
				}
			}

			if i == syscall.FD_SETSIZE {
				log.Printf("too many clients\n")
				ioutil.Close(connfd)
				client[i] = -1
				continue
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

		for i := 1; i <= maxi; i++ {
			connfd := client[i]
			if connfd < 0 {
				continue
			}

			if fdset.IsSet(&rset, connfd) {
				args := &proto.Args{}
				size := binary.Size(*args)
				recvbuf := make([]byte, 1024)

				var err error
				for tn, rn := 0, 0; tn < size && err == nil; tn += rn {
					rn, err = ioutil.Read(connfd, recvbuf)
					if err != nil {
						log.Printf("read failed: %v\n", err)
						ioutil.Close(connfd)
						fdset.Clear(&allset, connfd)
						client[i] = -1
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
					ioutil.Close(connfd)
					fdset.Clear(&allset, connfd)
					client[i] = -1
					continue
				}

				ret := &proto.Result{Sum: args.Args1 + args.Args2}
				buf := bytes.NewBuffer([]byte{})
				if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
					log.Printf("binary write failed: %v\n", err)
					ioutil.Close(connfd)
					fdset.Clear(&allset, connfd)
					client[i] = -1
					continue
				}

				_, err = ioutil.Write(connfd, buf.Bytes())
				if err != nil {
					log.Printf("write failed: %v\n", err)
					syscall.Close(connfd)
					fdset.Clear(&allset, connfd)
					client[i] = -1
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
