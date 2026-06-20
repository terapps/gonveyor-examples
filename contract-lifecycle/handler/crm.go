package handler

import (
	"context"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func SyncCrm(_ context.Context, in st.SyncCrmInput) (st.SyncCrmOutput, error) {
	slog.Info("syncing to CRM",
		"entity_type", in.EntityType,
		"entity_id", in.EntityID,
		"metadata", in.Metadata,
		"docs", len(in.DocURLs),
	)
	return st.SyncCrmOutput{}, nil
}
