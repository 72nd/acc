package util

import (
	"fmt"
	"regexp"
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

func NewMonyFromParse(value string) (Money, error) {
	re := regexp.MustCompile(`^(\d*)\.(\d{2})\s([A-z]{3})$`)
	if !re.MatchString(value) {
		return Money{}, fmt.Errorf("given string \"%s\" doesn't match format \"USD 00000.00\"", value)
	}
	rsl := re.FindStringSubmatch(value)
	if len(rsl) != 4 {
		return Money{}, fmt.Errorf("regex submatch of string \"%s\" returned array with length != 4", value)
	}
	part1, err := strconv.ParseInt(rsl[1], 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("couldn't parse \"%s\" as number (int64)", rsl[1])
	}
	part2, err := strconv.ParseInt(rsl[2], 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("couldn't parse \"%s\" as number (int64)", rsl[2])
	}
	amount := part1*100 + part2

	return Money{money.New(amount, rsl[3])}, nil
}

func NewMonyFromDotNotation(value, code string) (Money, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return Money{}, fmt.Errorf("couldn't parse \"%s\" as a amount, format xxxx.xx", value)
	}
	part1, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("couldn't parse \"%s\" as a amount as it contains not only numbers", value)
	}
	part2, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("couldn't parse \"%s\" as a amount as it contains not only numbers", value)
	}
	if part2/100 > 0 {
		return Money{}, fmt.Errorf("couldn't parse \"%s\" as a amount, as only two digits are allowed after the point", value)
	}
	amount := part1*100 + part2
	return NewMoney(amount, code), nil
}

// NewMoneyFromFloat returns a new Money object from the given string.
// As floating point numbers are not a good idea for money this shouldn't
// be used. This method is only here for the Bimpf import.
func NewMoneyFromFloat(value float64, currency string) Money {
	amount := int64(value * 100)
	return NewMoney(amount, currency)
}

func (m *Money) UnmarshalYAML(value *yaml.Node) error {
	money, err := NewMonyFromParse(value.Value)
	if err != nil {
		return err
	}
	m.Money = money.Money
	return nil
}

func (m Money) Value() string {
	part1 := m.Amount() / 100
	part2 := m.Amount() - part1*100
	if part2 < 10 {
		return fmt.Sprintf("%d.0%d %s", part1, part2, m.Currency().Code)
	}
	return fmt.Sprintf("%d.%d %s", part1, part2, m.Currency().Code)
}

func (m Money) MarshalYAML() (interface{}, error) {
	return m.Value(), nil
}

func (m Money) String() string {
	return m.Display()
}
