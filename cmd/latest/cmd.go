package latest

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

func NewLatestCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "latest",
		Short: "Show latest stable version",
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
		sdkmPlugin        = sdkmPlugins.GetPluginByID(cmd)
		flagRebuildCache  = sdkmCmd.IsFlagRebuildCache(cmd)
		flagOnlyInstalled = !sdkmCmd.IsFlagWithUninstalled(cmd)

		sdkVersion sdkmSDKVersion.SDKVersion
		err        error
	)

	if len(args) == 0 {
		sdkVersion, err = sdkmPlugin.LatestVersion(cmd.Context(), flagRebuildCache, flagOnlyInstalled)
	} else {
		sdkVersion, err = sdkmPlugin.LatestVersionByPrefix(cmd.Context(), flagRebuildCache, flagOnlyInstalled, args[_idxArgVersion])
	}

	if err != nil {
		itbasisCoreCmd.Fatal(cmd, err)
	}

	cmd.Println(sdkVersion.PrintWithOptions(false, true, true))
}
