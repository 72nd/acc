// Schema contains the description of the fonts structure of acc.
package schema

import (
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"path"
)

const DefaultAccFile = "acc.json"

var DefaultProjectFiles = []string{
	DefaultAccFile,
	DefaultExpensesFile,
	DefaultInvoicesFile,
	DefaultPartiesFile,
	DefaultBankStatementFile}

// Acc represents an entry point into the fonts and also provides general information.
type Acc struct {
	// Company contains the information about the organisation which uses acc.
	Company               Party         `yaml:"company" default:""`
	ExpensesFilePath      string        `yaml:"expensesFilePath" default:"expenses.json"`
	InvoicesFilePath      string        `yaml:"invoicesFilePath" default:"invoices.json"`
	PartiesFilePath       string        `yaml:"partiesFilePath" default:"parties.json"`
	BankStatementFilePath string        `yaml:"bankStatementFilePath" default:"bank.json"`
	Expenses              Expenses      `yaml:"-"`
	Invoices              Invoices      `yaml:"-"`
	Parties               Parties       `yaml:"-"`
	BankStatement         BankStatement `yaml:"-"`
	fileName              string        `yaml:"-"`
}

// NewAcc returns a new Acc element with the default values.
func NewAcc(useDefaults bool) *Acc {
	acc := &Acc{}
	if err := defaults.Set(acc); err != nil {
		logrus.Fatal(err)
	}
	acc.Company = NewCompanyParty(useDefaults)
	return acc
}

// NewProject creates a new acc project in the given folder path.
func NewProject(folderPath string, doSave, useDefaults bool) Acc {
	acc := Acc{
		Company:               NewCompanyParty(useDefaults),
		ExpensesFilePath:      DefaultExpensesFile,
		InvoicesFilePath:      DefaultInvoicesFile,
		PartiesFilePath:       DefaultPartiesFile,
		BankStatementFilePath: DefaultBankStatementFile,
		fileName:              DefaultAccFile,
	}
	exp := NewExpenses()
	inv := NewInvoices()
	prt := NewParties()
	stm := NewBankStatement()

	if doSave {
		acc.Save(path.Join(folderPath, DefaultAccFile), true)
		exp.Save(path.Join(folderPath, DefaultExpensesFile), true)
		inv.Save(path.Join(folderPath, DefaultInvoicesFile), true)
		prt.Save(path.Join(folderPath, DefaultPartiesFile), true)
		stm.Save(path.Join(folderPath, DefaultBankStatementFile), true)
	}

	return acc
}

// OpenAcc opens a Acc saved in the json file given by the path.
func OpenAcc(path string) Acc {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	acc := Acc{}
	if err := json.Unmarshal(raw, &acc); err != nil {
		logrus.Fatal(err)
	}
	return acc
}

// OpenProject reads first the Acc file and then tries to open all linked files.
func OpenProject(path string) Acc {
	acc := OpenAcc(path)
	acc.Expenses = OpenExpenses(acc.ExpensesFilePath)
	acc.Invoices = OpenInvoices(acc.InvoicesFilePath)
	acc.Parties = OpenParties(acc.PartiesFilePath)
	acc.BankStatement = OpenBankStatement(acc.BankStatementFilePath)
	return acc
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (a Acc) Save(path string, indented bool) {
	SaveToYaml(a, path, indented)
}

// SaveProject saves all files linked in the Acc config.
func (a Acc) SaveProject(pth string, indented bool) {
	a.Save(path.Join(pth, a.fileName), indented)
	a.Expenses.Save(path.Join(pth, a.ExpensesFilePath), indented)
	a.Invoices.Save(path.Join(pth, a.InvoicesFilePath), indented)
	a.Parties.Save(path.Join(pth, a.PartiesFilePath), indented)
	a.BankStatement.Save(path.Join(pth, a.BankStatementFilePath), indented)
}

// Type returns a string with the type name of the element.
func (a Acc) Type() string {
	return "Acc-Main"
}

// String returns a human readable representation of the element.
func (a Acc) String() string {
	return fmt.Sprintf("for company %s", a.Company.Name)
}

// Conditions returns the validation conditions.
func (a Acc) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: a.Company.Name == "",
			Message:   "company name is not set (Company.Name is empty)",
		},
		{
			Condition: a.Company.Street == "",
			Message:   "company street name is not set (Company.Street is empty)",
		},
		{
			Condition: a.Company.StreetNr == 0,
			Message:   "company street number is not set (Company.StreetNr is 0)",
		},
		{
			Condition: a.Company.Place == "",
			Message:   "company place is not set (Company.Place is empty)",
		},
		{
			Condition: a.Company.PostalCode == 0,
			Message:   "company postal code is not set (Company.PostalCode is 0)",
		},
		{
			Condition: a.ExpensesFilePath == "",
			Message: "path to expenses file is not set (ExpensesFilePath is empty)",
		},
		{
			Condition: a.InvoicesFilePath == "",
			Message: "path to invoices file is not set (InvoicesFilePath is empty)",
		},
		{
			Condition: a.PartiesFilePath == "",
			Message: "path to parties file is not set (PartiesFilePath is empty)",
		},
		{
			Condition: a.BankStatementFilePath == "",
			Message: "path to bank statement file is not set (BankStatementFilePath is empty)",
		},
	}
}

// Validate the element and return the result.
func (a Acc) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(a)}
}


func (a Acc) ValidateProject() util.ValidateResults {
	results := a.Validate()
	for i := range a.Expenses {
		results = append(results, util.Check(a.Expenses[i]))
	}
	return results
}

// ValidateAndReportProject validates the acc project files and saves the report to the given path.
func (a Acc) ValidateAndReportProject(path string) {
	rpt := util.Report{
		Title:           "Acc Project Validation Report",
		ColumnTitles:    []string{"type", "element", "reason"},
		ValidateResults: a.ValidateProject(),
	}
	rpt.Write(path)
}

