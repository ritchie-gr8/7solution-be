package main

import (
	"os"

	"github.com/ritchie-gr8/7solution-be/internal/config"
	databases "github.com/ritchie-gr8/7solution-be/internal/database"
	"github.com/ritchie-gr8/7solution-be/internal/servers"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	}

	return os.Args[1]
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.DB())
	defer databases.DbDisconnect(db)

	server := servers.NewServer(cfg, db)
	server.Start()
}
