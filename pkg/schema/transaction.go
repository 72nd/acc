package schema

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/72nd/acc/pkg/util"
	"github.com/creasty/defaults"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

type JournalMode int

const (
	UnknownJournalMode JournalMode = iota
	ManualJournalMode
	AutoJournalMode
)

// Transaction represents a single transaction of a bank statement.
type Transaction struct {
	Id                 string               `yaml:"id" default:""`
	Identifier         string               `yaml:"identifier" default:""`
	Amount             util.Money           `yaml:"amount" default:"-" query:"amount"`
	Description        string               `yaml:"description" default:""`
	TransactionType    util.TransactionType `yaml:"transactionType" default:"0"`
	AssociatedParty    Ref                  `yaml:"associatedPartyId" default:"" query:"customer,employee"`
	AssociatedDocument Ref                  `yaml:"associatedDocumentId" default:"" query:"expense,invoice,misc"`
	Date               string               `yaml:"date" default:""`
	JournalMode        JournalMode          `yaml:"journalMode" default:"0"`
}

func NewTransaction() Transaction {
	trn := Transaction{}
	if err := defaults.Set(&trn); err != nil {
		logrus.Fatal(err)
	}
	trn.Amount = util.NewMoney(1000, "CHF")
	return trn
}

func NewTransactionWithUuid() Transaction {
	trn := NewTransaction()
	trn.Id = GetUuid()
	return trn
}

func InteractiveNewTransaction(s Statement, currency string) Transaction {
	trn := NewTransactionWithUuid()
	trn.Identifier = util.AskString(
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.GetIdentifiables(), DefaultTransactionPrefix),
	)
	trn.Description = util.AskString(
		"Description",
		"Some information about the transaction",
		"",
	)
	trn.TransactionType = util.TransactionType(util.AskIntFromList(
		"Transaction Type",
		"",
		util.SearchItems{
			util.SearchItem{
				Name:        "Incoming transaction",
				Value:       int(util.DebitTransaction),
				SearchValue: "1 Incoming Transaction",
			},
			util.SearchItem{
				Name:        "Outgoing transaction",
				Value:       int(util.CreditTransaction),
				SearchValue: "2 Outgoing Transaction",
			},
		}))
	trn.Date = util.AskDate(
		"Date",
		"Transaction date",
		time.Now())
	trn.Amount = util.AskMoney(
		"Amount",
		"",
		util.NewMoney(2342, currency),
		currency,
	)
	trn.JournalMode = JournalMode(util.AskIntFromList(
		"Journal Mode",
		"choose how journal entry will be generated for this transaction",
		util.SearchItems{
			{
				Name:  "Manual Mode",
				Value: int(ManualJournalMode),
			},
			{
				Name:  "Auto Mode",
				Value: int(AutoJournalMode),
			}}))
	return trn
}

func (t Transaction) AssistedCompletion(s Schema, doAll, autoMode, askSkip, documentsOnly bool) Transaction {
	tmp := t
	if autoMode {
		t.JournalMode = AutoJournalMode
	}
	if !doAll && util.Check(t).Valid() {
		fmt.Printf("%s %s\n", aurora.BrightMagenta(aurora.Bold("Skip transaction:")), aurora.BrightMagenta(t.String()))
		return t
	}
	fmt.Printf("%s %s %s\n", aurora.BrightMagenta(aurora.Bold("Optimize transaction:")), aurora.BrightMagenta(t.String()), aurora.BrightMagenta(t.Description))
	if !doAll && askSkip && t.Id != "" && t.Identifier != "" {
		skip := util.AskBool(
			"Skip",
			"Skip this non valid entry?",
			false)
		if skip {
			return t
		}
	}
	if !documentsOnly {
		if t.Id == "" {
			t.SetId()
		}
		identifier := SuggestNextIdentifier(s.Statement.GetIdentifiables(), DefaultTransactionPrefix)
		if t.Identifier == "" {
			t.Identifier = util.AskString(
				"Identifier",
				"Unique human readable identifier",
				identifier)
		}
		t.Description = util.AskString(
			"Description",
			"Description of the transaction",
			t.Description)
		t.JournalMode = JournalMode(util.AskIntFromList(
			"Journal Mode",
			"choose how journal entry will be generated for this transaction",
			util.SearchItems{
				{
					Name:  "Manual Mode",
					Value: int(ManualJournalMode),
				},
				{
					Name:  "Auto Mode",
					Value: int(AutoJournalMode),
				}}))
		if t.AssociatedParty.Empty() && t.JournalMode == AutoJournalMode {
			parties := append(s.Parties.CustomersSearchItems(), s.Parties.EmployeesSearchItems()...)
			suggestion, err := t.parseAssociatedParty(t.Description, s.Parties)
			if err == nil && util.AskForConformation(fmt.Sprintf("Use \"%s\" as associeted third party?", suggestion.String())) {
				t.AssociatedParty = NewRef(suggestion.GetId())
			} else {
				var pty interface{}
				value, pty := util.AskStringFromSearchWithNew(
					"Associated Party",
					"customer/employee which is originator/recipient of the transaction",
					parties,
					InteractiveNewGenericParty,
					s)
				t.AssociatedParty = NewRef(value)
				if pty != nil {
					value, ok := pty.(Party)
					if !ok {
						logrus.Fatal("returned party has invalid type")
					}
					if value.PartyType == CustomerType {
						s.Parties.Customers = append(s.Parties.Customers, value)
					} else if value.PartyType == EmployeeType {
						s.Parties.Employees = append(s.Parties.Employees, value)
					}
					t.AssociatedParty = NewRef(value.Id)
				}
			}
		}
	}

	document, err := t.parseAssociatedDocument(s.Expenses, s.Invoices)
	if err == nil && util.AskForConformation(fmt.Sprintf("Use \"%s\" as associated document?", document.String())) {
		t.AssociatedDocument = NewRef(document.GetId())
	} else {
		docs := append(s.Expenses.SearchItems(), s.Invoices.SearchItems(s)...)
		t.AssociatedDocument = NewRef(util.AskStringFromSearch(
			"Associated Document",
			"couldn't find associated document, manual search",
			docs))
	}

	strategy := util.AskForStategy()
	switch strategy {
	case util.RedoStrategy:
		t.AssistedCompletion(s, doAll, autoMode, askSkip, documentsOnly)
	case util.SkipStrategy:
		return tmp
	}
	return t
}

