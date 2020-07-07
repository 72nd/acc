package util

import (
	"testing"

	"gopkg.in/yaml.v3"
)

type TestStruct struct {
	Amount Money `yaml:"amount"`
}

func TestUnmarshalYAML(t *testing.T) {
	given := `amount: 23.42`
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
	given := NewMoney(2342, "CHF")
	b, _ := yaml.Marshal(&given)
	t.Log(string(b))
}
