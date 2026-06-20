package handler

import (
	"context"
	"log/slog"
	"time"

	st "github.com/terapps/gonveyor-examples/simple/stations"
)

func SendWelcome(_ context.Context, in st.WelcomeInput) (st.WelcomeOutput, error) {
	slog.Info("sending welcome email", "user_id", in.UserID, "email", in.Email)
	return st.WelcomeOutput{SentAt: time.Now().Format(time.RFC3339)}, nil
}
