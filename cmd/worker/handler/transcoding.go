package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	tst "github.com/terapps/gonveyor-examples/transcoding/stations"
)

func HandleDownload(_ context.Context, in tst.DownloadInput) (tst.DownloadOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("downloading asset", "asset_id", in.AssetID, "source", in.SourceURL)
	return tst.DownloadOutput{AssetID: in.AssetID, LocalURL: fmt.Sprintf("local://%s", in.AssetID)}, nil
}

func HandleTranscode(_ context.Context, in tst.TranscodeInput) (tst.TranscodeOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("transcoding video", "asset_id", in.AssetID)
	return tst.TranscodeOutput{VideoURL: fmt.Sprintf("cdn://video/%s.mp4", in.AssetID)}, nil
}

func HandleThumbnail(_ context.Context, in tst.ThumbnailInput) (tst.ThumbnailOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("generating thumbnail", "asset_id", in.AssetID)
	return tst.ThumbnailOutput{ThumbURL: fmt.Sprintf("cdn://thumb/%s.jpg", in.AssetID)}, nil
}

func HandleExtractAudio(_ context.Context, in tst.ExtractAudioInput) (tst.ExtractAudioOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("extracting audio", "asset_id", in.AssetID)
	return tst.ExtractAudioOutput{AudioURL: fmt.Sprintf("cdn://audio/%s.mp3", in.AssetID)}, nil
}

func HandlePackage(_ context.Context, in tst.PackageInput) (tst.PackageOutput, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	slog.Info("packaging bundle", "asset_id", in.AssetID, "video", in.VideoURL, "thumb", in.ThumbURL, "audio", in.AudioURL)
	return tst.PackageOutput{BundleURL: fmt.Sprintf("cdn://bundle/%s.zip", in.AssetID)}, nil
}
