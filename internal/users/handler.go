package users

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ritchie-gr8/7solution-be/pkg/response"
)

type IUserHandler interface {
	GetUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userHandler struct {
	service IUserService
}

func NewUserHandler(service IUserService) IUserHandler {
	return &userHandler{service: service}
}

func (uh *userHandler) GetUsers(c *fiber.Ctx) error {
	users, err := uh.service.GetUsers(c)
	if err != nil {
		return response.NewResponse(c).Error(fiber.StatusInternalServerError, "", "An unexpected error occurred while retrieving users.").Response()
	}
	return response.NewResponse(c).Success(fiber.StatusOK, users).Response()
}

func (uh *userHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.NewResponse(c).Error(fiber.StatusBadRequest, id, "id is required").Response()
	}

	user, err := uh.service.GetUserById(c, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return response.NewResponse(c).Error(fiber.StatusNotFound, id, "The requested user was not found.").Response()
		default:
			return response.NewResponse(c).Error(fiber.StatusInternalServerError, id, "An unexpected error occurred while retrieving the user.").Response()
		}
	}
	return response.NewResponse(c).Success(fiber.StatusOK, user).Response()
}

func (uh *userHandler) CreateUser(c *fiber.Ctx) error {
	var userReq CreateUserRequest
	if err := c.BodyParser(&userReq); err != nil {
		return response.NewResponse(c).Error(fiber.StatusBadRequest, "", err.Error()).Response()
	}

	user, err := uh.service.CreateUser(c, userReq)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			return response.NewResponse(c).Error(fiber.StatusConflict, "", "A user with this email already exists.").Response()
		default:
			return response.NewResponse(c).Error(fiber.StatusInternalServerError, "", "An unexpected error occurred while creating the user.").Response()
		}
	}
	return response.NewResponse(c).Success(fiber.StatusOK, user).Response()
}

func (uh *userHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.NewResponse(c).Error(fiber.StatusBadRequest, id, "id is required").Response()
	}

	tokenUser := c.Locals("userId").(string)
	if tokenUser != id {
		return response.NewResponse(c).Error(fiber.StatusUnauthorized, id, "Unauthorized").Response()
	}

	var userReq UpdateUserRequest
	if err := c.BodyParser(&userReq); err != nil {
		return response.NewResponse(c).Error(fiber.StatusBadRequest, id, err.Error()).Response()
	}

	updatedUser, err := uh.service.UpdateUser(c, id, userReq)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return response.NewResponse(c).Error(fiber.StatusNotFound, id, "The user you are trying to update was not found.").Response()
		case errors.Is(err, ErrEmailAlreadyExists):
			return response.NewResponse(c).Error(fiber.StatusConflict, id, "Cannot update user: email already exists.").Response()
		case errors.Is(err, ErrUpdateFailed):
			return response.NewResponse(c).Error(fiber.StatusInternalServerError, id, "Failed to update user due to an internal error.").Response()
		default:
			return response.NewResponse(c).Error(fiber.StatusInternalServerError, id, "An unexpected error occurred while updating the user.").Response()
		}
	}
	return response.NewResponse(c).Success(fiber.StatusOK, updatedUser).Response()
}

func (uh *userHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.NewResponse(c).Error(fiber.StatusBadRequest, id, "id is required").Response()
	}

	tokenUser := c.Locals("userId").(string)
	if tokenUser != id {
		return response.NewResponse(c).Error(fiber.StatusUnauthorized, id, "Unauthorized").Response()
	}

	err := uh.service.DeleteUser(c, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return response.NewResponse(c).Error(fiber.StatusNotFound, id, "The user you are trying to delete was not found.").Response()
		case errors.Is(err, ErrDeleteFailed):
			return response.NewResponse(c).Error(fiber.StatusInternalServerError, id, "Failed to delete user due to an internal error.").Response()
		default:
			return response.NewResponse(c).Error(fiber.StatusInternalServerError, id, "An unexpected error occurred while deleting the user.").Response()
		}
	}
	return response.NewResponse(c).Success(fiber.StatusOK, fmt.Sprintf("User with id %s deleted successfully", id)).Response()
}

func (uh *userHandler) Login(c *fiber.Ctx) error {
	var loginReq LoginUserRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return response.NewResponse(c).Error(fiber.StatusBadRequest, "", err.Error()).Response()
	}

	user, err := uh.service.Login(c, loginReq)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return response.NewResponse(c).Error(fiber.StatusUnauthorized, "", "Invalid email or password provided.").Response()
		case errors.Is(err, ErrUserNotFound):
			return response.NewResponse(c).Error(fiber.StatusNotFound, "", "User not found.").Response()
		default:
			return response.NewResponse(c).Error(fiber.StatusInternalServerError, "", "An unexpected error occurred during login.").Response()
		}
	}
	return response.NewResponse(c).Success(fiber.StatusOK, user).Response()
}
