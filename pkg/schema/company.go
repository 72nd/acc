package schema

import (
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
)

type Company struct {
	Name       string `yaml:"name" default:"Fortuna Inc."`
	Street     string `yaml:"street" default:"Main Street"`
	StreetNr   int    `yaml:"streetNr" default:"1"`
	Place      string `yaml:"place" default:"Zurich"`
	PostalCode int    `yaml:"postalCode" default:"8000"`
	Phone      string `yaml:"phone" default:"+41 78 000 00 00"`
	Mail       string `yaml:"mail" default:"info@fortuna.com"`
	Url        string `yaml:"url" default:"https://fortuna.com"`
	Logo       string `yaml:"logo" default:"/path/to/logo.png"`
}

func NewCompany() Company {
	cmp := Company{}
	if err := defaults.Set(&cmp); err != nil {
		logrus.Fatal(err)
	}
	return cmp
}
