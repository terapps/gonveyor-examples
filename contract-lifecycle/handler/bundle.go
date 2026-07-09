package handler

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func BundleContractDocs(_ context.Context, in st.BundleContractDocsInput) (st.BundleContractDocsOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("bundling contract docs", "contract_id", in.ContractID, "count", len(in.DocURLs))
	return st.BundleContractDocsOutput{
		ContractID:  in.ContractID,
		ClientEmail: in.ClientEmail,
		DocURLs:     in.DocURLs,
	}, nil
}
