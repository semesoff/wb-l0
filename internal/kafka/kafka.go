package kafka

import (
	"sync"
	"wb-l0/config"
	"wb-l0/internal/db/db"
	"wb-l0/internal/kafka/consumer"
	"wb-l0/internal/kafka/producer"
)

type KafkaManager struct {
	Consumer consumer.Consumer
	Producer producer.Producer
}

func NewKafkaManager(db *db.DatabaseProvider) *KafkaManager {
	return &KafkaManager{
		Consumer: consumer.NewConsumer(db),
		Producer: producer.NewProducer(),
	}
}

func (k *KafkaManager) StartKafkaServices(cfg *config.Kafka) {
	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(2)

	go func() {
		defer wg.Done()
		k.Consumer.Start(&ch, cfg)
	}()

	go func() {
		defer wg.Done()
		<-ch
		k.Producer.Start(cfg)
	}()

	wg.Wait()
}
