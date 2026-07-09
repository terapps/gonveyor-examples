package handler

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func SyncCrm(_ context.Context, in st.SyncCrmInput) (st.SyncCrmOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("syncing to CRM",
		"entity_type", in.EntityType,
		"entity_id", in.EntityID,
		"metadata", in.Metadata,
		"docs", len(in.DocURLs),
	)
	return st.SyncCrmOutput{}, nil
}
