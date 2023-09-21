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
	infra.WithLog().WithDB().WithKafkaProducer().WithElasticClient()

	// init cmd
	baseCli := cmdbase.NewBaseCommand(*infra).PopulateCommand()

	sort.Sort(cli.CommandsByName(baseCli.Cmd.Commands))
	sort.Sort(cli.FlagsByName(baseCli.Cmd.Flags))
	err := baseCli.Cmd.Run(os.Args)
	if err != nil {
		infra.Log.Error(err)
	}
}
