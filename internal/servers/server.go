package servers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ritchie-gr8/7solution-be/internal/auth"
	"github.com/ritchie-gr8/7solution-be/internal/config"
	"github.com/ritchie-gr8/7solution-be/internal/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type IServer interface {
	Start()
	GetServer() *server
}

type server struct {
	app    *fiber.App
	db     *mongo.Client
	cfg    config.IConfig
	cancel context.CancelFunc
}

func NewServer(cfg config.IConfig, db *mongo.Client) IServer {
	return &server{
		db:  db,
		cfg: cfg,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) Start() {
	s.app.Use(logger.New(logger.Config{
		Format:     "[${time}] | Status: ${status} | ${method} | Path: '${path}' | IP: ${ip} | Execution Time: ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Set up router groups
	v1 := s.app.Group("/v1")
	modules := InitModule(v1, s)
	modules.HealthModule()
	modules.UserModule()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	userRepo := users.NewUserRepository(s.db)
	userSvc := users.NewUserService(userRepo, auth.NewJWTAuthenticator(
		string(s.cfg.Jwt().SecretKey()),
		s.cfg.App().Name(),
		s.cfg.App().Name(),
		time.Duration(s.cfg.Jwt().AccessExpiresAt())*time.Second))
	StartUserCountMonitor(ctx, userSvc)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	done := make(chan bool, 1)

	go func() {
		log.Printf("server running on %s", s.cfg.App().Url())
		if err := s.app.Listen(s.cfg.App().Url()); err != nil {
			log.Printf("server error: %v", err)
		}
		done <- true
	}()

	select {
	case <-c:
		log.Println("shutting down server...")
		if s.cancel != nil {
			s.cancel()
		}
		s.app.Shutdown()
		<-done
	case <-done:
		log.Println("server stopped")
	}
}
