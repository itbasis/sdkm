package plugin

import (
	"context"
	"fmt"
	"io"

	itbasisCoreEnv "github.com/itbasis/go-tools-core/env"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=base-plugin.mock.go
type BasePlugin interface {
	fmt.GoStringer

	GetSDKDir() string
	GetSDKVersionDir(pluginID ID, version string) string
	HasInstalled(pluginID ID, version string) bool

	PrepareEnvironment(overrideEnv itbasisCoreEnv.Map, envKeys ...string) itbasisCoreEnv.Map

	Exec(ctx context.Context, cli string, overrideEnv map[string]string, stdIn io.Reader, stdOut, stdErr io.Writer, args []string) error
}
