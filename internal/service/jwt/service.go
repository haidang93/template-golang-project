package jwt

import (
	"errors"
	"time"

	"github.com/example/internal/i18n"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secretKey string
	issuer    string
}

type JWTClaims struct {
	UserID   *string `json:"id"`
	UserType *string `json:"type"`
	jwt.RegisteredClaims
}

func NewJwtService(secretKey string) *JwtService {
	return &JwtService{
		secretKey: secretKey,
		issuer:    "viewfinder",
	}
}

func (js *JwtService) GenerateToken(userID *string, userType *string) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    js.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(js.secretKey))
}

func (js *JwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(i18n.KEY_AUTHORIZATION_UNEXPECTED_SIGNING_METHOD)
		}
		return []byte(js.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New(i18n.KEY_AUTHORIZATION_INVALID_TOKEN)
	}

	return claims, nil
}
