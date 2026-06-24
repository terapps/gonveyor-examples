package handler

import (
	"context"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func BundleContractDocs(_ context.Context, in st.BundleContractDocsInput) (st.BundleContractDocsOutput, error) {
	slog.Info("bundling contract docs", "contract_id", in.ContractID, "count", len(in.DocURLs))
	return st.BundleContractDocsOutput{
		ContractID:  in.ContractID,
		ClientEmail: in.ClientEmail,
		DocURLs:     in.DocURLs,
	}, nil
}
