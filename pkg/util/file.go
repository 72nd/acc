package util

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func RelativeAssetPath(absolute, asset string) string {
	if asset == "" {
		return ""
	}
	if !filepath.IsAbs(asset) {
		asset = AbsolutePathWithWD(asset)
	}
	rsl, err := filepath.Rel(absolute, asset)
	if err != nil {
		logrus.Errorf("couldn't determine relative path \"%s\" in respective of \"%s\": %s", asset, absolute, err)
		return asset
	}
	return rsl
}

// AbsolutePathWithWD returns the absolute path of a given path in relation to the
// current working dir. This is mainly used to normalize user inputs (like for assets).
// If the path is already absolute it will be returned without any change.
func AbsolutePathWithWD(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	wd, err := os.Getwd()
	if err != nil {
		logrus.Fatal("working directory not found: ", err)
	}
	return filepath.Join(wd, path)
}
