package golang

import (
	"context"
	"sort"
	"strings"

	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
)

func (receiver *goPlugin) ListAllVersions(ctx context.Context, rebuildCache bool) (sdkmSDKVersion.SdkVersionList, error) {
	var sdkList, err = receiver.sdkVersions.AllVersions(ctx, rebuildCache) //nolint:wrapcheck // TODO
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO  c
	}

	for sdkVersion := range sdkList.Seq() {
		receiver.enrichSDKVersion(sdkVersion)
	}

	sort.Sort(sdkList)

	return sdkList, nil
}

func (receiver *goPlugin) ListAllVersionsByPrefix(ctx context.Context, rebuildCache bool, prefix string) (sdkmSDKVersion.SdkVersionList, error) {
	var allVersions, err = receiver.ListAllVersions(ctx, rebuildCache)
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO
	}

	if prefix == "" {
		return allVersions, nil
	}

	var sdkVersionList = sdkmSDKVersion.NewSdkVersionList()

	for sdkVersion := range allVersions.Seq() {
		if strings.HasPrefix(sdkVersion.GetId(), prefix) {
			sdkVersionList.Add(sdkVersion)
		}
	}

	return sdkVersionList, nil
}
