package main

import (
	"github.com/binance-converter/backend/internal/service"
	userDbPostgres "github.com/binance-converter/backend/internal/storage/user_db/postgres"
	"github.com/binance-converter/backend/internal/transport/grpc"
	"github.com/binance-converter/backend/internal/transport/grpc/handler"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func main() {
	logger := logrus.New()

	ctx := context.Background()

	userDbConfig := userDbPostgres.Config{}

	postgresDb, err := userDbPostgres.NewPostgresDB(ctx, userDbConfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("error connect to postgres database")
	}

	transaction := userDbPostgres.NewTransaction(postgresDb)

	userDb := userDbPostgres.NewUserDB(postgresDb, transaction)

	authService := service.NewAuth(userDb)
	converterService := service.NewConverter(nil, userDb)
	currencyService := service.NewCurrency(nil, userDb)

	auth := handler.NewAuthHandler(authService)
	converter := handler.NewConverterHandler(converterService)
	currencies := handler.NewCurrenciesHandler(currencyService)
	exchangePlot := handler.NewExchangePlotHandler(nil)

	grpcServer := grpc.NewServer(logger, auth, converter, currencies, exchangePlot, nil)

	grpcServer.ListenAndServe(9000)
}
