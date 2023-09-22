package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"marketplace-svc/app"
	"marketplace-svc/app/model/base"
	"marketplace-svc/app/model/request"
	"marketplace-svc/helper/queue"
	"net/http"
	"strings"
	"sync"
	"time"
)

type OrderCreateNotify struct {
	Topic string
	Infra app.Infra
}

func NewOrderCreateNotify(infra app.Infra) IWorker {
	return &OrderCreateNotify{
		Infra: infra,
		Topic: base.TOPIC_ORDER_CREATE_NOTIFIY,
	}
}

func (cp OrderCreateNotify) Cmd() *cli.Command {
	return &cli.Command{
		Name:  cp.Topic,
		Usage: "worker " + cp.Topic + " --indices=1",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "indices", Aliases: []string{"i"}, Value: 1},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("running worker " + cp.Topic + " with indices " + fmt.Sprint(c.Int("indices")))
			return cp.Subscriber(c.Int("indices"))
		},
	}
}

func (cp OrderCreateNotify) Subscriber(indices int) error {

	// publish task
	cp.publishTask()

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

func (cp OrderCreateNotify) HandlerSubscriber(msg *kafka.Message) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in main.go", r)
		}
	}()

	fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	cp.Infra.Log.Info("Message Received: " + string(msg.Value))

	orderCancelKafka := request.OrderCancelKafka{}
	err := json.Unmarshal(msg.Value, &orderCancelKafka)
	if err != nil {
		fmt.Printf("Error decoding JSON from Kafka: %s", err)
	}

	orderCancel := request.OrderCancel{}
	err = json.Unmarshal([]byte(orderCancelKafka.Body), &orderCancel)
	if err != nil {
		fmt.Printf("Error decoding JSON to Struct: %s", err)
	}

	//taskCancelMinutes := orderCancel.PaymentExpiredMinutes

	if orderCancel.QueueType == "notify" {
		go cp.initPushData(cp.Infra.Config.KalcareAPI.PostInterval, 1, []byte(orderCancelKafka.Body), []byte(orderCancel.NotifyPayload))
	}
}

func (cp OrderCreateNotify) publishTask() {
	var (
		token string
		err   error
	)
	cp.Infra.Log.Info("Start processing task")
	cfg := cp.Infra.Config.KalcareAPI

	if cfg.ClientID == "" {
		panic("ClientID is required")
	}

	requestBody, err := json.Marshal(map[string]string{
		"client_id": cfg.ClientID,
	})
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(cfg.Server+cfg.EndpointAuth, "application/json", bytes.NewBuffer(requestBody))
	defer resp.Body.Close()

	var loginResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginResult)

	cp.Infra.Log.Info("RQ: " + cfg.Server + cfg.EndpointAuth)

	token = fmt.Sprintf("%v", loginResult["data"].(map[string]interface{})["record"].(map[string]interface{})["token"])
	client := &http.Client{}

	respQ, err := http.NewRequest("GET", cfg.Server+cfg.EndpointQueue+"s", nil)
	respQ.Header.Add("Authorization", "Bearer "+token)
	reqQ, err := client.Do(respQ)
	defer reqQ.Body.Close()

	var tempQueueResult map[string]interface{}
	json.NewDecoder(reqQ.Body).Decode(&tempQueueResult)

	byteData, _ := json.Marshal(tempQueueResult["data"].(map[string]interface{})["records"])
	orderCancel := request.OrderCancelArr{}
	err = json.Unmarshal(byteData, &orderCancel)
	if err != nil {
		cp.Infra.Log.Info("Error decoding JSON: " + fmt.Sprint(err))
	}
	for _, data := range orderCancel {
		cp.Infra.Log.Info("queueType: " + data.QueueType + "Queue Received: " + data.OrderNo + " " + data.OrderStatus + "Queue Data: " + fmt.Sprint(data))

		valueData, _ := json.Marshal(data)
		var orderCancelKafka = request.OrderCancelKafka{Body: string(valueData), Properties: []string{}, Headers: []string{}}
		kafkaData, _ := json.Marshal(orderCancelKafka)
		err := cp.Infra.KafkaProducer.Publish(base.TOPIC_ORDER_CREATE_NOTIFIY, kafkaData)
		if err != nil {
			fmt.Printf("Failed to publish message: %v\n", err)
		}
	}
}

func (cp OrderCreateNotify) initCancelOrderKafka(data []byte, minutes int) {
	time.AfterFunc(time.Duration(minutes)*time.Minute, func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in initCancelOrderKafka", r)
			}
		}()
		cp.Infra.Log.Info("info Execute Order : " + string(data))

		err := cp.Infra.KafkaProducer.Publish(base.TOPIC_ORDER_CREATE_NOTIFIY, data)
		if err != nil {
			fmt.Printf("Failed to publish message: %v\n", err)
		}
	})
}

