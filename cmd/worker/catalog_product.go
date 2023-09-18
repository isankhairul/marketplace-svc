package worker

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/urfave/cli/v2"
	"marketplace-svc/app"
	"marketplace-svc/app/model/base"
	"marketplace-svc/helper/queue"
	"sync"
)

type CatalogProduct struct {
	Topic string
	Infra app.Infra
}

func NewCatalogProduct(infra app.Infra) IWorker {
	return &CatalogProduct{
		Infra: infra,
		Topic: base.TOPIC_CATALOG_PRODUCT,
	}
}

func (cp CatalogProduct) Cmd() *cli.Command {
	return &cli.Command{
		Name:  cp.Topic,
		Usage: "worker " + cp.Topic + " --indices=",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "indices", Aliases: []string{"i"}, Value: 1},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("running worker " + cp.Topic + " with indices " + fmt.Sprint(c.Int("indices")))
			return cp.Subscriber(c.Int("indices"))
		},
	}
}

func (cp CatalogProduct) Subscriber(indices int) error {
	var wg sync.WaitGroup

	for i := 0; i < indices; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			initConsumer, err := queue.NewKafkaConsumer(*cp.Infra.Config, cp.Infra.Log, cp.Topic)
			if err != nil {
				cp.Infra.Log.Error(err)
				return
			}
			initConsumer.Subscribe(cp.Topic, cp.HandlerSubscriber)
		}()
	}
	wg.Wait()
	return nil
}

func (cp CatalogProduct) HandlerSubscriber(km *kafka.Message) {
	fmt.Println("payload: ", string(km.Value))
}
