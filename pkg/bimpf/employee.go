package bimpf

import (
	"fmt"
	"gitlab.com/72th/acc/pkg"
)

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
func (e Employee) Conditions() pkg.Conditions {
	return pkg.Conditions{
		{
			Condition: e.Id < 1,
			Message:   "id is not set (id < 1)",
		},
		{
			Condition: e.SbId == "",
			Message:   "solutionsbüro id is not set",
		},
		{
			Condition: e.Name == "",
			Message:   "family name not set",
		},
	}
}

// Validate the element and return the result.
func (e Employee) Validate() []pkg.ValidateResult {
	return []pkg.ValidateResult{pkg.Check(e)}
}