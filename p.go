//go:build linux

package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	f, err := os.OpenFile("./test.txt.4M", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open failed. err(%v)\n\n", err)
	}

	// discard 0-2M 的空间，预期实际物理占用能减少 2M
	err = punchHoleLinux(f, 0, 2*1024*1024)
	if err != nil {
		fmt.Println("punch hole failed")
	}

	fmt.Println("punch hole success.")
}

func punchHoleLinux(file *os.File, offset int64, size int64) error {
	return syscall.Fallocate(int(file.Fd()), 0x1|0x2, offset, size)
}
