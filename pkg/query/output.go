package query

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"gitlab.com/72th/acc/pkg/util"
	"gopkg.in/yaml.v2"
)

type OutputMode int

const (
	YamlMode OutputMode = iota
	TableMode
)

type Outputs []Output

func OutputsFromElements(a schema.Acc, ele []Element) Outputs {
	var rsl []Output
	for i := range ele {
		rsl = append(rsl, NewOutput(a, ele[i]))
	}
	return rsl
}

func (o Outputs) PPKeyValue(a schema.Acc, mode OutputMode) {
	for i := range o {
		o[i].PPKeyValue(&a, mode)
	}
}

type Output struct {
	Header  []string
	Element Element
}

func NewOutput(a schema.Acc, ele Element) Output {
	return Output{
		Header:  []string{"Key", "Value"},
		Element: ele,
	}
}

func (o Output) PPKeyValue(a *schema.Acc, mode OutputMode) {
	switch mode {
	case YamlMode:
		fmt.Print(o.yamlKeyValue(a))
	case TableMode:
		fmt.Print(o.tableKeyValue(a))
	default:
		logrus.Fatalf("illegal output mode \"%d\"", mode)
	}
}

func (o Output) yamlKeyValue(a *schema.Acc) string {
	data := make(map[string]string)
	for i := range o.Element {
		if a != nil {
			data[o.Element[i].Key] = o.Element[i].RenderValue(*a)
		} else {
			data[o.Element[i].Key] = o.Element[i].Value
		}
	}
	yml, err := yaml.Marshal(data)
	if err != nil {
		logrus.Fatalf("error while marshaling \"%+v\": %s", data, err)
	}
	return string(yml)
}

func (o Output) tableKeyValue(a *schema.Acc) string {
	termWidth, err := util.TerminalWidth()
	if err != nil {
		logrus.Warnf("%s using 80 as default instead", err)
		termWidth = 80
	}
	fmt.Println(termWidth)
	valueWidth := int(termWidth) - o.Element.MaxKeyLength() - 7

	tblStr := &bytes.Buffer{}
	tbl := tablewriter.NewWriter(tblStr)
	tbl.SetHeader(o.Header)
	tbl.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	tbl.SetColumnColor(tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor})
	for i := range o.Element {
		value := o.Element[i].Value
		if a != nil {
			value = o.Element[i].RenderValue(*a)
		}
		tbl.Append([]string{o.Element[i].Key, multiline(value, valueWidth)})
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
