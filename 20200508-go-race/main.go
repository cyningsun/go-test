package main

import (
	"fmt"
	"github.com/cyningsun/go-test/20200508-go-race/cache"
	proto "github.com/gogo/protobuf/proto"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const Concurrency = 2
const MaxLen = 100

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterLength = len(letterBytes)

func RandLengthString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	c := cache.NewPersonCache()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	for i:=0;i<Concurrency;i++ {
		go func(n int) {
			for {
				select {
				case <-ch:
					break
				default:
					r := rand.Uint64() % cache.Max
					key := strconv.FormatUint(r, 10)
					p, ok := c.Get(key)
					if !ok {
						continue
					}
					p.Name = proto.String(RandLengthString(rand.Int()%letterLength))
					time.Sleep(time.Nanosecond)
					p.Address = proto.String(RandLengthString(rand.Int()%letterLength))
					fmt.Println(n, "Writting")
				}
			}
		}(i)
	}

	for i:=0;i<Concurrency;i++{
		go func(n int) {
			for {
				select {
				case <-ch:
					break
				default:
					r :=  rand.Uint64() % cache.Max
					key := strconv.FormatUint(r,10)
					p,ok := c.Get(key)
					if !ok {
						continue
					}
					b,_ := proto.Marshal(p)
					fmt.Println(n, "marshaling", string(b))
				}
			}
		}(i)
	}

	fmt.Println(<-ch)
}