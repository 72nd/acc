package util

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

const GermanLayout = "02-01-2006"

type SearchItems []SearchItem

func (s SearchItems) Match(search string) SearchItems {
	var result SearchItems
	for i := range s {
		if fuzzy.MatchFold(search, s[i].SearchValue) {
			result = append(result, s[i])
		}
	}
	return result
}

func (s SearchItems) ByIndex(index int) (*SearchItem, error) {
	if index < 0 || len(s) <= index+1 {
		return nil, errors.New(fmt.Sprintf("no item found for index %d", index))
	}
	return &s[index], nil
}

type SearchItem struct {
	Name        string
	Type        string
	Value       interface{}
	SearchValue string
}

type Strategy int

const (
	AcceptStrategy Strategy = iota
	RedoStrategy
	SkipStrategy
)

func AskString(name, desc, defaultValue string) string {
	input := simplePrompt(name, "string", desc, defaultValue)
	if input == "" {
		return defaultValue
	}
	return input
}

func AskStringFromSearch(name, desc string, searchItems SearchItems) string {
	result, _ := searchPrompt(name, desc, searchItems, false, nil, nil)
	value, ok := result.(string)
	if !ok {
		logrus.Fatalf("could not convert %+v to string", result)
	}
	return value
}

func AskStringFromSearchWithNew(name, desc string, searchItems SearchItems, newFunction func(arg interface{}) interface{}, arg interface{}) (value string, newElement interface{}) {
	result, newElement := searchPrompt(name, desc, searchItems, false, newFunction, arg)
	value, ok := result.(string)
	if !ok {
		logrus.Fatalf("could not convert %+v to string", result)
	}
	return value, newElement
}

func AskStringFromListSearch(name, desc string, searchItems SearchItems) string {
	result, _ := searchPrompt(name, desc, searchItems, true, nil, nil)
	value, ok := result.(string)
	if !ok {
		logrus.Fatalf("could not convert %+v to string", result)
	}
	return value
}

func AskInt(name, desc string, defaultValue int) int {
	input := simplePrompt(name, "int", desc, fmt.Sprintf("%d", defaultValue))
	if input == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(input)
	if err != nil {
		logrus.Warn("could not parse input as a number (int)")
		return AskInt(name, desc, defaultValue)
	}
	return value
}

func AskIntFromList(name, desc string, searchItems SearchItems) int {
	header(name, "selection", desc, fmt.Sprintf("Select a item between 1 and %d", len(searchItems)), false)
	listItems(searchItems)
	fmt.Print("--> ")
	input := getInput()
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(searchItems) {
		logrus.Error("invalid input, try again")
		return AskIntFromList(name, desc, searchItems)
	}
	value, ok := searchItems[index-1].Value.(int)
	if !ok {
		logrus.Fatalf("value %+v in search item with index %d is not an int", searchItems[index-1], index)
	}
	return value
}

func AskIntFromListSearch(name, desc string, searchItems SearchItems) int {
	result, _ := searchPrompt(name, desc, searchItems, true, nil, nil)
	value, ok := result.(int)
	if !ok {
		logrus.Fatalf("could not convert %+v to int", result)
	}
	return value
}

func AskBool(name, desc string, defaultValue bool) bool {
	input := simplePrompt(name, "bool", desc, strconv.FormatBool(defaultValue))
	if input == "" {
		return defaultValue
	}
	value, err := parseBool(input)
	if err != nil {
		logrus.Warn(err)
		return AskBool(name, desc, defaultValue)
	}
	return value
}

func AskForConformation(question string) bool {
	fmt.Printf("%s (Y/n)\n--> ", aurora.BrightCyan(question))
	input := getInput()
	value, err := parseBool(input)
	if err != nil {
		logrus.Warn(err)
		return AskForConformation(question)
	}
	return value
}

func AskFloat(name, desc string, defaultValue float64) float64 {
	input := simplePrompt(name, "float", desc, fmt.Sprintf("%.2f", defaultValue))
	if input == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(input, 32)
	if err != nil {
		logrus.Warn("Could not parse input as a floating number (float)")
		return AskFloat(name, desc, defaultValue)
	}
	return value
}

func AskMoney(name, desc string, defaultValue Money, currency string) Money {
	input := simplePrompt(name, "money", desc, defaultValue.Display())
	if input == "" {
		return defaultValue
	}
	value, err := NewMonyFromDotNotation(input, currency)
	if err != nil {
		logrus.Warn(err)
		return AskMoney(name, desc, defaultValue, currency)
	}
	return value
}

func AskDate(name, desc string, defaultValue time.Time) string {
	possibleLayouts := []string{
		GermanLayout,
		"01.02.2006",
		"2006-01-02",
		"2006.02.01",
	}
	header(name, "DD-MM-YYYY", desc, fmt.Sprintf("Enter for empty, 'T' for today (%s)", defaultValue.Format(GermanLayout)), true)
	input := getInput()
	if input == "" {
		return ""
	}
	if input == "T" {
		return defaultValue.Format(DateFormat)
	}
	success := false
	var value time.Time
	for i := range possibleLayouts {
		var err error
		value, err = time.Parse(possibleLayouts[i], input)
		if err == nil {
			success = true
			break
		}
	}
	if !success {
		logrus.Warnf("Could not parse input as date with format: %s", DateFormat)
		return AskDate(name, desc, defaultValue)
	}
	return value.Format(DateFormat)
}

func AskForStategy() Strategy {
	ok := AskForConformation("Were your entries correct?")
	if !ok {
		for {
			strategy := AskIntFromList(
				"Strategy",
				"how do you want to resolve this situation?",
				SearchItems{
					{
						Name:  "Redo",
						Value: 1,
					},
					{
						Name:  "Skip",
						Value: 2,
					},
				})
			switch strategy {
			case 1:
				return RedoStrategy
			case 2:
				return SkipStrategy
			default:
				logrus.Error("invalid input, try again")
			}
		}
	}
	return AcceptStrategy
}
