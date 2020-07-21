package schema

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/72nd/acc/pkg/util"
	"github.com/creasty/defaults"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

const DefaultExpensesFile = "expenses.yaml"
const DefaultExpensePrefix = "e-"

// Expenses is a slice of multiple expenses.
type Expenses []Expense

// NewExpenses returns a new Expense slice with the one Expense in it.
func NewExpenses(useDefaults bool) Expenses {
	if useDefaults {
		return Expenses{NewExpense()}
	}
	return Expenses{}
}

// OpenExpenses opens a Expenses saved in the YAML file given by the path.
func OpenExpenses(path string) Expenses {
	var exp Expenses
	util.OpenYaml(&exp, path, "expenses")
	return exp
}

// Save writes the element as a YAML to the given path.
// If s is not nil, the id fields will be commented with the linked element.
func (e Expenses) Save(s *Schema, path string) {
	/*
		if s != nil {
			e.CommentNodes(*s)
		}
	*/
	util.SaveToYaml(e, path, "expenses")
}

// SetId sets a unique id to all elements in the slice.
func (e Expenses) SetId() {
	for i := range e {
		e[i].SetId()
	}
}

func (e Expenses) ExpenseById(id string) (*Expense, error) {
	for i := range e {
		if e[i].Id == id {
			return &e[i], nil
		}
	}
	return nil, fmt.Errorf("no expense for id \"%s\" found", id)
}

func (e Expenses) ExpenseByIdent(ident string) (*Expense, error) {
	for i := range e {
		if e[i].Identifier == ident {
			return &e[i], nil
		}
	}
	return nil, fmt.Errorf("no expense for identifier \"%s\" found", ident)
}

func (e Expenses) GetIdentifiables() []Identifiable {
	exp := make([]Identifiable, len(e))
	for i := range e {
		exp[i] = e[i]
	}
	return exp
}

func (e Expenses) SearchItems() util.SearchItems {
	result := make(util.SearchItems, len(e))
	for i := range e {
		result[i] = e[i].SearchItem()
	}
	return result
}

func (e Expenses) Validate() util.ValidateResults {
	var rsl util.ValidateResults
	for i := range e {
		rsl = append(rsl, util.Check(e[i]))
	}
	return rsl
}

func (e Expenses) Filter(from *time.Time, to *time.Time, identifier string) (Expenses, error) {
	var result Expenses
	for i := range e {
		ok, err := e[i].Match(from, to, identifier)
		if err != nil {
			logrus.Errorf("error while matching \"%s\": %s", e[i].String(), err)
		}
		if ok {
			result = append(result, e[i])
		}
	}
	return result, nil
}

func (e Expenses) AssistedCompletion(s *Schema, doAll, autoSave, openAttachment, retainFocus bool) {
	first := true
	for i := range e {
		if !first {
			fmt.Println()
		} else {
			first = false
		}
		e[i] = e[i].AssistedCompletion(s, doAll, autoSave, openAttachment, retainFocus)
		if autoSave {
			s.Save()
		}
	}
}

func (e Expenses) Repopulate(s Schema) {
	for i := range e {
		e[i].Repopulate(s)
	}
}

func (e Expenses) Len() int {
	return len(e)
}

func (e Expenses) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Expenses) Less(i, j int) bool {
	di, _ := time.Parse(util.DateFormat, e[i].DateOfAccrual)
	dj, _ := time.Parse(util.DateFormat, e[j].DateOfAccrual)
	return di.Before(dj)
}

// SortByYear returns a map with the given elements sorted by year.
func (e Expenses) SortByYear() map[int]Expenses {
	rsl := make(map[int]Expenses)
	for i := range e {
		if e[i].GetDate() == nil {
			rsl[0] = append(rsl[0], e[i])
		}
		year := e[i].GetDate().Year()
		rsl[year] = append(rsl[year], e[i])
	}
	return rsl
}

/*
func (e Expenses) CommentNodes(s Schema) Expenses {
	rsl := make(Expenses, len(e))
	for i := range e {
		rsl[i] = e[i].CommentNodes(s)
	}
	return rsl
}
*/

