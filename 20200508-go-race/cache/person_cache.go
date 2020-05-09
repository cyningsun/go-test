package cache

import (
	"github.com/cyningsun/go-test/20200508-go-race/pb"
	gocache "github.com/patrickmn/go-cache"
	proto "github.com/golang/protobuf/proto"
	"math/rand"
	"strconv"
	"time"
)


type PersonCache struct {
	c *gocache.Cache
}

func NewPersonCache() *PersonCache {
	one := &PersonCache{c:gocache.New(time.Minute, time.Hour)}
	go one.load()
	return one
}

func (p *PersonCache) load() {
	r :=  rand.Uint64() % 100
	key := strconv.FormatUint(r,10)
	newOne := &pb.Person{
		Id:                   proto.Uint64(10),
		Name:                 proto.String("init"),
		Age:                  proto.Uint32(rand.Uint32()),
	}
	p.c.Set(key, newOne, time.Minute)
}

func (p *PersonCache) Get(key string) (*pb.Person,bool) {
	ret, ok := p.c.Get(key)
	if !ok {
		return nil, false
	}
	return ret.(*pb.Person),true
}



