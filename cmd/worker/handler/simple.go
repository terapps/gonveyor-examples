package handler

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor-examples/simple"
)

func HandleWelcome(_ context.Context, in simple.WelcomeInput) (simple.WelcomeOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("sending welcome email", "user_id", in.UserID, "email", in.Email)
	return simple.WelcomeOutput{SentAt: time.Now().Format(time.RFC3339)}, nil
}
