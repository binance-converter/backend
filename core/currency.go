package core

import "errors"

type CurrencyType int32

const (
	CurrencyTypeCrypto   CurrencyType = 0
	ECurrencyTypeClassic CurrencyType = 1
)

type CurrencyCode string
type CurrencyBank string

type FullCurrency struct {
	CurrencyType CurrencyType
	CurrencyCode CurrencyCode
	BankCode     CurrencyBank
}

var (
	ErrorCurrencyInvalidCurrencyType = errors.New("invalid currency type")
	ErrorCurrencyInvalidCurrencyCode = errors.New("invalid currency code")
	ErrorCurrencyInvalidBankCode     = errors.New("invalid bank code")
	ErrorCurrencyInternal            = errors.New("internal error")
)