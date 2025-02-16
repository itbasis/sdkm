package golang

import (
	"context"
	"log/slog"
	"os"
	"path"

	"github.com/itbasis/go-tools-core/env"
	itbasisCoreLog "github.com/itbasis/go-tools-core/log"
	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
	"golang.org/x/tools/godoc/vfs"
)

const (
	EnvGoBin  = "GOBIN"
	EnvGoPath = string(vfs.RootTypeGoPath)
	EnvGoRoot = string(vfs.RootTypeGoRoot)
	EnvPath   = "PATH"

	EnvSdkmOriginPrefix = "SDKM_ORIGIN_"
)

func (receiver *goPlugin) Env(ctx context.Context, rebuildCache bool, baseDir string) (env.Map, error) {
	sdkVersion, errCurrent := receiver.Current(ctx, rebuildCache, baseDir)

	if errCurrent != nil {
		return map[string]string{}, errCurrent
	}

	return receiver.EnvByVersion(ctx, sdkVersion.ID)
}

func (receiver *goPlugin) EnvByVersion(_ context.Context, version string) (env.Map, error) {
	var (
		goCacheDir = receiver.getGoCacheDir(version)
		goBin      = path.Join(goCacheDir, "bin")
	)

	var envs = receiver.prepareEnvMap()
	//
	envs[EnvGoRoot] = receiver.basePlugin.GetSDKVersionDir(pluginGoConsts.PluginID, version)
	envs[EnvGoPath] = goCacheDir
	envs[EnvGoBin] = goBin
	envs[EnvPath] = itbasisCoreOs.AddBeforePath(
		envs[EnvPath],
		path.Join(receiver.basePlugin.GetSDKVersionDir(pluginGoConsts.PluginID, version), "bin"),
		goBin,
		itbasisCoreOs.ExecutableDir(),
	)

	slog.Debug("envs", itbasisCoreLog.SlogAttrMap("envs", envs))

	return envs, nil
}

func (receiver *goPlugin) prepareEnvMap() env.Map {
	var envKeys = []string{EnvPath, EnvGoRoot, EnvGoPath, EnvGoBin}

	var result = make(env.Map, len(envKeys))

	for _, key := range envKeys {
		var (
			sdkmOriginKey   = EnvSdkmOriginPrefix + key
			sdkmOriginValue = os.Getenv(sdkmOriginKey)
			value           = os.Getenv(key)
		)

		if sdkmOriginValue != "" {
			result[sdkmOriginKey] = sdkmOriginValue
		} else {
			result[sdkmOriginKey] = value
		}

		result[key] = value
	}

	return result
}
