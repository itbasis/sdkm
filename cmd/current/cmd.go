package current

import (
	itbasisCoreCmd "github.com/itbasis/go-tools-core/cmd"
	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	sdkmCmd "github.com/itbasis/go-tools-sdkm/internal/cmd"
	sdkmPlugins "github.com/itbasis/go-tools-sdkm/plugins"
	"github.com/spf13/cobra"
)

func NewCurrentCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "current",
		Short: "Display current version",
	}

	sdkmCmd.InitFlagRebuildCache(cmd.PersistentFlags())
	sdkmCmd.InitFlagWithUninstalled(cmd.PersistentFlags())

	sdkmPlugins.AddPluginsAsSubCommands(
		cmd, func(cmdChild *cobra.Command) {
			cmdChild.Args = cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs)
			cmdChild.Run = _run
		},
	)

	return cmd
}

func _run(cmd *cobra.Command, _ []string) {
	var (
		sdkmPlugin        = sdkmPlugins.GetPluginByID(cmd)
		flagRebuildCache  = sdkmCmd.IsFlagRebuildCache(cmd)
		flagOnlyInstalled = !sdkmCmd.IsFlagWithUninstalled(cmd)
		sdkVersion, err   = sdkmPlugin.Current(cmd.Context(), flagRebuildCache, flagOnlyInstalled, itbasisCoreOs.Pwd())
	)

	if err != nil {
		itbasisCoreCmd.Fatal(cmd, err)
	}

	cmd.Println(sdkVersion.PrintWithOptions(false, false, false))
}
