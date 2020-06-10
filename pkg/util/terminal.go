package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

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
	} else if input == "N" && newFunction != nil {
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

type winSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func TerminalWidth() (int, error) {
	w := &winSize{}
	code, _, errNr := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(w)))

	if int(code) == -1 {
		return 0, fmt.Errorf("syscall to determine terminal width failed with error code \"%d\"", errNr)
	}
	return int(w.Col), nil
}
