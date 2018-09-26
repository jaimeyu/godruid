package messaging

import (
	"context"

	"github.com/accedian/adh-gather/gather"

	"github.com/accedian/adh-gather/logger"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer

	topicName string
}

func CreateKafkaWriter(topicName string) *KafkaProducer {
	cfg := gather.GetConfig()

	result := &KafkaProducer{}

	result.topicName = topicName

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{cfg.GetString(gather.CK_kafka_broker.String())},
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	})

	result.writer = w

	return result
}

func (p *KafkaProducer) Destroy() {
	if err := p.writer.Close(); err != nil {
		logger.Log.Errorf("Unable to close Kafka Producer: %s", err.Error())
	}
}

func (p *KafkaProducer) WriteMessage(messageKey string, messageBytes []byte) error {

	if err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(messageKey),
			Value: messageBytes,
		},
	); err != nil {
		return err
	}

	return nil
}
