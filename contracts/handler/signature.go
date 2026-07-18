package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/contracts/stations"
)

func InitiateSignature(_ context.Context, in st.InitiateSignatureInput) (st.InitiateSignatureOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("initiating signature process", "quote_id", in.QuoteID, "docs", len(in.DocURLs))
	processID := fmt.Sprintf("sig-%s", in.QuoteID)
	return st.InitiateSignatureOutput{
		ProcessID:    processID,
		SignatureURL: fmt.Sprintf("https://sign.example.com/process/%s", processID),
	}, nil
}
