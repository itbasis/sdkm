package golang

import (
	"path"

	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	sdkmCache "github.com/itbasis/go-tools-sdkm/internal/cache"
	cacheFileStorage "github.com/itbasis/go-tools-sdkm/internal/cache/storage/file-storage"
	sdkmCmd "github.com/itbasis/go-tools-sdkm/internal/cmd"
	pluginGoConsts "github.com/itbasis/go-tools-sdkm/internal/plugins/golang/consts"
	"github.com/spf13/cobra"
)

func CmdExtPlugin(cmd *cobra.Command) {
	flags := cmd.Flags()

	sdkmCmd.InitFlagCacheRootDir(flags)
	sdkmCmd.InitFlagRebuildCache(flags)
}

func (receiver *goPlugin) InitProcessCommandFlags(cmd *cobra.Command) {
	if cacheRootDir := sdkmCmd.GetCacheRootDir(cmd); cacheRootDir == "" {
		receiver.goCacheRootDir = path.Join(itbasisCoreOs.UserHomeDir(), ".cache", string(pluginGoConsts.PluginID))
	} else {
		receiver.goCacheRootDir = path.Join(cacheRootDir, string(pluginGoConsts.PluginID))
	}

	receiver.sdkVersions = receiver.sdkVersions.WithCache(
		sdkmCache.NewCache().
			WithExternalStore(cacheFileStorage.NewFileCacheStorage(
				path.Join(itbasisCoreOs.ExecutableDir(), ".cache"),
				pluginGoConsts.PluginID,
			)),
	)
}
