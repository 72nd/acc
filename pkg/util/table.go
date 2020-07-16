package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"unicode/utf8"
)

type TableRow []string

func (r TableRow) Render(widths map[int]int) string {
	rsl := ""
	for i := range r {
		rsl = fmt.Sprintf("%s│ %s%s", rsl, r[i], strings.Repeat(" ", widths[i]-utf8.RuneCountInString(r[i])+1))
	}
	rsl += "│"
	return rsl
}

type Table struct {
	Header TableRow
	Rows   []TableRow
}

func (t Table) Render() string {
	if err := t.validate(); err != nil {
		logrus.Fatal(err)
	}
	widths := t.maxCellWidths()

	rsl := t.renderHeader(widths)
	for i := range t.Rows {
		rsl += "\n" + t.Rows[i].Render(widths)
	}
	rsl += "\n" + renderBottomLine(widths) + "\n"
	return rsl
}

func (t Table) validate() error {
	length := len(t.Header)
	for i := range t.Rows {
		if length != len(t.Rows[i]) {
			return fmt.Errorf("row %d has another number of cells (%d) than in the header (%d): %s", i+1, len(t.Rows[i]), length, t.Rows[i])
		}
	}
	return nil
}

func (t Table) maxCellWidths() map[int]int {
	rsl := make(map[int]int)
	for i := range t.Header {
		rsl[i] = len(t.Header[i])
	}
	for i := range t.Rows {
		for j := range t.Rows[i] {
			if rsl[j] < len(t.Rows[i][j]) {
				rsl[j] = len(t.Rows[i][j])
			}
		}
	}
	return rsl
}

func (t Table) renderHeader(widths map[int]int) string {
	rsl := "┌"
	for i := range widths {
		rsl = fmt.Sprintf("%s%s┬", rsl, strings.Repeat("─", widths[i]+2))
	}
	for i := range t.Header {
		t.Header[i] = strings.ToUpper(t.Header[i])
	}
	rsl = strings.TrimSuffix(rsl, "┬")
	rsl = fmt.Sprintf("%s┐\n%s\n%s", rsl, t.Header.Render(widths), renderSepLine(widths))

	return rsl
}

func renderSepLine(widths map[int]int) string {
	rsl := "├"
	for i := range widths {
		rsl = fmt.Sprintf("%s%s┼", rsl, strings.Repeat("─", widths[i]+2))
	}
	rsl = strings.TrimSuffix(rsl, "┼")
	rsl += "┤"
	return rsl
}

func renderBottomLine(widths map[int]int) string {
	rsl := "└"
	for i := range widths {
		rsl = fmt.Sprintf("%s%s┴", rsl, strings.Repeat("─", widths[i]+2))
	}
	rsl = strings.TrimSuffix(rsl, "┴")
	rsl += "┘"
	return rsl
}
