package handler

import (
	"context"
	"fmt"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/contract-lifecycle/stations"
)

func GenerateDocument(_ context.Context, in st.DocumentInput) (st.DocumentOutput, error) {
	slog.Info("generating document", "entity_id", in.EntityID, "doc_type", in.DocType)
	return st.DocumentOutput{
		DocURL:  fmt.Sprintf("storage://%s/%s.pdf", in.EntityID, in.DocType),
		DocType: in.DocType,
	}, nil
}
