package handler

import (
	"errors"
	"github.com/binance-converter/backend-api/api/currencies"
	"github.com/binance-converter/backend/core"
	"github.com/sirupsen/logrus"
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
	GetMyCurrencies(ctx context.Context, currencyType *core.CurrencyType) ([]core.FullCurrency,
		error)
	DeleteCurrency(ctx context.Context, currencyType core.CurrencyCode) error
}

type CurrenciesHandler struct {
	currencies.UnimplementedCurrenciesServer
	service currenciesService
}

func NewCurrenciesHandler(service currenciesService) *CurrenciesHandler {
	return &CurrenciesHandler{service: service}
}

func (c *CurrenciesHandler) GetAvailableCurrencies(ctx context.Context,
	currencyType *currencies.CurrencyType) (*currencies.CurrencyCodes, error) {
	coreCurrencyType, err := convertProtoCurrencyTypeToCore(currencyType)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"currency_type": currencyType.GetType(),
			"error":         err.Error(),
		}).Error("error convert proto currency type to core")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	coreCurrencies, err := c.service.GetAvailableCurrencies(ctx, coreCurrencyType)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"currency_type": coreCurrencyType,
			"error":         err.Error(),
		}).Error("error get available currencies")
		switch err {
		case core.ErrorCurrencyInvalidCurrencyCode:
			return nil, status.Error(codes.Code(
				currencies.AdditionalErrorCode_INVALID_CURRENCY_CODE), err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertCoreCurrencyCodesToProto(coreCurrencies), nil
}

func (c *CurrenciesHandler) GetAvailableBankByCurrency(ctx context.Context,
	code *currencies.CurrencyCode) (*currencies.BankNames, error) {

	coreCode, err := convertProtoCurrencyCodeToCore(code)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"currency_code": code.CurrencyCode,
			"error":         err.Error(),
		}).Error("error convert proto currency type to core")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	banks, err := c.service.GetAvailableBankByCurrency(ctx, coreCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"currency_code": coreCode,
			"error":         err.Error(),
		}).Error("error get available banks")
		switch err {
		case core.ErrorCurrencyInvalidCurrencyCode:
			return nil, status.Error(codes.Code(
				currencies.AdditionalErrorCode_INVALID_CURRENCY_CODE), err.Error())
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
		logrus.WithFields(logrus.Fields{
			"currency": currency,
			"error":    err.Error(),
		}).Error("error convert proto full currency to core")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.service.SetCurrency(ctx, coreCurrency)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"currency": coreCurrency,
			"error":    err.Error(),
		}).Error("error set currency")
		switch err {
		case core.ErrorCurrencyInvalidCurrencyType:
			return nil, status.Error(codes.Code(
				currencies.AdditionalErrorCode_INVALID_CURRENCY_TYPE), err.Error())
		case core.ErrorCurrencyInvalidCurrencyCode:
			return nil, status.Error(codes.Code(
				currencies.AdditionalErrorCode_INVALID_CURRENCY_CODE), err.Error())
		case core.ErrorCurrencyInvalidBankCode:
			return nil, status.Error(codes.Code(currencies.AdditionalErrorCode_INVALID_BANK_CODE),
				err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (c *CurrenciesHandler) GetMyCurrencies(ctx context.Context,
	currencyType *currencies.CurrencyType) (*currencies.FullCurrencies, error) {

	coreCurrencyType := new(core.CurrencyType)
	var err error
	*coreCurrencyType, err = convertProtoCurrencyTypeToCore(currencyType)
	if err != nil {
		coreCurrencyType = nil
	}

	coreCurrencies, err := c.service.GetMyCurrencies(ctx, coreCurrencyType)

	if err != nil {
		switch err {
		case core.ErrorCurrencyInvalidCurrencyType:
			return nil, status.Error(codes.Code(
				currencies.AdditionalErrorCode_INVALID_CURRENCY_TYPE), err.Error())
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
	currency *currencies.CurrencyCode) (*emptypb.Empty, error) {

	coreCurrency, err := convertProtoCurrencyCodeToCore(currency)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.service.DeleteCurrency(ctx, coreCurrency)
	if err != nil {
		switch err {
		case core.ErrorCurrencyInvalidCurrencyType:
			return nil, status.Error(codes.Code(
				currencies.AdditionalErrorCode_INVALID_CURRENCY_TYPE), err.Error())
		case core.ErrorCurrencyInvalidCurrencyCode:
			return nil, status.Error(codes.Code(
				currencies.AdditionalErrorCode_INVALID_CURRENCY_CODE), err.Error())
		case core.ErrorCurrencyInvalidBankCode:
			return nil, status.Error(codes.Code(currencies.AdditionalErrorCode_INVALID_BANK_CODE),
				err.Error())
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
	if currencyType == nil {
		return core.CurrencyType(0), core.ErrorCurrencyEmptyInputArg
	}
	switch currencyType.Type {
	case currencies.ECurrencyType_CRYPTO:
		return core.CurrencyTypeCrypto, nil
	case currencies.ECurrencyType_CLASSIC:
		return core.CurrencyTypeClassic, nil
	}
	return 0, core.ErrorCurrencyInvalidCurrencyType
}

func convertCoreCurrencyTypeToProto(currencyType core.CurrencyType) (*currencies.CurrencyType,
	error) {
	switch currencyType {
	case core.CurrencyTypeCrypto:
		return &currencies.CurrencyType{
			Type: currencies.ECurrencyType_CRYPTO,
		}, nil
	case core.CurrencyTypeClassic:
		return &currencies.CurrencyType{
			Type: currencies.ECurrencyType_CLASSIC,
		}, nil
	}
	return nil, errors.New("error parsing currency type")
}

func convertProtoCurrencyCodeToCore(protoCurrencyCode *currencies.CurrencyCode) (
	core.CurrencyCode, error) {
	if protoCurrencyCode == nil {
		return "", core.ErrorCurrencyEmptyInputArg
	}
	currencyCode := core.CurrencyCode(protoCurrencyCode.CurrencyCode)
	return currencyCode, nil
}

func convertCoreCurrencyCodeToProto(coreCurrencyCode core.CurrencyCode) (
	currencyCode *currencies.CurrencyCode) {
	currencyCode = &currencies.CurrencyCode{}
	currencyCode.CurrencyCode = string(coreCurrencyCode)
	return currencyCode
}

func convertCoreCurrencyCodesToProto(coreCurrencies []core.CurrencyCode) (currencyCodes *currencies.
	CurrencyCodes) {
	currencyCodes = &currencies.CurrencyCodes{}
	for _, currency := range coreCurrencies {
		currencyCodes.CurrencyCodes = append(currencyCodes.CurrencyCodes,
			convertCoreCurrencyCodeToProto(currency))
	}
	return currencyCodes
}

func convertCoreCurrencyBankToProto(coreCurrencyBank core.CurrencyBank) (CurrencyBank *currencies.
	BankName) {
	CurrencyBank = &currencies.BankName{}
	CurrencyBank.BankName = string(coreCurrencyBank)
	return CurrencyBank
}

func convertProtoCurrencyBankToCore(protoCurrencyBank *currencies.
	BankName) (core.CurrencyBank, error) {
	if protoCurrencyBank == nil {
		return "", core.ErrorCurrencyEmptyInputArg
	}
	return core.CurrencyBank(protoCurrencyBank.BankName), nil
}

func convertCoreCurrencyBanksToProto(coreCurrencyBanks []core.CurrencyBank) (
	CurrencyBank *currencies.BankNames) {
	CurrencyBank = &currencies.BankNames{}
	for _, coreCurrencyBank := range coreCurrencyBanks {
		bankName := convertCoreCurrencyBankToProto(coreCurrencyBank)
		CurrencyBank.BankNames = append(CurrencyBank.BankNames,
			bankName)
	}
	return CurrencyBank
}

func convertProtoFullCurrencyToCore(protoCurrency *currencies.FullCurrency) (core.
	FullCurrency, error) {

	if protoCurrency == nil {
		return core.FullCurrency{}, core.ErrorCurrencyEmptyInputArg
	}

	coreCurrencyType, err := convertProtoCurrencyTypeToCore(protoCurrency.Type)
	if err != nil {
		return core.FullCurrency{}, err
	}

	currencyCode, err := convertProtoCurrencyCodeToCore(protoCurrency.CurrencyCode)
	if err != nil {
		return core.FullCurrency{}, err
	}

	bankCode, err := convertProtoCurrencyBankToCore(protoCurrency.BankName)
	if err != nil {
		return core.FullCurrency{}, err
	}

	return core.FullCurrency{
		CurrencyType: coreCurrencyType,
		CurrencyCode: currencyCode,
		BankCode:     bankCode,
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

func convertCoreFullCurrenciesToProto(coreCurrencies []core.FullCurrency) (protoCurrencies *currencies.
	FullCurrencies, err error) {
	protoCurrencies = &currencies.FullCurrencies{}
	for _, currency := range coreCurrencies {
		protoCurrency, err := convertCoreFullCurrencyToProto(currency)
		if err != nil {
			return nil, err
		}
		protoCurrencies.FullCurrencies = append(protoCurrencies.FullCurrencies, protoCurrency)
	}
	return protoCurrencies, nil
}
