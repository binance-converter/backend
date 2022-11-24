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
	SignUpUserByTelegram(ctx context.Context, data *core.ServiceSignUpUserByTelegramData) error
}

type AuthHandler struct {
	service AuthService

	auth.UnimplementedAuthServer
}

func NewAuthService(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (a *AuthHandler) SignUpUserByTelegram(ctx context.Context,
	request *auth.SignUpUserByTelegramRequest) (*emptypb.Empty, error) {

	err := a.service.SignUpUserByTelegram(ctx, &core.ServiceSignUpUserByTelegramData{
		ChatId:       request.ChatId,
		UserName:     request.UserName,
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		LanguageCode: request.LanguageCode,
	})

	if err != nil {
		switch err {
		case core.ErrorServiceAuthUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "fuck you")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
