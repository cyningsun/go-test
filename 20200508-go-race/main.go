package main

import (
	"fmt"
	"github.com/cyningsun/go-test/20200508-go-race/cache"
	proto "github.com/golang/protobuf/proto"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	c := cache.NewPersonCache()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func () {
		for {
			select {
			case <-ch:
				break
			default:
				r := rand.Uint64() % 100
				key := strconv.FormatUint(r, 10)
				p, ok := c.Get(key)
				if !ok {
					continue
				}
				p.Name = proto.String("reading....")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ch:
				break
			default:
				r :=  rand.Uint64() % 100
				key := strconv.FormatUint(r,10)
				p,ok := c.Get(key)
				if !ok {
					continue
				}
				proto.Marshal(p)
			}
		}
	}()

	fmt.Println(<-ch)
}