package query

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/util"
)

type SearchTerms []SearchTerm

func searchTermsFromUserInput(input string) SearchTerms {
	var rsl SearchTerms
	ele := separate(input, ",")
	for i := range ele {
		rsl = append(rsl, searchTermFromUserInput(ele[i]))
	}
	return rsl
}

type SearchTerm struct {
	Key   *regexp.Regexp
	Value *regexp.Regexp
}

func newSearchTerm(key, value string) SearchTerm {
	keyRe, err := regexp.Compile(key)
	if err != nil {
		logrus.Fatalf("error while parsing \"%s\" as key-regex from term \"%s:%s\"", key, key, value)
	}
	valRe, err := regexp.Compile(value)
	if err != nil {
		logrus.Fatalf("error while parsing \"%s\" as value-regex from term \"%s:%s\"", value, key, value)
	}
	return SearchTerm{
		Key:   keyRe,
		Value: valRe,
	}
}

func searchTermFromUserInput(input string) SearchTerm {
	ele := separate(input, ":")
	if len(ele) != 2 {
		logrus.Fatalf("input \"%s\" couldn't be parsed as KEY:VALUE, use \\: to escape colons inside your pattern", input)
	}
	return newSearchTerm(ele[0], ele[1])
}

func (s SearchTerm) matchKey(input string) bool {
	return s.Key.MatchString(input)
}

func (s SearchTerm) matchValue(input string) bool {
	return s.Value.MatchString(input)
}

type DateTerms []DateTerm

func dateTermsFromUserInput(input string) DateTerms {
	var rsl DateTerms
	ele := separate(input, ",")
	for i := range ele {
		rsl = append(rsl, dateTermFromUserInput(ele[i]))
	}
	return rsl

}

type DateTerm struct {
	Key  *regexp.Regexp
	From time.Time
	To   time.Time
}

func dateTermFromUserInput(input string) DateTerm {
	ele := separate(input, ":")
	if len(ele) != 3 {
		logrus.Fatalf("input \"%s\" couldn't be parsed as KEY:FROM:TO, use \\: to escape colons inside your pattern", input)
	}
	key, err := regexp.Compile(ele[0])
	if err != nil {
		logrus.Fatalf("error while parsing \"%s\" as key-regex from term \"%s:%s:%s\"", ele[0], ele[0], ele[1], ele[2])
	}
	from, err := time.Parse(util.DateFormat, ele[1])
	if err != nil {
		logrus.Fatalf("error while parsing \"%s\" as from-date from term \"%s:%s:%s\"", ele[1], ele[0], ele[1], ele[2])
	}
	to, err := time.Parse(util.DateFormat, ele[1])
	if err != nil {
		logrus.Fatalf("error while parsing \"%s\" as to-date from term \"%s:%s:%s\"", ele[2], ele[0], ele[1], ele[2])
	}
	return DateTerm{
		Key: key,
		From: from,
		To: to,
	}
}

func (d DateTerm) matchKey(input string) bool {
	return d.Key.MatchString(input)
}

func (d DateTerm) matchRange(input string) bool {
	date, err := time.Parse(util.DateFormat, input)
	if err != nil {
		logrus.Fatalf("\"%s\" couldn't be parsed as date: %s", input, err)
	}
	return !date.Before(d.From) && !date.After(d.To)
}


func separate(input, sep string) []string {
	const esc = "ESCAPE"
	if strings.Contains(input, esc) {
		logrus.Fatalf("match string may not contain \"%s\"", esc)
	}
	input = strings.Replace(input, fmt.Sprintf("\\%s", sep), esc, -1)
	ele := strings.Split(input, sep)
	for i := range ele {
		ele[i] = strings.Replace(ele[i], esc, sep, -1)
	}
	return ele
}
