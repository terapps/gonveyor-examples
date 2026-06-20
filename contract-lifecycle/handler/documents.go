package handler

import (
	"context"
	"fmt"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func GenerateQuoteDoc(_ context.Context, in st.GenerateQuoteDocInput) (st.GenerateQuoteDocOutput, error) {
	slog.Info("generating quote document", "quote_id", in.QuoteID, "doc_type", in.DocType)
	return st.GenerateQuoteDocOutput{
		DocURL:  fmt.Sprintf("storage://quotes/%s/%s.pdf", in.QuoteID, in.DocType),
		DocType: in.DocType,
	}, nil
}

func GenerateContractDoc(_ context.Context, in st.GenerateContractDocInput) (st.GenerateContractDocOutput, error) {
	slog.Info("generating contract document", "contract_id", in.ContractID, "doc_type", in.DocType)
	return st.GenerateContractDocOutput{
		DocURL:  fmt.Sprintf("storage://contracts/%s/%s.pdf", in.ContractID, in.DocType),
		DocType: in.DocType,
	}, nil
}
