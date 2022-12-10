package binance_api

import (
	"github.com/binance-converter/backend/core"
	binanceP2PApi "github.com/binance-converter/binance-p2p-api"
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
		return 0, core.ErrorBinanceApiInvalidConverterPair
	}

	var assets, fiat, tradeType string
	transAmount := float64(10000)
	var payTypes []string

	if converterPair.Currencies[0].CurrencyType == core.CurrencyTypeCrypto {
		assets = string(converterPair.Currencies[1].CurrencyCode)
		payTypes = []string{string(converterPair.Currencies[1].BankCode)}
		fiat = string(converterPair.Currencies[0].CurrencyType)
		tradeType = binanceP2PApi.OperationSell
	} else {
		assets = string(converterPair.Currencies[0].CurrencyCode)
		payTypes = []string{string(converterPair.Currencies[0].BankCode)}
		fiat = string(converterPair.Currencies[1].CurrencyType)
		tradeType = binanceP2PApi.OperationBuy
	}

	exchange, _, _, err := b.api.GetExchange(assets, fiat, payTypes, tradeType,
		transAmount)

	return core.Exchange(exchange), err
}
