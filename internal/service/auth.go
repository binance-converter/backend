package service

import (
	"github.com/binance-converter/backend/core"
	"golang.org/x/net/context"
)

type AuthDB interface {
	AddUser(ctx context.Context, user core.AddUser) (int, error)
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
