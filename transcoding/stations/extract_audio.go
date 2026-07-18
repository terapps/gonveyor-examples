package stations

import "github.com/terapps/gonveyor"

type ExtractAudioInput struct {
	AssetID  string
	LocalURL string
}
type ExtractAudioOutput struct {
	AudioURL string
}

var ExtractAudio = gonveyor.Define[ExtractAudioInput, ExtractAudioOutput]("extract_audio")
