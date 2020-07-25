package util

import (
	"fmt"
	"syscall"
	"unsafe"
)

func TerminalWidth() (int, error) {
	w := &winSize{}
	code, _, errNr := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(w)))

	if int(code) == -1 {
		return 0, fmt.Errorf("syscall to determine terminal width failed with error code \"%d\"", errNr)
	}
	return int(w.Col), nil
}
