package bimpf

import (
	"fmt"
	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
)

type Employees []Employee

func (e Employees) ById(id int) (*Employee, error) {
	for i := range e {
		if e[i].Id == id {
			return &e[i], nil
		}
	}
	return nil, fmt.Errorf("no employee in Bimpf dump for id «%d» found", id)
}

// Employee reassembles the structure of a Employee in a Bimpf json dump file.
type Employee struct {
	Id        int    `json:"id"`
	SbId      string `json:"sb_id"`
	FirstName string `json:"first_name"`
	Name      string `json:"name"`
	Username  string `json:"username"`
}

// Type returns a string with the type name of the element.
func (e Employee) Type() string {
	return "SB-Employee"
}

// String returns a human readable representation of the Employee.
func (e Employee) String() string {
	return fmt.Sprintf("%d/%s (%s %s)", e.Id, e.SbId, e.FirstName, e.Name)
}

// Conditions returns the validation conditions.
func (e Employee) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: e.Id < 1,
			Message:   "id is not set (id < 1)",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: e.SbId == "",
			Message:   "solutionsbüro id is not set",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: e.Name == "",
			Message:   "family name not set",
			Level:     util.BeforeImportFlaw,
		},
	}
}

// Validate the element and return the result.
func (e Employee) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(e)}
}

// Convert returns the employee as a acc Party.
func (e Employee) Convert() schema.Party {
	pty := schema.Party{
		Identifier: e.SbId,
		Name:       fmt.Sprintf("%s %s", e.FirstName, e.Name),
		Street:     "",
		StreetNr:   0,
		Place:      "",
		PostalCode: 0,
	}
	pty.SetId()
	return pty
}
