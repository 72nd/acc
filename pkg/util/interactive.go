package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		logrus.Fatalf("value %+v in search item with index %d is not an int", searchItems[index], index)
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

func AskDate(name, desc string, defaultValue time.Time) string {
	possibleLayouts := []string{
		GermanLayout,
		"01.02.2006",
		"2006-01-02",
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

func simplePrompt(name, typeName, desc, defaultValue string) string {
	header(name, typeName, desc, fmt.Sprintf("Enter for default (%s)", defaultValue), true)
	return getInput()
}

func simplePromptWithEmpty(name, typeName, desc, defaultValue string) (value string, empty bool) {
	header(name, typeName, desc, fmt.Sprintf("Enter for default (%s), 'E' for empty", defaultValue), true)
	input := getInput()
	if input == "E" {
		return "", true
	}
	return input, true
}

func searchPrompt(name, desc string, items SearchItems, showList bool, newFunction func(arg interface{}) interface{}, arg interface{}) (result interface{}, newElement interface{}) {
	fmt.Println(newFunction)
	functions := "'T(text)' for free text form, 'E' for empty"
	if showList {
		functions = fmt.Sprintf("%s, 'L(number)' for selecting by number", functions)
	}
	if newFunction != nil {
		functions = fmt.Sprintf("%s 'N' for new element", functions)
	}
	fmt.Printf("%s %s %s\n",
		aurora.BrightCyan(fmt.Sprintf("Search for a %s", aurora.Bold(name))),
		aurora.Yellow(fmt.Sprintf("(%s)", desc)),
		aurora.Green(functions),
	)
	if showList {
		listItems(items)
	}
	fmt.Print("--> ")
	input := getInput()
	if input == "T" {
		return AskString(name, desc, ""), nil
	} else if input == "E" {
		return "", nil
	} else if input == "L" && showList {
		index := AskInt("index", fmt.Sprintf("1 to %d", len(items)), 0)
		ele, err := items.ByIndex(index)
		if err != nil {
			logrus.Error("invalid input, try again")
			return searchPrompt(name, desc, items, showList, newFunction, arg)
		}
		return ele.Value, nil
	} else if strings.HasPrefix(input, "T") {
		if strings.HasPrefix(input, "T ") {
			return input[2:], nil
		}
		return input[1:], nil
	} else if strings.HasPrefix(input, "L") {
		input2 := input[1:]
		if strings.HasPrefix(input, "L ") {
			input2 = input[2:]
		}
		index, err := strconv.Atoi(input2)
		if err != nil {
			fmt.Println(aurora.BrightCyan("invalid input, try again"))
			return searchPrompt(name, desc, items, showList, newFunction, arg)
		}
		ele, err := items.ByIndex(index - 1)
		if err != nil {
			fmt.Println(aurora.BrightCyan("invalid input, try again"))
			return searchPrompt(name, desc, items, showList, newFunction, arg)
		}
		return ele.Value, nil
	} else if input == "N" {
		fmt.Print(aurora.BrightCyan(fmt.Sprintf("New %s: ", name)))
		return "", newFunction(arg)
	}

	matches := items.Match(input)
	if len(matches) == 0 {
		fmt.Println(aurora.BrightCyan("No entry found, try again"))
		return searchPrompt(name, desc, items, showList, newFunction, arg)
	}
	for {
		listItems(matches)
		fmt.Printf("%s ", aurora.BrightCyan("Choose item, 'S' to search again:"))
		input2 := getInput()
		if input2 == "S" {
			return searchPrompt(name, desc, items, showList, newFunction, arg)
		}

		value, err := strconv.Atoi(input2)
		if err != nil || value < 1 || value > len(items) {
			logrus.Error("invalid input, try again")
			continue
		}
		return matches[value-1].Value, nil
	}
}

func header(name, typeName, desc, options string, prompt bool) {
	var prm string
	if prompt {
		prm = "--> "
	}
	fmt.Printf("%s%s %s %s\n%s",
		aurora.BrightCyan(aurora.Bold(name)),
		aurora.BrightCyan(fmt.Sprintf(" (%s)", typeName)),
		aurora.Yellow(fmt.Sprintf("%s", desc)),
		aurora.Green(options),
		prm)
}

func listItems(items SearchItems) {
	for i := range items {
		var info string
		switch items[i].Value.(type) {
		case int:
			info = fmt.Sprintf("%d", items[i].Value)
		default:
			info = fmt.Sprintf("%s", items[i].Value)
		}
		if items[i].Type != "" {
			info = fmt.Sprintf("%s, %s", items[i].Type, info)
		}
		fmt.Printf("%s %s %s\n",
			aurora.Yellow(fmt.Sprintf("%d)", i+1)),
			items[i].Name,
			aurora.Green(fmt.Sprintf("(%s)", info)))
	}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1]
}

func parseBool(input string) (bool, error) {
	if input == "y" || input == "Y" || input == "1" || input == "true" || input == "t" {
		return true, nil
	} else if input == "n" || input == "N" || input == "0" || input == "false" || input == "f" {
		return false, nil
	}
	return false, errors.New("could not parse input as a boolean value (bool). Please use y/Y/t/1/true or n/N/f/0/false")
}
