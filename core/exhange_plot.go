package core

import (
	"errors"
	timeInterval "github.com/go-follow/time-interval"
	"image"
)

type TimeInterval timeInterval.Span

type PlotParams struct {
	ConverterPair ConverterPair
	TimeInterval  TimeInterval
}

type Plot image.RGBA

var (
	ErrorExchangePlotEmptyInputArg          = errors.New("empty input arguments")
	ErrorExchangePlotInvalidTimeInterval    = errors.New("invalid time interval")
	ErrorExchangePlotNoDataForTimeInterval  = errors.New("no data for time interval")
	ErrorExchangePlotInvalidConverterPair   = errors.New("invalid converter pair")
	ErrorExchangePlotCovertPairNotSupported = errors.New("converter pair not supported")
)
