package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type TranscodeInput struct {
	AssetID  string
	LocalURL string
}
type TranscodeOutput struct {
	VideoURL string
}

var Transcode = gonveyor.Define[TranscodeInput, TranscodeOutput]("transcode")

// --- Worker ---

func HandleTranscode(_ context.Context, in TranscodeInput) (TranscodeOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("transcoding video", "asset_id", in.AssetID)
	return TranscodeOutput{VideoURL: fmt.Sprintf("cdn://video/%s.mp4", in.AssetID)}, nil
}
