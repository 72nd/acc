package bimpf

import (
	"fmt"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
	"path"
)

// Projects is a slice of projects.
type Projects []Project

// Project reassembles the structure of a Project in a Bimpf json dump files.
type Project struct {
	Id           int        `json:"id"`
	SbId         string     `json:"sb_id"`
	Name         string     `json:"name"`
	NcFolderName string     `json:"nc_folder_name"`
	CustomerId   int        `json:"customer"`
	CustomerName string     `json:"customer_name"`
	IsArchived   bool       `json:"is_archived"`
	Quotes       []Document `json:"quotes"`
	Invoices     []Document `json:"invoices"`
	Reminders    []Document `json:"reminders"`
	Expenses     []Expense  `json:"expenses"`
}

// Type returns a string with the type name of the element.
func (p Project) Type() string {
	return "SB-Project"
}

// String returns a human readable representation of the element.
func (p Project) String() string {
	return fmt.Sprintf("%d/%s (%s) for customer %s", p.Id, p.SbId, p.Name, p.CustomerName)
}

// ShortDescription returns the SbId and the name as a string
func (p Project) ShortDescription() string {
	return fmt.Sprintf("%s: %s", p.SbId, p.Name)
}

// Conditions returns the validation conditions.
func (p Project) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: p.Id < 1,
			Message:   "id is not set (id < 1)",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: p.SbId == "",
			Message:   "solutionsbüro id is not set",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: p.Name == "",
			Message:   "name not set",
			Level:     util.BeforeImportFlaw,
		},
		{
			Condition: p.NcFolderName == "",
			Message:   "nextcloud folder not defined",
			Level:     util.BeforeImportFlaw,
		},
	}
}

// Validate the element and return the result.
func (p Project) Validate() []util.ValidateResult {
	var results []util.ValidateResult
	for i := range p.Quotes {
		results = append(results, util.Check(p.Quotes[i]))
	}
	for i := range p.Invoices {
		results = append(results, util.Check(p.Invoices[i]))
	}
	for i := range p.Reminders {
		results = append(results, util.Check(p.Reminders[i]))
	}
	return append(results, util.Check(p))
}

// Doc reassembles the structure of a Doc in a Bimpf json dump file.
type Document struct {
	Id           int    `json:"id"`
	SbId         string `json:"sb_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Path         string `json:"path"`
	CustomerId   int    `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	ProjectId    int    `json:"project_id"`
	ProjectName  string `json:"project_name"`
	SendDate     string `json:"send_date"`
	DaysOfGrace  int    `json:"days_of_grace"`
	IsStalled    bool   `json:"is_stalled"`
}

// ConvertAsInvoice converts the document to a Bimpf Invoice.
// Amount is set to -1 as Bimpf does not provide such information.
func (d Document) ConvertAsInvoice(pathPrefix, customerId, projectName string, parties schema.Parties) schema.Invoice {
	inv := schema.Invoice{
		Identifier:              d.SbId,
		Name:                    d.Name,
		Amount:                  -1,
		Path:                    path.Join(pathPrefix, d.Path),
		CustomerId:              customerId,
		SendDate:                d.SendDate,
		DateOfSettlement:        "",
		SettlementTransactionId: "",
		ProjectName:             projectName,
	}
	inv.SetId()
	return inv
}

// Type returns a string with the type name of the element.
func (d Document) Type() string {
	return "SB-Doc"
}

// String returns a human readable representation of the element.
func (d Document) String() string {
	return fmt.Sprintf("%d/%s (%s) for customer %s in project %s", d.Id, d.SbId, d.Name, d.CustomerName, d.ProjectName)
}

// Conditions returns the validation conditions.
func (d Document) Conditions() util.Conditions {
	return util.Conditions{
		{
			Condition: d.Id < 1,
			Message:   "id is not set (id < 1)",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: d.SbId == "",
			Message:   "solutionsbüro id is not set",
			Level:     util.FundamentalFlaw,
		},
		{
			Condition: d.Name == "",
			Message:   "name not set",
			Level:     util.BeforeImportFlaw,
		},
		{
			Condition: d.Path == "",
			Message:   "attachment path not specified",
			Level:     util.BeforeImportFlaw,
		},
		{
			Condition: d.SendDate == "",
			Message:   "send date not specified",
			Level:     util.BeforeImportFlaw,
		},
		{
			Condition: d.SendDate == "none",
			Message:   "send date not specified",
		}}
}

// Validate the element and return the result.
func (d Document) Validate() util.ValidateResults {
	return []util.ValidateResult{util.Check(d)}
}
