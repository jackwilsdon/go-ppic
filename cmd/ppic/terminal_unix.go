// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package main

import (
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
)

func isTerminal() bool {
	return terminal.IsTerminal(unix.Stdout)
}
