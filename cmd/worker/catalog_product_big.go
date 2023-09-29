package worker

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/urfave/cli/v2"
	"marketplace-svc/app"
	"marketplace-svc/app/model/base"
	"marketplace-svc/helper/queue"
)

type CatalogProductBig struct {
	Topic string
	Infra app.Infra
}

func NewCatalogProductBig(infra app.Infra) CatalogProductBig {
	return CatalogProductBig{
		Infra: infra,
		Topic: base.TOPIC_CATALOG_PRODUCT_BIG,
	}
}

func (cp CatalogProductBig) Cmd() *cli.Command {
	return &cli.Command{
		Name:  cp.Topic,
		Usage: "worker " + cp.Topic + " --indices=",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "indices", Aliases: []string{"i"}, Value: 1},
		},
		Action: func(c *cli.Context) error {
			return cp.Subscriber(c.Int("indices"))
		},
	}
}

func (cp CatalogProductBig) Subscriber(indices int) error {
	for i := 0; i < indices; i++ {
		initConsumer, err := queue.NewKafkaConsumer(*cp.Infra.Config, cp.Infra.Log, cp.Topic)
		if err != nil {
			cp.Infra.Log.Error(err)
			return err
		}
		initConsumer.Subscribe(cp.Topic, cp.HandlerSubscriber)
	}
	return nil
}

func (cp CatalogProductBig) HandlerSubscriber(km *kafka.Message) {
	//var payloadKafka base.PayloadKafka
	//json.Unmarshal(payload, &payloadKafka)

	fmt.Println("payload: ", string(km.Value))
}
