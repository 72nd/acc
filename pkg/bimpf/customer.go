package bimpf

import (
	"fmt"
	"gitlab.com/72th/acc/pkg/util"
)

// Customer reassembles the structure of a Customer TimeUnit in a Bimpf json dump file.
type Customer struct {
	Id           int       `json:"id"`
	SbId         string    `json:"sb_id"`
	Name         string    `json:"name"`
	Comment      string    `json:"comment"`
	NcFolderName string    `json:"nc_folder_name"`
	Recipient1   string    `json:"recipient_1"`
	Recipient2   string    `json:"recipient_2"`
	Recipient3   string    `json:"recipient_3"`
	Recipient4   string    `json:"recipient_4"`
	Email        string    `json:"email"`
	Projects     []Project `json:"projects"`
}

// Type returns a string with the type name of the element.
func (c Customer) Type() string {
	return "SB-Customer"
}

// String returns a human readable representation of the element.
func (c Customer) String() string {
	return fmt.Sprintf("%d/%s (%s)", c.Id, c.SbId, c.Name)
}

// Conditions returns the validation conditions.
func (c Customer) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: c.Id < 1,
			Message:   "id is not set (id < 1)",
		},
		{
			Condition: c.SbId == "",
			Message:   "solutionsbüro id is not set",
		},
		{
			Condition: c.Name == "",
			Message:   "name not set",
		},
		{
			Condition: c.NcFolderName == "",
			Message: "nextcloud folder not defined",
		},
	}
}

// Validate the element and return the result.
func (c Customer) Validate() util.ValidateResults {
	var results []util.ValidateResult
	for i := range c.Projects {
		results = append(results, util.Check(c.Projects[i]))
	}
	return append(results, util.Check(c))
}