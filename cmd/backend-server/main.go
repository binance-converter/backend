package main

import (
	"context"
	"github.com/binance-converter/backend/internal/service"
	userDbPostgres "github.com/binance-converter/backend/internal/storage/user_db/postgres"
	"github.com/binance-converter/backend/internal/transport/grpc"
	"github.com/binance-converter/backend/internal/transport/grpc/handler"
	"github.com/binance-converter/backend/pkg/binance_api"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/sirupsen/logrus"
)

type appConfig struct {
	Grpc struct {
		Port int `env:"GRPC_PORT"`
	}
	PostgresUserDb struct {
		Host     string `env:"POSTGRES_USER_DB_HOST"`
		Port     int    `env:"POSTGRES_USER_DB_PORT"`
		Username string `env:"POSTGRES_USER_DB_USER_NAME"`
		Password string `env:"POSTGRES_USER_DB_PASSWORD"`
		DBName   string `env:"POSTGRES_USER_DB_NAME"`
	}
}

func main() {
	setupLogs()

	cfg, err := initConfig()
	if err != nil {
		logrus.Warning(err)
	}

	logger := logrus.New()

	ctx := context.Background()

	userDbConfig := userDbPostgres.Config{
		Host:     cfg.PostgresUserDb.Host,
		Port:     cfg.PostgresUserDb.Port,
		Username: cfg.PostgresUserDb.Username,
		Password: cfg.PostgresUserDb.Password,
		DBName:   cfg.PostgresUserDb.DBName,
	}

	postgresDb, err := userDbPostgres.NewPostgresDB(ctx, userDbConfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("error connect to postgres database")
	}

	transaction := userDbPostgres.NewTransaction(postgresDb)

	userDb := userDbPostgres.NewUserDB(postgresDb, transaction)

	bApi := binance_api.NewBinanceApi()

	authService := service.NewAuth(userDb)
	converterService := service.NewConverter(bApi, userDb)
	currencyService := service.NewCurrency(userDb)

	auth := handler.NewAuthHandler(authService)
	converter := handler.NewConverterHandler(converterService)
	currencies := handler.NewCurrenciesHandler(currencyService)
	exchangePlot := handler.NewExchangePlotHandler(nil)

	grpcServer := grpc.NewServer(logger, auth, converter, currencies, exchangePlot, authService)

	err = grpcServer.ListenAndServe(cfg.Grpc.Port)
	if err != nil {
		logrus.Fatal(err)
	}
}

func setupLogs() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006.01.02 15:04:05",
		FullTimestamp:   true,
		DisableSorting:  true,
	})
}

func initConfig() (appConfig, error) {
	var cfg appConfig

	yamlFeeder := feeder.Yaml{Path: "config.yaml"}
	envFeeder := feeder.Env{}
	dotEnvFeeder := feeder.DotEnv{Path: ".env"}

	err := config.New().AddFeeder(envFeeder, dotEnvFeeder, yamlFeeder).AddStruct(&cfg).Feed()

	logrus.WithFields(logrus.Fields{
		"grpcPort": cfg.Grpc.Port,
		"dbHost":   cfg.PostgresUserDb.Host,
		"dbPort":   cfg.PostgresUserDb.Port,
		"dbUser":   cfg.PostgresUserDb.Username,
		"dbName":   cfg.PostgresUserDb.DBName,
	}).Info("app configs")

	return cfg, err
}
