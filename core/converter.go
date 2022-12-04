package core

import "errors"

type ConverterPair struct {
	Currencies []FullCurrency
}

type Exchange float32

type ThresholdConvertPair struct {
	ConverterPair ConverterPair
	Exchange      Exchange
}

var (
	ErrorConverterEmptyInputArg              = errors.New("empty input arguments")
	ErrorConverterInvalidConverterPair       = errors.New("invalid converter pair")
	ErrorConverterNotAuthorized              = errors.New("not authorized")
	ErrorConverterConverterPairAlreadyExists = errors.New("converter pair already exists")
	ErrorConverterConverterPairNotFound      = errors.New("converter pair not found")
)
