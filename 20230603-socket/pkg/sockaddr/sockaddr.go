package sockaddr

import (
	"net"
	"strconv"
	"strings"
	"syscall"
)

func Parse(addr string) (*syscall.SockaddrInet4, error) {
	s := strings.Split(addr, ":")

	ip4 := net.ParseIP(s[0])

	switch {
	case len(s[0]) != 0 && ip4 == nil:
		return nil, syscall.EINVAL
	case len(s[0]) != 0 && ip4 != nil:
		// do nothing
	case len(s[0]) == 0:
		ip4 = net.IPv4zero
	}

	port, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, err
	}

	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip4)
	return sa, nil
}
