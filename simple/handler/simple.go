package handler

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/simple/stations"
)

func SendWelcome(_ context.Context, in st.WelcomeInput) (st.WelcomeOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("sending welcome email", "user_id", in.UserID, "email", in.Email)
	return st.WelcomeOutput{SentAt: time.Now().Format(time.RFC3339)}, nil
}
