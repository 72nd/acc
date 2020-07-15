package query

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/72nd/acc/pkg/schema"
)

type ElementGroup []Element

func NewElements(ele interface{}) ElementGroup {
	var rsl []Element
	v := reflect.ValueOf(ele)
	if v.Kind() != reflect.Slice {
		logrus.Fatalf("\"%+v\" isn't a slice", ele)
	}
	for i := 0; i < v.Len(); i++ {
		rsl = append(rsl, NewElement(v.Index(i)))
	}
	return rsl
}

func (g ElementGroup) MatchTerm(terms SearchTerms, caseSensitive bool) ElementGroup {
	var rsl ElementGroup
	for i := range g {
		if g[i].Match(terms, caseSensitive) {
			rsl = append(rsl, g[i])
		}
	}
	return rsl
}

func (g ElementGroup) DateMatch(ranges DateTerms) ElementGroup {
	var rsl ElementGroup
	for i := range g {
		if g[i].DateMatch(ranges) {
			rsl = append(rsl, g[i])
		}
	}
	return rsl
}

func (g ElementGroup) Select(sel []string, caseSensitive bool) ElementGroup {
	rsl := make(ElementGroup, len(g))
	if !caseSensitive {
		for i := range sel {
			sel[i] = strings.ToLower(sel[i])
		}
	}
	for i := range g {
		rsl[i] = g[i].Select(sel, caseSensitive)
	}
	return rsl
}

type Element []KeyValue

func NewElement(v reflect.Value) Element {
	var rsl Element
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		rsl = append(rsl, KeyValue{
			Key:   t.Field(i).Name,
			Value: fmt.Sprint(v.Field(i)),
			Field: t.Field(i),
		})
	}
	return rsl
}

func (e Element) Match(terms SearchTerms, caseSensitive bool) bool {
	for i := range terms {
		for j := range e {
			if terms[i].matchKey(e[j].Key, caseSensitive) {
				if terms[i].matchValue(e[j].Value, caseSensitive) {
					return true
				}
			}
		}
	}
	return false
}

func (e Element) DateMatch(ranges DateTerms) bool {
	for i := range ranges {
		for j := range e {
			if ranges[i].matchKey(e[j].Key) {
				if !ranges[i].matchRange(e[j].Value) {
					return false
				}
			}
		}
	}
	return true
}

func (e Element) Select(sel []string, caseSensitive bool) Element {
	var rsl Element
	for i := range e {
		key := e[i].Key
		if !caseSensitive {
			key = strings.ToLower(key)
		}
		contains := false
		for j := range sel {
			if key == sel[j] {
				contains = true
			}
		}
		if contains {
			rsl = append(rsl, e[i])
		}
	}
	return rsl
}

func (e Element) MaxKeyLength() int {
	var rsl int
	for i := range e {
		if rsl < len(e[i].Key) {
			rsl = len(e[i].Key)
		}
	}
	return rsl
}

type KeyValue struct {
	Key   string
	Value string
	Field reflect.StructField
}

func (k KeyValue) RenderValue(s schema.Schema) string {
	switch k.Field.Tag.Get("query") {
	case "amount":
		amount, err := strconv.ParseFloat(k.Value, 64)
		if err != nil {
			return k.Value
		}
		return fmt.Sprintf("%s %.2f", s.JournalConfig.Currency, amount)
	case "customer":
		cst, err := s.Parties.CustomerById(k.Value)
		if err != nil {
			return fmt.Sprintf("%s (no such customer exists)", k.Value)
		}
		return fmt.Sprintf("%s, %s", k.Value, cst.Short())
	case "employee":
		emp, err := s.Parties.EmployeeById(k.Value)
		if err != nil {
			return fmt.Sprintf("%s (no such employee exists)", k.Value)
		}
		return fmt.Sprintf("%s, %s", k.Value, emp.Short())
	case "customer,employee":
		cst, err := s.Parties.CustomerById(k.Value)
		if err == nil {
			return fmt.Sprintf("%s, %s", k.Value, cst.Short())
		}
		emp, err := s.Parties.EmployeeById(k.Value)
		if err != nil {
			return fmt.Sprintf("%s (no such party exists)", k.Value)
		}
		return fmt.Sprintf("%s, %s", k.Value, emp.Short())
	case "expense":
		exp, err := s.Expenses.ExpenseById(k.Value)
		if err != nil {
			return fmt.Sprintf("%s (no such expense exists)", k.Value)
		}
		return fmt.Sprintf("%s, %s", k.Value, exp.Short())
	case "invoice":
		inv, err := s.Invoices.InvoiceById(k.Value)
		if err != nil {
			return fmt.Sprintf("%s (no such invoice exists)", k.Value)
		}
		return fmt.Sprintf("%s, %s", k.Value, inv.Short())
	case "expense,invoice":
		exp, err := s.Expenses.ExpenseById(k.Value)
		if err == nil {
			return fmt.Sprintf("%s, %s", k.Value, exp.Short())
		}
		inv, err := s.Invoices.InvoiceById(k.Value)
		if err != nil {
			return fmt.Sprintf("%s (no such expense/invoice exists)", k.Value)
		}
		return fmt.Sprintf("%s, %s", k.Value, inv.Short())
	case "transaction":
		trn, err := s.Statement.TransactionById(k.Value)
		if err != nil {
			return fmt.Sprintf("%s (no such transaction exists)", k.Value)
		}
		return fmt.Sprintf("%s (%s)", k.Value, trn.Date)
	}

	return k.Value
}
