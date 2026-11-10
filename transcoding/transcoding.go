package transcoding

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/transcoding/stations"
)

// Transcoding: download → [transcode, thumbnail, extract_audio] → package
//
//	                  ┌──> transcode    ──┐
//	download ─────────┼──> thumbnail    ──┼──> package
//	                  └──> extract_audio──┘
var Transcoding = gonveyor.New("transcoding",
	st.Download,
	st.Transcode.Intake(st.Download, func(o st.DownloadOutput, in *st.TranscodeInput) {
		in.AssetID = o.AssetID
		in.LocalURL = o.LocalURL
	}),
	st.Thumbnail.Intake(st.Download, func(o st.DownloadOutput, in *st.ThumbnailInput) {
		in.AssetID = o.AssetID
		in.LocalURL = o.LocalURL
	}),
	st.ExtractAudio.Intake(st.Download, func(o st.DownloadOutput, in *st.ExtractAudioInput) {
		in.AssetID = o.AssetID
		in.LocalURL = o.LocalURL
	}),
	st.Package.
		Intake(st.Transcode, func(o st.TranscodeOutput, in *st.PackageInput) {
			in.VideoURL = o.VideoURL
		}).
		Intake(st.Thumbnail, func(o st.ThumbnailOutput, in *st.PackageInput) {
			in.ThumbURL = o.ThumbURL
		}).
		Intake(st.ExtractAudio, func(o st.ExtractAudioOutput, in *st.PackageInput) {
			in.AudioURL = o.AudioURL
		}),
)

var Template = gonveyor.NewLaunchTemplate(Transcoding, func(p st.DownloadInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{st.Download.Seed(p)}
})

func Manifest(assetID, sourceURL string) (gonveyor.BlueprintManifest, error) {
	return Template.Manifest(st.DownloadInput{AssetID: assetID, SourceURL: sourceURL})
}
