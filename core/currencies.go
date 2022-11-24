package core

import "errors"

type CurrencyType int32

const (
	CurrencyTypeCrypto   CurrencyType = 0
	ECurrencyTypeClassic CurrencyType = 1
)

type CurrencyCode string
type CurrencyBank string

var (
	ErrorInvalidCurrencyType = errors.New("invalid currency type")
	ErrorInvalidCurrencyCode = errors.New("invalid currency code")
	ErrorInvalidBankCode     = errors.New("invalid bank code")
	ErrorInternal            = errors.New("internal error")
)

type FullCurrency struct {
	CurrencyType CurrencyType
	CurrencyCode CurrencyCode
	BankCode     CurrencyBank
}
