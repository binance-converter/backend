package handler

import (
	"github.com/binance-converter/backend-api/api/auth"
	"github.com/binance-converter/backend/core"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	SignUpUserByTelegram(ctx context.Context, data core.ServiceSignUpUserByTelegramData) error
}

type AuthHandler struct {
	service AuthService

	auth.UnimplementedAuthServer
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (a *AuthHandler) SignUpUserByTelegram(ctx context.Context,
	request *auth.SignUpUserByTelegramRequest) (*emptypb.Empty, error) {

	coreRequest, err := convertProtoSignUpUserByTelegramToCore(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = a.service.SignUpUserByTelegram(ctx, coreRequest)
	if err != nil {
		switch err {
		case core.ErrorAuthServiceAuthUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "fuck you")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

func convertProtoSignUpUserByTelegramToCore(protoRequest *auth.SignUpUserByTelegramRequest) (
	core.ServiceSignUpUserByTelegramData, error) {
	if protoRequest == nil {
		return core.ServiceSignUpUserByTelegramData{}, core.ErrorAuthServiceEmptyInputArg
	}
	return core.ServiceSignUpUserByTelegramData{
		ChatId:       protoRequest.ChatId,
		UserName:     protoRequest.UserName,
		FirstName:    protoRequest.FirstName,
		LastName:     protoRequest.LastName,
		LanguageCode: protoRequest.LanguageCode,
	}, nil
}
