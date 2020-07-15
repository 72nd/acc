package schema

import (
	"fmt"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"gitlab.com/72nd/acc/pkg/util"
)

const DefaultInvoicesFile = "invoices.yaml"
const DefaultInvoicesPrefix = "i-"

// Invoices is a slice of invoices.
type Invoices []Invoice

// NewInvoices returns a new Expense slice with the one Expense in it.
func NewInvoices(useDefaults bool) Invoices {
	if useDefaults {
		return []Invoice{NewInvoice()}
	}
	return []Invoice{}
}

// OpenInvoices opens a Expenses saved in the json file given by the path.
func OpenInvoices(path string) Invoices {
	var inv Invoices
	util.OpenYaml(&inv, path, "invoices")
	return inv
}

// Save writes the element as a json to the given path.
// Indented states whether «prettify» the json output.
func (i *Invoices) Save(path string) {
	util.SaveToYaml(i, path, "invoices")
}

// SetId sets a unique id to all elements in the slice.
func (i Invoices) SetId() {
	for j := range i {
		i[j].SetId()
	}
}

func (i Invoices) InvoiceById(id string) (*Invoice, error) {
	for j := range i {
		if i[j].Id == id {
			return &i[j], nil
		}
	}
	return nil, fmt.Errorf("no invoice for id «%s» found", id)
}

func (i Invoices) GetIdentifiables() []Identifiable {
	ivs := make([]Identifiable, len(i))
	for j := range i {
		ivs[j] = i[j]
	}
	return ivs
}

func (i Invoices) SearchItems(s Schema) util.SearchItems {
	result := make(util.SearchItems, len(i))
	for j := range i {
		result[j] = i[j].SearchItem(s)
	}
	return result
}

func (i Invoices) Filter(from *time.Time, to *time.Time) (Invoices, error) {
	var result Invoices
	for j := range i {
		date, err := time.Parse(util.DateFormat, i[j].SendDate)
		if err != nil {
			return nil, fmt.Errorf("invoice \"%s\": %s", i[j].String(), err)
		}
		if from != nil && to != nil && (date.After(*from) || date.Equal(*from)) && (date.Before(*to) || date.Equal(*to)) {
			result = append(result, i[j])
		}
	}
	return result, nil
}

func (i Invoices) AssistedCompletion(s Schema, doAll, autoSave, openAttachment, retainFocus bool) {
	first := true
	for j := range i {
		if !first {
			fmt.Println()
		} else {
			first = false
		}
		i[j] = i[j].AssistedCompletion(s, doAll, openAttachment, retainFocus)
		if autoSave {
			s.Save()
		}
	}
}

func (i Invoices) Repopulate(s Schema) {
	for j := range i {
		i[j].Repopulate(s)
	}
}

func (i Invoices) Len() int {
	return len(i)
}

func (i Invoices) Swap(j, k int) {
	i[j], i[k] = i[k], i[j]
}

func (i Invoices) Less(j, k int) bool {
	dj, _ := time.Parse(util.DateFormat, i[j].SendDate)
	dk, _ := time.Parse(util.DateFormat, i[k].SendDate)
	return dj.Before(dk)
}

// Invoice represents an invoice sent to a customer for some services.
type Invoice struct {
	// Id is the internal unique identifier of the Expense.
	Id string `yaml:"id" default:"1"`
	// Value is a unique user chosen identifier, has to be the same in all source files (bank statements, bimpf dumps...).
	Identifier string `yaml:"identifier" default:"i-19-1"`
	// Name describes meaningful the kind of the Expense.
	Name string `yaml:"name" default:"Expense Name"`
	// Amount states the amount of the Expense.
	Amount float64 `yaml:"amount" default:"10.00" query:"amount"`
	// Path is the full path to the voucher utils.
	Path string `yaml:"path" default:"/path/to/file.utils" query:"path"`
	// Revoked invoices are disabled an no longer taken into account.
	Revoked bool `yaml:"revoked" default:"false"`
	// CustomerId refers to the customer the invoice was sent to.
	CustomerId string `yaml:"customerId" default:"" query:"customer"`
	// SendDate states the date, the invoice was sent to the customer.
	SendDate string `yaml:"sendDate" default:"2019-12-20"`
	// DateOfSettlement states the date the customer paid the outstanding amount.
	DateOfSettlement string `yaml:"dateOfSettlement" default:"2019-12-25"`
	// SettlementTransactionId refers to a possible bank transaction which settled the Expense for the company.
	SettlementTransactionId string `yaml:"settlementTransactionId" default:"" query:"transaction"`
	// ProjectName refers to the associated project of the expense. Depreciated.
	ProjectName string `yaml:"projectName" default:""`
	// ProjectId refers to the associated project.
	ProjectId string `yaml:"projectId" default:""`
}

