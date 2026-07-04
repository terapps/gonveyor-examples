package stations

import "github.com/terapps/gonveyor"

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

var Download     = gonveyor.Define[DownloadInput, DownloadOutput]("download")
var Transcode    = gonveyor.Define[TranscodeInput, TranscodeOutput]("transcode")
var Thumbnail    = gonveyor.Define[ThumbnailInput, ThumbnailOutput]("thumbnail")
var ExtractAudio = gonveyor.Define[ExtractAudioInput, ExtractAudioOutput]("extract_audio")
var Package      = gonveyor.Define[PackageInput, PackageOutput]("package")
