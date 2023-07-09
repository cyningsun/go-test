package fdset

import "syscall"

// Set adds fd to the set fds.
func Set(fds *syscall.FdSet, fd int) {
	fds.Bits[fd/syscall.FD_SETSIZE] |= (1 << (uintptr(fd) % syscall.FD_SETSIZE))
}

// Clear removes fd from the set fds.
func Clear(fds *syscall.FdSet, fd int) {
	fds.Bits[fd/syscall.FD_SETSIZE] &^= (1 << (uintptr(fd) % syscall.FD_SETSIZE))
}

// IsSet returns whether fd is in the set fds.
func IsSet(fds *syscall.FdSet, fd int) bool {
	return fds.Bits[fd/syscall.FD_SETSIZE]&(1<<(uintptr(fd)%syscall.FD_SETSIZE)) != 0
}

// Zero clears the set fds.
func Zero(fds *syscall.FdSet) {
	for i := range fds.Bits {
		fds.Bits[i] = 0
	}
}
