package tasks

import (
	"github.com/urfave/cli/v2"
	"marketplace-svc/app"
	taskselastic "marketplace-svc/cmd/tasks/elastic"
)

func GetTasksHandlerCommand(infra app.Infra) *cli.Command {
	// list subcommands task
	var arrSubCmd []*cli.Command
	arrSubCmd = append(arrSubCmd, taskselastic.NewCatalogProduct(infra).Cmd()...)

	return &cli.Command{
		Name:        "task",
		Aliases:     []string{"t"},
		Usage:       "task",
		Description: "Handling Task Marketplace",
		HelpName:    "task",
		Subcommands: arrSubCmd,
	}
}
