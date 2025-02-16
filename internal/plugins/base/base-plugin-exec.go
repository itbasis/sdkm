package base

import (
	"context"
	"io"

	itbasisCoreExec "github.com/itbasis/go-tools-core/exec"
	"github.com/pkg/errors"
)

func (receiver *basePlugin) Exec(
	ctx context.Context,
	cli string,
	overrideEnv map[string]string,
	stdIn io.Reader, stdOut, stdErr io.Writer,
	args []string,
) error {
	var envMap = receiver.PrepareEnvironment(overrideEnv)

	cmd, err := itbasisCoreExec.NewExecutable(
		ctx,
		cli,
		itbasisCoreExec.WithArgs(args...),
		itbasisCoreExec.WithCustomIn(stdIn),
		itbasisCoreExec.WithCustomOut(stdOut, stdErr),
		itbasisCoreExec.WithEnv(envMap),
	)
	if err != nil {
		return errors.Wrap(err, "error executing plugin")
	}

	if err := cmd.Execute(ctx); err != nil {
		return errors.Wrap(err, itbasisCoreExec.ErrFailedExecuteCommand.Error())
	}

	return nil
}
