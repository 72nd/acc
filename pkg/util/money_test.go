package util

import (
	"testing"

	"gopkg.in/yaml.v3"
)

type TestStruct struct {
	Amount Money `yaml:"amount"`
}

func TestNewMonyFromParse(t *testing.T) {
	m1, err := NewMonyFromParse("CHF 23.42")
	if err != nil {
		t.Error(err)
	}
	checkMoney(t, m1, 2342, "CHF")

	m2, err := NewMonyFromParse("CHF 123.42")
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
	given := `amount: CHF 23.42`
	expected := TestStruct{Amount: NewMoney(2342, "CHF")}

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
	given := NewMoney(12342, "CHF")
	// expected := "amount: CHF 123.42\n       "
	b, _ := yaml.Marshal(TestStruct{given})
	t.Log(string(b))
	/*
		if string(b) != expected {
			t.Errorf("should be \"%s\" but is \"%s\"", expected, string(b))
		} */
}
