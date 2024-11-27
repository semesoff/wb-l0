package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"wb-l0/config"
	"wb-l0/internal/cache"
	"wb-l0/internal/db/db"
	"wb-l0/internal/models/order"
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

	close(*ch)
	log.Println("Consumer is started.")
	wg.Done()

	//defer func() {
	//	if err := reader.Close(); err != nil {
	//		log.Printf("error while closing reader: %v\n", err)
	//	}
	//}()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		log.Printf("Received message: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
		if orderData, err := encodeMessage(msg); err == nil {
			if err := (*kc.db).AddOrder(orderData); err != nil {
				log.Printf("error while adding order: %v\n", err)
			} else {
				log.Printf("message with key %s is added to database\n", msg.Key)
			}
			if dataBytes, err := prepareJSON(orderData); err != nil {
				log.Printf("error while adding order to cache: %v\n", err)
			} else {
				if err := (*kc.r).AddCache(orderData.OrderUid, dataBytes); err != nil {
					log.Printf("error while adding order to cache: %v\n", err)
				} else {
					log.Printf("message with key %s is added to cache\n", msg.Key)
				}
			}
		} else {
			log.Printf("error while encoding message: %v\n", err)
		}
	}
}

func encodeMessage(msg kafka.Message) (order.Order, error) {
	var orderData order.Order
	err := json.Unmarshal(msg.Value, &orderData)
	if err != nil {
		return order.Order{}, fmt.Errorf("error while unmarshalling message: %v", err)
	}
	return orderData, nil
}

func prepareJSON(orderData order.Order) ([]byte, error) {
	messageBytes, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}
	return messageBytes, nil
}
