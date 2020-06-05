// Schema contains the description of the utils structure of acc.
package schema

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
)

const DefaultAccFile = "acc.yaml"

var DefaultProjectFiles = []string{
	DefaultAccFile,
	DefaultExpensesFile,
	DefaultInvoicesFile,
	DefaultPartiesFile,
	DefaultBankStatementFile,
}

// Acc represents an entry point into the utils and also provides general information.
type Acc struct {
	// Company contains the information about the organisation which uses acc.
	Company               Company       `yaml:"company" default:""`
	JournalConfig         JournalConfig `yaml:"journalConfig" default:""`
	ProjectMode           bool          `yaml:"projectMode" default:"false"`
	ExpensesFilePath      string        `yaml:"expensesFilePath" default:"expenses.yaml"`
	InvoicesFilePath      string        `yaml:"invoicesFilePath" default:"invoices.yaml"`
	PartiesFilePath       string        `yaml:"partiesFilePath" default:"parties.yaml"`
	BankStatementFilePath string        `yaml:"bankStatementFilePath" default:"bank.yaml"`
	Expenses              Expenses      `yaml:"-"`
	Invoices              Invoices      `yaml:"-"`
	Parties               Parties       `yaml:"-"`
	BankStatement         BankStatement `yaml:"-"`
	fileName              string        `yaml:"-"`
	projectFolder         string        `yaml:"-"`
}

// NewProject creates a new acc project in the given folder path.
func NewProject(folderPath, logo string, doSave, interactive bool) Acc {
	var cmp Company
	var jrc JournalConfig
	if interactive {
		cmp = InteractiveNewCompany(logo)
		jrc = InteractiveNewJournalConfig()
	} else {
		cmp = NewCompany(logo)
		jrc = NewJournalConfig()
	}
	acc := Acc{
		Company:               cmp,
		JournalConfig:         jrc,
		ExpensesFilePath:      DefaultExpensesFile,
		InvoicesFilePath:      DefaultInvoicesFile,
		PartiesFilePath:       DefaultPartiesFile,
		BankStatementFilePath: DefaultBankStatementFile,
		fileName:              DefaultAccFile,
	}
	exp := NewExpenses(!interactive)
	inv := NewInvoices(!interactive)
	prt := NewParties(!interactive)
	stm := NewBankStatement(!interactive)

	if doSave {
		acc.Save(path.Join(folderPath, DefaultAccFile))
		exp.Save(path.Join(folderPath, DefaultExpensesFile))
		inv.Save(path.Join(folderPath, DefaultInvoicesFile))
		prt.Save(path.Join(folderPath, DefaultPartiesFile))
		stm.Save(path.Join(folderPath, DefaultBankStatementFile))
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
	if err := yaml.Unmarshal(raw, &acc); err != nil {
		logrus.Fatal(err)
	}
	acc.fileName = filepath.Base(path)
	acc.projectFolder = filepath.Dir(path)
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
func (a Acc) Save(path string) {
	SaveToYaml(a, path)
}

func (a Acc) SaveAtCurrent() {
	fmt.Println(path.Join(a.projectFolder, a.fileName))
	a.Save(path.Join(a.projectFolder, a.fileName))
}

func (a Acc) SaveProject() {
	a.SaveProjectToFolder(a.projectFolder)
}

// SaveProjectToFolder saves all files linked in the Acc config to the given folder.
func (a Acc) SaveProjectToFolder(pth string) {
	a.Save(path.Join(pth, a.fileName))
	a.Expenses.Save(path.Join(pth, a.ExpensesFilePath))
	a.Invoices.Save(path.Join(pth, a.InvoicesFilePath))
	a.Parties.Save(path.Join(pth, a.PartiesFilePath))
	a.BankStatement.Save(path.Join(pth, a.BankStatementFilePath))
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
			Condition: a.ExpensesFilePath == "",
			Message:   "path to expenses file is not set (ExpensesFilePath is empty)",
		},
		{
			Condition: a.InvoicesFilePath == "",
			Message:   "path to invoices file is not set (InvoicesFilePath is empty)",
		},
		{
			Condition: a.PartiesFilePath == "",
			Message:   "path to parties file is not set (PartiesFilePath is empty)",
		},
		{
			Condition: a.BankStatementFilePath == "",
			Message:   "path to bank statement file is not set (BankStatementFilePath is empty)",
		},
	}
}

// Validate the element and return the result.
func (a Acc) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(a)}
}

func (a Acc) ValidateProject() util.ValidateResults {
	results := a.Validate()
	results = append(results, util.Check(a.Company))
	for i := range a.Expenses {
		results = append(results, util.Check(a.Expenses[i]))
	}
	for i := range a.Invoices {
		results = append(results, util.Check(a.Invoices[i]))
	}
	for i := range a.Parties.Customers {
		results = append(results, util.Check(a.Parties.Customers[i]))
	}
	for i := range a.Parties.Employees {
		results = append(results, util.Check(a.Parties.Employees[i]))
	}
	results = append(results, a.BankStatement.Validate()...)
	return results
}

// ValidateAndReportProject validates the acc project files and saves the report to the given path.
func (a Acc) ValidateAndReportProject(path string) {
	rpt := util.Report{
		Title:           "Acc Validation Report",
		ColumnTitles:    []string{"type", "element", "reason"},
		ValidateResults: a.ValidateProject(),
	}
	rpt.Write(path)
}

func (a *Acc) Filter(types []string, from *time.Time, to *time.Time, suffix string, overwrite bool, identifier string) {
	expPath := appendSuffix(a.ExpensesFilePath, suffix)
	invPath := appendSuffix(a.InvoicesFilePath, suffix)
	if util.FileExist(expPath) || util.FileExist(invPath) && !overwrite {
		logrus.Warn("files already exist, use --force to overwrite")
		return
	}
	if util.Contains(types, "expenses") {
		var err error
		a.Expenses, err = a.Expenses.Filter(from, to, identifier)
		if err != nil {
			logrus.Fatal("error while filtering: ", err)
		}
		a.ExpensesFilePath = expPath
	}
	var err error
	a.Invoices, err = a.Invoices.Filter(from, to)
	if err != nil {
		logrus.Fatal("error while filtering: ", err)
	}
	a.InvoicesFilePath = invPath
}

func appendSuffix(file, suffix string) string {
	ext := path.Ext(file)
	return fmt.Sprintf("%s-%s%s", strings.TrimSuffix(file, ext), suffix, ext)
}
