package schema

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v3"
)

type Schema struct {
	Company             Company
	Expenses            Expenses
	Invoices            Invoices
	JournalConfig       JournalConfig
	MiscRecords         MiscRecords
	Parties             Parties
	Projects            Projects
	Statement           Statement
	AppendExpenseSuffix func(suffix string, overwrite bool)
	AppendInvoiceSuffix func(suffix string, overwrite bool)
	SaveFunc            func(s Schema)
	FileHashes          map[string]string
}

func (s Schema) Save() {
	s.SaveFunc(s)
}

func (s Schema) ValidateProject() util.ValidateResults {
	var rsl util.ValidateResults
	rsl = append(rsl, util.Check(s.Company))
	rsl = append(rsl, s.Expenses.Validate()...)
	rsl = append(rsl, s.Invoices.Validate()...)
	rsl = append(rsl, s.MiscRecords.Validate()...)
	rsl = append(rsl, s.Statement.Validate()...)
	rsl = append(rsl, s.Parties.Validate()...)
	rsl = append(rsl, s.Projects.Validate()...)
	rsl = append(rsl, s.Statement.Validate()...)
	return rsl
}

// ValidateAndReportProject validates the Schema and saves the report to the given path.
func (s Schema) ValidateAndReportProject(path string) {
	rpt := util.Report{
		Title:           "Schema Validation Report",
		ColumnTitles:    []string{"type", "element", "reason"},
		ValidateResults: s.ValidateProject(),
	}
	rpt.Write(path)
}

func (s *Schema) Filter(types []string, from *time.Time, to *time.Time, suffix string, overwrite bool, identifier string) {
	if util.Contains(types, "expenses") {
		var err error
		s.Expenses, err = s.Expenses.Filter(from, to, identifier)
		if err != nil {
			logrus.Fatal("error while filtering: ", err)
		}
		s.AppendExpenseSuffix(suffix, overwrite)
	}
	if util.Contains(types, "invoices") {
		var err error
		s.Invoices, err = s.Invoices.Filter(from, to)
		if err != nil {
			logrus.Fatal("error while filtering: ", err)
		}
		s.AppendInvoiceSuffix(suffix, overwrite)
	}

}

func (s Schema) FilterYear(year int) Schema {
	if year > 0 {
		from, to := util.DateRangeFromYear(year)
		s.Expenses, _ = s.Expenses.Filter(&from, &to, "")
		s.Invoices, _ = s.Invoices.Filter(&from, &to)
		s.Statement.Transactions, _ = s.Statement.FilterTransactions(&from, &to)
	}
	return s
}

// Identifiable describes types which are uniquely identifiable trough out the utils structure.
type Identifiable interface {
	GetId() string
	GetIdentifier() string
	String() string
}

func SuggestNextIdentifier(idt []Identifiable, prefix string) string {
	r := regexp.MustCompile(`(\d+)$`)
	max := 0
	for i := range idt {
		rsl := r.FindAllString(idt[i].GetIdentifier(), -1)
		if len(rsl) != 1 {
			continue
		}
		val, err := strconv.Atoi(rsl[0])
		if err != nil {
			logrus.Debugf("regex to find last number in identifier returned something else than a int (%+v), take a look: %s", rsl, err)
		}
		if max < val {
			max = val
		}
	}
	return fmt.Sprintf("%s%d", prefix, max+1)
}

// Completable describes types which automatically can resolve some missing information atomically.
// Example is the setting of a unique Id.
type Completable interface {
	// SetId sets a unique id (UUID-4) for the object.
	SetId()
}

func GetUuid() string {
	return uuid.Must(uuid.NewRandom()).String()
}

// OpenYaml does the same as util.OpenYaml but returns the hash of the file.
func OpenYamlHashed(data interface{}, path, dataType string) string {
	if path == "" {
		logrus.Fatalf("error reading %s file: given path is empty", dataType)
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("error reading %s file \"%s\": %s", dataType, path, err)
	}
	if err := yaml.Unmarshal(raw, data); err != nil {
		logrus.Fatalf("error converting (unmarshalling) %s data of file \"%s\" %s", dataType, path, err)
	}
	return hash(raw)
}

func SaveYamlOnChange(data interface{}, path, dataType, oldHash string) {
	var raw []byte
	var err error
	raw, err = yaml.Marshal(data)
	if err != nil {
		logrus.Fatalf("error converting (marshalling) %s data for file \"%s\" to YAML: %s", dataType, path, err)
	}
	if hash(raw) == oldHash {
		return
	}
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		logrus.Fatalf("error writing %s file \"%s\" %s", dataType, path, err)
	}
}

func hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return string(h.Sum(nil))
}
