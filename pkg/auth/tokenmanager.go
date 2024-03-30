package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidIdentityType = errors.New("invalid identity type")
)

type TokenManager interface {
	Generate(uid uint) (string, int64, error)
	Verify(tokenStr string) (uint, error)
}

type JWT struct {
	key    []byte
	expire time.Duration
}

func NewJWT(secretKey []byte, expiration time.Duration) TokenManager {
	return &JWT{
		key:    secretKey,
		expire: expiration,
	}
}

func (j *JWT) Generate(uid uint) (string, int64, error) {
	now := time.Now()
	expireAt := now.Add(j.expire).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uid,
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"exp": expireAt,
	})
	token, err := claims.SignedString(j.key)
	if err != nil {
		return "", 0, err
	}
	return token, expireAt, nil
}

func (j *JWT) Verify(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.key, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idStr, err := claims.GetSubject()
		if err != nil {
			return 0, err
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return 0, ErrInvalidIdentityType
		}
		return uint(id), nil
	}

	return 0, err
}
