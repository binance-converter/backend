package handler

import (
	"github.com/binance-converter/backend-api/api/converter"
	"github.com/binance-converter/backend/core"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ConverterService interface {
	GetAvailableConverterPairs(ctx context.Context) ([]core.ConverterPair, error)
	SetConvertPair(ctx context.Context, converterPair core.ConverterPair) error
	GetMyConvertPairs(ctx context.Context) ([]core.ConverterPair, error)
	SetThresholdConvertPair(ctx context.Context, threshold core.ThresholdConvertPair) error
	GetMyThresholdsConvertPairs(ctx context.Context) ([]core.ThresholdConvertPair, error)
	GetCurrentExchange(ctx context.Context, converterPair core.ConverterPair) (core.Exchange, error)
}

type ConverterHandler struct {
	converter.UnimplementedConverterServer
	service ConverterService
}

func NewConverterHandler(service ConverterService) *ConverterHandler {
	return &ConverterHandler{service: service}
}

func (c ConverterHandler) GetAvailableConverterPairs(ctx context.Context,
	empty *emptypb.Empty) (*converter.ConverterPairs, error) {

	pairs, err := c.service.GetAvailableConverterPairs(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error get available converter pairs")
		switch err {
		case core.ErrorConverterNotAuthorized:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	protoPairs, err := convertCoreConverterPairsToProto(pairs)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"pairs": pairs,
		}).Error("error convert core converter pairs to proto")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return protoPairs, nil
}

func (c ConverterHandler) SetConvertPair(ctx context.Context,
	pair *converter.ConverterPair) (*emptypb.Empty, error) {
	corePair, err := convertProtoConverterPairToCore(pair)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"pair":  pair,
		}).Error("error convert proto converter pairs to core")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.service.SetConvertPair(ctx, corePair)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"corePair": corePair,
		}).Error("error set converter pair")
		switch err {
		case core.ErrorConverterNotAuthorized:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case core.ErrorConverterInvalidConverterPair:
			return nil, status.Error(codes.Code(
				converter.AdditionalErrorCode_INVALID_CONVERTER_PAIR), err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

func (c ConverterHandler) GetMyConvertPairs(ctx context.Context,
	empty *emptypb.Empty) (*converter.ConverterPairs, error) {

	pairs, err := c.service.GetMyConvertPairs(ctx)
	if err != nil {
		switch err {
		case core.ErrorConverterNotAuthorized:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	protoPairs, err := convertCoreConverterPairsToProto(pairs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return protoPairs, nil
}

func (c ConverterHandler) SetThresholdConvertPairs(ctx context.Context,
	pair *converter.ThresholdConvertPair) (*emptypb.Empty, error) {
	corePair, err := convertProtoThresholdConverterPair(pair)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = c.service.SetThresholdConvertPair(ctx, corePair)
	if err != nil {
		switch err {
		case core.ErrorConverterNotAuthorized:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case core.ErrorConverterInvalidConverterPair:
			return nil, status.Error(codes.Code(
				converter.AdditionalErrorCode_INVALID_CONVERTER_PAIR), err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (c ConverterHandler) GetMyThresholdConvertPairs(ctx context.Context,
	empty *emptypb.Empty) (*converter.ThresholdConvertPairs, error) {
	threshold, err := c.service.GetMyThresholdsConvertPairs(ctx)
	if err != nil {
		switch err {
		case core.ErrorConverterNotAuthorized:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	protoPairs, err := convertCoreThresholdConverterPairsToProto(threshold)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return protoPairs, nil
}

func (c ConverterHandler) GetCurrentExchange(ctx context.Context,
	pair *converter.ConverterPair) (*converter.Exchange, error) {
	corePair, err := convertProtoConverterPairToCore(pair)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	exchange, err := c.service.GetCurrentExchange(ctx, corePair)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"corePair": corePair,
			"error":    err.Error(),
		}).Error("error get current exchange")
		switch err {
		case core.ErrorConverterNotAuthorized:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case core.ErrorConverterInvalidConverterPair:
			return nil, status.Error(codes.Code(
				converter.AdditionalErrorCode_INVALID_CONVERTER_PAIR), err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return convertCoreExchangeToProto(exchange), nil
}

// ------------------------------------------------------------------------------------------------
// helper functions

func convertCoreConverterPairToProto(corePair core.ConverterPair) (*converter.ConverterPair,
	error) {
	pair := &converter.ConverterPair{}
	for _, currency := range corePair.Currencies {
		protoCurrency, err := convertCoreFullCurrencyToProto(currency)
		if err != nil {
			return nil, err
		}
		pair.ConverterPair = append(pair.ConverterPair, protoCurrency)
	}
	return pair, nil
}

func convertCoreConverterPairsToProto(corePairs []core.ConverterPair) (*converter.ConverterPairs,
	error) {
	pairs := &converter.ConverterPairs{}
	for _, corePair := range corePairs {
		pair, err := convertCoreConverterPairToProto(corePair)
		if err != nil {
			return nil, err
		}
		pairs.ConverterPairs = append(pairs.ConverterPairs, pair)
	}

	return pairs, nil
}

func convertProtoConverterPairToCore(protoConverterPair *converter.ConverterPair) (core.
	ConverterPair, error) {
	if protoConverterPair == nil {
		return core.ConverterPair{}, core.ErrorConverterEmptyInputArg
	}
	coreConverterPair := core.ConverterPair{}
	for _, protoCurrency := range protoConverterPair.ConverterPair {
		coreCurrency, err := convertProtoFullCurrencyToCore(protoCurrency)
		if err != nil {
			return coreConverterPair, err
		}
		coreConverterPair.Currencies = append(coreConverterPair.Currencies, coreCurrency)
	}
	return coreConverterPair, nil
}

func convertCoreExchangeToProto(coreExchange core.Exchange) *converter.Exchange {
	return &converter.Exchange{
		Exchange: float32(coreExchange),
	}
}

func convertCoreThresholdConverterPairToProto(coreThreshold core.ThresholdConvertPair) (
	*converter.ThresholdConvertPair, error) {
	threshold := &converter.ThresholdConvertPair{
		Exchange: convertCoreExchangeToProto(coreThreshold.Exchange),
	}
	converterPair, err := convertCoreConverterPairToProto(coreThreshold.ConverterPair)
	if err != nil {
		return nil, err
	}
	threshold.ConverterPair = converterPair
	return threshold, nil
}

func convertCoreThresholdConverterPairsToProto(coreThreshold []core.ThresholdConvertPair) (
	*converter.ThresholdConvertPairs, error) {
	threshold := &converter.ThresholdConvertPairs{}
	for _, coreThresholdPair := range coreThreshold {
		converterPair, err := convertCoreThresholdConverterPairToProto(coreThresholdPair)
		if err != nil {
			return nil, err
		}
		threshold.ConverterPairs = append(threshold.ConverterPairs, converterPair)
	}
	return threshold, nil
}

func convertProtoExchangeToCore(protoExchange *converter.Exchange) (core.Exchange, error) {
	if protoExchange == nil {
		return core.Exchange(0), core.ErrorConverterEmptyInputArg
	}
	return core.Exchange(protoExchange.Exchange), nil
}

func convertProtoThresholdConverterPair(protoThreshold *converter.ThresholdConvertPair) (core.
	ThresholdConvertPair, error) {
	coreThreshold := core.ThresholdConvertPair{}

	var err error

	coreThreshold.Exchange, err = convertProtoExchangeToCore(protoThreshold.Exchange)
	if err != nil {
		return coreThreshold, err
	}

	coreThreshold.ConverterPair, err = convertProtoConverterPairToCore(protoThreshold.
		ConverterPair)

	return coreThreshold, err
}
