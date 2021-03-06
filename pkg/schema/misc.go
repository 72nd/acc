package schema

import (
	"fmt"
	"strings"
	"time"

	"github.com/72nd/acc/pkg/util"
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
)

const DefaultMiscRecordsFile = "misc.yaml"
const DefaultMiscRecordPrefix = "m-"

// MiscRecord is a collection of Miscellaneous Record elements.
type MiscRecords []MiscRecord

// NewMiscRecords returns an empty new MiscRecords collection.
func NewMiscRecords() MiscRecords {
	return MiscRecords{}
}

// OpenMiscRecords opens the MiscRecords saved in the YAML file given by the path.
func OpenMiscRecords(path string) MiscRecords {
	var mrc MiscRecords
	util.OpenYaml(&mrc, path, "misc-records")
	return mrc
}

// Save writes the element as YAML file to the given path.
func (m MiscRecords) Save(path string) {
	util.SaveToYaml(m, path, "misc-records")
}

// MiscRecordByRef returns the MiscRecord with the given id. If no record could be found
// an error will be returned.
func (m MiscRecords) MiscRecordByRef(ref Ref) (*MiscRecord, error) {
	for i := range m {
		if ref.Match(m[i]) {
			return &m[i], nil
		}
	}
	return nil, fmt.Errorf("no misc record for id \"%s\" found", ref.Id)
}

// MiscRecordByIdent returns the MiscRecord with the given identifier.
// If no record could be found an error will be returned.
func (m MiscRecords) MiscRecordByIdent(ident string) (*MiscRecord, error) {
	for i := range m {
		if m[i].Identifier == ident {
			return &m[i], nil
		}
	}
	return nil, fmt.Errorf("no misc record for ident \"%s\" found", ident)
}

// GetIdentifiables returns the a sliche of all identifiers. This is used for the
// identifier suggestion while interactively adding a new MiscRecord.
func (m MiscRecords) GetIdentifiables() []Identifiable {
	rsl := make([]Identifiable, len(m))
	for i := range m {
		rsl[i] = m[i]
	}
	return rsl
}

// SearchItems returns a searchable structure of the MiscRecords. So the user
// can search for MiscRecords in the interactive mode.
func (m MiscRecords) SearchItems() util.SearchItems {
	rsl := make(util.SearchItems, len(m))
	for i := range m {
		rsl[i] = m[i].SearchItem()
	}
	return rsl
}

// Validate all MiscRecords.
func (m MiscRecords) Validate() util.ValidateResults {
	var rsl util.ValidateResults
	for i := range m {
		rsl = append(rsl, util.Check(m[i]))
	}
	return rsl
}

// Repopulate all MiscRecords based on the transactions.
func (m MiscRecords) Repopulate(s Schema) {
	for i := range m {
		m[i].Repopulate(s)
	}
}

func (m MiscRecords) SetReferenceDestinations(trn []Identifiable) {
	for i := range m {
		m[i].Transaction.SetDestination(trn)
	}
}

// MiscRecord represents business records which are not invoices or expenses
// but still important for accounting. Example: A credit note from an insurance.
type MiscRecord struct {
	// Id is the internal unique identifier of the Miscellaneous Record.
	Id          string `yaml:"id" default:"1"`
	// Identifier is a unique user-chosen identifier for a misc-record, should be human readable.
	Identifier  string `yaml:"identifier" default:"m-19-1"`
	// Name of the Miscellaneous Record.
	Name        string `yaml:"name" default:""`
	// Path is the full file path to the associated business record.
	Path        string `yaml:"path" default:"/path/to/record.pdf" query:"path"`
	// Date represents the date the document arrived.
	Date        string `yaml:"date" default:"2019-12-20"`
	// Transaction refers to an optional transaction which was issued upon the arrival of the Miscellaneous Record.
	Transaction Ref    `yaml:"settlementTransactionId" default:"" query:"transaction"`
}

