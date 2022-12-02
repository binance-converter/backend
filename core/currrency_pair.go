package core

var USDT = FullCurrency{
	CurrencyType: CurrencyTypeCrypto,
	CurrencyCode: "USDT",
}

var RUB_TINKOFF = FullCurrency{
	CurrencyType: CurrencyTypeClassic,
	CurrencyCode: "RUB",
	BankCode:     "TinkoffNew",
}

var KZT_KASPI = FullCurrency{
	CurrencyType: CurrencyTypeClassic,
	CurrencyCode: "KZT",
	BankCode:     "KaspiBank",
}

var ConverterPairs = []ConverterPair{
	{
		Currencies: []FullCurrency{
			RUB_TINKOFF,
			USDT,
		},
	},
	{
		Currencies: []FullCurrency{
			USDT,
			KZT_KASPI,
		},
	},
}
