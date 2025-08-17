package ctxval

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type userInfoKey struct{}

func SetUserInfo(parent context.Context, val interface{ jwt.Claims }) context.Context {
	return context.WithValue(parent, &userInfoKey{}, val)
}

func GetUserInfo(ctx context.Context) (interface{ jwt.Claims }, bool) {
	if val := ctx.Value(&userInfoKey{}); val != nil {
		userInfo, ok := val.(interface{ jwt.Claims })
		if ok {
			return userInfo, ok
		}
	}

	return nil, false
}
