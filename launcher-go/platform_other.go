//go:build !windows

package main

import "fmt"

// Stand-ins so the core logic builds and tests on Linux; the real launcher is Windows-only.

func steamRoots() []string { return nil }
func gameRunning() bool    { return false }

func main() {
	fmt.Println("BL2 Arsenal Launcher is a Windows program.")
}
