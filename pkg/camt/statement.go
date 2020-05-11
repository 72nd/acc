package camt

import "gitlab.com/72th/acc/pkg/schema"

type BankToCustomerStatement struct {
	CamtPath string
}

func NewBankToCustomerStatement(path string) BankToCustomerStatement {
	return BankToCustomerStatement{
		CamtPath: path,
	}
}

func (s BankToCustomerStatement) Transactions() []schema.Transaction {

}