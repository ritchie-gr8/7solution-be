package servers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ritchie-gr8/7solution-be/internal/auth"
	"github.com/ritchie-gr8/7solution-be/internal/health"
	"github.com/ritchie-gr8/7solution-be/internal/middleware"
	"github.com/ritchie-gr8/7solution-be/internal/users"
)

type IModuleFactory interface {
	HealthModule()
	UserModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
}

func InitModule(r fiber.Router, server *server) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: server,
	}
}

func (m *moduleFactory) HealthModule() {
	healthHandler := health.NewMonitorHandler(m.server.cfg)
	m.router.Get("/health", healthHandler.HealthCheck)
}

func (m *moduleFactory) UserModule() {
	jwtAuth := auth.NewJWTAuthenticator(
		string(m.server.cfg.Jwt().SecretKey()),
		m.server.cfg.App().Name(),
		m.server.cfg.App().Name(),
		time.Duration(m.server.cfg.Jwt().AccessExpiresAt())*time.Second)
	userRepo := users.NewUserRepository(m.server.db)
	userSvc := users.NewUserService(userRepo, jwtAuth)
	userHandler := users.NewUserHandler(userSvc)

	userGroup := m.router.Group("/users")
	userGroup.Get("", userHandler.GetUsers)
	userGroup.Get("/:id", userHandler.GetUserById)
	userGroup.Post("", middleware.ValidateRequest(&users.CreateUserRequest{}), userHandler.CreateUser)
	userGroup.Put("/:id", middleware.ValidateToken(jwtAuth), middleware.ValidateRequest(&users.UpdateUserRequest{}), userHandler.UpdateUser)
	userGroup.Delete("/:id", middleware.ValidateToken(jwtAuth), userHandler.DeleteUser)
	userGroup.Post("/login", middleware.ValidateRequest(&users.LoginUserRequest{}), userHandler.Login)
}
