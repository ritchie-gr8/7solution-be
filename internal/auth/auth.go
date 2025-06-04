package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticator interface {
	GenerateToken(claims jwt.MapClaims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GenerateClaims(id primitive.ObjectID) jwt.MapClaims
}
