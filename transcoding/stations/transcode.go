package stations

import (
	"github.com/terapps/gonveyor"
)

type TranscodeInput struct {
	AssetID  string
	LocalURL string
}
type TranscodeOutput struct {
	VideoURL string
}

var Transcode = gonveyor.Define[TranscodeInput, TranscodeOutput]("transcode")
