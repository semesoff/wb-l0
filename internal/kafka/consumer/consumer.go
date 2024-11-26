package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"wb-l0/config"
	"wb-l0/internal/db/db"
)

type Consumer interface {
	Start(ch *chan struct{}, cfg *config.Kafka)
}

type KafkaConsumer struct {
	db *db.DatabaseProvider
}

func NewConsumer(db *db.DatabaseProvider) *KafkaConsumer {
	kc := &KafkaConsumer{db: db}
	return kc
}

func (kc *KafkaConsumer) Start(ch *chan struct{}, cfg *config.Kafka) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Topic:    cfg.Topic,
		GroupID:  "order-consumers",
		MaxBytes: 1024 * 1024 * 2,
	})

	close(*ch)
	log.Println("Consumer is started.")

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error while closing reader: %v\n", err)
		}
	}()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		fmt.Printf("Received message: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
	}
}