// NewInvoice returns a new Acc element with the default values.
func NewInvoice() Invoice {
	inv := Invoice{}
	if err := defaults.Set(&inv); err != nil {
		logrus.Fatal("error setting defaults: ", err)
	}
	return inv
}

func NewInvoiceWithUuid() Invoice {
	inv := NewInvoice()
	inv.Id = GetUuid()
	return inv
}

func InteractiveNewInvoice(s Schema, asset string) Invoice {
	inv := NewInvoiceWithUuid()
	inv.Identifier = util.AskString(
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.Invoices.GetIdentifiables(), DefaultInvoicesPrefix))
	inv.Name = util.AskString(
		"Name",
		"Name of the invoice",
		"Invoice for clingfilm")
	inv.Amount = util.AskFloat(
		"Amount",
		"How much is the outstanding balance",
		23.42)
	if asset == "" {
		inv.Path = util.AskString(
			"Asset",
			"Path to asset file (use --asset to set with flag)",
			"")
	} else {
		inv.Path = asset
	}
	inv.CustomerId = util.AskStringFromSearch(
		"Obliged Customer",
		"Customer which has to pay the invoice",
		s.Parties.CustomersSearchItems())
	inv.SendDate = util.AskDate(
		"Send Date",
		"Date the invoice was sent",
		time.Now())
	inv.DateOfSettlement = util.AskDate(
		"Date of settlement",
		"Date when invoice was paid",
		time.Now())
	inv.SettlementTransactionId = util.AskStringFromSearch(
		"Settlement Transaction",
		"Transaction which settled the invoice",
		s.Statement.TransactionSearchItems())
	inv.ProjectName = util.AskString(
		"Project Name",
		"Name of the associated project",
		"")
	inv.ProjectId = util.AskStringFromSearch(
		"Project",
		"Associated Project",
		s.Projects.SearchItems())

	return inv
}

func (i Invoice) AssistedCompletion(s Schema, doAll, openAttachment, retainFocus bool) Invoice {
	if !doAll && util.Check(i).Valid() {
		fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Skip invoice:")), aurora.BrightMagenta(i.String()))
		return i
	}
	tmp := i
	var ext util.External
	if i.Path != "" && openAttachment {
		ext = util.NewExternal(i.Path, retainFocus)
		ext.Open()
	}
	fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Optimize invoice:")), aurora.BrightMagenta(i.String()))
	if i.Amount <= 0.0 {
		i.Amount = util.AskFloat(
			"Amount",
			"How much is the outstanding balance",
			0.0)
	}

	strategy := util.AskForStategy()
	switch strategy {
	case util.RedoStrategy:
		inv := i.AssistedCompletion(s, doAll, openAttachment, retainFocus)
		ext.Close()
		return inv
	case util.SkipStrategy:
		return tmp
	}
	ext.Close()
	return i
}

func (i *Invoice) Repopulate(s Schema) {
	trn, err := s.Statement.TransactionForDocument(i.Id)
	if err != nil {
		logrus.Warnf("there is no transaction for invoice \"%s\" associated", i.String())
		return
	}
	i.DateOfSettlement = trn.Date
	i.SettlementTransactionId = trn.Id
	fmt.Println(i)
}

