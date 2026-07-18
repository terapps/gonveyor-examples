package transcoding

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type ThumbnailInput struct {
	AssetID  string
	LocalURL string
}
type ThumbnailOutput struct {
	ThumbURL string
}

var Thumbnail = gonveyor.Define[ThumbnailInput, ThumbnailOutput]("thumbnail")

// --- Worker ---

func HandleThumbnail(_ context.Context, in ThumbnailInput) (ThumbnailOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("generating thumbnail", "asset_id", in.AssetID)
	return ThumbnailOutput{ThumbURL: fmt.Sprintf("cdn://thumb/%s.jpg", in.AssetID)}, nil
}
