package messaging

import (
	"context"
	"fmt"

	"github.com/accedian/adh-gather/logger"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader

	topicName string
}

func CreateKafkaReader(topicName string) *KafkaConsumer {
	result := &KafkaConsumer{}

	result.topicName = topicName

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   result.topicName,

		Partition: 0,
	})

	result.reader = r

	return result
}

func (c *KafkaConsumer) Destroy() {
	if err := c.reader.Close(); err != nil {
		logger.Log.Errorf("Unable to close Kafka Consumer: %s", err.Error())
	}
}

func (c *KafkaConsumer) ReadMessage(action func([]byte) bool) ([]byte, error) {
	ctx := context.Background()
	m, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

	if action != nil {
		if action(m.Value) {

			c.reader.CommitMessages(ctx, m)

			logger.Log.Debugf("Successfully read message and completed action for: %s", string(m.Value))
			return m.Value, nil
		}

		logger.Log.Debugf("Successfully read message but could not complete action for: %s", string(m.Value))
		return nil, fmt.Errorf("Unable to complete action required by: %s", string(m.Value))
	}

	c.reader.CommitMessages(ctx, m)

	logger.Log.Debugf("Successfully read message: %s", string(m.Value))
	return m.Value, nil
}
