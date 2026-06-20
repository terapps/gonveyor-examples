package blueprint

import (
	"github.com/terapps/gonveyor"
	bp "github.com/terapps/gonveyor/blueprint"
	"github.com/terapps/gonveyor/ledger"
	st "github.com/terapps/gonveyor-examples/transcoding/stations"
)

// Transcoding: download → [transcode, thumbnail, extract_audio] → package
//
//	                  ┌──> transcode    ──┐
//	download ─────────┼──> thumbnail    ──┼──> package
//	                  └──> extract_audio──┘
var Transcoding = bp.New("transcoding",
	st.Download,
	bp.Wire(st.Transcode,
		gonveyor.Intake(st.Download, func(o st.DownloadOutput, in *st.TranscodeInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	bp.Wire(st.Thumbnail,
		gonveyor.Intake(st.Download, func(o st.DownloadOutput, in *st.ThumbnailInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	bp.Wire(st.ExtractAudio,
		gonveyor.Intake(st.Download, func(o st.DownloadOutput, in *st.ExtractAudioInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	bp.Wire(st.Package,
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

func Manifest(assetID, sourceURL string) (ledger.BlueprintManifest, error) {
	return Transcoding.Manifest(
		gonveyor.Seed(st.Download, st.DownloadInput{AssetID: assetID, SourceURL: sourceURL}),
	)
}
