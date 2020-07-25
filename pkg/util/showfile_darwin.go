package util

import "github.com/sirupsen/logrus"

type External struct{}

func NewExternal(path string, retainFocus bool) External {
	return External{}
}

func (e *External) Open() {
	logrus.Warn("external open is not implemented for the MacOS/darwin platform")
}

func (e External) Close() {
	logrus.Warn("external close is not implemented for the MacOS/darwin platform")
}
