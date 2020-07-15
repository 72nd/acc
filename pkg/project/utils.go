package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

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

// getFoldersInPath returns all folders as absolute path in the given folder.
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

// getMatchingFilesInPath matches all files and folders in the given path and matches
// the found files against the given regex. All matching elements will be returned
// as absolute path.
func getMatchingFilesInPath(path string, re *regexp.Regexp) []string {
	var rsl []string
	elements, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.Fatal("error reading dir: ", err)
	}
	for i := range elements {
		if re.MatchString(elements[i].Name()) {
			rsl = append(rsl, filepath.Join(path, elements[i].Name()))
		}
	}
	return rsl
}

// absolutPath returns the absolute path of the relativePath in relation to the folder location.
func absolutPath(folder, relativePath string) string {
	return filepath.Clean(filepath.Join(folder, relativePath))
}

// relativePath returns the relative path in relation to the folder (basepath).
// In the case of the error the absolute path is returned instead and the user is
// presented with an error message.
//
// This is legit for the acc use case. The function will only be used to link assets
// relative to the data files (like customer.yaml, project.yaml etc.) and thus enabling
// the compatibility across different users on different system. In a case of the error
// the user is asked to resolve this unlikely situation manually. If this doesn't happen
// another user will prompted an error thus the data coherency is guaranteed.
func relativePath(folder, absolutPath string) string {
	rsl, err := filepath.Rel(folder, absolutPath)
	if err != nil {
		logrus.Errorf("path \"%s\" couldn't be made relative to folder \"%s\", please resolve this manually", absolutPath, folder)
		return absolutPath
	}
	return rsl
}

// createNonExistingDir creates a directory with the given path. If the folder already exists, nothing will be done.
func createNonExistingDir(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}
	if err := os.Mkdir(path, 0777); err != nil {
		logrus.Fatal("error while creating folder: ", err)
	}
}

func relativeAssetPath(absolute, asset string) string {
	rsl, err := filepath.Rel(absolute, asset)
	if err != nil {
		logrus.Errorf("couldn't determine relative path \"%s\" in respective of \"%s\": %s", asset, absolute, err)
		return asset
	}
	return rsl
}


