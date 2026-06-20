package handler

import (
	"context"
	"fmt"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func CreateContract(_ context.Context, in st.CreateContractInput) (st.CreateContractOutput, error) {
	slog.Info("creating contract",
		"quote_id", in.QuoteID,
		"signature_id", in.SignatureID,
		"txn_id", in.TxnID,
		"amount", in.Amount,
	)
	return st.CreateContractOutput{
		ContractID: fmt.Sprintf("ctr-%s", in.QuoteID),
	}, nil
}
