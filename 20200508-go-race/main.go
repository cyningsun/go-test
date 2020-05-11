package main

import (
	"fmt"
	"github.com/cyningsun/go-test/20200508-go-race/cache"
	proto "github.com/gogo/protobuf/proto"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const MaxLen = 100

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterLength = len(letterBytes)

var (
	c = cache.NewPersonCache()
)

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Write(w http.ResponseWriter, req *http.Request) {
	r := rand.Uint64() % cache.Max
	key := strconv.FormatUint(r, 10)
	p, ok := c.Get(key)
	if !ok {
		return
	}
	p.Name = proto.String(randString(rand.Int()%letterLength))
	time.Sleep(time.Nanosecond)
	p.Address = proto.String(randString(rand.Int()%letterLength))
}

func Read(w http.ResponseWriter, req *http.Request) {
	r :=  rand.Uint64() % cache.Max
	key := strconv.FormatUint(r,10)
	p,ok := c.Get(key)
	if !ok {
		return
	}
	b,_ := proto.Marshal(p)
	w.Write(b)
}

func main() {
	http.HandleFunc("/read", Read)
	http.HandleFunc("/write", Write)
	fmt.Println("server is listening on 8080")
	http.ListenAndServe(":8080", nil)
}