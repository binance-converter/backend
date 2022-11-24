package handler

import (
	"bytes"
	"errors"
	"github.com/binance-converter/backend-api/api/exchange_plot"
	"github.com/binance-converter/backend/core"
	timeInterval "github.com/go-follow/time-interval"
	"github.com/openlyinc/pointy"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image"
	"image/png"
)

var (
	ErrorExchangePlotHandlerInputArgIsNull = errors.New("error input arg is null")
)

type exchangePlotService interface {
	GetExchangePlot(ctx context.Context, params core.PlotParams) (core.Plot, error)
}

type ExchangePlotHandler struct {
	exchange_plot.UnimplementedExchangePlotServer
	service exchangePlotService
}

func NewExchangePlotService() *ExchangePlotHandler {
	return &ExchangePlotHandler{}
}

func (e *ExchangePlotHandler) GetExchangePlot(ctx context.Context,
	params *exchange_plot.PlotParams) (*exchange_plot.Plot, error) {
	corePlotParams, err := convertProtoPlotParamsToCore(params)
	if err != nil {
		return nil, status.Error(codes.Code(
			exchange_plot.AdditionalErrorCode_INVALID_TIME_INTERVAL), err.Error())
	}

	plot, err := e.service.GetExchangePlot(ctx, corePlotParams)
	if err != nil {
		switch err {
		case core.ErrorExchangePlotCovertPairNotSupported:
			return nil, status.Error(codes.Code(
				exchange_plot.AdditionalErrorCode_NOT_SUPPORTED_CONVERTER_PAIR), err.Error())
		case core.ErrorExchangePlotInvalidConverterPair:
			return nil, status.Error(codes.Code(
				exchange_plot.AdditionalErrorCode_INVALID_CONVERTER_PAIR), err.Error())
		case core.ErrorExchangePlotInvalidTimeInterval:
			return nil, status.Error(codes.Code(
				exchange_plot.AdditionalErrorCode_INVALID_TIME_INTERVAL), err.Error())
		case core.ErrorExchangePlotNoDataForTimeInterval:
			return nil, status.Error(codes.Code(
				exchange_plot.AdditionalErrorCode_NO_DATA_FOR_TIME_INTERVAL), err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	protoPlot, err := convertCorePlotToProto(plot)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return protoPlot, nil
}

func convertProtoTimeIntervalToCore(protoTimeInterval *exchange_plot.TimeInterval) (core.
	TimeInterval, error) {
	start, err := timeInterval.New(protoTimeInterval.Start.AsTime(),
		protoTimeInterval.Start.AsTime())

	return core.TimeInterval(start), err
}

func convertProtoPlotParamsToCore(protoPlotParams *exchange_plot.PlotParams) (core.PlotParams,
	error) {
	if protoPlotParams == nil {
		return core.PlotParams{}, ErrorExchangePlotHandlerInputArgIsNull
	}

	coreConverterPair, err := convertProtoConverterPairToCore(protoPlotParams.Pair)
	if err != nil {
		return core.PlotParams{}, err
	}

	coreTimeInterval, err := convertProtoTimeIntervalToCore(protoPlotParams.Interval)
	if err != nil {
		return core.PlotParams{}, err
	}

	plotParams := core.PlotParams{
		ConverterPair: coreConverterPair,
		TimeInterval:  coreTimeInterval,
	}

	return plotParams, nil
}

func convertCorePlotToProto(plot core.Plot) (*exchange_plot.Plot, error) {
	raw := bytes.Buffer{}

	err := png.Encode(&raw, pointy.Pointer(image.RGBA(plot)))
	if err != nil {
		return nil, err
	}

	return &exchange_plot.Plot{
		Image: raw.Bytes(),
	}, nil
}
