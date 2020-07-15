package distributed

import (
	"gitlab.com/72nd/acc/pkg/schema"
)

// CustomersToSave is a slice of CustomerToSave's.
type CustomersToSave []CustomerToSave

// CustomerToSave serves as the information source to save all customers and their
// projects in the repository structure.
type CustomerToSave struct {
	Customer     schema.Party
	ProjectFiles ProjectFiles
}

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

// ProjectFile is only used in distributed mode to store all invoices and expenses of a given
// project in the same file in the project folder. This data structure is exclusively used
// to store data in distributed mode and is not used internally.
type ProjectFile struct {
	Project  schema.Project  `yaml:"project"`
	Expenses schema.Expenses `yaml:"expenses"`
	Invoices schema.Invoices `yaml:"invoices"`
}

// AbsolutePaths takes the location of the project folder and changes the relative paths
// of the assets to the correct absolute path.
func (p ProjectFile) AbsolutePaths(prjPath string) ProjectFile {
	for i := range p.Expenses {
		p.Expenses[i].Path = absolutPath(prjPath, p.Expenses[i].Path)
	}
	for i := range p.Invoices {
		p.Invoices[i].Path = absolutPath(prjPath, p.Invoices[i].Path)
	}
	return p
}
