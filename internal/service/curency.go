package service

import (
	"github.com/binance-converter/backend/core"
	"golang.org/x/net/context"
)

type currencyUserDb interface {
	addCurrency(ctx context.Context, userId int, currency core.FullCurrency) error
	getCurrencies(ctx context.Context, userId int, currencyType *core.CurrencyType) ([]core.
		FullCurrency, error)
	deleteCurrency(ctx context.Context, userId int, currency core.CurrencyCode) error
}

type currencyBinanceApi interface {
	getAvailableClassicCurrencies(ctx context.Context) ([]core.CurrencyCode, error)
	getAvailableBanks(ctx context.Context, currency core.CurrencyCode) ([]core.CurrencyBank, error)
	getAvailableCryptoCurrencies(ctx context.Context) ([]core.CurrencyCode, error)
}

type Currency struct {
	userDb     currencyUserDb
	binanceApi currencyBinanceApi
}

func (c Currency) GetAvailableCurrencies(ctx context.Context,
	currencyType core.CurrencyType) (currencies []core.CurrencyCode, err error) {

	switch currencyType {
	case core.CurrencyTypeClassic:
		currencies, err = c.binanceApi.getAvailableClassicCurrencies(ctx)
		if err != nil {
			return nil, err
		}
		break
	case core.CurrencyTypeCrypto:
		currencies, err = c.binanceApi.getAvailableCryptoCurrencies(ctx)
		if err != nil {
			return nil, err
		}
		break
	}
	return currencies, nil
}

func (c Currency) GetAvailableBankByCurrency(ctx context.Context,
	currencyCode core.CurrencyCode) (banks []core.CurrencyBank, err error) {
	banks, err = c.binanceApi.getAvailableBanks(ctx, currencyCode)
	return banks, err
}

func (c Currency) SetCurrency(ctx context.Context, currency core.FullCurrency) error {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return core.ErrorCurrencyNotAuthorized
	}
	err = c.userDb.addCurrency(ctx, userId, currency)
	return err
}

func (c Currency) GetMyCurrencies(ctx context.Context,
	currencyType *core.CurrencyType) ([]core.FullCurrency, error) {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return nil, core.ErrorCurrencyNotAuthorized
	}
	currencies, err := c.userDb.getCurrencies(ctx, userId, currencyType)

	return currencies, err
}

func (c Currency) DeleteCurrency(ctx context.Context, currency core.CurrencyCode) error {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return core.ErrorCurrencyNotAuthorized
	}

	err = c.userDb.deleteCurrency(ctx, userId, currency)
	return err
}
