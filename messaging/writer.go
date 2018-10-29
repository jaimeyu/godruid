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

	w.Stats()

	logger.Log.Debugf("Created kafka writer for topic %s", topicName)

	return result
}

func (p *KafkaProducer) Destroy() {
	if err := p.writer.Close(); err != nil {
		logger.Log.Errorf("Unable to close Kafka Producer on topic %s: %s", p.topicName, err.Error())
	}
}

func (p *KafkaProducer) WriteMessage(messageKey string, messageBytes []byte) error {

	defer logger.Log.Debugf("Dumping write stats for topic %s: %+v", p.topicName, p.writer.Stats())

	logger.Log.Debugf("Writing message to topic %s with key %s", p.topicName, messageKey)

	if err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(messageKey),
			Value: messageBytes,
		},
	); err != nil {
		logger.Log.Errorf("Message could not be written to topic %s with key %s due to: %s", p.topicName, messageKey, err)
		return err
	}

	logger.Log.Debugf("Message successfully written message %s to topic %s with key %s", string(messageBytes), p.topicName, messageKey)

	return nil
}
