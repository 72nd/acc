package bimpf

import "gitlab.com/72th/acc/pkg"

// Dump reassembles the structure of a Bimpf json dump file.
type Dump struct {
	TimeUnits []TimeUnit `json:"time_units"`
	Expenses  []Expense  `json:"expenses"`
	Customers []Customer `json:"customers"`
	Employees []Employee `json:"employees"`
}

// Validate the element and return the result.
func (d Dump) Validate() []pkg.ValidateResult {
	var results []pkg.ValidateResult
	for i := range d.TimeUnits {
		results = append(results, pkg.Check(d.TimeUnits[i]))
	}
	for i := range d.Expenses {
		results = append(results, pkg.Check(d.Expenses[i]))
	}
	for i := range d.Customers {
		results = append(results, pkg.Check(d.Customers[i]))
	}
	for i := range d.Employees {
		results = append(results, pkg.Check(d.Employees[i]))
	}
	return results
}
