package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type Checkable interface {
	Validate() ValidateResults
}

// Validatable types can be validated.
type Validatable interface {
	Type() string
	String() string
	Conditions() Conditions
}

// Check the Validatable and return the results.
func Check(v Validatable) ValidateResult {
	var messages []string
	for i := range v.Conditions() {
		if v.Conditions()[i].Condition {
			messages = append(messages, v.Conditions()[i].Message)
		}
	}
	return ValidateResult{
		Element:  v,
		Messages: messages,
	}
}

// ValidateResult contains the result of a validation check.
type ValidateResult struct {
	Element  Validatable
	Messages []string
}

// Valid returns true if there were no validation problems.
func (v ValidateResult) Valid() bool {
	return len(v.Messages) == 0
}

// Log logs all validation
func (v ValidateResult) Log() {
	for i := range v.Messages {
		logrus.WithFields(logrus.Fields{
			"type":   v.Element.Type(),
			"name":   v.Element.String(),
			"reason": v.Messages[i],
		})
	}
}

func (v ValidateResult) String() string {
	if len(v.Messages) == 0 {
		return ""
	}
	result := fmt.Sprintf("%s %s: %s", strings.ToUpper(v.Element.Type()), v.Element.String(), v.Messages[0])
	for i := 1; i < len(v.Messages); i++ {
		result = fmt.Sprintf("%s\n%s %s: %s", result, strings.ToUpper(v.Element.Type()), v.Element.String(), v.Messages[0])
	}
	return result
}

func (v ValidateResult) TableRows() []TableRow {
	if len(v.Messages) == 0 {
		return []TableRow{}
	}
	var results []TableRow
	for i := range v.Messages{
		results = append(results, []string{v.Element.Type(), v.Element.String(), v.Messages[i]})
	}
	return results
}

type ValidateResults []ValidateResult

func (v ValidateResults) String() string {
	var result string
	for i := range v {
		tmp := v[i].String()
		if tmp != "" {
			result += "\n" + tmp
		}
	}
	return result
}

func (v ValidateResults) TableRows() []TableRow {
	var results []TableRow
	for i := range v {
		results = append(results, v[i].TableRows()...)
	}
	return results
}

type Conditions []struct {
	Condition bool
	Message   string
}
