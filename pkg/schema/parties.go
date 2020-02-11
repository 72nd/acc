package schema

import (
	"encoding/json"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const DefaultPartiesFile = "parties.json"

type Parties struct {
	Employees []Party `json:"employees"`
	Customers []Party `json:"customers"`
}

// NewParties returns a new Parties struct with the one Expense in it.
func NewParties() Parties {
	return Parties{
		Employees: []Party{NewParty()},
		Customers: []Party{NewParty()},
	}
}

// OpenParties opens a Parties element saved in the json file given by the path.
func OpenParties(path string) Parties {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	pty := Parties{}
	if err := json.Unmarshal(raw, &pty); err != nil {
		logrus.Fatal(err)
	}
	return pty
}

// Save writes the element as a json to the given path.
func (p Parties) Save(path string) {
	raw, err := json.Marshal(p)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		logrus.Fatal(err)
	}
}

// SetId sets a unique id to all elements in the struct.
func (p Parties) SetId() {
	for i := range p.Employees {
		p.Employees[i].SetId()
	}
	for i := range p.Customers {
		p.Customers[i].SetId()
	}
}

// Party represents some person or organisation.
type Party struct {
	// Id is the internal unique identifier of the Expense.
	Id string `json:"id" default:""`
	// Identifier is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `json:"identifier" default:"c-1"`
	Name       string `json:"name" default:"Max Mustermann"`
	Street     string `json:"street" default:"Main Street"`
	StreetNr   int    `json:"streetNr" default:"1"`
	Place      string `json:"place" default:"Zurich"`
	PostalCode int    `json:"postalCode" default:"8000"`
}

// NewParty returns a new Party with the default values.
func NewParty() Party {
	pty := Party{}
	if err := defaults.Set(&pty); err != nil {
		logrus.Fatal(err)
	}
	return pty
}

// NewCompanyParty returns a new default company Party.
func NewCompanyParty() Party {
	return Party{
		Name:       "Fantasia Company",
		Street:     "Main Street",
		StreetNr:   10,
		Place:      "Zurich",
		PostalCode: 8000,
	}
}

// GetId returns the unique id of the element.
func (p Party) GetId() string {
	return p.Id
}

// SetId generates a unique id for the element if there isn't already one defined.
func (p *Party) SetId() {
	if p.Id != "" {
		return
	}
	p.Id = uuid.Must(uuid.NewRandom()).String()
}
