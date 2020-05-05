package schema

import (
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"io/ioutil"
	"strings"
)

const DefaultExpensesFile = "expenses.json"

// Expenses is a slice of multiple expenses.
type Expenses []Expense

// NewExpenses returns a new Expense slice with the one Expense in it.
func NewExpenses() Expenses {
	return Expenses{NewExpense()}
}

// OpenExpenses opens a Expenses saved in the json file given by the path.
func OpenExpenses(path string) Expenses {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	exp := Expenses{}
	if err := json.Unmarshal(raw, &exp); err != nil {
		logrus.Fatal(err)
	}
	return exp
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (e Expenses) Save(path string, indented bool) {
	SaveToJson(e, path, indented)
}

// SetId sets a unique id to all elements in the slice.
func (e Expenses) SetId() {
	for i := range e {
		e[i].SetId()
	}
}

// Expense represents a payment done by the company or a third party to assure the ongoing of the business.
type Expense struct {
	// Id is the internal unique identifier of the Expense.
	Id string `json:"id" default:"1"`
	// Identifier is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `json:"identifier" default:"e-19-1"`
	// Name describes meaningful the kind of the Expense.
	Name string `json:"name" default:"Expense Name"`
	// Amount states the amount of the Expense.
	Amount float64 `json:"amount" default:"10.00"`
	// Path is the full path to the voucher pdf.
	Path string `json:"path" default:"/path/to/file.pdf"`
	// DateOfAccrual represents the day the obligation emerged.
	DateOfAccrual string `json:"dateOfAccrual" default:"2019-12-20"`
	// Billable states if the costs for the Expense will be forwarded to the customer.
	Billable bool `json:"billable" default:"false"`
	// ObligedCustomerId refers to the customer which have to pay the Expense.
	ObligedCustomerId string `json:"obligedCustomerId" default:""`
	// AdvancedByThirdParty states if a third party (employee, etc.) advanced the payment of this expense for the company.
	AdvancedByThirdParty bool `json:"advancedByThirdParty" default:"false"`
	// AdvancePartyId refers to the third party which advanced the payment.
	AdvancedThirdPartyId string `json:"advancedThirdPartyId" default:""`
	// DateOfSettlement states the date of the settlement of the expense (the company has not to take further actions).
	DateOfSettlement string `json:"dateOfSettlement" default:"2019-12-25"`
	// SettlementTransactionId refers to a possible bank transaction which settled the Expense for the company.
	SettlementTransactionId string `json:"settlementTransactionId" default:""`
	// ProjectName refers to the associated project of the expense.
	ProjectName string `json:"projectName" default:""`
}

// NewExpense returns a new Expense element with the default values.
func NewExpense() Expense {
	exp := Expense{}
	if err := defaults.Set(&exp); err != nil {
		logrus.Fatal(err)
	}
	return exp
}

// GetId returns the unique id of the element.
func (e Expense) GetId() string {
	return e.Id
}

// SetId generates a unique id for the element if there isn't already one defined.
func (e *Expense) SetId() {
	if e.Id != "" {
		return
	}
	e.Id = uuid.Must(uuid.NewRandom()).String()
}

// Type returns a string with the type name of the element.
func (e Expense) Type() string {
	return ""
}

// String returns a human readable representation of the element.
func (e Expense) String() string {
	return fmt.Sprintf("")
}

// FileString returns the file name for exporting the expense as a document.
func (e Expense) FileString() string {
	result := e.Identifier
	result = strings.ReplaceAll(result, " ", "-")
	result = strings.ReplaceAll(result, ".", "-")
	return result
}

// Conditions returns the validation conditions.
func (e Expense) Conditions() util.Conditions {
	return util.Conditions{

	}
}

// Validate the element and return the result.
func (e Expense) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(e)}
}