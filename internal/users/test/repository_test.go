package test

import (
	"context"
	"errors" // Added for errors.Is
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ritchie-gr8/7solution-be/internal/users"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct {
	findOneFunc          func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult
	findFunc             func(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error)
	insertOneFunc        func(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	findOneAndUpdateFunc func(ctx context.Context, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	deleteOneFunc        func(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	countDocumentsFunc   func(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)
}

func (m *MockCollection) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
	if m.findOneFunc != nil {
		return m.findOneFunc(ctx, filter, opts...)
	}
	return nil
}

func (m *MockCollection) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	if m.findFunc != nil {
		return m.findFunc(ctx, filter, opts...)
	}
	return nil, nil
}

func (m *MockCollection) InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if m.insertOneFunc != nil {
		return m.insertOneFunc(ctx, document, opts...)
	}
	return nil, nil
}

func (m *MockCollection) FindOneAndUpdate(ctx context.Context, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	if m.findOneAndUpdateFunc != nil {
		return m.findOneAndUpdateFunc(ctx, filter, update, opts...)
	}
	return nil
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if m.deleteOneFunc != nil {
		return m.deleteOneFunc(ctx, filter, opts...)
	}
	return nil, nil
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	if m.countDocumentsFunc != nil {
		return m.countDocumentsFunc(ctx, filter, opts...)
	}
	return 0, nil
}

func createFiberCtx() *fiber.Ctx {
	app := fiber.New()
	return app.AcquireCtx(&fasthttp.RequestCtx{})
}

