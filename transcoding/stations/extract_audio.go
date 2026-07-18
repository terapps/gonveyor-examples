package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type ExtractAudioInput struct {
	AssetID  string
	LocalURL string
}
type ExtractAudioOutput struct {
	AudioURL string
}

var ExtractAudio = gonveyor.Define[ExtractAudioInput, ExtractAudioOutput]("extract_audio")

// --- Worker ---

func HandleExtractAudio(_ context.Context, in ExtractAudioInput) (ExtractAudioOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("extracting audio", "asset_id", in.AssetID)
	return ExtractAudioOutput{AudioURL: fmt.Sprintf("cdn://audio/%s.mp3", in.AssetID)}, nil
}
