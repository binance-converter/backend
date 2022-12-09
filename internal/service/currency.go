package service

import (
	"github.com/binance-converter/backend/core"
	"golang.org/x/net/context"
)

type CurrencyUserDb interface {
	AddUserCurrency(ctx context.Context, userId int, currency core.FullCurrency) (int, error)
	GetUserCurrencies(ctx context.Context, userId int, currencyType *core.CurrencyType) ([]core.
		FullCurrency, error)
	DeleteUserCurrency(ctx context.Context, userId int, currency core.CurrencyCode) error
	GetAvailableClassicCurrencies(ctx context.Context) ([]core.CurrencyCode, error)
	GetAvailableBanks(ctx context.Context, currency core.CurrencyCode) ([]core.CurrencyBank, error)
	GetAvailableCryptoCurrencies(ctx context.Context) ([]core.CurrencyCode, error)
}

type Currency struct {
	userDb CurrencyUserDb
}

func (c Currency) GetAvailableCurrencies(ctx context.Context,
	currencyType core.CurrencyType) (currencies []core.CurrencyCode, err error) {

	switch currencyType {
	case core.CurrencyTypeClassic:
		currencies, err = c.userDb.GetAvailableClassicCurrencies(ctx)
		if err != nil {
			return nil, err
		}
		break
	case core.CurrencyTypeCrypto:
		currencies, err = c.userDb.GetAvailableCryptoCurrencies(ctx)
		if err != nil {
			return nil, err
		}
		break
	}
	return currencies, nil
}

func (c Currency) GetAvailableBankByCurrency(ctx context.Context,
	currencyCode core.CurrencyCode) (banks []core.CurrencyBank, err error) {
	banks, err = c.userDb.GetAvailableBanks(ctx, currencyCode)
	return banks, err
}

func (c Currency) SetCurrency(ctx context.Context, currency core.FullCurrency) error {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return core.ErrorCurrencyNotAuthorized
	}
	// TODO: add validate currency
	_, err = c.userDb.AddUserCurrency(ctx, userId, currency)
	return err
}

func (c Currency) GetMyCurrencies(ctx context.Context,
	currencyType *core.CurrencyType) ([]core.FullCurrency, error) {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return nil, core.ErrorCurrencyNotAuthorized
	}
	currencies, err := c.userDb.GetUserCurrencies(ctx, userId, currencyType)

	return currencies, err
}

func (c Currency) DeleteCurrency(ctx context.Context, currency core.CurrencyCode) error {
	userId, err := core.ContextGetUserId(ctx)
	if err != nil {
		return core.ErrorCurrencyNotAuthorized
	}

	err = c.userDb.DeleteUserCurrency(ctx, userId, currency)
	return err
}
