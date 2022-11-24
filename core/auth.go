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
	ErrorAuthServiceEmptyInputArg         = errors.New("empty input arguments")
	ErrorAuthServiceAuthUserAlreadyExists = errors.New("error user already exists")
	ErrorAuthInternalError                = errors.New("internal error")
)
