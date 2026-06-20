package handler

import (
	"context"
	"fmt"
	"log/slog"

	bp "github.com/terapps/gonveyor-examples/blueprint"
)

func Download(_ context.Context, in bp.DownloadInput) (bp.DownloadOutput, error) {
	slog.Info("downloading asset", "asset_id", in.AssetID, "source", in.SourceURL)
	return bp.DownloadOutput{AssetID: in.AssetID, LocalURL: fmt.Sprintf("local://%s", in.AssetID)}, nil
}

func Transcode(_ context.Context, in bp.TranscodeInput) (bp.TranscodeOutput, error) {
	slog.Info("transcoding video", "asset_id", in.AssetID)
	return bp.TranscodeOutput{VideoURL: fmt.Sprintf("cdn://video/%s.mp4", in.AssetID)}, nil
}

func Thumbnail(_ context.Context, in bp.ThumbnailInput) (bp.ThumbnailOutput, error) {
	slog.Info("generating thumbnail", "asset_id", in.AssetID)
	return bp.ThumbnailOutput{ThumbURL: fmt.Sprintf("cdn://thumb/%s.jpg", in.AssetID)}, nil
}

func ExtractAudio(_ context.Context, in bp.ExtractAudioInput) (bp.ExtractAudioOutput, error) {
	slog.Info("extracting audio", "asset_id", in.AssetID)
	return bp.ExtractAudioOutput{AudioURL: fmt.Sprintf("cdn://audio/%s.mp3", in.AssetID)}, nil
}

func Package(_ context.Context, in bp.PackageInput) (bp.PackageOutput, error) {
	slog.Info("packaging bundle", "asset_id", in.AssetID, "video", in.VideoURL, "thumb", in.ThumbURL, "audio", in.AudioURL)
	return bp.PackageOutput{BundleURL: fmt.Sprintf("cdn://bundle/%s.zip", in.AssetID)}, nil
}
