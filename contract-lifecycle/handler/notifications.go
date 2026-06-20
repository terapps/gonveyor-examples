package handler

import (
	"context"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func SendQuoteEmail(_ context.Context, in st.SendQuoteEmailInput) (st.SendQuoteEmailOutput, error) {
	slog.Info("sending quote email", "to", in.ClientEmail, "signature_url", in.SignatureURL)
	return st.SendQuoteEmailOutput{}, nil
}

func SendContractEmail(_ context.Context, in st.SendContractEmailInput) (st.SendContractEmailOutput, error) {
	slog.Info("sending contract email", "to", in.ClientEmail, "docs", len(in.DocURLs))
	return st.SendContractEmailOutput{}, nil
}
