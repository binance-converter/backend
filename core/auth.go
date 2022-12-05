package core

import "errors"

type AddUser struct {
	ChatId       *int64
	UserName     *string
	FirstName    *string
	LastName     *string
	LanguageCode *string
}

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
	ErrorAuthServiceInternalError         = errors.New("internal error")
	ErrorAuthServiceUserNotFound          = errors.New("user not found")
)
