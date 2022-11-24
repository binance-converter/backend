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
	ErrorConverterInvalidConverterPair = errors.New("invalid converter pair")
)
