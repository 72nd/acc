// Provides functionality to validate and import Bimpf json dumps.
package bimpf

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
	"io/ioutil"
	"time"
)

// Dump reassembles the structure of a Bimpf json dump file.
type Dump struct {
	TimeUnits []TimeUnit `json:"time_units"`
	Customers Customers  `json:"customers"`
	Employees Employees `json:"employees"`
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

// ValidateAndReport validates the bimpf dump and saves the report to the given path.
func (d Dump) ValidateAndReport(path string) {
	table := util.Table{
		Header: []string{"type", "element", "reason"},
		Rows:   d.Validate().TableRows(),
	}
	output := fmt.Sprintf("BIMPF DUMP JSON VALIDATION REPORT\n---------------------------------\nGenerated at: %s\n\n%s", time.Now().String(), table.Render())
	if err := ioutil.WriteFile(path, []byte(output), 0644); err != nil {
		logrus.Fatal(err)
	}
}

// Convert returns the bimpf dump as an Acc struct. Needs a project path and a Nextcloud Bimpf folder path.
func (d Dump) Convert(outputFolder, bimpfFolder string) schema.Acc {
	acc := schema.NewProject(outputFolder, false, false)
	acc.Parties.Customers = make([]schema.Party, len(d.Customers))
	for i := range d.Customers {
		acc.Parties.Customers[i] = d.Customers[i].Convert()
	}
	acc.Parties.Employees = make([]schema.Party, len(d.Employees))
	for i := range d.Employees {
		acc.Parties.Employees[i] = d.Employees[i].Convert()
	}
	acc.Expenses = d.Customers.ConvertExpenses(bimpfFolder, acc.Parties, d.Employees)
	acc.Invoices  = d.Customers.ConvertInvoices(bimpfFolder, acc.Parties)

	return acc
}
