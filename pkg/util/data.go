package util

type TransactionType int

const (
	CreditTransaction TransactionType = iota // Incoming transaction
	DebitTransaction                         // Outgoing transaction
)

const DateFormat = "2006-01-02"

