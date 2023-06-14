package ioutil

import (
	"os"
	"syscall"
)

func Socket(family, sotype, proto int) (fd int, err error) {
	fd, err = syscall.Socket(family, sotype, proto)
	if err != nil {
		return -1, os.NewSyscallError("socket", err)
	}

	return fd, nil
}

func Connect(fd int, sa syscall.Sockaddr) (err error) {
	err = syscall.Connect(fd, sa)
	switch err {
	case syscall.EINPROGRESS, syscall.EALREADY, syscall.EINTR:
		return nil
	case nil, syscall.EISCONN:
		return nil
	default:
		return os.NewSyscallError("connect", err)
	}
}

func SetNonblock(fd int, nonblocking bool) (err error) {
	err = syscall.SetNonblock(fd, nonblocking)
	if err != nil {
		return os.NewSyscallError("setnonblock", err)
	}

	return nil
}

func Bind(fd int, sa syscall.Sockaddr) (err error) {
	err = syscall.Bind(fd, sa)
	if err != nil {
		return os.NewSyscallError("bind", err)
	}

	return nil
}

func Read(fd int, b []byte) (n int, err error) {
	n, err = syscall.Read(fd, b)
	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EINTR {
			return 0, nil
		}

		return n, os.NewSyscallError("read", err)
	}
	return n, nil
}

func Write(fd int, b []byte) (n int, err error) {
	n, err = syscall.Write(fd, b)
	if err != nil {
		if err == syscall.EAGAIN {
			return 0, nil
		}

		return n, os.NewSyscallError("write", err)
	}
	return n, nil
}
