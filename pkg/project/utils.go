package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// repositoryPath tries to get the `ACC_REPOSITORY` environment variable.
// If not set the current working directory will be used.
func repositoryPath() string {
	env := os.Getenv("ACC_REPOSITORY")
	if env == "" {
		path, err := os.Getwd()
		if err != nil {
			logrus.Fatal("\"ACC_REPOSITORY\" is not set and couldn't determine working directory: ", err)
		}
		logrus.Warnf("the enviroment variable \"ACC_REPOSITORY\" is not set, will use current working dir (%s)", path)
		return path
	}
	return env
}

// folderName takes a string and returns version without spaces, umlauts etc.
func folderName(name string) string {
	name = strings.ToLower(name)
	r := strings.NewReplacer(
		"ä", "ae",
		"ö", "oe",
		"ü", "ue",
		"à", "a",
		"é", "e",
		"è", "e",
		" ", "-",
		"_", "-",
		".", "-")
	return r.Replace(name)
}

// getFoldersInPath returns all folders in the given folder.
func getFoldersInPath(path string) []string {
	var rsl []string
	elements, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.Fatal("error reading dir: ", err)
	}
	for i := range elements {
		p := filepath.Join(path, elements[i].Name())
		if stat, _ := os.Stat(p); stat.IsDir() {
			rsl = append(rsl, p)
		}
	}
	return rsl
}
