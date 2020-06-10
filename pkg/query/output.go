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

func (o Outputs) PPKeyValue(a schema.Acc, mode OutputMode, render bool) {
	for i := range o {
		o[i].PPKeyValue(&a, mode, render)
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

func (o Output) PPKeyValue(a *schema.Acc, mode OutputMode, render bool) {
	switch mode {
	case YamlMode:
		fmt.Print(o.yamlKeyValue(a, render))
	case TableMode:
		fmt.Print(o.tableKeyValue(a, render))
	default:
		logrus.Fatalf("illegal output mode \"%d\"", mode)
	}
}

func (o Output) yamlKeyValue(a *schema.Acc, render bool) string {
	data := make(map[string]string)
	for i := range o.Element {
		if a != nil && render {
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

func (o Output) tableKeyValue(a *schema.Acc, render bool) string {
	termWidth, err := util.TerminalWidth()
	if err != nil && render {
		logrus.Warnf("%s using 80 as default instead", err)
		termWidth = 80
	}
	valueWidth := int(termWidth) - o.Element.MaxKeyLength() - 7

	tblStr := &bytes.Buffer{}
	tbl := tablewriter.NewWriter(tblStr)
	tbl.SetHeader(o.Header)
	tbl.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	tbl.SetColumnColor(tablewriter.Colors{tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor})
	tbl.SetAutoWrapText(false)
	for i := range o.Element {
		value := o.Element[i].Value
		if a != nil {
			value = o.Element[i].RenderValue(*a)
		}
		ele := []string{o.Element[i].Key, multiline(value, valueWidth)}
		fmt.Println(ele)
		tbl.Append(ele)
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
