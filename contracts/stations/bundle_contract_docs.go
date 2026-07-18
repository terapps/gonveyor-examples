package stations

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type BundleContractDocsInput struct {
	ContractID  string
	ClientEmail string
	DocURLs     []string
}
type BundleContractDocsOutput struct {
	ContractID  string
	ClientEmail string
	DocURLs     []string
}

var BundleContractDocs = gonveyor.Define[BundleContractDocsInput, BundleContractDocsOutput]("bundle_contract_docs")

// --- Worker ---

func HandleBundleContractDocs(_ context.Context, in BundleContractDocsInput) (BundleContractDocsOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("bundling contract docs", "contract_id", in.ContractID, "count", len(in.DocURLs))
	return BundleContractDocsOutput{ContractID: in.ContractID, ClientEmail: in.ClientEmail, DocURLs: in.DocURLs}, nil
}
