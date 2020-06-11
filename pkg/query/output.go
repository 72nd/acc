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

func OutputsFromElements(s schema.Schema, ele []Element) Outputs {
	var rsl []Output
	for i := range ele {
		rsl = append(rsl, NewOutput(s, ele[i]))
	}
	return rsl
}

func (o Outputs) PPKeyValue(s schema.Schema, mode OutputMode, render, openAttachment bool) {
	for i := range o {
		o[i].PPKeyValue(&s, mode, render)
		if openAttachment {
			o[i].OpenAttachment()
		}
	}
}

type Output struct {
	Header  []string
	Element Element
}

func NewOutput(s schema.Schema, ele Element) Output {
	return Output{
		Header:  []string{"Key", "Value"},
		Element: ele,
	}
}

func (o Output) PPKeyValue(a *schema.Schema, mode OutputMode, render bool) {
	switch mode {
	case YamlMode:
		fmt.Print(o.yamlKeyValue(a, render))
	case TableMode:
		fmt.Print(o.tableKeyValue(a, render))
	default:
		logrus.Fatalf("illegal output mode \"%d\"", mode)
	}
}

func (o Output) OpenAttachment() {
	for i := range o.Element {
		if o.Element[i].Field.Tag.Get("query") == "path" {
			ext := util.NewExternal(o.Element[i].Value, false)
			ext.Open()
		}
	}
}

func (o Output) yamlKeyValue(a *schema.Schema, render bool) string {
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

func (o Output) tableKeyValue(a *schema.Schema, render bool) string {
	termWidth, err := util.TerminalWidth()
	if err != nil && render {
		logrus.Warnf("%s using 80 as default instead", err)
		termWidth = 80
	}
	valueWidth := termWidth - o.Element.MaxKeyLength() - 7

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
