package handler

import (
	"context"
	"fmt"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func InitiateSignature(_ context.Context, in st.InitiateSignatureInput) (st.InitiateSignatureOutput, error) {
	slog.Info("initiating signature process", "quote_id", in.QuoteID, "docs", len(in.DocURLs))
	processID := fmt.Sprintf("sig-%s", in.QuoteID)
	return st.InitiateSignatureOutput{
		ProcessID:   processID,
		SignatureURL: fmt.Sprintf("https://sign.example.com/process/%s", processID),
	}, nil
}
