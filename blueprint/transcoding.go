package blueprint

import (
	"github.com/terapps/gonveyor"
	"github.com/terapps/gonveyor/blueprint"
	"github.com/terapps/gonveyor/ledger"
)

// --- types ---

type DownloadInput struct {
	AssetID   string
	SourceURL string
}
type DownloadOutput struct {
	AssetID  string
	LocalURL string
}

type TranscodeInput struct {
	AssetID  string
	LocalURL string
}
type TranscodeOutput struct {
	VideoURL string
}

type ThumbnailInput struct {
	AssetID  string
	LocalURL string
}
type ThumbnailOutput struct {
	ThumbURL string
}

type ExtractAudioInput struct {
	AssetID  string
	LocalURL string
}
type ExtractAudioOutput struct {
	AudioURL string
}

type PackageInput struct {
	AssetID  string
	VideoURL string
	ThumbURL string
	AudioURL string
}
type PackageOutput struct {
	BundleURL string
}

// --- stations ---

var Download     = blueprint.Define[DownloadInput, DownloadOutput]("download")
var Transcode    = blueprint.Define[TranscodeInput, TranscodeOutput]("transcode")
var Thumbnail    = blueprint.Define[ThumbnailInput, ThumbnailOutput]("thumbnail")
var ExtractAudio = blueprint.Define[ExtractAudioInput, ExtractAudioOutput]("extract_audio")
var Package      = blueprint.Define[PackageInput, PackageOutput]("package")

// Transcoding: download → [transcode, thumbnail, extract_audio] → package
//
//	                  ┌──> transcode    ──┐
//	download ─────────┼──> thumbnail    ──┼──> package
//	                  └──> extract_audio──┘
var Transcoding = blueprint.New("transcoding",
	Download,
	blueprint.Wire(Transcode,
		gonveyor.Intake(Download, func(o DownloadOutput, in *TranscodeInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	blueprint.Wire(Thumbnail,
		gonveyor.Intake(Download, func(o DownloadOutput, in *ThumbnailInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	blueprint.Wire(ExtractAudio,
		gonveyor.Intake(Download, func(o DownloadOutput, in *ExtractAudioInput) {
			in.AssetID = o.AssetID
			in.LocalURL = o.LocalURL
		}),
	),
	blueprint.Wire(Package,
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

func TranscodingManifest(assetID, sourceURL string) (ledger.BlueprintManifest, error) {
	return Transcoding.Manifest(
		gonveyor.Seed(Download, DownloadInput{AssetID: assetID, SourceURL: sourceURL}),
	)
}
