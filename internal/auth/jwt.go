package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTAuthenticator struct {
	secret    string
	audience  string
	issuer    string
	expiresAt time.Duration
}

func NewJWTAuthenticator(secret, audience, issuer string, expiresAt time.Duration) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:    secret,
		audience:  audience,
		issuer:    issuer,
		expiresAt: expiresAt,
	}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.audience),
		jwt.WithIssuer(a.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
}

func (a *JWTAuthenticator) GenerateClaims(id primitive.ObjectID) jwt.MapClaims {
	return jwt.MapClaims{
		"sub": id.Hex(),
		"exp": time.Now().Add(a.expiresAt).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": a.issuer,
		"aud": a.audience,
	}
}
