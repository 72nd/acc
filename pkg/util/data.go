package util

type TransactionType int

const (
	CreditTransaction TransactionType = iota // Incoming transaction
	DebitTransaction                         // Outgoing transaction
)


