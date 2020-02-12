// Schema contains the description of the data structure of acc.
package schema

import (
	"encoding/json"
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
)

const DefaultAccFile = "acc.json"

var DefaultProjectFiles = []string{
	DefaultAccFile,
	DefaultExpensesFile,
	DefaultInvoicesFile,
	DefaultPartiesFile,
	DefaultBankStatementFile}

// Acc represents an entry point into the data and also provides general information.
type Acc struct {
	// Company contains the information about the organisation which uses acc.
	Company               Party         `json:"company" default:""`
	ExpensesFilePath      string        `json:"expensesFilePath" default:"expenses.json"`
	InvoicesFilePath      string        `json:"invoicesFilePath" default:"invoices.json"`
	PartiesFilePath       string        `json:"partiesFilePath" default:"parties.json"`
	BankStatementFilePath string        `json:"bankStatementFilePath" default:"bank.json"`
	Expenses              Expenses      `json:"-"`
	Invoices              Invoices      `json:"-"`
	Parties               Parties       `json:"-"`
	BankStatement         BankStatement `json:"-"`
	fileName              string        `json:"-"`
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
	SaveToJson(a, path, indented)
}

// SaveProject saves all files linked in the Acc config.
func (a Acc) SaveProject(pth string, indented bool) {
	a.Save(path.Join(pth, a.fileName), indented)
	a.Expenses.Save(path.Join(pth, a.ExpensesFilePath), indented)
	a.Invoices.Save(path.Join(pth, a.InvoicesFilePath), indented)
	a.Parties.Save(path.Join(pth, a.PartiesFilePath), indented)
	a.BankStatement.Save(path.Join(pth, a.BankStatementFilePath), indented)
}
