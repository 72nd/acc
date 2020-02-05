package bimpf

import (
	"fmt"
	"gitlab.com/72th/acc/pkg"
)

// Expense reassembles the structure of a Expense in a Bimpf json dump file.
type Expense struct {
	Id                 int     `json:"id"`
	SbId               string  `json:"sb_id"`
	Name               string  `json:"name"`
	Comment            string  `json:"comment"`
	Path               string  `json:"path"`
	Amount             float64 `json:"amount"`
	DateOfAccrual      string  `json:"date_of_accrual"`
	AdvancedByEmployee bool    `json:"advanced_by_employee"`
	PaidByCustomer     bool    `json:"paid_by_customer"`
	IsPaid             bool    `json:"is_paid"`
	Billable           bool    `json:"billable"`
	EmployeeId         int     `json:"employee_id"`
}

// Type returns a string with the type name of the element.
func (e Expense) Type() string {
	return "SB-Expense"
}

// String returns a human readable representation of the element.
func (e Expense) String() string {
	return fmt.Sprintf("%d/%s (%s) for employee %d", e.Id, e.SbId, e.Name, e.EmployeeId)
}

// Conditions returns the validation conditions.
func (e Expense) Conditions() pkg.Conditions {
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
			Message:   "name not set",
		},
		{
			Condition: e.Path == "",
			Message:   "attachment path not specified",
		},
		{
			Condition: e.EmployeeId < 1,
			Message:   "employee id not set (id < 1)",
		},
		{
			Condition: e.DateOfAccrual == "",
			Message:   "date of accrual not set",
		},
		{
			Condition: e.Amount <= 0,
			Message: "amount not set (amount <= 0)",
		},
	}
}

// Validate the element and return the result.
func (e Expense) Validate() []pkg.ValidateResult {
	return []pkg.ValidateResult{pkg.Check(e)}
}
