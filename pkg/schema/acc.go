// Schema contains the description of the data structure of acc.
package schema

import (
	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
)

// Acc represents an entry point into the data and also provides general information.
type Acc struct {
	// Company contains the information about the organisation which uses acc.
	Company               Party  `json:"company" default:""`
	ExpensesFilePath      string `json:"expensesFilePath" default:"expenses.json"`
	InvoicesFilePath      string `json:"invoicesFilePath" default:"invoices.json"`
	PartiesFilePath       string `json:"partiesFilePath" default:"parties.json"`
	BankStatementFilePath string `json:"bankStatementFilePath" default:"bank.json"`
}

// NewAcc returns a new Acc element with the default values.
func NewAcc() *Acc {
	acc := &Acc{}
	if err := defaults.Set(acc); err != nil {
		logrus.Fatal(err)
	}
	acc.Company = Party{
		Name:       "Fantasia Company",
		Street:     "Main Street",
		StreetNr:   10,
		Place:      "Zurich",
		PostalCode: 8000,
	}
	return acc
}

// NewProject creates a new acc project in the given folder path.
func NewProject(folderPath string) {

}
