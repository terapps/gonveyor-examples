package stations

import "github.com/terapps/gonveyor"

type DownloadInput struct {
	AssetID   string `validate:"required"`
	SourceURL string `validate:"required,url"`
}
type DownloadOutput struct {
	AssetID  string
	LocalURL string
}

var Download = gonveyor.Define[DownloadInput, DownloadOutput]("download")
