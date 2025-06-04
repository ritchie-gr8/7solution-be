package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchie-gr8/7solution-be/internal/auth"
	"github.com/ritchie-gr8/7solution-be/pkg/response"
)

func ValidateToken(jwtAuth auth.IAuthenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.NewResponse(c).Error(fiber.StatusUnauthorized, "", "Unauthorized").Response()
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.NewResponse(c).Error(fiber.StatusUnauthorized, "", "Invalid token format").Response()
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return response.NewResponse(c).Error(fiber.StatusUnauthorized, "", "Missing token").Response()
		}

		jwtToken, err := jwtAuth.ValidateToken(token)
		if err != nil {
			return response.NewResponse(c).Error(fiber.StatusUnauthorized, "", "Invalid token").Response()
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return response.NewResponse(c).Error(fiber.StatusUnauthorized, "", "Invalid token claims").Response()
		}

		userId := claims["sub"].(string)
		c.Locals("userId", userId)

		return c.Next()
	}
}
