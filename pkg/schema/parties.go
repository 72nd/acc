package schema

import "github.com/google/uuid"

// Employees represents a slice of parties as employees.
type Employees []Party

func (e Employees) SetId() {
	for i := range e {
		e[i].SetId()
	}
}

// Customer represents a customer of the company.
type Customer struct {
	Party
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
