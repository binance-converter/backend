package service

import (
	"github.com/binance-converter/backend/core"
	"golang.org/x/net/context"
)

type AuthDB interface {
	AddUser(ctx context.Context, user core.AddUser) (int, error)
	ValidateUser(ctx context.Context, chatId int) (int, error)
}

type Auth struct {
	db AuthDB
}

func NewAuth(db AuthDB) *Auth {
	return &Auth{
		db: db,
	}
}

func (a *Auth) SignUpUserByTelegram(ctx context.Context,
	data core.ServiceSignUpUserByTelegramData) error {
	addUser := convertServiceSignUpUserByTelegramDataToAddUser(data)
	_, err := a.db.AddUser(ctx, addUser)
	return err
}

func (a *Auth) ValidateUserByChatId(ctx context.Context, chatId int) (int, error) {
	userId, err := a.db.ValidateUser(ctx, chatId)
	return userId, err
}

func convertServiceSignUpUserByTelegramDataToAddUser(data core.
	ServiceSignUpUserByTelegramData) core.AddUser {
	return core.AddUser{
		ChatId:       &data.ChatId,
		UserName:     &data.UserName,
		FirstName:    &data.FirstName,
		LastName:     &data.LastName,
		LanguageCode: &data.LanguageCode,
	}
}
