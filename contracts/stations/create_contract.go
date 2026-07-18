package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type CreateContractInput struct {
	QuoteID     string
	ClientEmail string
	SignatureID string
	TxnID       string
	Amount      float64
}
type CreateContractOutput struct {
	ContractID  string
	ClientEmail string
}

var CreateContract = gonveyor.Define[CreateContractInput, CreateContractOutput]("create_contract")

// --- Worker ---

func HandleCreateContract(_ context.Context, in CreateContractInput) (CreateContractOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("creating contract", "quote_id", in.QuoteID, "signature_id", in.SignatureID, "txn_id", in.TxnID, "amount", in.Amount)
	return CreateContractOutput{ContractID: fmt.Sprintf("ctr-%s", in.QuoteID)}, nil
}
