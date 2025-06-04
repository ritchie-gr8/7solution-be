package servers

import (
	"context"
	"log"
	"time"

	"github.com/ritchie-gr8/7solution-be/internal/users"
)

func StartUserCountMonitor(ctx context.Context, userService users.IUserService) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				countAndLogUsers(userService)
			}
		}
	}()

	log.Println("User count monitor started")
}

func countAndLogUsers(userService users.IUserService) {
	ctx := context.Background()
	count, err := userService.CountUsers(ctx)
	if err != nil {
		log.Printf("Failed to count users: %v", err)
	}
	log.Printf("Current user count: %d", count)
}
