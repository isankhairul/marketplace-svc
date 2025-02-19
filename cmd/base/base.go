package cmdbase

import (
	"marketplace-svc/app"
	"marketplace-svc/cmd/tasks"
	"marketplace-svc/cmd/worker"

	"github.com/urfave/cli/v2"
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
	baseCommand.Cmd.Authors = []*cli.Author{{Name: "Khairul Ihksan", Email: "khairul.ihksan@klikdokter.com"}}
	baseCommand.Cmd.Version = "1.0.0"

	baseCommand.Cmd.CommandNotFound = func(c *cli.Context, command string) {
		infra.Log.Info("No matching command " + command)
		//_ = cli.ShowAppHelp(c)
	}

	return baseCommand
}

func (base *BaseCommand) PopulateCommand() *BaseCommand {
	cmdWorker := worker.GetWorkerHandlerCommand(base.Infra)
	cmdTask := tasks.GetTasksHandlerCommand(base.Infra)
	arrCmd := []*cli.Command{cmdWorker, cmdTask}
	base.Cmd.Commands = arrCmd

	return base
}
