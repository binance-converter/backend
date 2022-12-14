package binance_api

import (
	"github.com/binance-converter/backend/core"
	binanceP2PApi "github.com/binance-converter/binance-p2p-api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type BinanceApi struct {
	api binanceP2PApi.BinanceP2PApi
}

func NewBinanceApi() *BinanceApi {
	return &BinanceApi{}
}

func (b *BinanceApi) GetExchange(ctx context.Context,
	converterPair core.ConverterPair) (core.Exchange, error) {
	if len(converterPair.Currencies) != 2 {
		return 0, core.ErrorBinanceApiInvalidConverterPair
	}

	if converterPair.Currencies[0].CurrencyType == converterPair.Currencies[1].CurrencyType {
		logrus.WithFields(logrus.Fields{
			"currency1": converterPair.Currencies[0],
			"currency2": converterPair.Currencies[1],
		}).Error("Currencies is equal")
		return 0, core.ErrorBinanceApiInvalidConverterPair
	}

	var assets, fiat, tradeType string
	transAmount := float64(10000)
	var payTypes []string

	if converterPair.Currencies[0].CurrencyType == core.CurrencyTypeCrypto {
		assets = string(converterPair.Currencies[0].CurrencyCode)
		payTypes = []string{string(converterPair.Currencies[1].BankCode)}
		fiat = string(converterPair.Currencies[1].CurrencyCode)
		tradeType = binanceP2PApi.OperationSell
	} else {
		assets = string(converterPair.Currencies[1].CurrencyCode)
		payTypes = []string{string(converterPair.Currencies[0].BankCode)}
		fiat = string(converterPair.Currencies[0].CurrencyCode)
		tradeType = binanceP2PApi.OperationBuy
	}

	exchange, _, _, err := b.api.GetExchange(assets, fiat, payTypes, tradeType,
		transAmount)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"assets":      assets,
			"fiat":        fiat,
			"payTypes":    payTypes,
			"tradeType":   tradeType,
			"transAmount": transAmount,
			"err":         err,
		}).Error("Error get exchange")
	}

	return core.Exchange(exchange), err
}
