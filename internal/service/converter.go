package service

import (
	"github.com/binance-converter/backend/core"
	"golang.org/x/net/context"
)

type converterBinanceApi interface {
	getExchange(ctx context.Context, converterPair core.ConverterPair) (core.Exchange, error)
}

type converterUserDb interface {
	SetConverterPair(ctx context.Context, userId int, converterPair core.ConverterPair) error
	GetConverterPairs(ctx context.Context, userId int) ([]core.ConverterPair, error)
	SetThresholdConvertPair(ctx context.Context, userId int,
		threshold core.ThresholdConvertPair) error
	GetThresholdConvertPair(ctx context.Context, userId int) ([]core.ThresholdConvertPair, error)
}

type Converter struct {
	binanceApi converterBinanceApi
	UserDb     converterUserDb
}

func (c *Converter) GetAvailableConverterPairs(ctx context.Context) ([]core.ConverterPair, error) {
	converterPairs := core.ConverterPairs
	for _, valF := range core.ConverterPairs {
		for _, valS := range core.ConverterPairs {
			converterPair, err := c.makeSecondLevelPair(valF, valS)
			if err == nil {
				converterPairs = append(converterPairs, converterPair)
			}

			converterPair, err = c.makeSecondLevelPair(valS, valF)
			if err == nil {
				converterPairs = append(converterPairs, converterPair)
			}
		}
	}
	return converterPairs, nil
}

func (c *Converter) SetConvertPair(ctx context.Context, converterPair core.ConverterPair) error {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return core.ErrorConverterNotAuthorized
	}

	return c.UserDb.SetConverterPair(ctx, userId, converterPair)
}

func (c *Converter) GetMyConvertPairs(ctx context.Context) ([]core.ConverterPair, error) {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return nil, core.ErrorConverterNotAuthorized
	}
	converterPairs, err := c.UserDb.GetConverterPairs(ctx, userId)
	if err != nil {
		return nil, err
	}
	return converterPairs, nil
}

func (c *Converter) SetThresholdConvertPair(ctx context.Context,
	threshold core.ThresholdConvertPair) error {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return core.ErrorConverterNotAuthorized
	}

	return c.UserDb.SetThresholdConvertPair(ctx, userId, threshold)
}

func (c *Converter) GetMyThresholdsConvertPairs(ctx context.Context) ([]core.ThresholdConvertPair,
	error) {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return nil, core.ErrorConverterNotAuthorized
	}

	thresholds, err := c.UserDb.GetThresholdConvertPair(ctx, userId)
	if err != nil {
		return nil, err
	}
	return thresholds, nil
}

func (c *Converter) GetCurrentExchange(ctx context.Context,
	converterPair core.ConverterPair) (core.Exchange, error) {
	exchange, err := c.binanceApi.getExchange(ctx, converterPair)
	if err != nil {
		return core.Exchange(0), err
	}
	return exchange, nil
}

func (c *Converter) makeSecondLevelPair(first core.ConverterPair,
	second core.ConverterPair) (core.ConverterPair, error) {
	if first.Currencies[1] != second.Currencies[0] {
		return core.ConverterPair{},
			core.ErrorConverterInvalidConverterPair
	}

	return core.ConverterPair{
		Currencies: []core.FullCurrency{
			first.Currencies[0],
			first.Currencies[1],
			second.Currencies[1],
		},
	}, nil
}