func TestGetUsers(t *testing.T) {
	t.Run("Successfully get all users", func(t *testing.T) {
		user1 := users.User{
			ID:        primitive.NewObjectID(),
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "password1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		user2 := users.User{
			ID:        primitive.NewObjectID(),
			Name:      "User 2",
			Email:     "user2@example.com",
			Password:  "password2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		usersList := []users.User{user1, user2}

		mockCurrsor, err := mongo.NewCursorFromDocuments(bson.A{user1, user2}, nil, nil)
		if err != nil {
			t.Fatalf("Failed to create mock cursor: %v", err)
		}

		mockColl := &MockCollection{
			findFunc: func(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error) {
				filterMap, ok := filter.(bson.D)
				if !ok {
					t.Fatalf("Expected filter to be bson.D, got %T", filter)
				}

				if len(filterMap) != 0 {
					t.Fatalf("Expected empty filter for GetUsers, got: %v", filterMap)
				}

				return mockCurrsor, nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		result, err := repo.GetUsers(ctx)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != len(usersList) {
			t.Fatalf("Expected %d users, got %d", len(usersList), len(result))
		}

		for i, user := range result {
			if user.Email != usersList[i].Email {
				t.Errorf("Expected user email %s, got %s", usersList[i].Email, user.Email)
			}
			if user.Name != usersList[i].Name {
				t.Errorf("Expected user name %s, got %s", usersList[i].Name, user.Name)
			}
		}
	})

	t.Run("Error get all users", func(t *testing.T) {
		mockColl := &MockCollection{
			findFunc: func(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error) {
				return nil, mongo.CommandError{Message: "Database error", Code: 123}
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		result, err := repo.GetUsers(ctx)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestGetUserById(t *testing.T) {
	t.Run("User found", func(t *testing.T) {
		id := "507f1f77bcf86cd799439011"
		objectID, _ := primitive.ObjectIDFromHex(id)
		expectedUser := users.User{
			ID:        objectID,
			Name:      "Test User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				filterMap, ok := filter.(bson.M)
				if !ok {
					t.Fatalf("Expected filter to be bson.M, got %T", filter)
				}

				filterID, ok := filterMap["_id"]
				if !ok || filterID != objectID {
					t.Fatalf("Expected filter to have _id: %v, got: %v", objectID, filterMap)
				}

				return mongo.NewSingleResultFromDocument(expectedUser, nil, nil)
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()
		user, err := repo.GetUserById(ctx, id)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.ID != objectID {
			t.Errorf("Expected user ID %v, got %v", objectID, user.ID)
		}

		if user.Name != expectedUser.Name {
			t.Errorf("Expected user name %s, got %s", expectedUser.Name, user.Name)
		}

		if user.Email != expectedUser.Email {
			t.Errorf("Expected user email %s, got %s", expectedUser.Email, user.Email)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		id := "507f1f77bcf86cd799439012"
		objectID, _ := primitive.ObjectIDFromHex(id)

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {

				filterMap, ok := filter.(bson.M)
				if !ok {
					t.Fatalf("Expected filter to be bson.M, got %T", filter)
				}

				filterID, ok := filterMap["_id"]
				if !ok || filterID != objectID {
					t.Fatalf("Expected filter to have _id: %v, got: %v", objectID, filterMap)
				}

				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()
		user, err := repo.GetUserById(ctx, id)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrUserNotFound) {
			t.Errorf("Expected users.ErrUserNotFound, got: %v", err)
		}

		if user != nil {
			t.Errorf("Expected nil user, got: %v", user)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		id := "invalid-id"

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				t.Fatal("FindOne should not be called with invalid ID")
				return nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()
		user, err := repo.GetUserById(ctx, id)

		if err == nil {
			t.Error("Expected error for invalid ID, got nil")
		}

		if err == mongo.ErrNoDocuments {
			t.Error("Expected ObjectID parsing error, not ErrNoDocuments")
		}

		if user != nil {
			t.Errorf("Expected nil user for invalid ID, got: %v", user)
		}
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("User found by email", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := users.User{
			ID:        primitive.NewObjectID(),
			Name:      "Test User",
			Email:     email,
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				filterMap, ok := filter.(bson.M)
				if !ok {
					t.Fatalf("Expected filter to be bson.M, got %T", filter)
				}

				filterEmail, ok := filterMap["email"]
				if !ok || filterEmail != email {
					t.Fatalf("Expected filter to have email: %v, got: %v", email, filterMap)
				}

				return mongo.NewSingleResultFromDocument(expectedUser, nil, nil)
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		user, err := repo.GetUserByEmail(ctx, email)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.Email != email {
			t.Errorf("Expected user email %s, got %s", email, user.Email)
		}
	})

	t.Run("User not found by email", func(t *testing.T) {
		email := "notfound@example.com"

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		user, err := repo.GetUserByEmail(ctx, email)

		if !errors.Is(err, users.ErrUserNotFound) {
			t.Errorf("Expected users.ErrUserNotFound, got: %v", err)
		}

		if user != nil {
			t.Errorf("Expected nil user, got: %v", user)
		}
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Successful creation", func(t *testing.T) {
		insertedID := primitive.NewObjectID()

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
			insertOneFunc: func(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
				user, ok := document.(users.User)
				if !ok {
					t.Fatalf("Expected document to be users.User, got %T", document)
				}

				if user.Name == "" {
					t.Error("Expected user name to be set")
				}
				if user.Email == "" {
					t.Error("Expected user email to be set")
				}
				if user.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set")
				}
				if user.UpdatedAt.IsZero() {
					t.Error("Expected UpdatedAt to be set")
				}

				return &mongo.InsertOneResult{InsertedID: insertedID}, nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		createReq := users.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "hashedpassword",
		}

		user, err := repo.CreateUser(ctx, createReq)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.Name != createReq.Name {
			t.Errorf("Expected name %s, got %s", createReq.Name, user.Name)
		}

		if user.Email != createReq.Email {
			t.Errorf("Expected email %s, got %s", createReq.Email, user.Email)
		}

		if user.ID != insertedID {
			t.Errorf("Expected ID %v, got %v", insertedID, user.ID)
		}
	})

	t.Run("Insert failure", func(t *testing.T) {
		expectedError := users.ErrInsertFailed

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
			insertOneFunc: func(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
				return nil, expectedError
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		createReq := users.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "hashedpassword",
		}

		user, err := repo.CreateUser(ctx, createReq)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrInsertFailed) {
			t.Errorf("Expected users.ErrInsertFailed, got: %v", err)
		}

		if user != nil {
			t.Errorf("Expected nil user, got: %v", user)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("User updated successfully", func(t *testing.T) {
		userID := primitive.NewObjectID()
		updatedUser := users.User{
			ID:        userID,
			Name:      "Updated Name",
			Email:     "updated@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
			findOneAndUpdateFunc: func(ctx context.Context, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
				filterDoc, ok := filter.(bson.M)
				if !ok {
					t.Fatalf("Expected filter to be bson.M, got %T", filter)
				}

				filterID, ok := filterDoc["_id"].(primitive.ObjectID)
				if !ok {
					t.Fatalf("Expected filter to contain _id as ObjectID, got %T", filterDoc["_id"])
				}

				if filterID != userID {
					t.Errorf("Expected filter ID %v, got %v", userID, filterID)
				}

				updateDoc, ok := update.(bson.M)
				if !ok {
					t.Fatalf("Expected update to be bson.M, got %T", update)
				}

				setDoc, ok := updateDoc["$set"].(bson.M)
				if !ok {
					t.Fatalf("Expected update to contain $set operation, got %T", updateDoc["$set"])
				}

				if setDoc["name"] != "Updated Name" {
					t.Errorf("Expected name to be Updated Name, got %v", setDoc["name"])
				}

				if setDoc["email"] != "updated@example.com" {
					t.Errorf("Expected email to be updated@example.com, got %v", setDoc["email"])
				}

				_, hasUpdatedAt := setDoc["updatedAt"]
				if !hasUpdatedAt {
					t.Error("Expected updatedAt to be set")
				}

				return mongo.NewSingleResultFromDocument(updatedUser, nil, nil)
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		updateReq := users.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		updatedUserResult, err := repo.UpdateUser(ctx, userID.Hex(), updateReq)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if updatedUserResult == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if updatedUserResult.ID != userID {
			t.Errorf("Expected ID %v, got %v", userID, updatedUserResult.ID)
		}

		if updatedUserResult.Name != updateReq.Name {
			t.Errorf("Expected name %s, got %s", updateReq.Name, updatedUserResult.Name)
		}

		if updatedUserResult.Email != updateReq.Email {
			t.Errorf("Expected email %s, got %s", updateReq.Email, updatedUserResult.Email)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		userID := primitive.NewObjectID()

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
			findOneAndUpdateFunc: func(ctx context.Context, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		updateReq := users.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		updatedUser, err := repo.UpdateUser(ctx, userID.Hex(), updateReq)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrUserNotFound) {
			t.Errorf("Expected users.ErrUserNotFound, got: %v", err)
		}

		if updatedUser != nil {
			t.Errorf("Expected nil user, got: %v", updatedUser)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		invalidID := "invalid-id"

		mockColl := &MockCollection{
			findOneAndUpdateFunc: func(ctx context.Context, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
				t.Fatal("FindOneAndUpdate should not be called with invalid ID")
				return nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		updateReq := users.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		updatedUser, err := repo.UpdateUser(ctx, invalidID, updateReq)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrInvalidID) {
			t.Errorf("Expected users.ErrInvalidID, got: %v", err)
		}

		if updatedUser != nil {
			t.Errorf("Expected nil user, got: %v", updatedUser)
		}
	})

	t.Run("Update failed", func(t *testing.T) {
		userID := primitive.NewObjectID()
		updatedUser := users.User{
			ID:        userID,
			Name:      "Updated Name",
			Email:     "updated@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockColl := &MockCollection{
			findOneFunc: func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
			},
			findOneAndUpdateFunc: func(ctx context.Context, filter any, update any, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(updatedUser, users.ErrUpdateFailed, nil)
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		updateReq := users.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		updatedUserResult, err := repo.UpdateUser(ctx, userID.Hex(), updateReq)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrUpdateFailed) {
			t.Errorf("Expected users.ErrUpdateFailed, got: %v", err)
		}

		if updatedUserResult != nil {
			t.Errorf("Expected nil user, got: %v", updatedUserResult)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Delete user successfully", func(t *testing.T) {
		userID := primitive.NewObjectID()

		mockColl := &MockCollection{
			deleteOneFunc: func(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
				filterDoc, ok := filter.(bson.M)
				if !ok {
					t.Fatalf("Expected filter to be bson.M, got %T", filter)
				}

				filterID, ok := filterDoc["_id"].(primitive.ObjectID)
				if !ok {
					t.Fatalf("Expected filter to contain _id as ObjectID, got %T", filterDoc["_id"])
				}

				if filterID != userID {
					t.Errorf("Expected filter ID %v, got %v", userID, filterID)
				}

				return &mongo.DeleteResult{DeletedCount: 1}, nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		err := repo.DeleteUser(ctx, userID.Hex())

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		userID := primitive.NewObjectID()

		mockColl := &MockCollection{
			deleteOneFunc: func(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
				return &mongo.DeleteResult{DeletedCount: 0}, nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		err := repo.DeleteUser(ctx, userID.Hex())

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrUserNotFound) {
			t.Errorf("Expected users.ErrUserNotFound, got: %v", err)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		invalidID := "invalid-id"

		mockColl := &MockCollection{
			deleteOneFunc: func(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
				t.Fatal("DeleteOne should not be called with invalid ID")
				return nil, nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		err := repo.DeleteUser(ctx, invalidID)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrInvalidID) {
			t.Errorf("Expected users.ErrInvalidID, got: %v", err)
		}
	})

	t.Run("Delete failed", func(t *testing.T) {
		userID := primitive.NewObjectID()
		expectedError := users.ErrDeleteFailed

		mockColl := &MockCollection{
			deleteOneFunc: func(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
				return nil, expectedError
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)
		ctx := createFiberCtx()

		err := repo.DeleteUser(ctx, userID.Hex())

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !errors.Is(err, users.ErrDeleteFailed) {
			t.Errorf("Expected users.ErrDeleteFailed, got: %v", err)
		}
	})
}

func TestCountUsers(t *testing.T) {
	t.Run("Count users successfully", func(t *testing.T) {
		mockColl := &MockCollection{
			countDocumentsFunc: func(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
				return 10, nil
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)

		count, err := repo.CountUsers(context.Background())

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if count != 10 {
			t.Errorf("Expected count 10, got %d", count)
		}
	})

	t.Run("Count failed", func(t *testing.T) {
		mockColl := &MockCollection{
			countDocumentsFunc: func(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
				return 0, mongo.CommandError{Message: "Database error", Code: 123}
			},
		}

		repo := users.NewUserRepositoryWithCollection(mockColl)

		count, err := repo.CountUsers(context.Background())

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if count != 0 {
			t.Errorf("Expected count 0, got %d", count)
		}
	})
}
