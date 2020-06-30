package project

import "gitlab.com/72th/acc/pkg/schema"

// ProjectFiles is a slice of ProjectFile's.
type ProjectFiles []ProjectFile

// Projects returns all projects of all ProjectFile's.
func (p ProjectFiles) Projects() schema.Projects {
	var rsl schema.Projects
	for i := range p {
		rsl = append(rsl, p[i].Project)
	}
	return rsl
}

// Expenses returns all expenses of all ProjectFile's.
func (p ProjectFiles) Expenses() schema.Expenses {
	var rsl schema.Expenses
	for i := range p {
		rsl = append(rsl, p[i].Expenses...)
	}
	return rsl
}

// Invoices returns all Invoices of all ProjectFile's.
func (p ProjectFiles) Invoices() schema.Invoices {
	var rsl schema.Invoices
	for i := range p {
		rsl = append(rsl, p[i].Invoices...)
	}
	return rsl
}

// ProjectFile is only used in project mode to store all invoices and expenses of a given
// project in the same file in the project folder. This data structure is exclusively used
// to store data in project mode and is not used internally.
type ProjectFile struct {
	Project  schema.Project  `yaml:"project"`
	Expenses schema.Expenses `yaml:"expenses"`
	Invoices schema.Invoices `yaml:"invoices"`
}

