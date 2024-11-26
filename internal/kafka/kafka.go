package kafka

import (
	"sync"
	"wb-l0/internal/kafka/consumer"
	"wb-l0/internal/kafka/producer"
)

type KafkaManager struct {
	Consumer consumer.Consumer
	Producer producer.Producer
}

func NewKafkaManager() *KafkaManager {
	return &KafkaManager{
		Consumer: consumer.NewConsumer(),
		Producer: producer.NewProducer(),
	}
}

func (k *KafkaManager) StartKafkaServices() {
	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(2)

	go func() {
		defer wg.Done()
		k.Consumer.Start(&ch)
	}()

	go func() {
		defer wg.Done()
		<-ch
		k.Producer.Start()
	}()

	wg.Wait()
}
