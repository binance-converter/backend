package core

import (
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	UserIdCtx = "userId"
)

var allContextValues = []string{UserIdCtx}

var (
	ErrorContextErrorGettingUserIdFromContext = errors.New("error getting user id from context")
)

func ContextGetUserId(ctx context.Context) (int, error) {
	userId := ctx.Value(UserIdCtx)
	if userId == nil {
		return 0, ErrorContextErrorGettingUserIdFromContext
	}
	id, ok := userId.(int)
	if !ok {
		return 0, ErrorContextErrorGettingUserIdFromContext
	}
	return id, nil
}

func ContextAddUserId(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, UserIdCtx, userId)
}

func LogContext(ctx context.Context) *logrus.Fields {
	fields := make(logrus.Fields)
	for _, val := range allContextValues {
		fields[val] = ctx.Value(val)
	}
	return &fields
}
