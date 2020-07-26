package bimpf

import (
	"fmt"
	"os"
	"path"

	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
	"github.com/sirupsen/logrus"
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
	AdvancedByEmployee bool    `json:"advance_by_employee"`
	PaidByCustomer     bool    `json:"paid_by_customer"`
	IsPaid             bool    `json:"is_paid"`
	Billable           bool    `json:"billable"`
	EmployeeId         int     `json:"advanced_employee"`
}

// Type returns a string with the type name of the element.
func (e Expense) Type() string {
	return "SB-Expense"
}

// String returns a human readable representation of the element.
func (e Expense) String() string {
	return fmt.Sprintf("%d/%s (%s)", e.Id, e.SbId, e.Name)
}

// Conditions returns the validation conditions.
func (e Expense) Conditions() util.Conditions {
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
			Message:   "name not set",
			Level:     util.BeforeImportFlaw,
		},
		{
			Condition: e.Path == "",
			Message:   "attachment path not specified", Level: util.BeforeImportFlaw,
		},
		{
			Condition: func() bool {
				if _, err := os.Stat(e.Path); os.IsNotExist(err) {
					return false
				}
				return true
			}(),
			Message: "attachment path doesn't exist",
			Level:   util.BeforeImportFlaw,
		},
		{
			Condition: e.EmployeeId < 1,
			Message:   "employee id not set (id < 1)",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: e.DateOfAccrual == "",
			Message:   "date of accrual not set",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: e.Amount <= 0,
			Message:   "amount not set (amount <= 0)",
			Level:     util.BeforeMergeFlaw,
		},
	}
}

// Validate the element and return the result.
func (e Expense) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(e)}
}

// Convert returns the bimpf Expense as a acc Expense.
// As the path in the Bimpf Expense structure is not absolute (is relative to project folder in a nextcloud folder), a folder prefix is needed.
func (e Expense) Convert(pathPrefix, obligedCustomerId, project string, parties schema.Parties, bimpfEmployees Employees) schema.Expense {
	exp := schema.Expense{
		Identifier:              e.SbId,
		Name:                    e.Name,
		Amount:                  util.NewMoneyFromFloat(e.Amount, "CHF"),
		Path:                    path.Join(pathPrefix, e.Path),
		DateOfAccrual:           e.DateOfAccrual,
		Billable:                e.Billable,
		ObligedCustomer:         schema.NewRef(obligedCustomerId),
		AdvancedByThirdParty:    e.AdvancedByEmployee,
		AdvancedThirdParty:      schema.NewRef(e.getAdvancedPartyId(parties, bimpfEmployees)),
		DateOfSettlement:        "",
		SettlementTransaction: "",
		ProjectName:             project,
	}
	exp.SetId()
	return exp
}

// getAdvancedPartyId tries to find the Acc Id associated advancing employee.
// In Bimpf only employees can advance payments.
func (e Expense) getAdvancedPartyId(parties schema.Parties, bimpfEmployees Employees) string {
	if e.AdvancedByEmployee {
		bimpfEmployee, err := bimpfEmployees.ById(e.EmployeeId)
		if err != nil {
			logrus.Warn(err)
		}
		employee, err := parties.EmployeeByIdentifier(bimpfEmployee.SbId)
		if err != nil {
			logrus.Warn(err)
		}
		return employee.Id
	}
	return ""
}
