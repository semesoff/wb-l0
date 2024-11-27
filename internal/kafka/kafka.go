package kafka

import (
	"sync"
	"wb-l0/config"
	"wb-l0/internal/cache"
	"wb-l0/internal/db/db"
	"wb-l0/internal/kafka/consumer"
	"wb-l0/internal/kafka/producer"
)

type KafkaManager struct {
	Consumer consumer.Consumer
	Producer producer.Producer
}

func NewKafkaManager(db *db.DatabaseProvider, r *cache.RedisProvider) *KafkaManager {
	return &KafkaManager{
		Consumer: consumer.NewConsumer(db, r),
		Producer: producer.NewProducer(),
	}
}

func (k *KafkaManager) StartKafkaServices(cfg *config.Kafka) {
	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(2)

	go func() {
		// wg done in consumer start
		k.Consumer.Start(&ch, cfg, &wg)
	}()

	go func() {
		// wg done in producer start
		<-ch
		k.Producer.Start(cfg, &wg)
	}()

	wg.Wait()
}
