// +build !aix,!darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package main

func isTerminal() bool {
	return false
}
