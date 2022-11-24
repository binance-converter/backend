package core

import "errors"

type ServiceSignUpUserByTelegramData struct {
	ChatId       int64
	UserName     string
	FirstName    string
	LastName     string
	LanguageCode string
}

var (
	ErrorAuthServiceAuthUserAlreadyExists = errors.New("error user already exists")
	ErrorAuthInternalError                = errors.New("internal error")
)
