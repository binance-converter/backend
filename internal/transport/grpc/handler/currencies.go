package handler

import (
	"errors"
	"github.com/binance-converter/backend-api/api/currencies"
	"github.com/binance-converter/backend/core"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type currenciesService interface {
	GetAvailableCurrencies(ctx context.Context, currencyType core.CurrencyType) ([]core.
		CurrencyCode, error)
	GetAvailableBankByCurrency(ctx context.Context, currencyCode core.CurrencyCode) ([]core.
		CurrencyBank, error)
	SetCurrency(ctx context.Context, currency core.FullCurrency) error
	GetMyCurrencies(ctx context.Context, currencyType core.CurrencyType) ([]core.FullCurrency,
		error)
	DeleteCurrency(ctx context.Context, currencyType core.FullCurrency) error
}

type CurrenciesHandler struct {
	currencies.UnimplementedCurrenciesServer
	service currenciesService
}

func (c *CurrenciesHandler) GetAvailableCurrencies(ctx context.Context,
	currencyType *currencies.CurrencyType) (*currencies.CurrencyCodes, error) {
	coreCurrencyCode, err := convertProtoCurrencyTypeToCore(currencyType)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	coreCurrencies, err := c.service.GetAvailableCurrencies(ctx, coreCurrencyCode)
	if err != nil {
		switch err {
		case core.ErrorInvalidCurrencyCode:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertCoreCurrencyCodesToProto(coreCurrencies), nil
}

func (c *CurrenciesHandler) GetAvailableBankByCurrency(ctx context.Context,
	code *currencies.CurrencyCode) (*currencies.BankNames, error) {

	banks, err := c.service.GetAvailableBankByCurrency(ctx, convertProtoCurrencyCodeToCore(code))
	if err != nil {
		switch err {
		case core.ErrorInvalidCurrencyCode:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return convertCoreCurrencyBanksToProto(banks), nil
}

func (c *CurrenciesHandler) SetCurrency(ctx context.Context,
	currency *currencies.FullCurrency) (*emptypb.Empty, error) {
	coreCurrency, err := convertProtoFullCurrencyToCore(currency)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.service.SetCurrency(ctx, coreCurrency)

	if err != nil {
		switch err {
		case core.ErrorInvalidCurrencyType:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case core.ErrorInvalidCurrencyCode:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case core.ErrorInvalidBankCode:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (c *CurrenciesHandler) GetMyCurrencies(ctx context.Context,
	currencyType *currencies.CurrencyType) (*currencies.FullCurrencies, error) {
	coreCurrencyType, err := convertProtoCurrencyTypeToCore(currencyType)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	coreCurrencies, err := c.service.GetMyCurrencies(ctx, coreCurrencyType)

	if err != nil {
		switch err {
		case core.ErrorInvalidCurrencyType:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	protoCurrencies, err := convertCoreFullCurrenciesToProto(coreCurrencies)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return protoCurrencies, nil
}

func (c *CurrenciesHandler) DeleteCurrency(ctx context.Context,
	currency *currencies.FullCurrency) (*emptypb.Empty, error) {

	coreCurrency, err := convertProtoFullCurrencyToCore(currency)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.service.DeleteCurrency(ctx, coreCurrency)
	if err != nil {
		switch err {
		case core.ErrorInvalidCurrencyType:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case core.ErrorInvalidCurrencyCode:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case core.ErrorInvalidBankCode:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

// ------------------------------------------------------------------------------------------------
// helper functions

func convertProtoCurrencyTypeToCore(currencyType *currencies.CurrencyType) (core.CurrencyType,
	error) {
	switch currencyType.Type {
	case currencies.ECurrencyType_CRYPTO:
		return core.CurrencyTypeCrypto, nil
	case currencies.ECurrencyType_CLASSIC:
		return core.ECurrencyTypeClassic, nil
	}
	return 0, errors.New("error parsing currency type")
}

func convertCoreCurrencyTypeToProto(currencyType core.CurrencyType) (*currencies.CurrencyType,
	error) {
	switch currencyType {
	case core.CurrencyTypeCrypto:
		return &currencies.CurrencyType{
			Type: currencies.ECurrencyType_CRYPTO,
		}, nil
	case core.ECurrencyTypeClassic:
		return &currencies.CurrencyType{
			Type: currencies.ECurrencyType_CLASSIC,
		}, nil
	}
	return nil, errors.New("error parsing currency type")
}

func convertProtoCurrencyCodeToCore(protoCurrencyCode *currencies.CurrencyCode) (
	currencyCode core.CurrencyCode) {
	currencyCode = core.CurrencyCode(protoCurrencyCode.CurrencyCode)
	return currencyCode
}

func convertCoreCurrencyCodeToProto(coreCurrencyCode core.CurrencyCode) (
	currencyCode *currencies.CurrencyCode) {
	currencyCode.CurrencyCode = string(coreCurrencyCode)
	return currencyCode
}

func convertCoreCurrencyCodesToProto(coreCurrencies []core.CurrencyCode) (currencyCodes *currencies.
	CurrencyCodes) {
	for _, currency := range coreCurrencies {
		currencyCodes.CurrencyCodes = append(currencyCodes.CurrencyCodes,
			convertCoreCurrencyCodeToProto(currency))
	}
	return currencyCodes
}

func convertCoreCurrencyBankToProto(coreCurrencyBank core.CurrencyBank) (CurrencyBank *currencies.
	BankName) {
	CurrencyBank.BankName = string(coreCurrencyBank)
	return CurrencyBank
}

func convertProtoCurrencyBankToCore(protoCurrencyBank *currencies.
	BankName) core.CurrencyBank {
	return core.CurrencyBank(protoCurrencyBank.BankName)
}

func convertCoreCurrencyBanksToProto(coreCurrencyBanks []core.CurrencyBank) (
	CurrencyBank *currencies.BankNames) {
	for _, coreCurrencyBank := range coreCurrencyBanks {
		bankName := convertCoreCurrencyBankToProto(coreCurrencyBank)
		CurrencyBank.BankNames = append(CurrencyBank.BankNames,
			bankName)
	}
	return CurrencyBank
}

func convertProtoFullCurrencyToCore(protoCurrency *currencies.FullCurrency) (core.
	FullCurrency, error) {
	coreCurrencyType, err := convertProtoCurrencyTypeToCore(protoCurrency.Type)
	if err != nil {
		return core.FullCurrency{}, err
	}

	return core.FullCurrency{
		CurrencyType: coreCurrencyType,
		CurrencyCode: convertProtoCurrencyCodeToCore(protoCurrency.CurrencyCode),
		BankCode:     convertProtoCurrencyBankToCore(protoCurrency.BankName),
	}, nil
}

func convertCoreFullCurrencyToProto(coreCurrency core.FullCurrency) (*currencies.
	FullCurrency, error) {
	currencyType, err := convertCoreCurrencyTypeToProto(coreCurrency.CurrencyType)
	if err != nil {
		return nil, err
	}
	return &currencies.FullCurrency{
		Type:         currencyType,
		CurrencyCode: convertCoreCurrencyCodeToProto(coreCurrency.CurrencyCode),
		BankName:     convertCoreCurrencyBankToProto(coreCurrency.BankCode),
	}, nil
}

func convertCoreFullCurrenciesToProto(coreCurrencies []core.FullCurrency) (currencies *currencies.
	FullCurrencies, err error) {
	for _, currency := range coreCurrencies {
		protoCurrency, err := convertCoreFullCurrencyToProto(currency)
		if err != nil {
			return nil, err
		}
		currencies.FullCurrencies = append(currencies.FullCurrencies, protoCurrency)
	}
	return currencies, nil
}
