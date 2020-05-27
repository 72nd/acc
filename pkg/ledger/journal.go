package ledger

import (
	"gitlab.com/72th/acc/pkg/schema"
)

const HLedgerDateFormat = "2006-01-02"

type Journal []Entry 

func JournalFromStatement(acc schema.Acc) Journal {
	return Journal{}
}
