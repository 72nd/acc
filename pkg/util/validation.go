package util

import (
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/sirupsen/logrus"
)

// FlawLevel states the importance of a validation error.
// The levels also states in which step the validation should be fixed.
type FlawLevel int

const (
	// UndefinedFlaw was not defined by the programmer.
	UndefinedFlaw FlawLevel = iota
	// FundamentalFlaw has to be fixed right now.
	FundamentalFlaw
	// BeforeImportFlaw should be fixed before import into Acc.
	BeforeImportFlaw
	// BeforeMergeFlaw should be fixed before using any merging functions.
	BeforeMergeFlaw
	// BeforeExportFlaw should be fixed before using export functions.
	BeforeExportFlaw
)

// String returns a human readable representation for the flaw level.
func (l FlawLevel) String() string {
	switch l {
	case FundamentalFlaw:
		return "fundamental"
	case BeforeImportFlaw:
		return "before import"
	case BeforeMergeFlaw:
		return "before merge"
	case BeforeExportFlaw:
		return "before export"
	default:
		return "undefined"
	}
}

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
	var conditions Conditions
	for i := range v.Conditions() {
		if v.Conditions()[i].Condition {
			conditions = append(conditions, v.Conditions()[i])
		}
	}
	return ValidateResult{
		Element:    v,
		Conditions: conditions,
	}
}

// ValidateResult contains the result of a validation check.
type ValidateResult struct {
	Element    Validatable
	Conditions Conditions
}

// Valid returns true if there were no validation problems.
func (v ValidateResult) Valid() bool {
	return len(v.Conditions) == 0
}

// Log logs all validation
func (v ValidateResult) Log() {
	for i := range v.Conditions {
		logrus.WithFields(logrus.Fields{
			"type":   v.Element.Type(),
			"name":   v.Element.String(),
			"reason": v.Conditions[i].Message,
			"level":  v.Conditions[i].Level.String(),
		})
	}
}

func (v ValidateResult) String() string {
	if len(v.Conditions) == 0 {
		return ""
	}
	result := fmt.Sprintf("%s %s: %+v", strings.ToUpper(v.Element.Type()), v.Element.String(), v.Conditions[0])
	for i := 1; i < len(v.Conditions); i++ {
		result = fmt.Sprintf("%s\n%s %s: %s (%s)", result, strings.ToUpper(v.Element.Type()), v.Element.String(), v.Conditions[0].Message, v.Conditions[0].Level)
	}
	return result
}

func (v ValidateResult) TableRows() [][]string {
	if len(v.Conditions) == 0 {
		return [][]string{}
	}
	var results [][]string
	for i := range v.Conditions {
		results = append(results, []string{v.Element.Type(), v.Element.String(), v.Conditions[i].Message, v.Conditions[i].Level.String()})
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

func (v ValidateResults) TableRows() [][]string {
	var results [][]string
	for i := range v {
		results = append(results, v[i].TableRows()...)
	}
	return results
}

type Conditions []struct {
	Condition bool
	Message   string
	Level     FlawLevel
}

// FileExist returns whether a file at a given path exists or not.
func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ValidDate(format, date string) bool {
	_, err := time.Parse(format, date); if err != nil {
		return false
	}
	return true
}
