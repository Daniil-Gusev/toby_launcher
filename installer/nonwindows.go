//go:build !windows

package main

import "fmt"

func installWindows() error {
	return fmt.Errorf("installWindows is not supported on this platform")
}

func uninstallWindows() error {
	return fmt.Errorf("uninstallWindows is not supported on this platform")
}

func isAdmin() bool {
	return false
}
