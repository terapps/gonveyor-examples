package stations

import (
	"github.com/terapps/gonveyor"
)

type ThumbnailInput struct {
	AssetID  string
	LocalURL string
}
type ThumbnailOutput struct {
	ThumbURL string
}

var Thumbnail = gonveyor.Define[ThumbnailInput, ThumbnailOutput]("thumbnail")
