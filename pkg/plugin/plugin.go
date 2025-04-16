package plugin

import (
	"context"
	"io"

	"github.com/itbasis/go-tools-core/env"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
)

//nolint:interfacebloat // TODO
//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=plugin.mock.go
type SDKMPlugin interface {
	WithVersions(versions sdkmSDKVersion.SDKVersions) SDKMPlugin
	// WithBasePlugin(basePlugin BasePlugin) SDKMPlugin
	// WithDownloader(downloader) SDKMPlugin

	ListAllVersions(ctx context.Context, rebuildCache bool) (sdkmSDKVersion.SdkVersionList, error)
	ListAllVersionsByPrefix(ctx context.Context, rebuildCache bool, prefix string) (sdkmSDKVersion.SdkVersionList, error)

	LatestVersion(ctx context.Context, rebuildCache, onlyInstalled bool) (sdkmSDKVersion.SDKVersion, error)
	LatestVersionByPrefix(ctx context.Context, rebuildCache, onlyInstalled bool, prefix string) (sdkmSDKVersion.SDKVersion, error)

	Current(ctx context.Context, rebuildCache, onlyInstalled bool, baseDir string) (sdkmSDKVersion.SDKVersion, error)

	Install(ctx context.Context, rebuildCache bool, baseDir string) error
	InstallVersion(ctx context.Context, version string) error

	Env(ctx context.Context, rebuildCache, onlyInstalled bool, baseDir string) (env.Map, error)
	EnvByVersion(ctx context.Context, version string) (env.Map, error)

	Exec(ctx context.Context, rebuildCache bool, baseDir string, stdIn io.Reader, stdOut, stdErr io.Writer, args []string) error
}
