package stations

import (
	"github.com/terapps/gonveyor"
)

type PackageInput struct {
	AssetID  string
	VideoURL string
	ThumbURL string
	AudioURL string
}
type PackageOutput struct {
	BundleURL string
}

var Package = gonveyor.Define[PackageInput, PackageOutput]("package")