// NewMiscRecord returns a new MiscRecord element with the default values.
func NewMiscRecord() MiscRecord {
	mrc := MiscRecord{}
	if err := defaults.Set(&mrc); err != nil {
		logrus.Fatal("error setting defaults for misc record: ", err)
	}
	mrc.Id = GetUuid()
	return mrc
}

// InteractiveNewMiscRecord returns a new MiscRecord based on the user input.
func InteractiveNewMiscRecord(s Schema, asset string) MiscRecord {
	mrc := NewMiscRecord()
	mrc.Identifier = util.AskString(
		"Identifier",
		"Unique human readable identifier",
		SuggestNextIdentifier(s.MiscRecords.GetIdentifiables(), DefaultMiscRecordPrefix))
	mrc.Name = util.AskString(
		"Name",
		"Name of the record",
		"AHV record")
	if asset == "" {
		mrc.Path = util.AskString(
			"Asset",
			"Path to asset file (use --asset to set with flag)",
			"")
		mrc.Date = util.AskDate(
			"Date",
			"Receiving date of the record",
			time.Now())
	} else {
		mrc.Path = asset
	}
	return mrc
}

// SetId generates a unique id for the element if there isn't already one defined.
func (m *MiscRecord) SetId() {
	if m.Id != "" {
		return
	}
	m.Id = GetUuid()
}

// Repopulate MiscRecord based on the transactions.
func (m *MiscRecord) Repopulate(s Schema) {
	trn, err := s.Statement.TransactionForDocument(m.Id)
	if err != nil {
		logrus.Warnf("there is no transaction for expense \"%s\" associated", m.String())
		return
	}
	m.Transaction = NewRef(trn.Id)
}

// GetId returns the id of the MiscRecord.
func (m MiscRecord) GetId() string {
	return m.Id
}

// GetIdentifier return the identifier of the MiscRecord.
func (m MiscRecord) GetIdentifier() string {
	return m.Identifier
}

// SearchItem returns a searchable representation of the MiscRecord.
func (m MiscRecord) SearchItem() util.SearchItem {
	return util.SearchItem{
		Name:        fmt.Sprintf("%s (%s)", m.Name, m.Identifier),
		Type:        m.Type(),
		Value:       m.Id,
		SearchValue: fmt.Sprintf("%s %s", m.Identifier, m.Name)}
}

// String returns a human readable representation of the element.
func (m MiscRecord) String() string {
	return fmt.Sprintf("misc record %s (%s)", m.Name, m.Identifier)
}

// Short returns a short representation of the element.
func (m MiscRecord) Short() string {
	return fmt.Sprintf("%s (%s)", m.Name, m.Identifier)
}

// Type returns a string with the type name of the element.
func (m MiscRecord) Type() string {
	return "MiscRecord"
}

// FileString returns the file name for exporting the misc record as a document.
func (m MiscRecord) FileString() string {
	rsl := m.Identifier
	rsl = strings.ReplaceAll(rsl, " ", "-")
	rsl = strings.ReplaceAll(rsl, ".", "-")
	return rsl
}

// Conditions returns the validation conditions.
func (m MiscRecord) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: m.Id == "",
			Message:   "unique identifier not set (Id is empty)",
		},
		{
			Condition: m.Identifier == "",
			Message:   "human readable identifier not set (Identifier is empty)",
		},
		{
			Condition: m.Name == "",
			Message:   "name not set (Name is empty)",
		},
		{
			Condition: !util.FileExist(m.Path),
			Message:   fmt.Sprintf("business record document at \"%s\" not found", m.Path),
		},
		{
			Condition: m.Date != "" && util.ValidDate(util.DateFormat, m.Date),
			Message:   fmt.Sprintf("string \"%s\" could not be parsed with format YYYY-MM-DD", m.Date),
		},
	}
}

// GetDate returns the date of the misc-records as a time.Time struct.
func (m MiscRecord) GetDate() *time.Time {
	date, err := time.Parse(util.DateFormat, m.Date)
	if err != nil {
		return nil
	}
	return &date
}
