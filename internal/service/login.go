package service

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/ffauzann/loan-service/internal/constant"
	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, req *model.LoginRequest) (res *model.LoginResponse, err error) {
	// Get user.
	user, err := s.repository.db.GetUserByOneOfIdentifier(ctx, req.UserId)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Validate user status.
	if user.Status != constant.UserStatusActive {
		return nil, constant.ErrUserIsNotActive
	}

	// Validate password.
	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, constant.ErrInvalidUsernamePassword
		}
		util.LogContext(ctx).Error(err.Error())
		return
	}

	var token model.Token
	// Generate access_token
	token.AccessToken, err = s.generateToken(ctx, &model.GenerateTokenRequest{
		User:      user,
		TokenType: constant.TokenTypeAccess,
		Extended:  false,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Generate refresh_token
	token.RefreshToken, err = s.generateToken(ctx, &model.GenerateTokenRequest{
		User:      user,
		TokenType: constant.TokenTypeRefresh,
		Extended:  req.RememberMe,
	})
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct response.
	res = &model.LoginResponse{
		Token: token,
	}

	return
}

func (s *service) generateToken(ctx context.Context, req *model.GenerateTokenRequest) (token string, err error) { //nolint
	// Preserve variables.
	var (
		iss, exp, pk string
		d            time.Duration
	)

	// Determine which config to fetch based on its type.
	switch req.TokenType {
	case constant.TokenTypeAccess:
		// signingKey = s.config.Jwt.AccessToken.SigningKey
		iss = s.config.Jwt.AccessToken.Iss
		exp = s.config.Jwt.AccessToken.Exp
		pk = s.config.Jwt.AsymmetricKeys[0].PrivateKey

		d, err = time.ParseDuration(exp)
		if err != nil {
			util.LogContext(ctx).Error(err.Error())
			return
		}
	case constant.TokenTypeRefresh:
		// signingKey = s.config.Jwt.RefreshToken.SigningKey
		iss = s.config.Jwt.RefreshToken.Iss
		exp = s.config.Jwt.RefreshToken.Exp
		pk = s.config.Jwt.AsymmetricKeys[0].PrivateKey

		if req.Extended {
			exp = s.config.Jwt.RefreshToken.ExtendedExp
		}

		d, err = time.ParseDuration(exp)
		if err != nil {
			util.LogContext(ctx).Error(err.Error())
			return
		}
	}

	// Decode base64.
	bPrivateKey, err := base64.StdEncoding.DecodeString(pk)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Construct base claims.
	claims := model.Claims{
		UserId:      req.User.Id,
		Name:        req.User.Name,
		Email:       req.User.Email,
		PhoneNumber: req.User.PhoneNumber,
		RoleId:      req.User.RoleId,
		TokenType:   req.TokenType,
		Extended:    req.Extended,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   req.User.Email,
			Issuer:    iss,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
		},
	}

	// Create token & sign.
	// t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = s.config.Jwt.AsymmetricKeys[0].Kid

	// Parse private key PEM to struct.
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(bPrivateKey)
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return
	}

	// Sign token and return.
	return t.SignedString(privateKey)
}
