package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type DownloadInput struct {
	AssetID   string `validate:"required"`
	SourceURL string `validate:"required,url"`
}
type DownloadOutput struct {
	AssetID  string
	LocalURL string
}

var Download = gonveyor.Define[DownloadInput, DownloadOutput]("download")

// --- Worker ---

func HandleDownload(_ context.Context, in DownloadInput) (DownloadOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("downloading asset", "asset_id", in.AssetID, "source", in.SourceURL)
	return DownloadOutput{AssetID: in.AssetID, LocalURL: fmt.Sprintf("local://%s", in.AssetID)}, nil
}
