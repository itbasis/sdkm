package golang

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	itbasisCoreLog "github.com/itbasis/go-tools-core/log"
	"github.com/itbasis/go-tools-sdkm/internal/plugins/golang/modfile"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
)

func (receiver *goPlugin) Current(ctx context.Context, rebuildCache, onlyInstalled bool, baseDir string) (sdkmSDKVersion.SDKVersion, error) {
	goModFile, errGoModFile := modfile.ReadGoModFile(baseDir)
	if errGoModFile != nil {
		slog.Error("Failed to read go.mod file", itbasisCoreLog.SlogAttrError(errGoModFile))

		return nil, errGoModFile //nolint:wrapcheck // TODO
	}

	var (
		sdkVersion sdkmSDKVersion.SDKVersion
		err        error
	)

	if toolchain := goModFile.Toolchain; toolchain != nil {
		sdkVersion, err = receiver.LatestVersionByPrefix(ctx, rebuildCache, onlyInstalled, toolchain.Name[2:])
	} else {
		var prefix = goModFile.Go.Version

		if strings.HasSuffix(prefix, ".0") && strings.Count(prefix, ".") == 2 {
			prefix = prefix[0:len(prefix) - 2]
		}

		sdkVersion, err = receiver.LatestVersionByPrefix(ctx, rebuildCache, onlyInstalled, prefix)
	}

	slog.Debug(fmt.Sprintf("sdkVersion: %++v", sdkVersion))

	if err != nil {
		return nil, err
	}

	return sdkVersion, nil
}
