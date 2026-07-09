package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func InitiatePayment(_ context.Context, in st.InitiatePaymentInput) (st.InitiatePaymentOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("initiating payment process",
		"quote_id", in.QuoteID,
		"client_email", in.ClientEmail,
		"amount", in.Amount,
	)
	processID := fmt.Sprintf("pay-%s", in.QuoteID)
	return st.InitiatePaymentOutput{
		ProcessID:  processID,
		PaymentURL: fmt.Sprintf("https://pay.example.com/process/%s", processID),
	}, nil
}
