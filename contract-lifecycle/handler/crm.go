package handler

import (
	"context"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func SyncCrmQuote(_ context.Context, in st.SyncCrmQuoteInput) (st.SyncCrmQuoteOutput, error) {
	slog.Info("syncing quote to CRM", "quote_id", in.QuoteID, "process_id", in.ProcessID)
	return st.SyncCrmQuoteOutput{}, nil
}

func SyncCrmContract(_ context.Context, in st.SyncCrmContractInput) (st.SyncCrmContractOutput, error) {
	slog.Info("syncing contract to CRM", "contract_id", in.ContractID, "docs", len(in.DocURLs))
	return st.SyncCrmContractOutput{}, nil
}
