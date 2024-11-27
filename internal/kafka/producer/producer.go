package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"sync"
	"wb-l0/config"
	"wb-l0/internal/models/order"
)

type Producer interface {
	Start(cfg *config.Kafka, wg *sync.WaitGroup)
}

type KafkaProducer struct{}

func NewProducer() *KafkaProducer {
	return &KafkaProducer{}
}

func prepareJSON() ([]byte, error) {
	data, err := os.ReadFile("data/order.json")
	if err != nil {
		return nil, err
	}

	var orderData order.Order
	err = json.Unmarshal(data, &orderData)
	if err != nil {
		return nil, err
	}

	messageBytes, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}
	return messageBytes, nil
}

func (kp *KafkaProducer) Start(cfg *config.Kafka, wg *sync.WaitGroup) {
	writer := kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	log.Println("Producer is started.")
	wg.Done()

	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("error while closing reader: %v\n", err)
		}
	}()

	for i := 1; i <= 5; i++ {
		msg, err := prepareJSON()
		if err != nil {
			log.Println(err)
			continue
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(fmt.Sprintf("Key %d", i)),
			Value: msg,
		})

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Sent message to [%d]: %s\n", i, msg)
	}
}
