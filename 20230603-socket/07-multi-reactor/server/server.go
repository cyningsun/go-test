//go:build linux && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"sync"
	"syscall"

	"github.com/cyningsun/go-test/20230603-socket/pkg/ioutil"
	"github.com/cyningsun/go-test/20230603-socket/pkg/proto"
	"github.com/cyningsun/go-test/20230603-socket/pkg/sockaddr"
)

const (
	MAX_OPEN = 1024
)

var addr string

type Hdlr func(int)

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

	r, err := newReactorMgr(listenfd, 4)
	if err != nil {
		log.Printf("new reactor manager failed: %v\n", err)
		return
	}
	defer r.Close()

	if err := r.Run(); err != nil {
		log.Printf("reactor manager run failed: %v\n", err)
		return
	}
}

type ReactorMgr struct {
	reactors []*Reactor
}

func newReactorMgr(listenfd int, reactorNums int) (*ReactorMgr, error) {
	var err error
	srs := make([]*Reactor, reactorNums)
	for i := 0; i < reactorNums; i++ {
		srs[i], err = newReactor()
		if err != nil {
			return nil, err
		}
	}

	rm := &ReactorMgr{
		reactors: srs,
	}

	reactor, _ := rm.Pick(listenfd)
	reactor.OnAccept(listenfd, rm.acceptor)

	return rm, nil
}

func (rm *ReactorMgr) Run() error {
	var wg sync.WaitGroup
	wg.Add(len(rm.reactors))

	for _, r := range rm.reactors {
		go func(iwg *sync.WaitGroup, ir *Reactor) {
			defer iwg.Done()

			for {
				if err := ir.EpollWait(); err != nil {
					log.Printf("epoll wait failed: %v\n", err)
					return
				}
			}
		}(&wg, r)
	}

	wg.Wait()

	return nil
}

func (rm *ReactorMgr) Close() error {
	for _, sr := range rm.reactors {
		if err := sr.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (rm *ReactorMgr) acceptor(listenfd int) {
	_, id := rm.Pick(listenfd)
	prefix := fmt.Sprintf("r:%d", id)

	connfd, _, err := ioutil.Accept(listenfd)
	if err != nil {
		log.Printf("accept failed: %v\n", err)
		return
	}

	log.Printf("[%s]accepted a connection", prefix)

	if err := ioutil.SetNonblock(connfd, true); err != nil {
		log.Printf("set nonblock failed: %v\n", err)
		return
	}

	reactor, idx := rm.Pick(connfd)
	reactor.OnAccept(connfd, rm.handler)

	log.Printf("[%s]dispatched to subreactor %d", prefix, idx)
}

// Pick up a subreactor to handle the connection
func (rm *ReactorMgr) Pick(fd int) (*Reactor, int) {
	idx := fd % len(rm.reactors)
	subReactor := rm.reactors[idx]
	return subReactor, idx
}

func (rm *ReactorMgr) handler(listenfd int) {
	reactor, id := rm.Pick(listenfd)
	prefix := fmt.Sprintf("r:%d", id)

	args := &proto.Args{}
	size := binary.Size(*args)
	recvbuf := make([]byte, 1024)

	var err error
	tn, rn := 0, 0
	for tn, rn = 0, 0; tn < size && err == nil; tn += rn {
		rn, err = ioutil.Read(listenfd, recvbuf)
		if err != nil {
			log.Printf("[%s]read failed: %v\n", prefix, err)
			break
		}

		if rn <= 0 {
			break
		}
	}

	if tn == 0 || err == ioutil.ECONNRESET {
		log.Printf("[%s]connection reset by peer", prefix)
		reactor.OnClose(listenfd)
		ioutil.Close(listenfd)
		return
	}

	if err != nil {
		reactor.OnClose(listenfd)
		ioutil.Close(listenfd)
		return
	}

	if err := binary.Read(bytes.NewBuffer(recvbuf[:size]), binary.BigEndian, args); err != nil {
		log.Printf("[%s]binary read failed: %v\n", prefix, err)
		return
	}

	ret := &proto.Result{Sum: args.Args1 + args.Args2}
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, ret); err != nil {
		log.Printf("[%s]binary write failed: %v\n", prefix, err)
		return
	}

	_, err = ioutil.Write(listenfd, buf.Bytes())
	if err != nil {
		log.Printf("[%s]write failed: %v\n", prefix, err)
		return
	}
}

type Reactor struct {
	epfd  int
	hdlrs sync.Map
}

func newReactor() (*Reactor, error) {
	epfd, err := ioutil.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &Reactor{
		epfd: epfd,
	}, nil
}

func (r *Reactor) EpollWait() error {
	events := make([]syscall.EpollEvent, MAX_OPEN)

	nready, err := ioutil.EpollWait(r.epfd, events, -1)
	if err != nil {
		log.Printf("epoll wait failed: %v\n", err)
		return err
	}

	for i := 0; i < nready; i++ {
		found, ok := r.hdlrs.Load(events[i].Fd)
		if !ok {
			log.Printf("handler not found for fd %d\n", events[i].Fd)
			continue
		}

		hdlr, ok := found.(Hdlr)
		if !ok {
			log.Printf("handler type not match for fd %d\n", events[i].Fd)
		}

		hdlr(int(events[i].Fd))
	}

	return nil
}

func (r *Reactor) OnAccept(fd int, hdlr Hdlr) {
	if err := ioutil.EpollCtl(r.epfd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Fd:     int32(fd),
		Events: syscall.EPOLLIN,
	}); err != nil {
		log.Printf("epoll ctl failed: %v\n", err)
		return
	}

	r.hdlrs.Store(int32(fd), hdlr)
}

func (r *Reactor) OnClose(fd int) {
	if err := ioutil.EpollCtl(r.epfd, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
		log.Printf("epoll ctl failed: %v\n", err)
		return
	}

	r.hdlrs.Delete(int32(fd))
}

func (r *Reactor) Close() error {
	return ioutil.Close(r.epfd)
}
