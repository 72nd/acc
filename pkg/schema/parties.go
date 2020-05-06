package schema

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const DefaultPartiesFile = "parties.yaml"

type Parties struct {
	Employees []Party `yaml:"employees"`
	Customers []Party `yaml:"customers"`
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
	if err := yaml.Unmarshal(raw, &pty); err != nil {
		logrus.Fatal(err)
	}
	return pty
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (p Parties) Save(path string) {
	SaveToYaml(p, path)
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

// EmployeeByIdentifier returns a Employee if there is one with the given identifier. Otherwise an error will be returned.
func (p Parties) EmployeeByIdentifier(ident string) (*Party, error) {
	for i := range p.Employees {
		if p.Employees[i].Identifier == ident {
			return &p.Employees[i], nil
		}
	}
	return nil, fmt.Errorf("no employee for identifier «%s» found", ident)
}

// CustomerByIdentifier returns a Customer if there is one with the given identifier. Otherwise an error will be returned.
func (p Parties) CustomerByIdentifier(ident string) (*Party, error) {
	for i := range p.Customers {
		if p.Customers[i].Identifier == ident {
			return &p.Customers[i], nil
		}
	}
	return nil, fmt.Errorf("no customer for identifier «%s» found", ident)
}

// Party represents some person or organisation.
type Party struct {
	// Id is the internal unique identifier of the Expense.
	Id string `yaml:"id" default:""`
	// Identifier is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `yaml:"identifier" default:"c-1"`
	Name       string `yaml:"name" default:"Max Mustermann"`
	Street     string `yaml:"street" default:"Main Street"`
	StreetNr   int    `yaml:"streetNr" default:"1"`
	Place      string `yaml:"place" default:"Zurich"`
	PostalCode int    `yaml:"postalCode" default:"8000"`
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
func NewCompanyParty(useDefaults bool) Party {
	if useDefaults {
		return Party{
			Name:       "Fantasia Company",
			Street:     "Main Street",
			StreetNr:   10,
			Place:      "Zurich",
			PostalCode: 8000,
		}
	}
	return Party{}
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

// Type returns a string with the type name of the element.
func (p Party) Type() string {
	return ""
}

// String returns a human readable representation of the element.
func (p Party) String() string {
	return fmt.Sprintf("")
}

// Conditions returns the validation conditions.
func (p Party) Conditions() util.Conditions {
	return util.Conditions{

	}
}

// Validate the element and return the result.
func (p Party) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(p)}
}
