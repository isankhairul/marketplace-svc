package cmdbase

import (
	"github.com/urfave/cli/v2"
	"marketplace-svc/app"
	"marketplace-svc/cmd/worker"
)

type BaseCommand struct {
	Cmd   *cli.App
	Infra app.Infra
}

func NewBaseCommand(infra app.Infra) *BaseCommand {
	baseCommand := &BaseCommand{
		Cmd:   cli.NewApp(),
		Infra: infra,
	}

	baseCommand.Cmd.Name = "Command execution for Marketplace CLI"
	baseCommand.Cmd.Usage = "Run task by command CLI for Membership"
	baseCommand.Cmd.Authors = []*cli.Author{{"Khairul Ihksan", "khairul.ihksan@klikdokter.com"}}
	baseCommand.Cmd.Version = "1.0.0"

	baseCommand.Cmd.CommandNotFound = func(c *cli.Context, command string) {
		infra.Log.Info("No matching command " + command)
		cli.ShowAppHelp(c)
	}

	return baseCommand
}

func (base *BaseCommand) PopulateCommand() *BaseCommand {
	cmdWorker := worker.GetWorkerHandlerCommand(base.Infra)
	arrCmd := []*cli.Command{cmdWorker}
	base.Cmd.Commands = arrCmd

	return base
}