func (cp OrderCreateNotify) initPushData(postInterval int, currentInterval int, messageData []byte, payload []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in initPushData", r)
		}
	}()
	cfg := cp.Infra.Config.KalcareAPI

	requestBody, err := json.Marshal(map[string]string{
		"client_id": cfg.ClientID,
	})
	if err != nil {
		cp.Infra.Log.Error(errors.New("Error Marshal Token Auth: " + err.Error()))
		return
	}
	resp, err := http.Post(cfg.Server+cfg.EndpointAuth, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		cp.Infra.Log.Error(errors.New("Error Token Auth: " + err.Error()))
		return
	}

	var loginResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginResult)

	token := fmt.Sprintf("%v", loginResult["data"].(map[string]interface{})["record"].(map[string]interface{})["token"])

	var m map[string]interface{}
	err = json.Unmarshal(messageData, &m)
	notifyCode := fmt.Sprintf("%v", m["notify_code"])
	merchantId := fmt.Sprintf("%v", m["merchant_id"])
	if m["merchant_id"] == nil {
		merchantId = ""
	}

	client := &http.Client{}

	respQ, err := http.NewRequest("GET", cfg.Server+cfg.EndpointWebhook+"s?module=order-create&third_party="+notifyCode+"&merchant_id="+merchantId, nil)
	respQ.Header.Add("Authorization", "Bearer "+token)
	reqQ, err := client.Do(respQ)

	if err != nil {
		return
	}

	var tempRegisteredWebhookResult map[string]interface{}
	json.NewDecoder(reqQ.Body).Decode(&tempRegisteredWebhookResult)

	byteData, _ := json.Marshal(tempRegisteredWebhookResult["data"].(map[string]interface{})["records"])
	registeredWebhook := request.RegisteredWebhookArr{}
	err = json.Unmarshal(byteData, &registeredWebhook)
	if err != nil {
		cp.Infra.Log.Error(errors.New("Error decoding JSON Push Data: " + err.Error()))
		return
	}
	for _, data := range registeredWebhook {
		respP, err := http.NewRequest("POST", data.Url, bytes.NewBuffer(payload))

		respHQ, err := http.NewRequest("GET", fmt.Sprintf("%s/%d/create-header?merchant_id=%s",
			cfg.Server+cfg.EndpointWebhook, data.ID, merchantId), nil)

		respHQ.Header.Add("Authorization", "Bearer "+token)
		reqHQ, err := client.Do(respHQ)
		defer reqHQ.Body.Close()

		var tempHeaderResult map[string]interface{}
		json.NewDecoder(reqHQ.Body).Decode(&tempHeaderResult)

		byteHeaderData, _ := json.Marshal(tempHeaderResult["data"])
		dataHeader := request.AuthorizationData{}
		err = json.Unmarshal(byteHeaderData, &dataHeader)
		if err != nil {
			cp.Infra.Log.Error(errors.New("Error decoding JSON: " + err.Error()))
		}
		for _, head := range dataHeader.Data {
			log.Println(head)
			s := strings.SplitN(head, "=", 2)
			respP.Header.Add(s[0], s[1])
		}

		reqP, err := client.Do(respP)
		if err != nil {
			cp.Infra.Log.Error(errors.New("Error Posting Data: " + err.Error()))
			return
		}

		tempBodyResult, err := ioutil.ReadAll(reqP.Body)
		if err != nil {
			cp.Infra.Log.Error(errors.New("Error Posting Data: " + err.Error()))
			return
		}
		bodyString := string(tempBodyResult)

		statusCode := 200
		statusMessage := "success"
		if reqP.StatusCode == http.StatusOK {
			statusCode = reqP.StatusCode
			if statusCode != 200 {
				statusMessage = "failed"
			}
		} else {
			statusCode = 500
			statusMessage = "failed"
		}

		log.Println(reqP)
		cp.Infra.Log.Info("Posting Data to : " + data.Url + " Header: " + strings.Join(dataHeader.Data, ";") + " Payload: " + string(payload) + " Status Code: " + reqP.Status + " Body: " + bodyString)

		if statusCode != 200 && currentInterval < postInterval {
			cp.Infra.Log.Info(fmt.Sprintf(" Re-attempt Post Order %d/%d", currentInterval+1, postInterval))
			time.AfterFunc(time.Duration(cfg.PostMinutes)*time.Minute, func() {
				go cp.initPushData(postInterval, currentInterval+1, messageData, payload)
			})
			statusMessage = "failed"
		}
		if currentInterval > 0 {
			m["post_order"] = 1
			m["post_order_data"] = string(payload)
			m["post_order_header"] = strings.Join(dataHeader.Data, ";")
			m["post_order_endpoint"] = data.Url
			m["action_type"] = data.Code
			m["action_module"] = data.ModuleIntegration
			m["third_party"] = data.ThirdParty
			m["status_message"] = statusMessage
			m["message_result"] = bodyString
			newMessageData, err := json.Marshal(m)
			respD, err := http.NewRequest("DELETE", cfg.Server+cfg.EndpointQueue, bytes.NewBuffer(newMessageData))
			respD.Header.Add("Authorization", "Bearer "+token)
			respD.Header.Add("Content-Type", "application/json")
			reqD, err := client.Do(respD)
			if err != nil {
				cp.Infra.Log.Error(errors.New("Error Remove Queue: " + err.Error()))
				return
			}
			log.Println(reqD)
			bodyDeleteBytes, err := ioutil.ReadAll(reqD.Body)
			if err != nil {
				cp.Infra.Log.Error(errors.New("Error Remove Queue Resp: " + err.Error()))
				return
			}

			cp.Infra.Log.Info("Delete Temp Data: " + cfg.Server + cfg.EndpointQueue + " Payload: " + string(newMessageData))
			cp.Infra.Log.Info("Delete Status Code: " + reqD.Status + " Body: " + string(bodyDeleteBytes))
		}
	}
}
