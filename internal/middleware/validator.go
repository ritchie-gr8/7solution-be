package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ritchie-gr8/7solution-be/pkg/response"
)

var validate = validator.New()

func ValidateRequest(model any) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(model); err != nil {
			return response.NewResponse(c).Error(fiber.StatusBadRequest, "", "Invalid request body").Response()
		}

		if err := validate.Struct(model); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errorMsg := formatValidationErrors(validationErrors)
			return response.NewResponse(c).Error(fiber.StatusBadRequest, "", errorMsg).Response()
		}

		return c.Next()
	}
}

func formatValidationErrors(errs validator.ValidationErrors) string {
	if len(errs) > 0 {
		err := errs[0]
		return "Validation failed on field '" + err.Field() + "', condition: " + err.Tag()
	}
	return "Validation failed"
}
