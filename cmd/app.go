package cmd

import (
	"context"
	"log"

	itbasisCoreApp "github.com/itbasis/go-tools-core/app"
	itbasisCoreCmd "github.com/itbasis/go-tools-core/cmd"
	sdkmCmdCurrent "github.com/itbasis/go-tools-sdkm/cmd/current"
	sdkmCmdEnv "github.com/itbasis/go-tools-sdkm/cmd/env"
	sdkmCmdExec "github.com/itbasis/go-tools-sdkm/cmd/exec"
	sdkmCmdInstall "github.com/itbasis/go-tools-sdkm/cmd/install"
	sdkmCmdLatest "github.com/itbasis/go-tools-sdkm/cmd/latest"
	sdkmCmdList "github.com/itbasis/go-tools-sdkm/cmd/list"
	sdkmCmdPlugins "github.com/itbasis/go-tools-sdkm/cmd/plugins"
	sdkmCmdReshim "github.com/itbasis/go-tools-sdkm/cmd/reshim"
	sdkmCmd "github.com/itbasis/go-tools-sdkm/internal/cmd"
)

func InitApp(ctx context.Context) *itbasisCoreApp.App {
	var cmdRoot, err = itbasisCoreCmd.InitDefaultCmdRoot(
		ctx,
		"SDK Manager",
		itbasisCoreCmd.WithFlagDebug(sdkmCmd.FlagDebugName, itbasisCoreCmd.DefaultFlagDebugDescription, true, true),
	)

	if err != nil {
		log.Fatal(err)
	}

	sdkmCmd.InitFlagSdkRootDir(cmdRoot.PersistentFlags())

	cmdRoot.AddCommand(
		sdkmCmdPlugins.NewPluginsCommand(),
		sdkmCmdList.NewListCommand(),
		sdkmCmdLatest.NewLatestCommand(),
		sdkmCmdCurrent.NewCurrentCommand(),
		sdkmCmdInstall.NewInstallCommand(),
		sdkmCmdEnv.NewEnvCommand(),
		sdkmCmdExec.NewExecCommand(),
		sdkmCmdReshim.NewReshimCommand(),
	)

	return itbasisCoreApp.NewApp(cmdRoot)
}
