package versions

import (
	"context"
	"log/slog"
	"sort"

	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
)

func (receiver *versions) AllVersions(ctx context.Context, rebuildCache bool) (sdkmSDKVersion.SdkVersionList, error) {
	slog.Debug("get all versions", slog.Bool("rebuildCache", rebuildCache))

	if rebuildCache || !receiver.cache.Valid(ctx) {
		receiver.updateCache(ctx, true, true, true)
	}

	if result := receiver.filterCacheVersions(ctx); result.Len() > 0 {
		return result, nil
	}

	slog.Debug("Trying to force refresh cache to search for versions")

	receiver.updateCache(ctx, true, true, true)

	return receiver.filterCacheVersions(ctx), nil
}

func (receiver *versions) filterCacheVersions(ctx context.Context) sdkmSDKVersion.SdkVersionList {
	var sdkVersionList = sdkmSDKVersion.NewSdkVersionList()

	for _, versionType := range []sdkmSDKVersion.VersionType{sdkmSDKVersion.TypeStable, sdkmSDKVersion.TypeUnstable, sdkmSDKVersion.TypeArchived} {
		v := receiver.cache.Load(ctx, versionType)
		if v.IsEmpty() {
			continue
		}

		sdkVersionList.Add(v.AsList()...)
	}

	sort.Sort(sdkVersionList)

	return sdkVersionList
}
