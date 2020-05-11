package schema

import (
	"bufio"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const DefaultExpensesFile = "expenses.yaml"
const DefaultExpensePrefix = "e-"

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
		logrus.Fatalf("error while reading file %s: %s", path, err)
	}
	exp := Expenses{}
	if err := yaml.Unmarshal(raw, &exp); err != nil {
		logrus.Fatalf("error reading (unmarshalling) YAML file %s: %s", path, err)
	}
	return exp
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (e Expenses) Save(path string) {
	SaveToYaml(e, path)
}

// SetId sets a unique id to all elements in the slice.
func (e Expenses) SetId() {
	for i := range e {
		e[i].SetId()
	}
}

func (e Expenses) GetIdentifiables() []Identifiable {
	exp := make([]Identifiable, len(e))
	for i := range e {
		exp[i] = e[i]
	}
	return exp
}

func (e Expenses) GetExpenseCategories() util.SearchItems {
	var result util.SearchItems
	for i := range e {
		if e[i].ExpenseCategory == "" {
			continue
		}
		existing := false
		for j := range result {
			if e[i].ExpenseCategory == result[j].Identifier {
				existing = true
			}
		}
		if !existing {
			result = append(result, util.SearchItem{
				Name:       e[i].ExpenseCategory,
				Identifier: e[i].ExpenseCategory,
				Value:      e[i].ExpenseCategory,
			})
		}
	}
	return result
}

// Expense represents a payment done by the company or a third party to assure the ongoing of the business.
type Expense struct {
	// Id is the internal unique identifier of the Expense.
	Id string `yaml:"id" default:"1"`
	// Identifier is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `yaml:"identifier" default:"e-19-1"`
	// Name describes meaningful the kind of the Expense.
	Name string `yaml:"name" default:"Expense Name"`
	// Amount states the amount of the Expense.
	Amount float64 `yaml:"amount" default:"10.00"`
	// Path is the full path to the voucher utils.
	Path string `yaml:"path" default:"/path/to/file.utils"`
	// DateOfAccrual represents the day the obligation emerged.
	DateOfAccrual string `yaml:"dateOfAccrual" default:"2019-12-20"`
	// Billable states if the costs for the Expense will be forwarded to the customer.
	Billable bool `yaml:"billable" default:"false"`
	// ObligedCustomerId refers to the customer which have to pay the Expense.
	ObligedCustomerId string `yaml:"obligedCustomerId" default:""`
	// AdvancedByThirdParty states if a third party (employee, etc.) advanced the payment of this expense for the company.
	AdvancedByThirdParty bool `yaml:"advancedByThirdParty" default:"false"`
	// AdvancePartyId refers to the third party which advanced the payment.
	AdvancedThirdPartyId string `yaml:"advancedThirdPartyId" default:""`
	// DateOfSettlement states the date of the settlement of the expense (the company has not to take further actions).
	DateOfSettlement string `yaml:"dateOfSettlement" default:"2019-12-25"`
	// SettlementTransactionId refers to a possible bank transaction which settled the Expense for the company.
	SettlementTransactionId string `yaml:"settlementTransactionId" default:""`
	// ExpenseCategory gives additional info for the categorization of the expense in the journal.
	ExpenseCategory string `yaml:"expenseCategory" default:""`
	// ProjectName refers to the associated project of the expense.
	ProjectName string `yaml:"projectName" default:""`
}

// NewExpense returns a new Expense element with the default values.
func NewExpense() Expense {
	exp := Expense{}
	if err := defaults.Set(&exp); err != nil {
		logrus.Fatal(err)
	}
	return exp
}

func NewExpenseWithUuid() Expense {
	exp := NewExpense()
	exp.Id = GetUuid()
	return exp
}

// InteractiveNewExpense returns a new Expense based on the user input.
func InteractiveNewExpense(a Acc, asset string) Expense {
	reader := bufio.NewReader(os.Stdin)
	exp := NewExpenseWithUuid()
	exp.Identifier = util.AskString(
		reader,
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(a.Expenses.GetIdentifiables(), DefaultExpensePrefix),
	)
	exp.ExpenseCategory = util.AskStringFromListSearch(
		reader,
		"Expense Category",
		"Used for journal generation",
		a.Expenses.GetExpenseCategories(),
	)
	exp.Name = util.AskString(reader,
		"Name",
		"Name of the expense",
		"HAL 9000",
	)
	exp.Amount = util.AskFloat(reader,
		"Amount",
		"How much did you spend?",
		23.42,
	)
	if asset == "" {
		exp.Path = util.AskString(
			reader,
			"Asset",
			"Path to asset file (use --asset to set with flag)",
			"",
		)
	} else {
		exp.Path = asset
	}
	exp.DateOfAccrual = util.AskDate(reader,
		"Date of Accrual",
		"Date when the obligation accrued",
		time.Now(),
	)
	exp.Billable = util.AskBool(
		reader,
		"Billable?",
		"Is expense billable to customer?",
		false,
	)
	if exp.Billable {
		exp.ObligedCustomerId = util.AskStringFromSearch(
			reader,
			"Obliged Customer",
			"Customer which has to pay this expense",
			a.Parties.CustomersSearchItems(),
		)
	} else {
		exp.ObligedCustomerId = ""
	}
	exp.AdvancedByThirdParty = util.AskBool(
		reader,
		"Advanced?",
		"Was this expense advanced by some third party (ex: employee)?",
		false,
	)
	if exp.AdvancedByThirdParty {
		exp.AdvancedThirdPartyId = util.AskStringFromSearch(
			reader,
			"Advanced party",
			"Employee which advanced the expense",
			a.Parties.EmployeesSearchItems(),
		)
	}
	exp.DateOfSettlement = util.AskDate(
		reader,
		"Date of settlement",
		"Date when the obligation was settelt for the company",
		time.Now(),
	)
	exp.ProjectName = util.AskString(
		reader,
		"Project Name",
		"Name of the associated project",
		"",
	)
	return exp
}

// GetId returns the unique id of the element.
func (e Expense) GetId() string {
	return e.Id
}

func (e Expense) GetIdentifier() string {
	return e.Identifier
}

// SetId generates a unique id for the element if there isn't already one defined.
func (e *Expense) SetId() {
	if e.Id != "" {
		return
	}
	e.Id = GetUuid()
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
