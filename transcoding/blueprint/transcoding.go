package blueprint

import (
	"github.com/terapps/gonveyor"
	st "github.com/terapps/gonveyor-examples/transcoding/stations"
	"github.com/terapps/gonveyor/ledger"
)

// Transcoding: download → [transcode, thumbnail, extract_audio] → package
//
//	                  ┌──> transcode    ──┐
//	download ─────────┼──> thumbnail    ──┼──> package
//	                  └──> extract_audio──┘
var Transcoding = gonveyor.New("transcoding",
	st.Download,
	gonveyor.Wire(st.Transcode,
		gonveyor.Intake(st.Download, func(o st.DownloadOutput, in *st.TranscodeInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	gonveyor.Wire(st.Thumbnail,
		gonveyor.Intake(st.Download, func(o st.DownloadOutput, in *st.ThumbnailInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	gonveyor.Wire(st.ExtractAudio,
		gonveyor.Intake(st.Download, func(o st.DownloadOutput, in *st.ExtractAudioInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	gonveyor.Wire(st.Package,
		gonveyor.Intake(st.Transcode, func(o st.TranscodeOutput, in *st.PackageInput) {
			in.VideoURL = o.VideoURL
		}),
		gonveyor.Intake(st.Thumbnail, func(o st.ThumbnailOutput, in *st.PackageInput) {
			in.ThumbURL = o.ThumbURL
		}),
		gonveyor.Intake(st.ExtractAudio, func(o st.ExtractAudioOutput, in *st.PackageInput) {
			in.AudioURL = o.AudioURL
		}),
	),
)

var Launcher = gonveyor.NewManifestBuilder(Transcoding, func(p st.DownloadInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(st.Download, p)}
})

func Manifest(assetID, sourceURL string) (ledger.BlueprintManifest, error) {
	return Launcher.Manifest(st.DownloadInput{AssetID: assetID, SourceURL: sourceURL})
}
