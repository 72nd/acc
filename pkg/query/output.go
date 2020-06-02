package query

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type OutputMode int

const (
	YamlMode OutputMode = iota
	TableMode
)

type Output struct {
	Header []string
	Data   [][]string
}

func (o Output) PPKeyValue(mode OutputMode) {
	if !o.validateForKeyValue() {
		logrus.Fatalf("output with header \"%+v\" couldn't be printed as key value as it contains more than two columns", o.Header)
	}
	switch mode {
	case YamlMode:
		fmt.Print(o.yamlKeyValue())
	case TableMode:
		fmt.Print(o.tableKeyValue())
	default:
		logrus.Fatalf("illegal output mode \"%d\"", mode)
	}
}

func (o Output) yamlKeyValue() string {
	data := make(map[string]string)
	for i := range o.Data {
		data[o.Data[i][0]] = o.Data[i][1]
	}
	yml, err := yaml.Marshal(data)
	if err != nil {
		logrus.Fatalf("error while marshaling \"%+v\": %s", data, err)
	}
	return string(yml)
}

func (o Output) tableKeyValue() string {
	tblStr := &bytes.Buffer{}
	tbl := tablewriter.NewWriter(tblStr)
	tbl.SetHeader(o.Header)
	tbl.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	tbl.SetColumnColor(tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor})
	for i := range o.Data {
		tbl.Append([]string{o.Data[i][0], multiline(o.Data[i][1], 60)})
	}
	tbl.Render()
	return tblStr.String()

}

func (o Output) validateForKeyValue() bool {
	if len(o.Data) != 0 && len(o.Data[0]) != len(o.Header) {
		return false
	}
	return true
}

func tableFormat(ele interface{}) string {
	tblStr := &bytes.Buffer{}
	tbl := tablewriter.NewWriter(tblStr)
	tbl.SetHeader([]string{"Key", "Value"})
	tbl.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	tbl.SetColumnColor(tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor})

	v := reflect.ValueOf(ele)
	t := reflect.TypeOf(ele)
	for i := 0; i < t.NumField(); i++ {
		tbl.Append([]string{t.Field(i).Name, multiline(fmt.Sprint(v.Field(i)), 60)})
	}
	tbl.Render()
	return tblStr.String()
}

func multiline(text string, width int) string {
	for i := width; i < len(text); i += width {
		text = text[:i] + "\n" + text[i:]
	}
	return text
}
