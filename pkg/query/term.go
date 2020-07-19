package query

import (
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/72nd/acc/pkg/util"
)

type SearchTerms []SearchTerm

func searchTermsFromUserInput(input string, caseSensitive bool) SearchTerms {
	var rsl SearchTerms
	ele := util.EscapedSplit(input, ",")
	for i := range ele {
		rsl = append(rsl, searchTermFromUserInput(ele[i], caseSensitive))
	}
	return rsl
}

type SearchTerm struct {
	Key   *regexp.Regexp
	Value *regexp.Regexp
}

func newSearchTerm(key, value string, caseSensitive bool) SearchTerm {
	if !caseSensitive {
		key = strings.ToLower(key)
		value = strings.ToLower(value)
	}
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

func searchTermFromUserInput(input string, caseSensitive bool) SearchTerm {
	ele := util.EscapedSplit(input, ":")
	if len(ele) != 2 {
		logrus.Fatalf("input \"%s\" couldn't be parsed as KEY:VALUE, use \\: to escape colons inside your pattern", input)
	}
	return newSearchTerm(ele[0], ele[1], caseSensitive)
}

func (s SearchTerm) matchKey(input string, caseSensitive bool) bool {
	if !caseSensitive {
		input = strings.ToLower(input)
	}
	return s.Key.MatchString(input)
}

func (s SearchTerm) matchValue(input string, caseSensitive bool) bool {
	if !caseSensitive {
		input = strings.ToLower(input)
	}
	rsl := s.Value.MatchString(input)
	return rsl
}

type DateTerms []DateTerm

func dateTermsFromUserInput(input string) DateTerms {
	var rsl DateTerms
	ele := util.EscapedSplit(input, ",")
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
	ele := util.EscapedSplit(input, ":")
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
		Key:  key,
		From: from,
		To:   to,
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
