package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type InitiateSignatureInput struct {
	QuoteID     string
	ClientEmail string
	DocURLs     []string
}
type InitiateSignatureOutput struct {
	ProcessID    string
	SignatureURL string
}

var InitiateSignature = gonveyor.Define[InitiateSignatureInput, InitiateSignatureOutput]("initiate_signature")

// --- Worker ---

func HandleInitiateSignature(_ context.Context, in InitiateSignatureInput) (InitiateSignatureOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("initiating signature process", "quote_id", in.QuoteID, "docs", len(in.DocURLs))
	processID := fmt.Sprintf("sig-%s", in.QuoteID)
	return InitiateSignatureOutput{ProcessID: processID, SignatureURL: fmt.Sprintf("https://sign.example.com/process/%s", processID)}, nil
}
