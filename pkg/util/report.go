package util

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"text/template"
	"unicode/utf8"
)

const reportTpl = `
{{ .Title }}
{{ .TitleLine }}

{{ .Table }}

The flaws levels states the moment, when a validation flaw should be fixed:
- FundamentalFlaw: Has to be fixed right now.
- BeforeImportFlaw: Should be fixed before import into Acc.
- BeforeMergeFlaw: Should be fixed before using any merging functions.
- BeforeExportFlaw: Should be fixed before using export functions.
`

type Report struct {
	Title           string
	ColumnTitles    []string
	ValidateResults ValidateResults
}

// Render returns the report as a string.
func (r Report) Render() string {
	data := map[string]interface{}{
		"Title":     strings.ToUpper(r.Title),
		"TitleLine": strings.Repeat("-", utf8.RuneCountInString(r.Title)),
		"Table":     r.renderTable(),
	}
	tpl := template.Must(template.New("report").Parse(reportTpl))
	rsl := strings.Builder{}
	if err := tpl.Execute(&rsl, data); err != nil {
		logrus.Fatal(err)
	}
	return rsl.String()
}

// renderTable renders the table and returns it as a string.
func (r Report) renderTable() string {
	tblStr := &bytes.Buffer{}
	tbl := tablewriter.NewWriter(tblStr)
	tbl.SetHeader([]string{"Type", "Name", "Reason", "Level"})
	tbl.SetColWidth(100)

	lines := r.ValidateResults.TableRows()
	for i := range lines {
		tbl.Append(lines[i])
	}

	tbl.Render()
	return tblStr.String()
}

// Write renders the Report and writes it to the given path.
func (r Report) Write(pth string) {
	output := r.Render()
	if err := ioutil.WriteFile(pth, []byte(output), 0644); err != nil {
		logrus.Fatal(err)
	}
}
