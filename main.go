package main

import (
	"context"

	"github.com/itbasis/go-tools-sdkm/cmd"
)

func main() {
	cmd.InitApp(context.Background()).Run()
}
