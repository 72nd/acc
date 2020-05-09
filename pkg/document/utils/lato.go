package utils

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/sirupsen/logrus"
)

func riceBox() *rice.Box {
	box, err := rice.FindBox("utils")
	if err != nil {
		logrus.Error("rice find box failed: ", err)
	}
	return box
}

func LatoHeavy() []byte {
	box := riceBox()
	data, err := box.Bytes("Lato-Heavy.ttf")
	if err != nil {
		logrus.Errorf("could not load lato heavy: ", err)
	}
	return data
}

func LatoRegular() []byte {
	box := riceBox()
	data, err := box.Bytes("Lato-Regular.ttf")
	if err != nil {
		logrus.Errorf("could not load lato regular: ", err)
	}
	return data
}
