package worker

import (
	"github.com/urfave/cli/v2"
	"marketplace-svc/app"
)

func GetWorkerHandlerCommand(infra app.Infra) *cli.Command {
	// list subcommands worker
	arrSubCmd := []*cli.Command{
		//NewCatalogProduct(infra).Cmd(),
		//NewCatalogProductBig(infra).Cmd(),
		NewOrderCreateNotify(infra).Cmd(),
	}
	return &cli.Command{
		Name:        "worker",
		Aliases:     []string{"w"},
		Usage:       "worker",
		Description: "Handling Worker Marketplace",
		HelpName:    "worker",
		Subcommands: arrSubCmd,
	}
}
