package util

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const ApplicationFolder = "/usr/share/applications/"

type External struct {
	Running     bool
	path        string
	cmd         *exec.Cmd
	retainFocus bool
}

func NewExternal(path string, retainFocus bool) External {
	return External{
		Running:     true,
		path:        path,
		retainFocus: retainFocus,
	}
}

func (e *External) Open() {
	var window string
	if e.retainFocus {
		var err error
		window, err = currentActiveWindow()
		if err != nil {
			logrus.Error("unable to identify active window: ", err)
		}
	}
	prg, err := e.application()
	if err != nil {
		logrus.Errorf("error while detect default application for \"%s\": %s", e.path, err)
	}
	e.cmd = exec.Command(prg, e.path)
	if err := e.cmd.Start(); err != nil {
		logrus.Errorf("could not open \"%s\" with %s", e.path, prg)
	}
	if window != "" {
		gainFocus(window)
	}
}

func (e External) Close() {
	if !e.Running {
		return
	}
	if err := e.cmd.Process.Kill(); err != nil {
		logrus.Warnf("could not kill the program from \"%s\"", e.path)
	}
}

func (e External) application() (string, error) {
	mime, err := exec.Command("xdg-mime", "query", "filetype", e.path).Output()
	if err != nil {
		return "", err
	}
	dft, err := exec.Command("xdg-mime", "query", "default", strings.Replace(string(mime), "\n", "", -1)).Output()
	if err != nil {
		return "", err
	}
	return extractExec(path.Join(ApplicationFolder, strings.Replace(string(dft), "\n", "", -1)))
}

func extractExec(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\nExec\=(.*?) .*\n`)
	rsl := re.FindStringSubmatch(string(data))
	if len(rsl) != 2 {
		return "", errors.New("could not parse desktop file")
	}
	return rsl[1], nil
}

func currentActiveWindow() (string, error) {
	output, err := exec.Command("xprop", "-root").Output()
	if err != nil {
		return "", err
	}
	rsl := regexp.MustCompile(`_NET_ACTIVE_WINDOW\(WINDOW\): .*#\s(.*)`).FindStringSubmatch(string(output))
	if len(rsl) != 2 {
		return "", errors.New("regex failed")
	}
	logrus.Debug("current window identified as ", rsl[1])
	return rsl[1], nil
}

func gainFocus(window string) {
	time.Sleep(500 * time.Millisecond)
	err := exec.Command("wmctrl", "-ia", window).Run()
	if err != nil {
		logrus.Error("error while retain focus with wmctrl: ", err)
	}
}
