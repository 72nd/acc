package schema

import (
	"bufio"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const DefaultPartiesFile = "parties.yaml"
const DefaultEmployeePrefix = "y-"
const DefaultCustomerPrefix = "c-"

type PartyType int

const (
	EmployeeType PartyType = iota
	CustomerType
)

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

func (p Parties) GetCustomerIdentifiables() []Identifiable {
	pty := make([]Identifiable, len(p.Customers))
	for i := range p.Customers {
		pty[i] = p.Customers[i]
	}
	return pty
}

func (p Parties) GetEmployeeIdentifiables() []Identifiable {
	pty := make([]Identifiable, len(p.Employees))
	for i := range p.Employees {
		pty[i] = p.Employees[i]
	}
	return pty
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


func (p Parties) EmployeeById(id string) (*Party, error) {
	for i := range p.Employees {
		if p.Employees[i].Id == id {
			return &p.Employees[i], nil
		}
	}
	return nil, fmt.Errorf("no employee for id «%s» found", id)
}

func (p Parties) CustomerById(id string) (*Party, error) {
	for i := range p.Customers {
		if p.Customers[i].Id == id {
			return &p.Customers[i], nil
		}
	}
	return nil, fmt.Errorf("no customer for id «%s» found", id)
}

func (p Parties) EmployeeStringById(id string) string {
	emp, err := p.EmployeeById(id)
	if err != nil {
		logrus.Error("no employee found: ", err)
		return "no employee for id"
	}
	return emp.String()
}

func (p Parties) CustomerStringById(id string) string {
	cst, err := p.CustomerById(id)
	if err != nil {
		logrus.Error("no  found: ", err)
		return "no customer for id"
	}
	return cst.String()
}

func (p Parties) EmployeesSearchItems() util.SearchItems {
	result := make(util.SearchItems, len(p.Employees))
	for i := range p.Employees {
		result[i] = p.Employees[i].SearchItem()
	}
	return result
}

func (p Parties) CustomersSearchItems() util.SearchItems {
	result := make(util.SearchItems, len(p.Customers))
	for i := range p.Customers {
		result[i] = p.Customers[i].SearchItem()
	}
	return result
}

// Party represents some person or organisation.
type Party struct {
	// Id is the internal unique identifier of the Expense.
	Id string `yaml:"id" default:""`
	// Identifier is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `yaml:"identifier" default:"?-1"`
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

func NewPartyWithUuid() Party {
	pty := NewParty()
	pty.Id = GetUuid()
	return pty
}

func InteractiveNewParty(partyType string) Party {
	reader := bufio.NewReader(os.Stdin)
	pty := NewPartyWithUuid()
	pty.Name = util.AskString(
		reader,
		"Name",
		fmt.Sprintf("Name of the %s", partyType),
		"Bimpf the first",
	)
	pty.Street = util.AskString(
		reader,
		"Street",
		fmt.Sprintf("Street of the %s", partyType),
		"Society Street",
	)
	pty.StreetNr = util.AskInt(
		reader,
		"Street Nr.",
		"Number of the street",
		49,
	)
	pty.Place = util.AskString(
		reader,
		"Place",
		fmt.Sprintf("Place/City of %s", partyType),
		"Zurich",
	)
	pty.PostalCode = util.AskInt(
		reader,
		"Postal Code",
		"Postal/ZIP Code",
		4223,
	)
	return pty
}

func InteractiveNewCustomer(a Acc) Party {
	reader := bufio.NewReader(os.Stdin)
	pty := InteractiveNewParty("Customer")
	pty.Identifier = util.AskString(
		reader,
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(a.Parties.GetCustomerIdentifiables(), DefaultCustomerPrefix),
	)
	return pty
}

func InteractiveNewEmployee(a Acc) Party {
	reader := bufio.NewReader(os.Stdin)
	pty := InteractiveNewParty("Employee")
	pty.Identifier = util.AskString(
		reader,
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(a.Parties.GetEmployeeIdentifiables(), DefaultEmployeePrefix),
	)
	return pty
}

func (p Party) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:       p.Name,
		Identifier: p.Id,
		Value:      fmt.Sprintf("%s %s %s", p.Name, p.Identifier, p.Place),
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

// Type returns a string with the type name of the element.
func (p Party) Type() string {
	return ""
}

// String returns a human readable representation of the element.
func (p Party) String() string {
	return fmt.Sprintf("%s (%s), %s", p.Name, p.Identifier, p.Place)
}

func (p Party) AddressLines() string {
	result := p.Name
	if p.Street != "" && p.StreetNr != 0 {
		result = fmt.Sprintf("%s\n%s %d", result, p.Street, p.StreetNr)
	}
	if p.PostalCode != 0 && p.Place != "" {
		result = fmt.Sprintf("%s\n%d %s", result, p.PostalCode, p.Place)
	}
	return result
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
