package contracts

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type DocumentInput struct {
	EntityID string // QuoteID or ContractID
	DocType  string // e.g. "proposal", "pricing", "contract", "annex_a"
}
type DocumentOutput struct {
	DocURL   string
	DocType  string
	EntityID string
}

// GenerateQuoteDoc and GenerateContractDoc share DocumentInput/Output and the same handler.
// They are distinct stations because a node key must be unique within a blueprint.
var GenerateQuoteDoc = gonveyor.Define[DocumentInput, DocumentOutput]("generate_quote_doc", gonveyor.WithRoute("tasks.document"))
var GenerateContractDoc = gonveyor.Define[DocumentInput, DocumentOutput]("generate_contract_doc", gonveyor.WithRoute("tasks.document"))

// --- Worker ---

func HandleDocument(_ context.Context, in DocumentInput) (DocumentOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("generating document", "entity_id", in.EntityID, "doc_type", in.DocType)
	return DocumentOutput{
		DocURL:   fmt.Sprintf("storage://%s/%s.pdf", in.EntityID, in.DocType),
		DocType:  in.DocType,
		EntityID: in.EntityID,
	}, nil
}
