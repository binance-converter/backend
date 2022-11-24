package grpc

import (
	"fmt"
	"github.com/binance-converter/backend-api/api/auth"
	"github.com/binance-converter/backend-api/api/converter"
	"github.com/binance-converter/backend-api/api/currencies"
	"github.com/binance-converter/backend-api/api/exchange_plot"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	Logger *logrus.Logger

	auth         auth.AuthServer
	converter    converter.ConverterServer
	currencies   currencies.CurrenciesServer
	exchangePlot exchange_plot.ExchangePlotServer

	srv *grpc.Server
}

func NewServer(logger *logrus.Logger, auth auth.AuthServer,
	converter converter.ConverterServer, currencies currencies.CurrenciesServer,
	exchangePlot exchange_plot.ExchangePlotServer) *Server {
	logrusLogger := logrus.NewEntry(logger)
	return &Server{
		Logger:       logger,
		auth:         auth,
		converter:    converter,
		currencies:   currencies,
		exchangePlot: exchangePlot,
		srv: grpc.NewServer(
			grpc.StreamInterceptor(
				grpc_middleware.ChainStreamServer(
					grpc_logrus.StreamServerInterceptor(logrusLogger),
					grpc_recovery.StreamServerInterceptor(),
				)),
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_logrus.UnaryServerInterceptor(logrusLogger),
					grpc_recovery.UnaryServerInterceptor(),
				)),
		),
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	auth.RegisterAuthServer(s.srv, s.auth)
	converter.RegisterConverterServer(s.srv, s.converter)
	currencies.RegisterCurrenciesServer(s.srv, s.currencies)
	exchange_plot.RegisterExchangePlotServer(s.srv, s.exchangePlot)

	if err := s.srv.Serve(lis); err != nil {
		return err
	}

	return nil
}
