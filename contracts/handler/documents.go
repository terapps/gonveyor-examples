package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	st "github.com/terapps/gonveyor-examples/contracts/stations"
)

func GenerateDocument(_ context.Context, in st.DocumentInput) (st.DocumentOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("generating document", "entity_id", in.EntityID, "doc_type", in.DocType)
	return st.DocumentOutput{
		DocURL:   fmt.Sprintf("storage://%s/%s.pdf", in.EntityID, in.DocType),
		DocType:  in.DocType,
		EntityID: in.EntityID,
	}, nil
}
