package taskselastic

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"marketplace-svc/app"
	elasticregistry "marketplace-svc/app/registry/elastic"
	elasticservice "marketplace-svc/app/service/elastic"
	"time"
)

type Product struct {
	Infra                 app.Infra
	ElasticProductService elasticservice.ElasticProductService
}

func NewCatalogProduct(infra app.Infra) Product {
	return Product{
		Infra:                 infra,
		ElasticProductService: elasticregistry.RegisterEsProductService(&infra),
	}
}

func (cp Product) Cmd() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "elastic-product:reindex",
			Usage: "task elastic-product:reindex",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "index", Usage: "product_store"},
				&cli.StringFlag{Name: "ids", Usage: "1,2,3,4"},
				&cli.IntFlag{Name: "store_id", Usage: "1", Value: 1},
				&cli.StringFlag{Name: "type", Usage: "id/sku"},
				&cli.IntFlag{Name: "merchant_id", Usage: "620"},
				&cli.BoolFlag{Name: "flush", Usage: "true/false", Value: false},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("elastic-product-reindex")
				// of Since method
				now := time.Now()
				_ = cp.ElasticProductService.Reindex(c.String("index"), c.String("ids"), c.Int("store_id"), c.String("type"), c.Int("merchant_id"), c.Bool("flush"))
				fmt.Println("time elapse: ", time.Since(now))
				return nil
			},
		},
		{
			Name:  "elastic-product:reindexByID",
			Usage: "task elastic-product:reindexByID",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "ids", Usage: "1,2,3,4", Required: true},
				&cli.IntFlag{Name: "store_id", Value: 1, Usage: "1"},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("elastic-product-reindex by ID")
				// of Since method
				now := time.Now()
				_ = cp.ElasticProductService.Reindex("", c.String("ids"), c.Int("store_id"), "id", 0, false)

				fmt.Println("time elapse: ", time.Since(now))
				return nil
			},
		},
	}
}
