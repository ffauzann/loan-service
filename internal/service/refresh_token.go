package service

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/constant"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/model"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

func (s *service) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (res *model.RefreshTokenResponse, err error) {
	// Extract and validate token string.
	// claims, ok := util.ExtractClaimsFromString(ctx, req.RefreshToken, s.config.Jwt.RefreshToken.SigningKey)
	// if !ok {
	// 	err = constant.ErrInvalidToken
	// 	util.LogContext(ctx).Error(err.Error())
	// 	return
	// }
	claims, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct user model from claims.
	user := &model.User{
		CommonModel: model.CommonModel{
			Id: claims.UserId,
		},
		Name:        claims.Name,
		Email:       claims.Email,
		PhoneNumber: claims.PhoneNumber,

		RoleId: claims.RoleId,
	}

	// Generate access_token.
	accessToken, err := s.generateToken(ctx, &model.GenerateTokenRequest{
		User:      user,
		TokenType: constant.TokenTypeAccess,
		Extended:  false,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Generate refresh_token.
	refreshToken, err := s.generateToken(ctx, &model.GenerateTokenRequest{
		User:      user,
		TokenType: constant.TokenTypeRefresh,
		Extended:  claims.Extended,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &model.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return
}

func (s *service) verifyRefreshToken(ctx context.Context, refreshTokenString string) (verifiedClaims *model.Claims, err error) { //nolint
	// Get JWKS.
	jwks, err := s.Jwks(ctx)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Parse unverified just to get `kid`.
	unverifiedToken, _, err := jwt.NewParser().ParseUnverified(refreshTokenString, &model.Claims{})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		err = constant.ErrInvalidToken
		return
	}

	// Find the right public key based on its `kid`.
	var publicKey *rsa.PublicKey
	for _, v := range jwks {
		if v.KeyID == unverifiedToken.Header["kid"] {
			var nBytes []byte
			nBytes, err = base64.RawURLEncoding.DecodeString(v.Modulus)
			if err != nil {
				util.LogContext(ctx).Error(err.Error())
				return
			}

			var eBytes []byte
			eBytes, err = base64.RawURLEncoding.DecodeString(v.Exponent)
			if err != nil {
				util.LogContext(ctx).Error(err.Error())
				return
			}

			publicKey = &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: int(new(big.Int).SetBytes(eBytes).Int64()),
			}

			break
		}
	}

	// Parse verified.
	verifiedClaims = &model.Claims{} // Preserve memory address.
	verifiedToken, err := jwt.NewParser(
		jwt.WithIssuedAt(),
		jwt.WithIssuer(s.config.Jwt.RefreshToken.Iss),
	).ParseWithClaims(refreshTokenString, verifiedClaims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("Invalid signing method. expected : %v | got : %s", jwt.SigningMethodRS256.Alg(), t.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		err = constant.ErrInvalidToken
		return
	}

	// Validate whether verified is nil.
	if verifiedToken == nil {
		util.LogContext(ctx).Error("Empty token")
		err = constant.ErrInvalidToken
		return
	}

	// Verify and return claims if it is valid.
	_, ok := verifiedToken.Claims.(interface{ jwt.Claims })
	if ok && verifiedToken.Valid {
		return verifiedClaims, nil
	}

	// Invalid token, return error.
	err = constant.ErrInvalidToken
	return
}
