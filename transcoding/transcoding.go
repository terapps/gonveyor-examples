package transcoding

import (
	"github.com/terapps/gonveyor"
)

// Transcoding: download → [transcode, thumbnail, extract_audio] → package
//
//	                  ┌──> transcode    ──┐
//	download ─────────┼──> thumbnail    ──┼──> package
//	                  └──> extract_audio──┘
var Transcoding = gonveyor.New("transcoding",
	Download,
	gonveyor.Wire(Transcode,
		gonveyor.Intake(Download, func(o DownloadOutput, in *TranscodeInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	gonveyor.Wire(Thumbnail,
		gonveyor.Intake(Download, func(o DownloadOutput, in *ThumbnailInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	gonveyor.Wire(ExtractAudio,
		gonveyor.Intake(Download, func(o DownloadOutput, in *ExtractAudioInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	gonveyor.Wire(Package,
		gonveyor.Intake(Transcode, func(o TranscodeOutput, in *PackageInput) {
			in.VideoURL = o.VideoURL
		}),
		gonveyor.Intake(Thumbnail, func(o ThumbnailOutput, in *PackageInput) {
			in.ThumbURL = o.ThumbURL
		}),
		gonveyor.Intake(ExtractAudio, func(o ExtractAudioOutput, in *PackageInput) {
			in.AudioURL = o.AudioURL
		}),
	),
)

var Template = gonveyor.NewLaunchTemplate(Transcoding, func(p DownloadInput) []gonveyor.ManifestOption {
	return []gonveyor.ManifestOption{gonveyor.Seed(Download, p)}
})

func Manifest(assetID, sourceURL string) (gonveyor.BlueprintManifest, error) {
	return Template.Manifest(DownloadInput{AssetID: assetID, SourceURL: sourceURL})
}
