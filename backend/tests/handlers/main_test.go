package handlers

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"gomor-e-commerce/internal/app"
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/tests/mocks"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database
var authMiddleware *auth.AuthMiddleware

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx, cancel := context.WithCancel(context.Background())

	DB = app.ConnectDB(ctx)
	if DB == nil {
		log.Fatal("Failed to connect to database in TestMain")
	}

	authClient := new(mocks.MockAuthClient)
	mockUserRepo := new(mocks.MockUsersRepository)
	// Mock Save to avoid panics when middleware tries to auto-create users from tokens
	mockUserRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	authMiddleware = auth.NewAuthMiddleware(authClient, mockUserRepo)
	adminToken := &auth.Token{
		UID: "user123",
		Claims: map[string]interface{}{
			"email":          "testadmin@example.com",
			"username":       "testadmin",
			"role":           "admin",
			"email_verified": true,
		},
	}
	userToken := &auth.Token{
		UID: "user123",
		Claims: map[string]interface{}{
			"email":          "testuser@example.com",
			"username":       "testuser",
			"role":           "user",
			"email_verified": true,
		},
	}

	authClient.On("VerifyIDToken", mock.Anything, "admin-token").Return(adminToken, nil).Maybe()
	authClient.On("VerifyIDToken", mock.Anything, "user-token").Return(userToken, nil).Maybe()

	code := m.Run()

	cancel()                          // This will trigger the disconnect goroutine in ConnectDB
	time.Sleep(10 * time.Millisecond) // Allow logs to print
	os.Exit(code)
}
