package main

import (
	"github.com/urfave/cli/v2"
	"marketplace-svc/app"
	cmdbase "marketplace-svc/cmd/base"
	"marketplace-svc/helper/config"
	"os"
	"sort"
)

func main() {
	cfg := config.Init()
	// pass to infra value
	infra := &app.Infra{
		Config: cfg,
	}
	// enable configuration
	infra.EnableLog().EnableDB().EnableKafkaProducer()

	// init cmd
	app := cmdbase.NewBaseCommand(*infra).PopulateCommand()

	sort.Sort(cli.CommandsByName(app.Cmd.Commands))
	sort.Sort(cli.FlagsByName(app.Cmd.Flags))
	err := app.Cmd.Run(os.Args)
	if err != nil {
		infra.Log.Error(err)
	}
}
