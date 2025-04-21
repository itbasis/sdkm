package env

import (
	itbasisCoreCmd "github.com/itbasis/go-tools-core/cmd"
	itbasisCoreEnv "github.com/itbasis/go-tools-core/env"
	itbasisCoreOs "github.com/itbasis/go-tools-core/os"
	sdkmCmd "github.com/itbasis/go-tools-sdkm/internal/cmd"
	sdkmPlugins "github.com/itbasis/go-tools-sdkm/plugins"
	"github.com/spf13/cobra"
)

const (
	_idxArgVersion = 0
)

func NewEnvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Displays environment variables inside the environment used for the plugin",
	}

	sdkmCmd.InitFlagRebuildCache(cmd.PersistentFlags())
	sdkmCmd.InitFlagWithUninstalled(cmd.PersistentFlags())

	sdkmPlugins.AddPluginsAsSubCommands(
		cmd, func(cmdChild *cobra.Command) {
			cmdChild.Use = itbasisCoreCmd.BuildUse(cmdChild.Use, sdkmCmd.UseArgVersion)
			cmdChild.ArgAliases = []string{sdkmCmd.ArgAliasVersion}
			cmdChild.Args = cobra.MatchAll(cobra.MaximumNArgs(2), cobra.OnlyValidArgs)
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
		envMap            map[string]string
		err               error
	)

	if len(args) == 0 {
		envMap, err = sdkmPlugin.Env(
			cmd.Context(),
			flagRebuildCache,
			flagOnlyInstalled,
			itbasisCoreOs.Pwd(),
		)
	} else {
		envMap, err = sdkmPlugin.EnvByVersion(cmd.Context(), args[_idxArgVersion])
	}

	if err != nil {
		itbasisCoreCmd.Fatal(cmd, err)
	}

	for _, env := range itbasisCoreEnv.MapToSlices(envMap) {
		cmd.Println(env)
	}
}
