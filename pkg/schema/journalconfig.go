package schema

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
)

type JournalConfig struct {
	BankAccount                string            `yaml:"bankAccount" default:"assets:Umlaufvermögen:Flüssige Mittel:Raiffeisenbank Bern"`
	ReceivableAccount          string            `yaml:"receivableAccount" default:"assets:Umlaufvermögen:Debitoren"`
	RevenueAccount             string            `yaml:"revenueAccount" default:"revenues:Betrieblicher Ertrag:Dienstleistungserlös"`
	PayableAccount             string            `yaml:"payableAccount" default:"liabilities:Kurzfristiges Fremdkapital:Kreditoren"`
	EmployeeLiabilitiesAccount string            `yaml:"employeeLiabilitiesAccount" default:"liabilities:Kurzfristiges Fremdkapital:Verbindlichkeiten gegenüber Genossenschaftler"`
	ExpenseCategories          ExpenseCategories `yaml:"expenseCategories" default:"[]"`
}

func NewJournalConfig() JournalConfig {
	jrc := JournalConfig{}
	if err := defaults.Set(&jrc); err != nil {
		logrus.Fatal("error setting defaults for journal config: ", err)
	}
	jrc.ExpenseCategories = ExpenseCategories{NewExpenseCategory()}
	return jrc
}

func InteractiveNewJournalConfig() JournalConfig {
	jrc := NewJournalConfig()
	jrc.BankAccount = util.AskString(
		"Bank Account",
		"Ledger account of your bank account",
		jrc.BankAccount)
	jrc.ReceivableAccount = util.AskString(
		"Receivable Account",
		"Ledger account for receivables (debitors)",
		jrc.ReceivableAccount)
	jrc.RevenueAccount = util.AskString(
		"Revenue Account",
		"Default ledger account for earnings",
		jrc.RevenueAccount)
	jrc.PayableAccount = util.AskString(
		"Payable Account",
		"Ledger account for payables",
		jrc.PayableAccount)
	jrc.EmployeeLiabilitiesAccount = util.AskString(
		"Emloyee Liabilities Account",
		"Ledger Account for unpaid liabilities against employees",
		jrc.EmployeeLiabilitiesAccount)
	jrc.ExpenseCategories = ExpenseCategories{}
	return jrc
}

func (JournalConfig) Type() string {
	return "Journal-Config"
}

func (c JournalConfig) String() string {
	return "journal config"
}

func (c JournalConfig) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: c.BankAccount == "",
			Message:   "bank account is not set (BankAccount is empty)",
		},
		{
			Condition: c.ReceivableAccount == "",
			Message:   "receivable account is not set (ReceivableAccount is empty)",
		},
		{
			Condition: c.RevenueAccount == "",
			Message:   "revenue account is not set (ReceivableAccount is empty)",
		},
		{
			Condition: c.PayableAccount == "",
			Message:   "payable account is not set (PayableAccount is empty)",
		},
		{
			Condition: c.EmployeeLiabilitiesAccount == "",
			Message:   "employee liabilities account is not set (EmployeeLiabilitiesAccount is empty",
		},
	}
}

func (c JournalConfig) Validate() util.ValidateResults {
	return append(util.ValidateResults{util.Check(c)}, c.ExpenseCategories.Validate()...)
}

type ExpenseCategories []ExpenseCategory

func InteractiveNewExpenseCategories(multiple bool) ExpenseCategories {
	cat := ExpenseCategories{InteractiveNewExpenseCategory()}
	if multiple && util.AskBool("Continue", "Add another expense categries?", false) {
		return append(cat, InteractiveNewExpenseCategories(multiple)...)
	}
	return cat
}

func (e ExpenseCategories) CategoryByName(name string) (*ExpenseCategory, error) {
	for i := range e {
		if e[i].Name == name {
			return &e[i], nil
		}
	}
	return nil, fmt.Errorf("no epense category for name «%s» found", name)
}

func (e ExpenseCategories) SearchItems() util.SearchItems {
	result := make(util.SearchItems, len(e))
	for i := range e {
		result[i] = e[i].SearchItem()
	}
	return result
}

func (e ExpenseCategories) Type() string {
	return "Expense-Categories"
}

func (e ExpenseCategories) String() string {
	return "expense categories"
}

func (e ExpenseCategories) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: func() bool {
				for i := range e {
					for j := i + 1; j < len(e); j++ {
						if e[i].Name == e[j].Name {
							return true
						}
					}
				}
				return false
			}(),
			Message: "same name is used for multiple expense categories",
		},
		{
			Condition: func() bool {
				for i := range e {
					for j := i + 1; j < len(e); j++ {
						if e[i].Account == e[j].Account {
							return true
						}
					}
				}
				return false
			}(),
			Message: "same account is used for multiple expense categories",
		},
	}
}

func (e ExpenseCategories) Validate() util.ValidateResults {
	result := util.ValidateResults{util.Check(e)}
	for i := range e {
		result = append(result, util.Check(e[i]))
	}
	return result
}

type ExpenseCategory struct {
	Name    string `yaml:"name" default:"Matrial Costs"`
	Account string `yaml:"account" default:"expenses:Betrieblicher Aufwand:Materialaufwand"`
}

func NewExpenseCategory() ExpenseCategory {
	cat := ExpenseCategory{}
	if err := defaults.Set(&cat); err != nil {
		logrus.Fatal("error setting defaults for expense category: ", err)
	}
	return cat
}

func InteractiveNewExpenseCategory() ExpenseCategory {
	cat := ExpenseCategory{}
	cat.Name = util.AskString(
		"Name",
		"Name of expense category",
		cat.Name)
	cat.Account = util.AskString(
		"Account",
		"Ledger account for expense category",
		cat.Account)
	return cat
}

func InteractiveNewGenericExpenseCategory() interface{} {
	return InteractiveNewExpenseCategory()
}

func (e ExpenseCategory) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:        e.Name,
		Value:       e.Name,
		SearchValue: fmt.Sprintf("%s %s", e.Name, e.Account)}
}

func (e ExpenseCategory) Type() string {
	return "Expense-Category"
}

func (e ExpenseCategory) String() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Account)
}

func (e ExpenseCategory) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: e.Name == "",
			Message:   "name is not set (Name is empty)",
		},
		{
			Condition: e.Account == "",
			Message:   "account is not set (Account is empty)",
		},
	}
}
