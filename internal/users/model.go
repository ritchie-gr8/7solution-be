package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password,omitempty" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID    primitive.ObjectID `json:"id,omitempty"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}

type UserResponseWithMessage struct {
	UserResponse
	Message string `json:"message"`
}

type UserResponseWithToken struct {
	UserResponse
	Token string `json:"token"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (u *User) ToResponseWithMessage(message string) *UserResponseWithMessage {
	return &UserResponseWithMessage{
		UserResponse: UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		},
		Message: message,
	}
}

func (u *User) ToResponseWithToken(token string) *UserResponseWithToken {
	return &UserResponseWithToken{
		UserResponse: UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		},
		Token: token,
	}
}
