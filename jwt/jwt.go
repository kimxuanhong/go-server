package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CustomClaims struct {
	UserInfo
	jwt.RegisteredClaims
}

type Jwt struct {
	secretKey string
	expIn     int
}

func NewJwt(configs ...*Config) *Jwt {
	cfg := GetConfig(configs...)
	return &Jwt{
		secretKey: cfg.SecretKey,
		expIn:     cfg.ExpIn,
	}
}

func (j *Jwt) IssueToken(ctx context.Context, user UserInfo) (string, error) {
	now := time.Now().UTC()
	claims := CustomClaims{
		UserInfo: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.expIn))),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *Jwt) ExpIn() int {
	return j.expIn
}

func (j *Jwt) Validate(tokenStr string) (*UserInfo, error) {
	var claims CustomClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims.UserInfo, nil
}