func (t Transaction) parseAssociatedDocument(expenses Expenses, invoices Invoices) (Identifiable, error) {
	r := regexp.MustCompile(`([ei]-(|.*-)(\d+))(\s|$)`)
	matches := r.FindAllStringSubmatch(t.Description, -1)
	if len(matches) != 1 || len(matches[0]) != 5 {
		return nil, fmt.Errorf("no document identifier found")
	}
	ident := matches[0][1]
	expense, err1 := expenses.ExpenseByIdent(ident)
	invoice, err2 := invoices.InvoiceByIdent(ident)
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("no expense or invoice for identifier %s found", ident)
	}
	if invoice == nil {
		return expense, nil
	}
	return invoice, nil
}

func (t Transaction) parseAssociatedParty(desc string, parties PartiesCollection) (Identifiable, error) {
	desc = strings.ToLower(desc)
	pty := append(parties.Customers, parties.Employees...)
	for i := range pty {
		if strings.Contains(desc, strings.ToLower(pty[i].Name)) {
			return pty[i], nil
		}
	}
	return nil, fmt.Errorf("no party for description found")
}

// GetId returns the unique id of the element.
func (t Transaction) GetId() string {
	return t.Id
}

func (t Transaction) GetIdentifier() string {
	return t.Identifier
}

// SetId generates a unique id for the element if there isn't already one defined.
func (t *Transaction) SetId() {
	if t.Id != "" {
		return
	}
	t.Id = uuid.Must(uuid.NewRandom()).String()
}

func (t Transaction) DateTime() time.Time {
	result, err := time.Parse(util.DateFormat, t.Date)
	if err != nil {
		logrus.Fatalf("could not parse «%s» as date with YYYY-MM-DD: %s", t.Date, err)
	}
	return result
}

// Type returns a string with the type name of the element.
func (Transaction) Type() string {
	return "Transaction"
}

// String returns a human readable representation of the element.
func (t Transaction) String() string {
	if t.TransactionType == util.CreditTransaction {
		return fmt.Sprintf("%s: received %s at %s", t.Identifier, t.Amount.Display(), t.Date)
	}
	return fmt.Sprintf("%s: paid %s at %s", t.Identifier, t.Amount.Display(), t.Date)
}

// Short returns a short representation of the element.
func (t Transaction) Short() string {
	return fmt.Sprintf("%s", t.Id)
}

// Conditions returns the validation conditions.
func (t Transaction) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: t.Id == "",
			Message:   "unique identifier not set (Id is empty)",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.Identifier == "",
			Message:   "human readable identifier not set (Identifier is empty)",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.Description == "",
			Message:   "no description set",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.TransactionType < 0 || t.TransactionType > 1,
			Message:   "transaction type not valid",
			Level:     util.BeforeExportFlaw,
		},
		{
			Condition: t.AssociatedParty.Empty() && t.JournalMode == AutoJournalMode,
			Message:   "no associated party set although auto journal mode is set",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: t.AssociatedDocument.Empty() && t.JournalMode == AutoJournalMode,
			Message:   "no associated document set although auto journal mode is set",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: t.Date == "",
			Message:   "date not set",
			Level:     util.BeforeMergeFlaw,
		},
		{
			Condition: t.Amount.Amount() <= 0,
			Message:   "amount is not set",
			Level:     util.BeforeMergeFlaw,
		},
	}
}

// Validate the element and return the result.
func (t Transaction) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(t)}
}

func (t *Transaction) Repopulate(s Schema) {
	docDone := false
	for i := range s.Expenses {
		if s.Expenses[i].SettlementTransaction.Match(t) {
			t.AssociatedDocument = NewRef(s.Expenses[i].SettlementTransaction.Id)
			docDone = true
			break
		}
	}
	if !docDone {
		for i := range s.Invoices {
			if s.Invoices[i].SettlementTransaction.Match(t) {
				t.AssociatedDocument = NewRef(s.Invoices[i].SettlementTransaction.Id)
				docDone = true
				break
			}
		}
	}
	if !docDone {
		for i := range s.MiscRecords {
			if s.MiscRecords[i].Transaction.Match(t) {
				t.AssociatedDocument = NewRef(s.MiscRecords[i].Transaction.Id)
				docDone = true
				break
			}
		}
	}
}

func (t Transaction) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:        t.Description,
		Value:       t.Id,
		SearchValue: t.Description,
	}
}

func (t Transaction) GetDate() *time.Time {
	date, err := time.Parse(util.DateFormat, t.Date)
	if err != nil {
		return nil
	}
	return &date
}
