package service

import (
	"errors"
	"time"

	"github.com/bestetufan/beste-store/config"
	"github.com/bestetufan/beste-store/internal/domain/entity"
	"github.com/dgrijalva/jwt-go"
)

type (
	JWTAuthService struct {
		cfg config.Config
	}

	JwtClaims struct {
		UserId int    `json:"id,omitempty"`
		Email  string `json:"email,omitempty"`
		jwt.StandardClaims
	}
)

func NewJWTAuthService(c config.Config) *JWTAuthService {
	return &JWTAuthService{
		cfg: c,
	}
}

func (s *JWTAuthService) VerifyToken(tokenString string) (bool, *JwtClaims) {
	claims := &JwtClaims{}
	token, _ := getTokenFromString(tokenString, s.cfg.JWTSecret, claims)
	if token.Valid {
		if e := claims.Valid(); e == nil {
			return true, claims
		}
	}
	return false, claims
}

func (s *JWTAuthService) CreateToken(user entity.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"iss":   s.cfg.JWTIss,
		"exp":   time.Now().Add(s.cfg.JWTExp).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("unable to create signed token")
	}

	return &tokenString, nil
}

func getTokenFromString(tokenString string, secret string, claims *JwtClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
}
