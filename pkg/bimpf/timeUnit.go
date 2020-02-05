package bimpf

import (
	"fmt"
	"gitlab.com/72th/acc/pkg/util"
)

// TimeUnit reassembles the structure of a TimeUnit in a Bimpf json dump file.
type TimeUnit struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	StartTime   string `json:"start_time"`
	EndDate     string `json:"end_date"`
	EndTime     string `json:"end_time"`
	Billable    bool   `json:"billable"`
}

// Type returns a string with the type name of the element.
func (t TimeUnit) Type() string {
	return "SB-TimeUnit"
}

// String returns a human readable representation of the element.
func (t TimeUnit) String() string {
	return fmt.Sprintf("%d (%s)", t.Id, t.Description)
}

// Conditions returns the validation conditions.
func (t TimeUnit) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: t.Id < 1,
			Message:   "id is not set (id < 1)",
		},
		{
			Condition: t.Description == "",
			Message:   "description not set",
		},
		{
			Condition: t.StartDate == "",
			Message:   "start date not set",
		},
		{
			Condition: t.StartTime == "",
			Message:   "start time not set",
		},
		{
			Condition: t.EndDate == "",
			Message:   "end date not set",
		},
		{
			Condition: t.EndTime == "",
			Message:   "end time not set",
		},
	}
}

// Validate the element and return the result.
func (t TimeUnit) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(t)}
}
