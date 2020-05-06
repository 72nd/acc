package util

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func AskString(reader *bufio.Reader, name, desc, defaultValue string) string {
	input := simplePrompt(reader, name, "string", desc, defaultValue)
	if input == "" {
		return defaultValue
	}
	return input
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

func simplePrompt(reader *bufio.Reader, name, typeName, desc, defaultValue string) string {
	fmt.Printf("%s%s %s %s\n--> ",
		aurora.BrightCyan(aurora.Bold(name)),
		aurora.BrightCyan(fmt.Sprintf(" (%s)", typeName)),
		aurora.Yellow(fmt.Sprintf("%s", desc)),
		aurora.Green(fmt.Sprintf("Enter for default (%s)", defaultValue)))
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}