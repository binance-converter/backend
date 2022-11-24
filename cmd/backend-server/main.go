package main

import (
	"github.com/binance-converter/backend/internal/transport/grpc"
	"github.com/binance-converter/backend/internal/transport/grpc/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	auth := handler.NewAuthHandler(nil)
	converter := handler.NewConverterHandler(nil)
	currencies := handler.NewCurrenciesHandler(nil)
	exchangePlot := handler.NewExchangePlotHandler(nil)

	grpcServer := grpc.NewServer(logger, auth, converter, currencies, exchangePlot)

	grpcServer.ListenAndServe(9000)
}