// Expense represents a payment done by the company or a third party to assure the ongoing of the business.
type Expense struct {
	// Id is the internal unique identifier of the Expense.
	Id string `yaml:"id" default:"1"`
	// Value is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `yaml:"identifier" default:"e-19-1"`
	// Name describes meaningful the kind of the Expense.
	Name string `yaml:"name" default:"Expense Name"`
	// Amount states the amount of the Expense.
	Amount util.Money `yaml:"amount" default:"" query:"amount"`
	// Path is the full path to the business record document.
	Path string `yaml:"path" default:"/path/to/expense.pdf" query:"path"`
	// DateOfAccrual represents the day the obligation emerged.
	DateOfAccrual string `yaml:"dateOfAccrual" default:"2019-12-20"`
	// Billable states if the costs for the Expense will be forwarded to the customer.
	Billable bool `yaml:"billable" default:"false"`
	// ObligedCustomerId refers to the customer which have to pay the Expense.
	ObligedCustomerId string `yaml:"obligedCustomerId" default:"" query:"customer"`
	// AdvancedByThirdParty states if a third party (employee, etc.) advanced the payment of this expense for the company.
	AdvancedByThirdParty bool `yaml:"advancedByThirdParty" default:"false"`
	// AdvancePartyId refers to the third party which advanced the payment.
	AdvancedThirdPartyId string `yaml:"advancedThirdPartyId" default:"" query:"emplyee"`
	// DateOfSettlement states the date of the settlement of the expense (the company has not to take further actions).
	DateOfSettlement string `yaml:"dateOfSettlement" default:"2019-12-25"`
	// SettlementTransactionId refers to a possible bank transaction which settled the Expense for the company.
	SettlementTransactionId string `yaml:"settlementTransactionId" default:"" query:"transaction"`
	// ExpenseCategory gives additional info for the categorization of the expense in the journal.
	ExpenseCategory string `yaml:"expenseCategory" default:""`
	// Debit Payment states whether the expense was directly paid with the main account debithether the expense was directly paid with the main account debit card.
	PaidWithDebit bool `yaml:"paidWithDebit" default:"false"`
	// Internal states whether this expense is for an internal purpose or not.
	Internal bool `yaml:"internal" default:"true"`
	// ProjectName refers to the associated project of the expense. Depreciated.
	ProjectName string `yaml:"projectName" default:""`
	// ProjectId refers to the associated project.
	ProjectId string `yaml:"projectId" default:""`
}

// NewExpense returns a new Expense element with the default values.
func NewExpense() Expense {
	exp := Expense{}
	exp.Id = GetUuid()
	if err := defaults.Set(&exp); err != nil {
		logrus.Fatal("error setting defaults: ", err)
	}
	exp.Amount = util.NewMoney(1000, "CHF")
	return exp
}

func NewExpenseWithUuid() Expense {
	exp := NewExpense()
	exp.Id = GetUuid()
	return exp
}

