package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type InitiatePaymentInput struct {
	QuoteID     string
	ClientEmail string
	Amount      float64
	DocURLs     []string
}
type InitiatePaymentOutput struct {
	ProcessID  string
	PaymentURL string
}

var InitiatePayment = gonveyor.Define[InitiatePaymentInput, InitiatePaymentOutput]("initiate_payment")

// --- Worker ---

func HandleInitiatePayment(_ context.Context, in InitiatePaymentInput) (InitiatePaymentOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("initiating payment process", "quote_id", in.QuoteID, "client_email", in.ClientEmail, "amount", in.Amount)
	processID := fmt.Sprintf("pay-%s", in.QuoteID)
	return InitiatePaymentOutput{ProcessID: processID, PaymentURL: fmt.Sprintf("https://pay.example.com/process/%s", processID)}, nil
}
