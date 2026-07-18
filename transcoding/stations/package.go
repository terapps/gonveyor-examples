package stations

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/terapps/gonveyor"
)

type PackageInput struct {
	AssetID  string
	VideoURL string
	ThumbURL string
	AudioURL string
}
type PackageOutput struct {
	BundleURL string
}

var Package = gonveyor.Define[PackageInput, PackageOutput]("package")

// --- Worker ---

func HandlePackage(_ context.Context, in PackageInput) (PackageOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("packaging bundle", "asset_id", in.AssetID, "video", in.VideoURL, "thumb", in.ThumbURL, "audio", in.AudioURL)
	return PackageOutput{BundleURL: fmt.Sprintf("cdn://bundle/%s.zip", in.AssetID)}, nil
}
