package util

import (
	"context"
	"fmt"

	"github.com/ffauzann/loan-service/client"
	"github.com/ffauzann/loan-service/internal/model"
	authCtx "github.com/ffauzann/loan-service/pkg/common/auth/jwt/ctxval"
	"github.com/golang-jwt/jwt/v5"
)

// Deprecated: Service is no longer using symmetric signature.
func ExtractClaimsFromString(ctx context.Context, strToken, signingKey string) (*model.Claims, bool) {
	hmacSecret := []byte(signingKey)

	token, err := jwt.ParseWithClaims(strToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		LogContext(ctx).Error(err.Error())
		return nil, false
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		LogContext(ctx).Error(fmt.Sprintf("%v", token.Valid))
		return claims, true
	}

	LogContext(ctx).Error(fmt.Sprintf("%v", token.Claims))
	return nil, false
}

func ClaimsFromContext(ctx context.Context) (claims *client.Claims, ok bool) {
	iClaims, ok := authCtx.GetUserInfo(ctx)
	if !ok {
		return
	}

	claims, ok = iClaims.(*client.Claims)
	return
}
