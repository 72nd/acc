package util

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type TestStruct struct {
	Amount Money `yaml:"amount"`
}

func TestNewMonyFromParse(t *testing.T) {
	m1, err := NewMonyFromParse("23.42 CHF")
	if err != nil {
		t.Error(err)
	}
	checkMoney(t, m1, 2342, "CHF")

	m2, err := NewMonyFromParse("123.42 CHF")
	if err != nil {
		t.Error(err)
	}
	checkMoney(t, m2, 12342, "CHF")
}

func checkMoney(t *testing.T, money Money, expectedAmount int64, expectedCode string) {
	if money.Amount() != expectedAmount {
		t.Errorf("amount should be \"%d\" but is \"%d\"", expectedAmount, money.Amount())
	}
	if money.Currency().Code != expectedCode {
		t.Errorf("currency code should be \"%s\" but is \"%s\"", expectedCode, money.Currency().Code)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	given := `amount: 2342.42 CHF`
	expected := TestStruct{Amount: NewMoney(234242, "CHF")}

	var s TestStruct
	err := yaml.Unmarshal([]byte(given), &s)
	if err != nil {
		t.Error(err)
	}
	if ok, _ := s.Amount.Equals(expected.Amount.Money); !ok {
		t.Errorf("%s is not equal %s", expected.Amount.Display(), s.Amount.Display())
	}
}

func TestMarshalYAML(t *testing.T) {
	given := NewMoney(112342, "CHF")
	expected := "amount: 1123.42 CHF"

	b, _ := yaml.Marshal(TestStruct{given})
	if strings.TrimSpace(string(b)) != expected {
		t.Errorf("should be \"%s\" but is \"%s\"", expected, string(b))
	}
}
