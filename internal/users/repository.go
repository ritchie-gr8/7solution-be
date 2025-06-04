package users

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollection interface {
	Find(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOneAndUpdate(ctx context.Context, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)
}

type IUserRepository interface {
	GetUsers(c *fiber.Ctx) ([]User, error)
	GetUserById(c *fiber.Ctx, id string) (*User, error)
	GetUserByEmail(c *fiber.Ctx, email string) (*User, error)
	CreateUser(c *fiber.Ctx, user CreateUserRequest) (*User, error)
	UpdateUser(c *fiber.Ctx, id string, user UpdateUserRequest) (*User, error)
	DeleteUser(c *fiber.Ctx, id string) error
	CountUsers(ctx context.Context) (int64, error)
}

type userRepository struct {
	collection MongoCollection
}

func NewUserRepository(db *mongo.Client) IUserRepository {
	database := db.Database("userdb")
	collection := database.Collection("users")
	return &userRepository{collection: collection}
}

func NewUserRepositoryWithCollection(collection MongoCollection) IUserRepository {
	return &userRepository{collection: collection}
}

func (r *userRepository) GetUsers(c *fiber.Ctx) ([]User, error) {
	cursor, err := r.collection.Find(c.Context(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c.Context())

	var users []User
	if err := cursor.All(c.Context(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserById(c *fiber.Ctx, id string) (*User, error) {
	var user User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	err = r.collection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(c *fiber.Ctx, email string) (*User, error) {
	var user User
	err := r.collection.FindOne(c.Context(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(c *fiber.Ctx, userReq CreateUserRequest) (*User, error) {
	if err := r.checkEmailUniqueness(c.Context(), userReq.Email); err != nil {
		return nil, err
	}

	user := User{
		Name:      userReq.Name,
		Email:     userReq.Email,
		Password:  userReq.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.collection.InsertOne(c.Context(), user)
	if err != nil {
		return nil, ErrInsertFailed
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (r *userRepository) UpdateUser(c *fiber.Ctx, id string, userReq UpdateUserRequest) (*User, error) {
	var user User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	if err := r.checkEmailUniqueness(c.Context(), userReq.Email, objectID); err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = r.collection.FindOneAndUpdate(
		c.Context(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"name":      userReq.Name,
			"email":     userReq.Email,
			"updatedAt": time.Now(),
		}},
		opts,
	).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}

		return nil, ErrUpdateFailed
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(c *fiber.Ctx, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	result, err := r.collection.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		return ErrDeleteFailed
	}

	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) checkEmailUniqueness(ctx context.Context, email string, excludeID ...primitive.ObjectID) error {
	filter := bson.M{"email": email}

	if len(excludeID) > 0 && !excludeID[0].IsZero() {
		filter = bson.M{
			"email": email,
			"_id":   bson.M{"$ne": excludeID[0]},
		}
	}

	var existingUser User
	err := r.collection.FindOne(ctx, filter).Decode(&existingUser)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	if existingUser.Email != "" {
		return ErrEmailAlreadyExists
	}

	return nil
}
