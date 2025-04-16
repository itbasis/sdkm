package list

import (
	itbasisCoreCmd "github.com/itbasis/go-tools-core/cmd"
	sdkmCmd "github.com/itbasis/go-tools-sdkm/internal/cmd"
	sdkmSDKVersion "github.com/itbasis/go-tools-sdkm/pkg/sdk-version"
	sdkmPlugins "github.com/itbasis/go-tools-sdkm/plugins"
	"github.com/spf13/cobra"
)

const (
	_idxArgVersion = 0
)

func newListAllCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "all",
		Short: "List all versions",
	}

	sdkmPlugins.AddPluginsAsSubCommands(
		cmd, func(cmdChild *cobra.Command) {
			cmdChild.Use = itbasisCoreCmd.BuildUse(cmdChild.Use, sdkmCmd.UseArgVersion)
			cmdChild.ArgAliases = []string{sdkmCmd.ArgAliasVersion}
			cmdChild.Args = cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs)
			cmdChild.Run = _run
		},
	)

	return cmd
}

func _run(cmd *cobra.Command, args []string) {
	var (
		plugin         = sdkmPlugins.GetPluginByID(cmd)
		sdkVersionList sdkmSDKVersion.SdkVersionList
		err            error
	)

	var flagRebuildCache = sdkmCmd.IsFlagRebuildCache(cmd)

	if len(args) == 0 {
		sdkVersionList, err = plugin.ListAllVersions(cmd.Context(), flagRebuildCache)
	} else {
		sdkVersionList, err = plugin.ListAllVersionsByPrefix(cmd.Context(), flagRebuildCache, args[_idxArgVersion])
	}

	if err != nil {
		itbasisCoreCmd.Fatal(cmd, err)
	}

	for sdkVersion := range sdkVersionList.Seq() {
		// TODO code smell
		cmd.Println(sdkVersion)
	}
}
