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

func SetsockoptInt(fd, level, opt int, value int) (err error) {
	err = syscall.SetsockoptInt(fd, level, opt, value)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
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

func Listen(fd, backlog int) (err error) {
	err = syscall.Listen(fd, backlog)
	if err != nil {
		return os.NewSyscallError("listen", err)
	}

	return nil
}

func Accept(fd int) (nfd int, sa syscall.Sockaddr, err error) {
	nfd, sa, err = syscall.Accept(fd)
	if err != nil {
		return -1, nil, os.NewSyscallError("accept", err)
	}

	return nfd, sa, nil
}

func Shutdown(fd, how int) (err error) {
	err = syscall.Shutdown(fd, how)
	if err != nil {
		return os.NewSyscallError("shutdown", err)
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

func Close(fd int) (err error) {
	err = syscall.Close(fd)
	if err != nil {
		return os.NewSyscallError("close", err)
	}

	return nil
}

func Select(nfd int, rset, wset, eset *syscall.FdSet, timeout *syscall.Timeval) (n int, err error) {
	n, err = syscall.Select(nfd, rset, wset, eset, timeout)
	if err != nil {
		if err == syscall.EINTR {
			return 0, nil
		}

		return n, os.NewSyscallError("select", err)
	}
	return n, nil
}

func EpollCreate1(flag int) (fd int, err error) {
	fd, err = syscall.EpollCreate1(flag)
	if err != nil {
		return -1, os.NewSyscallError("epoll_create1", err)
	}

	return fd, nil
}

func EpollCreate(size int) (fd int, err error) {
	fd, err = syscall.EpollCreate(size)
	if err != nil {
		return -1, os.NewSyscallError("epoll_create", err)
	}

	return fd, nil
}

func EpollCtl(epfd, op, fd int, event *syscall.EpollEvent) (err error) {
	err = syscall.EpollCtl(epfd, op, fd, event)
	if err != nil {
		return os.NewSyscallError("epoll_ctl", err)
	}

	return nil
}

func EpollWait(epfd int, events []syscall.EpollEvent, msec int) (n int, err error) {
	n, err = syscall.EpollWait(epfd, events, msec)
	if err != nil {
		if err == syscall.EINTR {
			return 0, nil
		}

		return n, os.NewSyscallError("epoll_wait", err)
	}
	return n, nil
}
