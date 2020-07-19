package bimpf

import (
	"fmt"
	"github.com/72nd/acc/pkg/schema"
	"github.com/72nd/acc/pkg/util"
	"github.com/sirupsen/logrus"
	"path"
	"regexp"
	"strconv"
)

// Customers represents a slice of Customer slices
type Customers []Customer

// ConvertExpenses cycles trough the Customers and returns all expenses found as a acc Expenses.
// As the expenses in Bimpf only contain a relative path to receipts the nextcloud prefix has to be provided.
func (c Customers) ConvertExpenses(folderPrefix string, parties schema.Parties, bimpfEmployees Employees) schema.Expenses {
	var exp schema.Expenses
	for i := range c {
		cst, err := parties.CustomerByIdentifier(c[i].SbId)
		if err != nil {
			logrus.Warn(err)
			continue
		}
		for j := range c[i].Projects {
			for k := range c[i].Projects[j].Expenses {
				folder := path.Join(folderPrefix, c[i].NcFolderName, c[i].Projects[j].NcFolderName)
				prjDesc := c[i].Projects[j].ShortDescription()
				ex := c[i].Projects[j].Expenses[k]
				exp = append(exp, ex.Convert(folder, cst.Id, prjDesc, parties, bimpfEmployees))
			}
		}
	}
	return exp
}

// ConvertInvoices cycles trough the Customers and returns all invoices found as a acc Invoice.
// As the expenses in Bimpf only contain a relative path to receipts the nextcloud prefix has to be provided.
func (c Customers) ConvertInvoices(folderPrefix string, parties schema.Parties) schema.Invoices {
	var inv []schema.Invoice
	for i := range c {
		cst, err := parties.CustomerByIdentifier(c[i].SbId)
		if err != nil {
			logrus.Warn(err)
			continue
		}
		for j := range c[i].Projects {
			for k := range c[i].Projects[j].Invoices {
				folder := path.Join(folderPrefix, c[i].NcFolderName, c[i].Projects[j].NcFolderName)
				prjDesc := c[i].Projects[j].ShortDescription()
				inv = append(inv, c[i].Projects[j].Invoices[k].ConvertAsInvoice(folder, cst.Id, prjDesc, parties))
			}
		}
	}
	return inv
}

// Customer reassembles the structure of a Customer TimeUnit in a Bimpf json dump file.
type Customer struct {
	Id           int       `json:"id"`
	SbId         string    `json:"sb_id"`
	Name         string    `json:"name"`
	Comment      string    `json:"comment"`
	NcFolderName string    `json:"nc_folder_name"`
	Recipient1   string    `json:"recipient_1"`
	Recipient2   string    `json:"recipient_2"`
	Recipient3   string    `json:"recipient_3"`
	Recipient4   string    `json:"recipient_4"`
	Email        string    `json:"email"`
	Projects     []Project `json:"projects"`
}

// Type returns a string with the type name of the element.
func (c Customer) Type() string {
	return "SB-Customer"
}

// String returns a human readable representation of the element.
func (c Customer) String() string {
	return fmt.Sprintf("%d/%s (%s)", c.Id, c.SbId, c.Name)
}

// Conditions returns the validation conditions.
func (c Customer) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: c.Id < 1,
			Message:   "id is not set (id < 1)",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: c.SbId == "",
			Message:   "solutionsbüro id is not set",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: c.Name == "",
			Message:   "name not set",
			Level:     util.BeforeImportFlaw,
		},
		{
			Condition: c.NcFolderName == "",
			Message:   "nextcloud folder not defined",
			Level:     util.BeforeImportFlaw,
		},
	}
}

// Validate the element and return the result.
func (c Customer) Validate() util.ValidateResults {
	var results []util.ValidateResult
	for i := range c.Projects {
		results = append(results, util.Check(c.Projects[i]))
	}
	return append(results, util.Check(c))
}

// Convert returns the customers as a acc Party.
func (c Customer) Convert() schema.Party {
	street, number, place, postal := c.parseAddress()
	pty := schema.Party{
		Identifier: c.SbId,
		Name:       c.Name,
		Street:     street,
		StreetNr:   number,
		Place:      place,
		PostalCode: postal,
	}
	pty.SetId()
	return pty
}

// parseAddress tries to parse the Solutionsbüro address structure.
// The result should to be manually validated.
func (c Customer) parseAddress() (street string, number int, place string, postal int) {
	lines := []string{
		c.Recipient1,
		c.Recipient2,
		c.Recipient3,
		c.Recipient4}
	for i := range lines {
		result := regexp.MustCompile(`^\s*(\d{4}) ([A-z]*)$`).FindStringSubmatch(lines[i])
		if len(result) != 3 {
			continue
		}
		postal, _ = strconv.Atoi(result[1])
		place = result[2]
		lines = append(lines[:i], lines[i+1:]...)
		break
	}
	for i := range lines {
		result := regexp.MustCompile(`^\s*([A-z]*) ([0-9]{1,3})$`).FindStringSubmatch(lines[i])
		if len(result) != 3 {
			continue
		}
		street = result[1]
		number, _ = strconv.Atoi(result[2])
	}
	return street, number, place, postal
}
