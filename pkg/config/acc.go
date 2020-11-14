package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/72nd/acc/pkg/distributed"
	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
	"github.com/sirupsen/logrus"
)

// DefaultConfigFile states the default name of the config/project file.
const DefaultConfigFile = "acc.yaml"

// DefaultProjectFiles contains the default names of all the project files of flat mode.
// This is used to check if a Acc project is present in a given folder.
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
	Currency            string               `yaml:"currency" default:"CHF"`
	DistributedMode     bool                 `yaml:"distributedMode" default:"false"`
	ExpensesFilePath    string               `yaml:"expensesFilePath" default:"expenses.yaml"`
	InvoicesFilePath    string               `yaml:"invoicesFilePath" default:"invoices.yaml"`
	MiscRecordsFilePath string               `yaml:"miscRecordsFilePath" default:"misc.yaml"`
	PartiesFilePath     string               `yaml:"partiesFilePath" default:"parties.yaml"`
	ProjectsFilePath    string               `yaml:"projectsFilePath" default:"projects.yaml"`
	StatementFilePath   string               `yaml:"statementFilePath" default:"bank.yaml"`
	FileName            string               `yaml:"-"`
}

// NewDistributedModeAcc acc takes a flat file acc configuration file and returns the
// new structure for a config file in project (aka. folder) mode.
func (a Acc) NewDistributedModeAcc(repoPath string) Acc {
	return Acc{
		Company:         a.Company,
		JournalConfig:   a.JournalConfig,
		Currency:        "CHF",
		DistributedMode: true,
		FileName:        filepath.Join(repoPath, DefaultConfigFile),
	}
}

// NewSchema creates a new acc project in the given folder path.
func NewSchema(folderPath, logo string, doSave, interactive, distMode bool) schema.Schema {
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
		DistributedMode:     distMode,
		Currency:            "CHF",
		ExpensesFilePath:    schema.DefaultExpensesFile,
		InvoicesFilePath:    schema.DefaultInvoicesFile,
		MiscRecordsFilePath: schema.DefaultMiscRecordsFile,
		PartiesFilePath:     schema.DefaultPartiesFile,
		ProjectsFilePath:    schema.DefaultProjectsFile,
		StatementFilePath:   schema.DefaultStatementFile,
		FileName:            DefaultConfigFile,
	}
	exp := schema.NewExpenses(!interactive)
	inv := schema.NewInvoices(!interactive)
	mrc := schema.NewMiscRecords()
	prt := schema.NewParties(!interactive)
	prj := schema.NewProjects()
	stm := schema.NewBankStatement(!interactive)

	if doSave && !distMode {
		acc.Save(filepath.Join(folderPath, DefaultConfigFile))
		exp.Save(nil, filepath.Join(folderPath, schema.DefaultExpensesFile))
		inv.Save(filepath.Join(folderPath, schema.DefaultInvoicesFile))
		mrc.Save(filepath.Join(folderPath, schema.DefaultMiscRecordsFile))
		prt.Save(filepath.Join(folderPath, schema.DefaultPartiesFile))
		prj.Save(filepath.Join(folderPath, schema.DefaultProjectsFile))
		stm.Save(filepath.Join(folderPath, schema.DefaultStatementFile))
	} else if doSave && distMode {
		acc = acc.NewDistributedModeAcc(folderPath)
		s := schema.Schema{
			Company:             cmp,
			Expenses:            exp,
			Invoices:            inv,
			JournalConfig:       jrc,
			Currency:            acc.Currency,
			MiscRecords:         mrc,
			Parties:             prt,
			Projects:            prj,
			Statement:           stm,
			AppendExpenseSuffix: acc.AppendExpensesSuffix,
			AppendInvoiceSuffix: acc.AppendInvoiceSuffix,
		}
		acc.Save(acc.FileName)
		distributed.Save(s, s.BaseFolder)
	}

	return schema.Schema{
		Company:             cmp,
		Expenses:            exp,
		Invoices:            inv,
		JournalConfig:       jrc,
		MiscRecords:         mrc,
		Parties:             prt,
		Projects:            prj,
		Statement:           stm,
		AppendExpenseSuffix: acc.AppendExpensesSuffix,
		AppendInvoiceSuffix: acc.AppendInvoiceSuffix,
	}
}

// OpenAcc opens a Acc saved in the json file given by the path.
func OpenAcc(path string) Acc {
	var acc Acc
	path = util.AbsolutePathWithWD(path)
	util.OpenYaml(&acc, path, "acc")
	acc.FileName = path
	return acc
}

// OpenSchema reads first the Acc file and then tries to open all linked files.
func OpenSchema(path string) schema.Schema {
	baseFolder := filepath.Dir(util.AbsolutePathWithWD(path))
	acc := OpenAcc(path)
	if acc.DistributedMode {
		return distributed.Open(baseFolder, acc.Company, acc.JournalConfig, acc.SaveSchema, acc.Currency)
	}
	return schema.Schema{
		Company:             acc.Company,
		Expenses:            schema.OpenExpenses(filepath.Join(baseFolder, acc.ExpensesFilePath)),
		Invoices:            schema.OpenInvoices(filepath.Join(baseFolder, acc.InvoicesFilePath)),
		JournalConfig:       acc.JournalConfig,
		Currency:            acc.Currency,
		MiscRecords:         schema.OpenMiscRecords(filepath.Join(baseFolder, acc.MiscRecordsFilePath)),
		Parties:             schema.OpenParties(filepath.Join(baseFolder, acc.PartiesFilePath)),
		Projects:            schema.OpenProjects(filepath.Join(baseFolder, acc.ProjectsFilePath)),
		Statement:           schema.OpenBankStatement(filepath.Join(baseFolder, acc.StatementFilePath)),
		AppendExpenseSuffix: acc.AppendExpensesSuffix,
		AppendInvoiceSuffix: acc.AppendInvoiceSuffix,
		BaseFolder:          baseFolder,
		SaveFunc:            acc.SaveSchema,
	}
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (a Acc) Save(path string) {
	util.SaveToYaml(a, path, "acc-config")
}

func (a Acc) SaveSchema(s schema.Schema) {
	if a.DistributedMode {
		a.Save(a.FileName)
		distributed.Save(s, s.BaseFolder)
		return
	}
	a.SaveSchemaToFolder(s)
}

// SaveProjectToFolder saves all files linked in the Acc config to the base folder defined in the schema.
func (a Acc) SaveSchemaToFolder(s schema.Schema) {
	a.Company = s.Company
	a.JournalConfig = s.JournalConfig
	a.Save(a.FileName)

	s.Expenses.Save(&s, filepath.Join(s.BaseFolder, a.ExpensesFilePath))
	s.Invoices.Save(filepath.Join(s.BaseFolder, a.InvoicesFilePath))
	s.MiscRecords.Save(filepath.Join(s.BaseFolder, a.MiscRecordsFilePath))
	s.Parties.Save(filepath.Join(s.BaseFolder, a.PartiesFilePath))
	s.Projects.Save(filepath.Join(s.BaseFolder, a.ProjectsFilePath))
	s.Statement.Save(filepath.Join(s.BaseFolder, a.StatementFilePath))
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
	ext := filepath.Ext(file)
	return fmt.Sprintf("%s-%s%s", strings.TrimSuffix(file, ext), suffix, ext)
}
