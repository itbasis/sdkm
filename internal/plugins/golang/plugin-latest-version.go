package golang

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	"github.com/pkg/errors"
)

func (receiver *goPlugin) LatestVersion(ctx context.Context, rebuildCache, onlyInstalled bool) (sdkmSDKVersion.SDKVersion, error) {
	slog.Debug(fmt.Sprintf("searching for latest version [onlyInstalled=%t]", onlyInstalled))

	sdkVersionList, err := receiver.sdkVersions.AllVersions(ctx, rebuildCache)
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO
	}

	if !onlyInstalled {
		slog.Debug("return first SDK in list")

		return sdkVersionList.First()
	}

	for sdkVersion := range sdkVersionList.Seq() {
		slog.Debug(fmt.Sprintf("check SDK [%s]: %t", sdkVersion.GetId(), sdkVersion.HasInstalled()))

		if sdkVersion.HasInstalled() {
			slog.Debug("return first installed SDK in list")

			return sdkVersion, nil
		}
	}

	return nil, sdkmSDKVersion.ErrSDKVersionNotFound
}

func (receiver *goPlugin) LatestVersionByPrefix(ctx context.Context, rebuildCache, onlyInstalled bool, prefix string) (sdkmSDKVersion.SDKVersion, error) {
	slog.Debug(fmt.Sprintf("searching for latest version [onlyInstalled=%t] by prefix: %s", onlyInstalled, prefix))

	if prefix == "" {
		return receiver.LatestVersion(ctx, rebuildCache, onlyInstalled)
	}

	sdkVersionList, err := receiver.ListAllVersions(ctx, rebuildCache)
	if err != nil {
		return nil, errors.WithMessage(sdkmSDKVersion.ErrSDKVersionNotFound, err.Error())
	}

	for sdkVersion := range sdkVersionList.Seq() {
		if strings.HasPrefix(sdkVersion.GetId(), prefix) && (!onlyInstalled || sdkVersion.HasInstalled()) {
			return sdkVersion, nil
		}
	}

	return nil, errors.WithMessagef(sdkmSDKVersion.ErrSDKVersionNotFound, "version by prefix %s", prefix)
}
