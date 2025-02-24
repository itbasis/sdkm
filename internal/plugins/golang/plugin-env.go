package golang

import (
	"context"
	"log/slog"
	"path"
	"strings"

	itbasisCoreEnv "github.com/itbasis/go-tools-core/env"
	itbasisCoreLog "github.com/itbasis/go-tools-core/log"
	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	pluginBase "github.com/itbasis/go-tools-sdkm/internal/plugins/base"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
)

func (receiver *goPlugin) Env(ctx context.Context, rebuildCache bool, baseDir string) (itbasisCoreEnv.Map, error) {
	sdkVersion, errCurrent := receiver.Current(ctx, rebuildCache, baseDir)

	if errCurrent != nil {
		return map[string]string{}, errCurrent
	}

	return receiver.EnvByVersion(ctx, sdkVersion.ID)
}

func (receiver *goPlugin) EnvByVersion(_ context.Context, version string) (itbasisCoreEnv.Map, error) {
	var (
		goCacheDir = receiver.getGoCacheDir(version)
		goBin      = path.Join(goCacheDir, "bin")
		goRootDir  = receiver.basePlugin.GetSDKVersionDir(pluginGoConsts.PluginID, version)

		envs = receiver.basePlugin.PrepareEnvironment(nil, itbasisCoreEnv.KeyGoRoot, itbasisCoreEnv.KeyGoPath, itbasisCoreEnv.KeyGoBin)
	)

	slog.Debug("envs", itbasisCoreLog.SlogAttrMap("envs", envs))

	var result = itbasisCoreEnv.Map{
		itbasisCoreEnv.KeyGoRoot: goRootDir,
		itbasisCoreEnv.KeyGoPath: goCacheDir,
		itbasisCoreEnv.KeyGoBin:  goBin,
		itbasisCoreEnv.KeyPath: itbasisCoreOs.FixPath(
			itbasisCoreOs.AddBeforePath(
				envs[itbasisCoreEnv.KeyPath],
				path.Join(goRootDir, "bin"),
				goBin,
				itbasisCoreOs.ExecutableDir(),
			),
		),
	}

	for k, v := range envs {
		if strings.HasPrefix(k, pluginBase.EnvSdkmOriginPrefix) {
			result[k] = v
		}
	}

	slog.Debug("envs", itbasisCoreLog.SlogAttrMap("result", result))

	return result, nil
}