func (i Invoice) SearchItem(s Schema) util.SearchItem {
	party := ""
	if i.CustomerId != "" {
		pty, err := s.Parties.CustomerById(i.CustomerId)
		if err != nil {
			logrus.Warn("while creating search items: ", err)
		} else {
			party = pty.Name
		}
	}
	return util.SearchItem{
		Name:        fmt.Sprintf("%s for customer %s, amount: %.2f", i.Name, party, i.Amount),
		Type:        i.Type(),
		Value:       i.Id,
		SearchValue: fmt.Sprintf("%s %s %s %s", i.Name, i.Identifier, i.ProjectName, party),
	}
}

func (i Invoices) Validate() util.ValidateResults {
	var rsl util.ValidateResults
	for j := range i {
		rsl = append(rsl, util.Check(i[j]))
	}
	return rsl
}

// GetId returns the unique id of the element.
func (i Invoice) GetId() string {
	return i.Id
}

func (i Invoice) GetIdentifier() string {
	return i.Identifier
}

// SetId generates a unique id for the element if there isn't already one defined.
func (i *Invoice) SetId() {
	if i.Id != "" {
		return
	}
	i.Id = uuid.Must(uuid.NewRandom()).String()
}

func (i Invoices) InvoiceByIdent(ident string) (*Invoice, error) {
	for j := range i {
		if i[j].Identifier == ident {
			return &i[j], nil
		}
	}
	return nil, fmt.Errorf("no invoice for identifier «%s» found", ident)
}

// Type returns a string with the type name of the element.
func (i Invoice) Type() string {
	return "Invoice"
}

// String returns a human readable representation of the element.
func (i Invoice) String() string {
	return fmt.Sprintf("%s (%s): %.2f", i.Name, i.Identifier, i.Amount)
}

// Short returns a short represenation of the element.
func (i Invoice) Short() string {
	return fmt.Sprintf("%s (%s)", i.Name, i.Identifier)
}

func (i Invoice) FileString() string {
	result := fmt.Sprintf("%s", i.Identifier)
	result = strings.ReplaceAll(result, " ", "-")
	result = strings.ReplaceAll(result, ".", "-")
	return result
}

// Conditions returns the validation conditions.
func (i Invoice) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: i.Id == "",
			Message:   "unique identifier not set (Id is empty)",
		},
		{
			Condition: i.Identifier == "",
			Message:   "human readable identifier not set (Identifier is empty)",
		},
		{
			Condition: i.Amount <= 0.0,
			Message:   "amount is not set (Amount is 0.0)",
		},
		{
			Condition: !util.FileExist(i.Path),
			Message:   fmt.Sprintf("business record document at «%s» not found", i.Path),
		},
		{
			Condition: i.CustomerId == "",
			Message:   "customer id is not set (CustomerId empty)",
		},
		{
			Condition: !util.ValidDate(util.DateFormat, i.SendDate),
			Message:   fmt.Sprintf("string «%s» could not be parsed with format YYYY-MM-DD", i.SendDate),
		},
		{
			Condition: i.DateOfSettlement != "" && !util.ValidDate(util.DateFormat, i.DateOfSettlement),
			Message:   fmt.Sprintf("string «%s» could not be parsed with format YYYY-MM-DD", i.DateOfSettlement),
		},
		{
			Condition: i.DateOfSettlement != "" && i.SettlementTransactionId == "",
			Message:   "although date of settlement is set, the corresponding transaction is empty (SettlementTransactionId is empty",
		},
		{
			Condition: i.ProjectName == "",
			Message:   "project name is not set (ProjectName is empty)",
		},
	}
}

func (i Invoice) SendDateTime() time.Time {
	result, err := time.Parse(util.DateFormat, i.SendDate)
	if err != nil {
		logrus.Fatalf("could not parse «%s» as date with YYYY-MM-DD: %s", i.SendDate, err)
	}
	return result
}

func (i Invoice) GetDate() *time.Time {
	date, err := time.Parse(util.DateFormat, i.SendDate)
	if err != nil {
		return nil
	}
	return &date
}
