package util

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
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
	Value       interface{}
	SearchValue string
}

func AskString(name, desc, defaultValue string) string {
	input := simplePrompt(name, "string", desc, defaultValue)
	if input == "" {
		return defaultValue
	}
	return input
}

func AskStringFromSearch(name, desc string, searchItems SearchItems) string {
	result := searchPrompt(name, desc, searchItems, false)
	value, ok := result.(string)
	if !ok {
		logrus.Fatalf("could not convert %+v to string", result)
	}
	return value
}

func AskStringFromListSearch(name, desc string, searchItems SearchItems) string {
	result := searchPrompt(name, desc, searchItems, true)
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
	header(name, "selection", desc, fmt.Sprintf("Select a item between 1 and %d", len(searchItems)))
	listItems(searchItems)
	fmt.Print("--> ")
	input := getInput()
	value, err := strconv.Atoi(input)
	if err != nil || value < 1 || value > len(searchItems) {
		logrus.Warn("invali")
	}
}

func AskIntFromListSearch(name, desc string, searchItems SearchItems) int {
	result := searchPrompt(name, desc, searchItems, true)
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
	input := simplePrompt(name, "int", desc, fmt.Sprintf("%.2f", defaultValue))
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
	header(name, "DD-MM-YYYY", desc, fmt.Sprintf("Enter for %s, 'T' for today (%s)", aurora.Underline("empty"), defaultValue.Format(GermanLayout)))
	input := getInput()
	if input == "" {
		return ""
	}
	if input == "T" {
		return defaultValue.Format(GermanLayout)
	}
	success := false
	var value time.Time
	for i := range possibleLayouts {
		var err error
		value, err = time.Parse(GermanLayout, possibleLayouts[i])
		if err == nil {
			success = true
			break
		}
	}
	if !success {
		logrus.Warnf("Could not parse input as date with format: %s", GermanLayout)
		return AskDate(name, desc, defaultValue)
	}
	return value.Format(GermanLayout)
}

func simplePrompt(name, typeName, desc, defaultValue string) string {
	header(name, typeName, desc, fmt.Sprintf("Enter for default (%s)", defaultValue))
	return getInput()
}

func simplePromptWithEmpty(name, typeName, desc, defaultValue string) (value string, empty bool) {
	header(name, typeName, desc, fmt.Sprintf("Enter for default (%s), 'E' for empty", defaultValue))
	input := getInput()
	if input == "E" {
		return "", true
	}
	return input, true
}

func searchPrompt(name, desc string, items SearchItems, showList bool) (result interface{}) {
	functions := "'T(text)' for free text form, 'E' for empty"
	if showList {
		functions = fmt.Sprintf("%s, 'L(number)' for selecting by number", functions)
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
		return AskString(name, desc, "")
	} else if input == "E" {
		return ""
	} else if input == "L" && showList {
		index := AskInt("index", fmt.Sprintf("1 to %d", len(items)), 0)
		ele, err := items.ByIndex(index)
		if err != nil {
			logrus.Error("invalid input, try again")
			return searchPrompt(name, desc, items, showList)
		}
		return ele.Value
	} else if strings.HasPrefix(input, "T") {
		if strings.HasPrefix(input, "T ") {
			return input[2:]
		}
		return input[1:]
	} else if strings.HasPrefix(input, "L") {
		input2 := input[1:]
		if strings.HasPrefix(input, "L ") {
			input2 = input[2:]
		}
		index, err := strconv.Atoi(input2)
		if err != nil {
			fmt.Println(aurora.BrightCyan("invalid input, try again"))
			return searchPrompt(name, desc, items, showList)
		}
		ele, err := items.ByIndex(index - 1)
		if err != nil {
			fmt.Println(aurora.BrightCyan("invalid input, try again"))
			return searchPrompt(name, desc, items, showList)
		}
		return ele.Value
	}

	matches := items.Match(input)
	if len(matches) == 0 {
		fmt.Println(aurora.BrightCyan("No entry found, try again"))
		return searchPrompt(name, desc, items, showList)
	}
	for {
		listItems(matches)
		fmt.Printf("%s ", aurora.BrightCyan("Choose item, 'S' to search again:"))
		input2 := getInput()
		input2 = strings.Replace(input2, "\n", "", -1)
		if input2 == "S" {
			return searchPrompt(name, desc, items, showList)
		}

		value, err := strconv.Atoi(input)
		if err != nil || value < 1 || value > len(items) {
			logrus.Error("invalid input, try again")
			continue
		}
		return items[value-1].Value
	}
}

func header(name, typeName, desc, options string) {
	fmt.Printf("%s%s %s %s\n--> ",
		aurora.BrightCyan(aurora.Bold(name)),
		aurora.BrightCyan(fmt.Sprintf(" (%s)", typeName)),
		aurora.Yellow(fmt.Sprintf("%s", desc)),
		aurora.Green(options))
}

func listItems(items SearchItems) {
	for i := range items {
		fmt.Printf("%s %s %s\n",
			aurora.Yellow(fmt.Sprintf("%d)", i+1)),
			items[i].Name,
			aurora.Green(fmt.Sprintf("(%s)", items[i].Value)))
	}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1]
}

func parseBool(input string) (bool, error) {
	if input == "y" || input == "1" || input == "true" || input == "t" {
		return true, nil
	} else if input == "n" || input == "0" || input == "false" || input == "f" {
		return false, nil
	}
	return false, errors.New("could not parse input as a boolean value (bool). Please use y/t/1/true or n/f/0/false")
}
