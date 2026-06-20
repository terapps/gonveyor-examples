package handler

import (
	"context"
	"fmt"
	"log/slog"

	st "github.com/terapps/gonveyor-examples/transcoding/stations"
)

func Download(_ context.Context, in st.DownloadInput) (st.DownloadOutput, error) {
	slog.Info("downloading asset", "asset_id", in.AssetID, "source", in.SourceURL)
	return st.DownloadOutput{AssetID: in.AssetID, LocalURL: fmt.Sprintf("local://%s", in.AssetID)}, nil
}

func Transcode(_ context.Context, in st.TranscodeInput) (st.TranscodeOutput, error) {
	slog.Info("transcoding video", "asset_id", in.AssetID)
	return st.TranscodeOutput{VideoURL: fmt.Sprintf("cdn://video/%s.mp4", in.AssetID)}, nil
}

func Thumbnail(_ context.Context, in st.ThumbnailInput) (st.ThumbnailOutput, error) {
	slog.Info("generating thumbnail", "asset_id", in.AssetID)
	return st.ThumbnailOutput{ThumbURL: fmt.Sprintf("cdn://thumb/%s.jpg", in.AssetID)}, nil
}

func ExtractAudio(_ context.Context, in st.ExtractAudioInput) (st.ExtractAudioOutput, error) {
	slog.Info("extracting audio", "asset_id", in.AssetID)
	return st.ExtractAudioOutput{AudioURL: fmt.Sprintf("cdn://audio/%s.mp3", in.AssetID)}, nil
}

func Package(_ context.Context, in st.PackageInput) (st.PackageOutput, error) {
	slog.Info("packaging bundle", "asset_id", in.AssetID, "video", in.VideoURL, "thumb", in.ThumbURL, "audio", in.AudioURL)
	return st.PackageOutput{BundleURL: fmt.Sprintf("cdn://bundle/%s.zip", in.AssetID)}, nil
}