// InteractiveNewExpense returns a new Expense based on the user input.
func InteractiveNewExpense(s *Schema, asset string) Expense {
	exp := NewExpenseWithUuid()
	exp.Identifier = util.AskString(
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.Expenses.GetIdentifiables(), DefaultExpensePrefix))
	exp.Name = util.AskString(
		"Name",
		"Name of the expense",
		"HAL 9000")
	exp.Amount = util.AskMoney(
		"Amount",
		"How much did you spend?",
		util.NewMoney(2342, s.Currency),
		s.Currency)
	if asset == "" {
		exp.Path = util.AskString(
			"Asset",
			"Path to asset file (use --asset to set with flag)",
			"")
	} else {
		exp.Path = asset
	}
	exp.DateOfAccrual = util.AskDate(
		"Date of Accrual",
		"Date when the obligation accrued",
		time.Now(),
	)
	exp.Billable = util.AskBool(
		"Billable?",
		"Is expense billable to customer?",
		false)
	if exp.Billable {
		exp.ObligedCustomerId = util.AskStringFromSearch(
			"Obliged Customer",
			"Customer which has to pay this expense",
			s.Parties.CustomersSearchItems())
	} else {
		exp.ObligedCustomerId = ""
	}
	exp.AdvancedByThirdParty = util.AskBool(
		"Advanced?",
		"Was this expense advanced by some third party (ex: employee)?",
		false)
	if exp.AdvancedByThirdParty {
		exp.AdvancedThirdPartyId = util.AskStringFromSearch(
			"Advanced party",
			"Employee which advanced the expense",
			s.Parties.EmployeesSearchItems())
	}
	exp.DateOfSettlement = util.AskDate(
		"Date of settlement",
		"Date when the obligation was settelt for the company",
		time.Now())
	var cat interface{}
	exp.ExpenseCategory, cat = util.AskStringFromSearchWithNew(
		"Expense Category",
		"Used for journal genertaion",
		s.JournalConfig.ExpenseCategories.SearchItems(),
		InteractiveNewGenericExpenseCategory,
		nil)
	if cat != nil {
		value, ok := cat.(ExpenseCategory)
		if !ok {
			logrus.Fatal("returned new expense category has different type")
		}
		s.JournalConfig.ExpenseCategories = append(s.JournalConfig.ExpenseCategories, value)
		exp.ExpenseCategory = value.Name
	}
	exp.PaidWithDebit = util.AskBool(
		"Paid with Debit",
		"Was this expense directly paid via the main account debit card?",
		false)
	exp.Internal = util.AskBool(
		"Internal",
		"Has this expense an internal prupose?",
		false)
	if !exp.Internal {
		exp.ProjectId = util.AskStringFromSearch(
			"Project",
			"Associated Project",
			s.Projects.SearchItems())
	}
	return exp
}

func (e Expense) AssistedCompletion(s *Schema, doAll, autoSave, openAttachment, retainFocus bool) Expense {
	if !doAll && util.Check(e).Valid() {
		fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Skip expense:")), aurora.BrightMagenta(e.String()))
		return e
	}
	tmp := e
	var ext util.External
	if e.Path != "" && openAttachment {
		ext = util.NewExternal(e.Path, retainFocus)
		ext.Open()
	}
	fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Optimize expense:")), aurora.BrightMagenta(e.String()))
	if e.AdvancedByThirdParty && e.AdvancedThirdPartyId == "" {
		e.AdvancedThirdPartyId = util.AskStringFromListSearch(
			"Advanced party",
			"Employee which advanced the expense",
			s.Parties.EmployeesSearchItems())
	}
	if e.ExpenseCategory == "" {
		var cat interface{}
		e.ExpenseCategory, cat = util.AskStringFromSearchWithNew(
			"Expense Category",
			"Used for journal generation",
			s.JournalConfig.ExpenseCategories.SearchItems(),
			InteractiveNewGenericExpenseCategory,
			nil)
		if cat != nil {
			value, ok := cat.(ExpenseCategory)
			if !ok {
				logrus.Fatal("returned new expense category has different type")
			}
			s.JournalConfig.ExpenseCategories = append(s.JournalConfig.ExpenseCategories, value)
			if autoSave {
				s.Save()
			}
			e.ExpenseCategory = value.Name
		}
	}
	/* Why should we do this?
	e.PaidWithDebit = util.AskBool(
		"Paid with Debit",
		"Was this expense directly paid via the main account debit card?",
		false)
	e.Internal = util.AskBool(
		"Internal",
		"Has this expense an internal prupose?",
		false)
	*/

	strategy := util.AskForStategy()
	switch strategy {
	case util.RedoStrategy:
		exp := e.AssistedCompletion(s, doAll, autoSave, openAttachment, retainFocus)
		ext.Close()
		return exp
	case util.SkipStrategy:
		return tmp
	}
	ext.Close()
	return e
}

