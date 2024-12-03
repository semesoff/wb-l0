package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"wb-l0/config"
	"wb-l0/internal/cache"
	"wb-l0/internal/db/db"
	"wb-l0/internal/utils"
)

type Consumer interface {
	Start(ch *chan struct{}, cfg *config.Kafka, wg *sync.WaitGroup)
}

type KafkaConsumer struct {
	r  *cache.RedisProvider
	db *db.DatabaseProvider
}

func NewConsumer(db *db.DatabaseProvider, r *cache.RedisProvider) *KafkaConsumer {
	kc := &KafkaConsumer{db: db, r: r}
	return kc
}

func (kc *KafkaConsumer) Start(ch *chan struct{}, cfg *config.Kafka, wg *sync.WaitGroup) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Topic:    cfg.Topic,
		GroupID:  "order-consumers",
		MaxBytes: 1024 * 1024 * 2,
	})
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {

		}
	}(reader)

	close(*ch)
	log.Println("Consumer is started.")
	wg.Done()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		log.Printf("Received message: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
		if orderData, err := utils.EncodeMessage(msg.Value); err == nil {
			if err := (*kc.db).AddOrder(orderData, msg.Value); err != nil {
				log.Printf("error while adding order: %v\n", err)
			} else {
				log.Printf("message with key %s is added to database\n", msg.Key)
			}
			if err := (*kc.r).SetCache(orderData.OrderUid, msg.Value); err != nil {
				log.Printf("error while adding order to cache: %v\n", err)
			} else {
				log.Printf("message with key %s is added to cache\n", msg.Key)
			}
		} else {
			log.Printf("error while encoding message: %v\n", err)
		}
	}
}
