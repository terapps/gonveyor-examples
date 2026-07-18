package contracts

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type SyncCrmInput struct {
	EntityType string // "quote" or "contract"
	EntityID   string
	Metadata   map[string]string
	DocURLs    []string
}
type SyncCrmOutput struct{}

var SyncCrmQuote = gonveyor.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_quote")
var SyncCrmContract = gonveyor.Define[SyncCrmInput, SyncCrmOutput]("sync_crm_contract")

// --- Worker ---

func HandleCrm(_ context.Context, in SyncCrmInput) (SyncCrmOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("syncing to CRM", "entity_type", in.EntityType, "entity_id", in.EntityID, "metadata", in.Metadata, "docs", len(in.DocURLs))
	return SyncCrmOutput{}, nil
}
