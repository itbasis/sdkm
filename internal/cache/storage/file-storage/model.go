package filestorage

import sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"

type model struct {
	Updated  updated                               `json:"updated"`
	Versions sdkmSDKVersion.MapSdkVersionGroupType `json:"versions"`
}
