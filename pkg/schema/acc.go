// Schema contains the description of the data structure of acc.
package schema

// Acc represents an entry point into the data and also provides general information.
type Acc struct {
	// Company contains the information about the organisation which uses acc.
	Company               Party  `json:"company" default:"[]"`
	ExpensesFilePath      string `json:"expensesFilePath" default:"expenses.json"`
	InvoicesFilePath      string `json:"invoicesFilePath" default:"invoices.json"`
	PartiesFilePath       string `json:"partiesFilePath" default:"parties.json"`
	BankStatementFilePath string `json:"bankStatementFilePath" default:"bank.json"`
}
