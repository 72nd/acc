package util

import (
	"bufio"
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const IsoLayout = "2006-01-02"

type SearchItems []SearchItem

func (s SearchItems) Match(search string) SearchItems {
	var result SearchItems
	for i := range s {
		if fuzzy.MatchFold(search, s[i].Value) {
			result = append(result, s[i])
		}
	}
	return result
}

type SearchItem struct {
	Name       string
	Identifier string
	Value      string
}

func AskString(reader *bufio.Reader, name, desc, defaultValue string) string {
	input := simplePrompt(reader, name, "string", desc, defaultValue)
	if input == "" {
		return defaultValue
	}
	return input
}

func AskStringFromSearch(reader *bufio.Reader, name, desc string, searchItems SearchItems) string {
	input := searchPrompt(reader, name, desc)
	rsl := searchItems.Match(input)
	identifier, redo := searchItemsPrompt(reader, rsl, true)
	if redo {
		return AskStringFromSearch(reader, name, desc, searchItems)
	}
	return identifier
}

func AskStringFromList(reader *bufio.Reader, name, desc string, showList bool, values map[string]string) string {
	return ""
}

/**
func SetStringFieldFromList(reader *bufio.Reader, field *string, name string, showList bool, possible []string) {
	if showList {
		fmt.Printf("%s%s %s\n", aurora.Bold(name), aurora.Bold(" (text)"), aurora.Italic("Possibilities:"))
	} else {
		fmt.Printf("%s%s %s\n--> ", aurora.Bold(name), aurora.Bold(" (text)"), aurora.Italic(fmt.Sprintf("default: «%s» (choose default with enter) type «l» to list possibilietes", *field)))
	}
}
*/

func AskInt(reader *bufio.Reader, name, desc string, defaultValue int) int {
	input := simplePrompt(reader, name, "int", desc, string(defaultValue))
	if input == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(input)
	if err != nil {
		logrus.Warn("Could not parse input as a number (int)")
		return AskInt(reader, name, desc, defaultValue)
	}
	return value
}

func AskIntWithHelp(reader *bufio.Reader, name string, defaultValue int, showList bool, possible map[int]string) int {
	if showList {
		fmt.Printf("%s%s %s", aurora.Bold(name), aurora.Bold(" (int)"), aurora.Italic("Possibilities:"))
		for value, explanation := range possible {
			fmt.Printf("\n%s %s", aurora.Bold(fmt.Sprintf("[%d]", value)), explanation)
		}
		fmt.Printf("\n--> ")
	} else {
		fmt.Printf("%s%s %s\n--> ", aurora.Bold(name), aurora.Bold(" (int)"), aurora.Italic(fmt.Sprintf("default: «%d» (choose default with enter) type «l» to list possibilietes", defaultValue)))
	}
	input, _ := reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if input == "" {
		return defaultValue
	} else if input == "l" {
		return AskIntWithHelp(reader, name, defaultValue, true, possible)
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		logrus.Warn("Could not parse input as a number (int)")
		return AskIntWithHelp(reader, name, defaultValue, showList, possible)
	}
	var exists bool
	for possibleValue := range possible {
		if value == possibleValue {
			exists = true
			break
		}
	}
	if !exists {
		logrus.Warn("This value is not valid, please see the list:")
		return AskIntWithHelp(reader, name, defaultValue, true, possible)
	}

	return value
}

func AskBool(reader *bufio.Reader, name, desc string, defaultValue bool) bool {
	input := simplePrompt(reader, name, "bool", desc, strconv.FormatBool(defaultValue))
	if input == "" {
		return defaultValue
	}

	if input == "y" || input == "1" || input == "true" {
		return true
	} else if input == "n" || input == "0" || input == "false" {
		return false
	} else {
		logrus.Warn("Could not parse input as a boolean value (bool). Please use y/1/true or n/0/false")
		return AskBool(reader, name, desc, defaultValue)
	}
}

func AskFloat(reader *bufio.Reader, name, desc string, defaultValue float64) float64 {
	input := simplePrompt(reader, name, "int", desc, fmt.Sprintf("%f", defaultValue))
	if input == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(input, 32)
	if err != nil {
		logrus.Warn("Could not parse input as a floating number (float)")
		return AskFloat(reader, name, desc, defaultValue)
	}
	return value
}

func AskDate(reader *bufio.Reader, name, desc string, defaultValue time.Time) string {
	input := simplePrompt(reader, name, IsoLayout, desc, defaultValue.Format(IsoLayout))
	if input == "" {
		return defaultValue.Format(IsoLayout)
	}
	value, err := time.Parse(IsoLayout, input)
	if err != nil {
		logrus.Warnf("Could not parse input as date with format: %s", IsoLayout)
		return AskDate(reader, name, desc, defaultValue)
	}
	return value.Format(IsoLayout)
}

func simplePrompt(reader *bufio.Reader, name, typeName, desc, defaultValue string) string {
	fmt.Printf("%s%s %s %s\n--> ",
		aurora.BrightCyan(aurora.Bold(name)),
		aurora.BrightCyan(fmt.Sprintf(" (%s)", typeName)),
		aurora.Yellow(fmt.Sprintf("%s", desc)),
		aurora.Green(fmt.Sprintf("Enter for default (%s)", defaultValue)))
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}

func searchPrompt(reader *bufio.Reader, name, desc string) string {
	fmt.Printf("%s %s: ",
		aurora.BrightCyan(fmt.Sprintf("Search for a %s", aurora.Bold(name))),
		aurora.Yellow(fmt.Sprintf("(%s)", desc)))
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}

func searchItemsPrompt(reader *bufio.Reader, items SearchItems, showList bool) (result string, back bool) {
	if len(items) == 0 {
		fmt.Printf("%s: ", aurora.BrightCyan("No entry found, 'T' for free text, any other key for search again"))
		input, _ := reader.ReadString('\n')
		if input == "T\n" {
			fmt.Print("--> ")
			input, _ := reader.ReadString('\n')
			return strings.Replace(input, "\n", "", -1), false
		}
		return "", true
	}
	if showList {
		for i := range items {
			fmt.Printf("%s %s %s\n",
				aurora.Yellow(fmt.Sprintf("%d)", i+1)),
				items[i].Name,
				aurora.Green(fmt.Sprintf("(%s)", items[i].Identifier)))
		}
	}
	fmt.Printf("%s ", aurora.BrightCyan("Choose item, 'S' to search again:"))
	input, _ := reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if input == "S" {
		return "", true
	}

	value, err := strconv.Atoi(input)
	if err != nil || value < 1 || value > len(items) {
		logrus.Error("invalid input, try again")
		return searchItemsPrompt(reader, items, false)
	}
	return items[value-1].Identifier, false
}
