package schema

import (
	"fmt"
	"github.com/72nd/acc/pkg/util"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
func NewParties(useDefault bool) Parties {
	if useDefault {
		return Parties{
			Employees: []Party{NewParty()},
			Customers: []Party{NewParty()},
		}
	}
	return Parties{
		Employees: []Party{},
		Customers: []Party{},
	}
}

// OpenParties opens a Parties element saved in the json file given by the path.
func OpenParties(path string) Parties {
	var prt Parties
	util.OpenYaml(&prt, path, "parties")
	return prt
}

// Save writes the element as a json to the given path.
func (p Parties) Save(path string) {
	util.SaveToYaml(p, path, "parties")
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

func (p Parties) CustomerByIdent(ident string) (*Party, error) {
	for i := range p.Customers {
		if p.Customers[i].Identifier == ident {
			return &p.Customers[i], nil
		}
	}
	return nil, fmt.Errorf("no customer for identifier «%s» found", ident)
}

func (p Parties) EmployeeByIdent(ident string) (*Party, error) {
	for i := range p.Employees {
		if p.Employees[i].Identifier == ident {
			return &p.Employees[i], nil
		}
	}
	return nil, fmt.Errorf("no employee for identifier «%s» found", ident)
}

func (p Parties) CustomerStringById(id string) string {
	if id == "" {
		return "no customer associated"
	}
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
		result[i] = p.Employees[i].SearchItem("Employe")
	}
	return result
}

func (p Parties) CustomersSearchItems() util.SearchItems {
	result := make(util.SearchItems, len(p.Customers))
	for i := range p.Customers {
		result[i] = p.Customers[i].SearchItem("Customer")
	}
	return result
}

func (p Parties) Validate() util.ValidateResults {
	var rsl util.ValidateResults
	for i := range p.Customers {
		rsl = append(rsl, util.Check(p.Customers[i]))
	}
	for i := range p.Employees {
		rsl = append(rsl, util.Check(p.Employees[i]))
	}
	return rsl
}

// Party represents some person or organisation.
type Party struct {
	// Id is the internal unique identifier of the Expense.
	Id string `yaml:"id" default:""`
	// Value is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string    `yaml:"identifier" default:"?-1"`
	Name       string    `yaml:"name" default:"Max Mustermann"`
	Street     string    `yaml:"street" default:"Main Street"`
	StreetNr   int       `yaml:"streetNr" default:"1"`
	PostalCode int       `yaml:"postalCode" default:"8000"`
	Place      string    `yaml:"place" default:"Zurich"`
	PartyType  PartyType `yaml:"partyType" default:"0"`
}

// NewParty returns a new Party with the default values.
func NewParty() Party {
	pty := Party{}
	pty.Id = GetUuid()
	if err := defaults.Set(&pty); err != nil {
		logrus.Fatal("error setting defaults: ", err)
	}
	return pty
}

func NewPartyWithUuid() Party {
	pty := NewParty()
	pty.Id = GetUuid()
	return pty
}

func InteractiveNewParty(partyType string) Party {
	pty := NewPartyWithUuid()
	pty.Name = util.AskString(
		"Name",
		fmt.Sprintf("Name of the %s", partyType),
		"Bimpf the first",
	)
	pty.Street = util.AskString(
		"Street",
		fmt.Sprintf("Street of the %s", partyType),
		"Society Street",
	)
	pty.StreetNr = util.AskInt(
		"Street Nr.",
		"Number of the street",
		49,
	)
	pty.PostalCode = util.AskInt(
		"Postal Code",
		"Postal/ZIP Code",
		4223,
	)
	pty.Place = util.AskString(
		"Place",
		fmt.Sprintf("Place/City of %s", partyType),
		"Zurich",
	)
	return pty
}

func InteractiveNewCustomer(s Schema) Party {
	pty := InteractiveNewParty("Customer")
	pty.Identifier = util.AskString(
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.Parties.GetCustomerIdentifiables(), DefaultCustomerPrefix))
	pty.PartyType = CustomerType
	return pty
}

func InteractiveNewEmployee(s Schema) Party {
	pty := InteractiveNewParty("Employee")
	pty.Identifier = util.AskString(
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.Parties.GetEmployeeIdentifiables(), DefaultEmployeePrefix))
	pty.PartyType = EmployeeType
	return pty
}

func InteractiveNewGenericParty(arg interface{}) interface{} {
	sel := util.AskIntFromList(
		"Type",
		"Choose type of party",
		util.SearchItems{
			{
				Name:  "Customer",
				Value: 1,
			},
			{
				Name:  "Employee",
				Value: 2,
			},
		})
	s, ok := arg.(Schema)
	if !ok {
		logrus.Fatalf("arg \"%s\" couldn't be parsed as Acc", arg)
	}
	switch sel {
	case 1:
		return InteractiveNewCustomer(s)
	case 2:
		return InteractiveNewEmployee(s)
	default:
		logrus.Fatal("invalid result form AskIntFromList")
	}
	return nil
}

func (p Party) SearchItem(typ string) util.SearchItem {
	return util.SearchItem{
		Name:        p.Name,
		Type:        typ,
		Value:       p.Id,
		SearchValue: fmt.Sprintf("%s %s %s", p.Name, p.Identifier, p.Place),
	}
}

// GetId returns the unique id of the element.
func (p Party) GetId() string {
	return p.Id
}

func (p Party) GetIdentifier() string {
	return p.Identifier
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
	return "Party"
}

// String returns a human readable representation of the element.
func (p Party) String() string {
	return fmt.Sprintf("%s (%s), %s", p.Name, p.Identifier, p.Place)
}

// Short returns a short representation of the element.
func (p Party) Short() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Identifier)
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
		{
			Condition: p.Id == "",
			Message:   "unique identifier not set (Id is empty)",
		},
		{
			Condition: p.Identifier == "",
			Message:   "human readable identifier not set (Identifier is empty)",
		},
		{
			Condition: p.Name == "",
			Message:   "name is not set (Name is empty)",
		},
		{
			Condition: p.Street == "",
			Message:   "street name is not set (Street is empty)",
		},
		{
			Condition: p.StreetNr == 0,
			Message:   "street number is not set (StreetNr is 0)",
		},
		{
			Condition: p.Place == "",
			Message:   "place is not set (Place is empty)",
		},
		{
			Condition: p.PostalCode == 0,
			Message:   "postal code is not set (PostalCode is 0)",
		},
	}
}
