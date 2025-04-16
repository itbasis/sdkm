package cmd

import (
	itbasisCoreCmd "github.com/itbasis/go-tools-core/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	_flagWithUninstalled = "with-uninstalled"
)

func InitFlagWithUninstalled(flags *pflag.FlagSet) {
	flags.Bool(_flagWithUninstalled, false, "check with uninstalled SDK")
}

func IsFlagWithUninstalled(cmd *cobra.Command) bool {
	flag, err := cmd.Flags().GetBool(_flagWithUninstalled)
	itbasisCoreCmd.RequireNoError(cmd, err)

	return flag
}