func (e *Expense) Repopulate(s Schema) {
	trn, err := s.Statement.TransactionForDocument(e.Id)
	if err != nil {
		logrus.Warnf("there is no transaction for expense \"%s\" associated", e.String())
		return
	}
	e.DateOfSettlement = trn.Date
	e.SettlementTransactionId = trn.Id
}

func (e Expense) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:        fmt.Sprintf("%s for %.2f", e.Name, e.Amount),
		Type:        e.Type(),
		Value:       e.Id,
		SearchValue: fmt.Sprintf("%s %s %s", e.Name, e.Identifier, e.ProjectName),
	}
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
	return "Expense"
}

// String returns a human readable representation of the element.
func (e Expense) String() string {
	return fmt.Sprintf("%s (%s): %s for %s", e.Name, e.Identifier, e.Amount.Display(), e.ProjectName)
}

// Short returns a short representation of the element.
func (e Expense) Short() string {
	return fmt.Sprintf("%s (%s)", e.Name, e.Identifier)
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
		{
			Condition: e.Id == "",
			Message:   "unique identifier not set (Id is empty)",
		},
		{
			Condition: e.Identifier == "",
			Message:   "human readable identifier not set (Identifier is empty)",
		},
		{
			Condition: e.Amount.Amount() == 0,
			Message:   "amount is not set (Amount is 0.0)",
		},
		{
			Condition: !util.FileExist(e.Path),
			Message:   fmt.Sprintf("business record document at «%s» not found", e.Path),
		},
		{
			Condition: !util.ValidDate(util.DateFormat, e.DateOfAccrual),
			Message:   fmt.Sprintf("string «%s» could not be parsed with format YYYY-MM-DD", e.DateOfAccrual),
		},
		{
			Condition: e.Billable && e.ObligedCustomerId == "",
			Message:   "although billable, no obliged customer is set (ObligedCustomerId is empty)",
		},
		{
			Condition: e.AdvancedByThirdParty && e.AdvancedThirdPartyId == "",
			Message:   "although advanced by third party, no third party id is set (AdvancedThirdPartyId is empty)",
		},
		{
			Condition: e.DateOfSettlement != "" && util.ValidDate(util.DateFormat, e.DateOfSettlement),
			Message:   fmt.Sprintf("string «%s» could not be parsed with format YYYY-MM-DD", e.DateOfSettlement),
		},
		{
			Condition: e.DateOfSettlement != "" && e.SettlementTransactionId == "",
			Message:   "although date of settlement is set, the corresponding transaction is empty (SettlementTransactionId is empty",
		},
		{
			Condition: e.ExpenseCategory == "",
			Message:   "expense category is not set (ExpenseCategory is empty)",
		},
		{
			Condition: !e.Internal && e.ProjectId == "",
			Message:   "altrough not an internal expense, project id is not set (ProjectId is empty)",
		},
	}
}

func (e Expense) AccrualDateTime() time.Time {
	result, err := time.Parse(util.DateFormat, e.DateOfAccrual)
	if err != nil {
		logrus.Fatalf("could not parse «%s» as date with YYYY-MM-DD: %s", e.DateOfAccrual, err)
	}
	return result
}

func (e Expense) Match(from *time.Time, to *time.Time, identifier string) (bool, error) {
	date, err := time.Parse(util.DateFormat, e.DateOfAccrual)
	if err != nil {
		return false, fmt.Errorf("expense \"%s\": %s", e.String(), err)
	}
	if from != nil && date.Before(*from) {
		return false, nil
	}
	if to != nil && date.After(*to) {
		return false, nil
	}
	re, err := regexp.Compile(identifier)
	if err != nil {
		return false, fmt.Errorf("error while parsing \"%s\" as regex for identifier", identifier)
	}
	if identifier != "" && !re.MatchString(e.Identifier) {
		return false, nil
	}
	return true, nil
}

func (e Expense) GetDate() *time.Time {
	date, err := time.Parse(util.DateFormat, e.DateOfAccrual)
	if err != nil {
		return nil
	}
	return &date
}
