package handler

import (
	"context"
	"log/slog"
	"time"

	bp "github.com/terapps/gonveyor-examples/blueprint"
)

func SendWelcome(_ context.Context, in bp.WelcomeInput) (bp.WelcomeOutput, error) {
	slog.Info("sending welcome email", "user_id", in.UserID, "email", in.Email)
	return bp.WelcomeOutput{SentAt: time.Now().Format(time.RFC3339)}, nil
}
