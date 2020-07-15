package schema

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gitlab.com/72nd/acc/pkg/util"
)

type Company struct {
	Name       string `yaml:"name" default:"Fortuna Inc."`
	Street     string `yaml:"street" default:"Main Street"`
	StreetNr   int    `yaml:"streetNr" default:"1"`
	PostalCode int    `yaml:"postalCode" default:"8000"`
	Place      string `yaml:"place" default:"Zurich"`
	Phone      string `yaml:"phone" default:"+41 78 000 00 00"`
	Mail       string `yaml:"mail" default:"info@fortuna.com"`
	Url        string `yaml:"url" default:"https://fortuna.com"`
	Logo       string `yaml:"logo" default:"/path/to/logo.png"`
}

func NewCompany(logo string) Company {
	cmp := Company{}
	if err := defaults.Set(&cmp); err != nil {
		logrus.Fatal("error setting defaults: ", err)
	}
	if logo != "" {
		cmp.Logo = logo
	}
	return cmp
}

func InteractiveNewCompany(logo string) Company {
	cmp := NewCompany(logo)
	cmp.Name = util.AskString(
		"Name",
		"Name of the company",
		"Fortuna Inc.",
	)
	cmp.Street = util.AskString(
		"Street",
		"Street of the company",
		"Society Street",
	)
	cmp.StreetNr = util.AskInt(
		"Street Nr.",
		"Number of the street",
		49,
	)
	cmp.PostalCode = util.AskInt(
		"Postal Code",
		"Postal/ZIP Code",
		4223,
	)
	cmp.Place = util.AskString(
		"Place",
		"place of the company",
		"Zurich",
	)
	cmp.Place = util.AskString(
		"Phone",
		"Phone number of the company",
		"+41 78 000 00 00")
	cmp.Mail = util.AskString(
		"Mail",
		"General mail address",
		"info@fortuna.com",
	)
	cmp.Url = util.AskString(
		"Url",
		"Website URL",
		"https://fortuna.com",
	)
	if logo == "" {
		cmp.Logo = util.AskString(
			"Logo",
			"Path to logo file (use --logo to set with flag)",
			"logo.png",
		)
	}
	return cmp
}

func (c Company) Type() string {
	return "Company"
}

func (c Company) String() string {
	return fmt.Sprintf("%s at %s %d", c.Name, c.Street, c.StreetNr)
}

func (c Company) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: c.Name == "",
			Message:   "name is not set (Name is empty)",
		},
		{
			Condition: c.Street == "",
			Message:   "street name is not set (Street is empty)",
		},
		{
			Condition: c.StreetNr == 0,
			Message:   "street number is not set (StreetNr is 0)",
		},
		{
			Condition: c.Place == "",
			Message:   "place is not set (Place is empty)",
		},
		{
			Condition: c.PostalCode == 0,
			Message:   "postal code is not set (PostalCode is 0)",
		},
		{
			Condition: c.Phone == "",
			Message:   "phone is not set (Phone is empty)",
		},
		{
			Condition: c.Mail == "",
			Message:   "mail is not set (Mail is empty)",
		},
		{
			Condition: c.Url == "",
			Message:   "url is not set (Url is empty)",
		},
		{
			Condition: !util.FileExist(c.Logo),
			Message:   fmt.Sprintf("logo at «%s» not found", c.Logo),
		},
	}
}
