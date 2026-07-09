package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func CreateContract(_ context.Context, in st.CreateContractInput) (st.CreateContractOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
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
