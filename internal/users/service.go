package users

import (
	"context"
	"errors" // Added for errors.Is

	"github.com/gofiber/fiber/v2"
	"github.com/ritchie-gr8/7solution-be/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUsers(c *fiber.Ctx) ([]*UserResponse, error)
	GetUserById(c *fiber.Ctx, id string) (*UserResponse, error)
	Login(c *fiber.Ctx, user LoginUserRequest) (*UserResponseWithToken, error)
	CreateUser(c *fiber.Ctx, user CreateUserRequest) (*UserResponseWithToken, error)
	UpdateUser(c *fiber.Ctx, id string, user UpdateUserRequest) (*UserResponseWithMessage, error)
	DeleteUser(c *fiber.Ctx, id string) error
	CountUsers(context context.Context) (int64, error)
}

type userService struct {
	repo IUserRepository
	jwt  auth.IAuthenticator
}

func NewUserService(repo IUserRepository, jwt auth.IAuthenticator) IUserService {
	return &userService{repo: repo, jwt: jwt}
}

func (s *userService) GetUsers(c *fiber.Ctx) ([]*UserResponse, error) {
	users, err := s.repo.GetUsers(c)
	if err != nil {
		return nil, err
	}
	var userResponses []*UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}
	return userResponses, nil
}

func (s *userService) GetUserById(c *fiber.Ctx, id string) (*UserResponse, error) {
	user, err := s.repo.GetUserById(c, id)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) CreateUser(c *fiber.Ctx, userReq CreateUserRequest) (*UserResponseWithToken, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userReq.Password = string(hashPassword)

	user, err := s.repo.CreateUser(c, userReq)
	if err != nil {
		return nil, err
	}

	claims := s.jwt.GenerateClaims(user.ID)

	token, err := s.jwt.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	return user.ToResponseWithToken(token), nil
}

func (s *userService) UpdateUser(c *fiber.Ctx, id string, userReq UpdateUserRequest) (*UserResponseWithMessage, error) {
	user, err := s.repo.UpdateUser(c, id, userReq)
	if err != nil {
		return nil, err
	}

	return user.ToResponseWithMessage("User updated successfully"), nil
}

func (s *userService) DeleteUser(c *fiber.Ctx, id string) error {
	return s.repo.DeleteUser(c, id)
}

func (s *userService) Login(c *fiber.Ctx, userReq LoginUserRequest) (*UserResponseWithToken, error) {
	user, err := s.repo.GetUserByEmail(c, userReq.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password)); err != nil {
		// If passwords don't match, return specific error
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		// For other bcrypt errors, return the original error
		return nil, err
	}

	claims := s.jwt.GenerateClaims(user.ID)

	token, err := s.jwt.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	return user.ToResponseWithToken(token), nil
}

func (s *userService) CountUsers(c context.Context) (int64, error) {
	return s.repo.CountUsers(c)
}
