package util

import "errors"

func TerminalWidth() (int, error) {
	return -1, errors.New("get terminal width not implemented for MacOS/darwin")
}
