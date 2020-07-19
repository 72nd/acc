// Provides functionality to validate and import Bimpf json dumps.
package bimpf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/72nd/acc/pkg/config"
	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
	"github.com/sirupsen/logrus"
)

// Dump reassembles the structure of a Bimpf json dump file.
type Dump struct {
	TimeUnits []TimeUnit `json:"time_units"`
	Customers Customers  `json:"customers"`
	Employees Employees  `json:"employees"`
	Expenses  []Expense  `json:"expenses"`
}

func OpenDump(path string) *Dump {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal(err)
	}
	dump := &Dump{}
	if err := json.Unmarshal(data, &dump); err != nil {
		logrus.Fatal(err)
	}
	return dump
}

// Validate the element and return the result.
func (d Dump) Validate() util.ValidateResults {
	var results []util.ValidateResult
	for i := range d.TimeUnits {
		results = append(results, util.Check(d.TimeUnits[i]))
	}
	for i := range d.Expenses {
		results = append(results, util.Check(d.Expenses[i]))
	}
	for i := range d.Customers {
		results = append(results, util.Check(d.Customers[i]))
	}
	for i := range d.Employees {
		results = append(results, util.Check(d.Employees[i]))
	}
	return results
}

// ValidateAndReportProject validates the bimpf dump and saves the report to the given path.
func (d Dump) ValidateAndReport(path string) {
	rpt := util.Report{
		Title:           "Bimpf Dump Json Validation Report",
		ColumnTitles:    []string{"type", "element", "reason"},
		ValidateResults: d.Validate(),
	}
	rpt.Write(path)
}

// Convert returns the bimpf dump as an Acc struct. Needs a project path and a Nextcloud Bimpf folder path.
func (d Dump) Convert(outputFolder, bimpfFolder string) schema.Schema {
	s := config.NewSchema(outputFolder, "", false, false, false)
	s.Parties.Customers = make([]schema.Party, len(d.Customers))
	for i := range d.Customers {
		s.Parties.Customers[i] = d.Customers[i].Convert()
	}
	s.Parties.Employees = make([]schema.Party, len(d.Employees))
	for i := range d.Employees {
		s.Parties.Employees[i] = d.Employees[i].Convert()
	}
	s.Expenses = d.Customers.ConvertExpenses(bimpfFolder, s.Parties, d.Employees)
	s.Invoices = d.Customers.ConvertInvoices(bimpfFolder, s.Parties)

	return s
}
