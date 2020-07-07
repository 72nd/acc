package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Rhymond/go-money"
	"gopkg.in/yaml.v3"
)

type Money struct {
	*money.Money
}

func NewMoney(amount int64, code string) Money {
	return Money{money.New(amount, code)}
}

func (m *Money) UnmarshalYAML(value *yaml.Node) error {
	parts := strings.Split(value.Value, ".")
	if len(parts) != 2 {
		return fmt.Errorf("couldn't parse \"%s\" as a amount, format xxxx.xx", value.Value)
	}
	part1, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse \"%s\" as a amount as it contains not only numbers", value.Value)
	}
	part2, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse \"%s\" as a amount as it contains not only numbers", value.Value)
	}
	if part2 / 100 > 0 {
		return fmt.Errorf("couldn't parse \"%s\" as a amount, as only two digits are allowed after the point", value.Value)
	}
	amount := part1 * 100 + part2
	m.Money = money.New(amount, "CHF")
	return nil
}

func (m Money) MarshalYAML() (interface{}, error) {
	return map[string]string{
		"amount":   strconv.FormatInt(m.Amount(), 10),
		"currency": m.Currency().Code}, nil
}
