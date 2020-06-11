package config

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
)

const DefaultConfigFile = "acc.yaml"

var DefaultProjectFiles = []string{
	DefaultConfigFile,
	schema.DefaultExpensesFile,
	schema.DefaultInvoicesFile,
	schema.DefaultMiscRecordsFile,
	schema.DefaultPartiesFile,
	schema.DefaultProjectsFile,
	schema.DefaultStatementFile,
}

// Acc represents an entry point into the utils and also provides general information.
type Acc struct {
	// Company contains the information about the organisation which uses acc.
	Company             schema.Company       `yaml:"company" default:""`
	JournalConfig       schema.JournalConfig `yaml:"journalConfig" default:""`
	ProjectMode         bool                 `yaml:"projectMode" default:"false"`
	ExpensesFilePath    string               `yaml:"expensesFilePath" default:"expenses.yaml"`
	InvoicesFilePath    string               `yaml:"invoicesFilePath" default:"invoices.yaml"`
	MiscRecordsFilePath string               `yaml:"miscRecordsFilePath" default:"misc.yaml"`
	PartiesFilePath     string               `yaml:"partiesFilePath" default:"parties.yaml"`
	ProjectsFilePath    string               `yaml:"projectsFilePath" default:"projects.yaml"`
	StatementFilePath   string               `yaml:"statementFilePath" default:"bank.yaml"`
	fileName            string               `yaml:"-"`
	projectFolder       string               `yaml:"-"`
}

// NewSchema creates a new acc project in the given folder path.
func NewSchema(folderPath, logo string, doSave, interactive bool) schema.Schema {
	var cmp schema.Company
	var jrc schema.JournalConfig
	if interactive {
		cmp = schema.InteractiveNewCompany(logo)
		jrc = schema.InteractiveNewJournalConfig()
	} else {
		cmp = schema.NewCompany(logo)
		jrc = schema.NewJournalConfig()
	}
	acc := Acc{
		Company:             cmp,
		JournalConfig:       jrc,
		ProjectMode:         false,
		ExpensesFilePath:    schema.DefaultExpensesFile,
		InvoicesFilePath:    schema.DefaultInvoicesFile,
		MiscRecordsFilePath: schema.DefaultMiscRecordsFile,
		PartiesFilePath:     schema.DefaultPartiesFile,
		ProjectsFilePath:    schema.DefaultProjectsFile,
		StatementFilePath:   schema.DefaultStatementFile,
		fileName:            DefaultConfigFile,
	}
	exp := schema.NewExpenses(!interactive)
	inv := schema.NewInvoices(!interactive)
	mrc := schema.NewMiscRecords()
	prt := schema.NewParties(!interactive)
	pry := schema.NewProjects()
	stm := schema.NewBankStatement(!interactive)

	if doSave {
		acc.Save(path.Join(folderPath, DefaultConfigFile))
		exp.Save(path.Join(folderPath, schema.DefaultExpensesFile))
		inv.Save(path.Join(folderPath, schema.DefaultInvoicesFile))
		mrc.Save(path.Join(folderPath, schema.DefaultMiscRecordsFile))
		prt.Save(path.Join(folderPath, schema.DefaultPartiesFile))
		pry.Save(path.Join(folderPath, schema.DefaultProjectsFile))
		stm.Save(path.Join(folderPath, schema.DefaultStatementFile))
	}

	return schema.Schema{
		Company:             cmp,
		Expenses:            exp,
		Invoices:            inv,
		JournalConfig:       jrc,
		MiscRecords:         mrc,
		Parties:             prt,
		Projects:            pry,
		Statement:           stm,
		AppendExpenseSuffix: acc.AppendExpensesSuffix,
		AppendInvoiceSuffix: acc.AppendInvoiceSuffix,
		SaveFunc:            acc.SaveSchema,
	}
}

// OpenAcc opens a Acc saved in the json file given by the path.
func OpenAcc(path string) Acc {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	acc := Acc{}
	if err := yaml.Unmarshal(raw, &acc); err != nil {
		logrus.Fatal("error unmarshalling: ", err)
	}
	acc.fileName = filepath.Base(path)
	acc.projectFolder = filepath.Dir(path)
	return acc
}

// OpenSchema reads first the Acc file and then tries to open all linked files.
func OpenSchema(path string) schema.Schema {
	acc := OpenAcc(path)
	return schema.Schema{
		Company:             acc.Company,
		Expenses:            schema.OpenExpenses(acc.ExpensesFilePath),
		Invoices:            schema.OpenInvoices(acc.InvoicesFilePath),
		JournalConfig:       acc.JournalConfig,
		MiscRecords:         schema.OpenMiscRecords(acc.MiscRecordsFilePath),
		Parties:             schema.OpenParties(acc.PartiesFilePath),
		Projects:            schema.OpenProjects(acc.ProjectsFilePath),
		Statement:           schema.OpenBankStatement(acc.StatementFilePath),
		AppendExpenseSuffix: acc.AppendExpensesSuffix,
		AppendInvoiceSuffix: acc.AppendInvoiceSuffix,
		SaveFunc:            acc.SaveSchema,
	}
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (a Acc) Save(path string) {
	schema.SaveToYaml(a, path)
}

func (a Acc) SaveSchema(s schema.Schema) {
	a.SaveSchemaToFolder(s, a.projectFolder)
}

// SaveProjectToFolder saves all files linked in the Acc config to the given folder.
func (a Acc) SaveSchemaToFolder(s schema.Schema, pth string) {
	fmt.Printf("journal cat: %d", len(s.JournalConfig.ExpenseCategories))

	a.Company = s.Company
	a.JournalConfig = s.JournalConfig
	a.Save(path.Join(pth, a.fileName))

	s.Expenses.Save(path.Join(pth, a.ExpensesFilePath))
	s.Invoices.Save(path.Join(pth, a.InvoicesFilePath))
	s.MiscRecords.Save(path.Join(pth, a.InvoicesFilePath))
	s.Parties.Save(path.Join(pth, a.PartiesFilePath))
	s.Projects.Save(path.Join(pth, a.ProjectsFilePath))
	s.Statement.Save(path.Join(pth, a.StatementFilePath))
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
			Condition: a.StatementFilePath == "",
			Message:   "path to bank statement file is not set (BankStatementFilePath is empty)",
		},
	}
}

// Validate the element and return the result.
func (a Acc) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(a)}
}

func (a *Acc) AppendExpensesSuffix(suffix string, overwrite bool) {
	path := appendSuffix(a.ExpensesFilePath, suffix)
	if util.FileExist(path) {
		logrus.Warn("filtered expenses file already exists, use --force to overwrite")
		return
	}
	a.ExpensesFilePath = path
}

func (a *Acc) AppendInvoiceSuffix(suffix string, overwrite bool) {
	path := appendSuffix(a.InvoicesFilePath, suffix)
	if util.FileExist(path) {
		logrus.Warn("filtered invoices file already exists, use --force to overwrite")
		return
	}
	a.InvoicesFilePath = path
}

func appendSuffix(file, suffix string) string {
	ext := path.Ext(file)
	return fmt.Sprintf("%s-%s%s", strings.TrimSuffix(file, ext), suffix, ext)
}
