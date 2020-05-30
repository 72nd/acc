package util

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type TransactionType int

const (
	CreditTransaction TransactionType = iota // Incoming transaction
	DebitTransaction                         // Outgoing transaction
)

const DateFormat = "2006-01-02"

func Contains(list []string, key string) bool {
	for i := range list {
		if list[i] == key {
			return true
		}
	}
	return false
}

type PrettyMode int

const (
	YamlMode PrettyMode = iota
	TableMode
)

func PrettyFormat(mode PrettyMode, ele interface{}) string {
	switch mode {
	case YamlMode:
		return yamlFormat(ele)
	case TableMode:
		return tableFormat(ele)
	default:
		return "UNDEFINED"
	}
}

func yamlFormat(ele interface{}) string {
	yml, err := yaml.Marshal(ele)
	if err != nil {
		logrus.Fatalf("error while marshaling \"%s\": %s", ele, err)
	}
	return string(yml)
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
