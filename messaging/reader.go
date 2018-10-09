package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/accedian/adh-gather/gather"

	"github.com/accedian/adh-gather/logger"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader

	topicName string
}

func CreateKafkaReader(topicName string, groupTag string) *KafkaConsumer {
	cfg := gather.GetConfig()
	result := &KafkaConsumer{}

	result.topicName = topicName

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.GetString(gather.CK_kafka_broker.String())},
		Topic:          result.topicName,
		GroupID:        result.topicName + "-" + groupTag,
		Partition:      0,
		CommitInterval: time.Second,
	})

	result.reader = r

	return result
}

func (c *KafkaConsumer) Destroy() {
	if err := c.reader.Close(); err != nil {
		logger.Log.Errorf("Unable to close Kafka Consumer for topic: %s", c.topicName, err.Error())
	}
}

func (c *KafkaConsumer) ReadMessage(action func([]byte) bool) ([]byte, error) {
	ctx := context.Background()
	m, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("message for topic %s at offset %d: %s = %s\n", c.topicName, m.Offset, string(m.Key), string(m.Value))

	if action != nil {
		if action(m.Value) {

			c.reader.CommitMessages(ctx, m)

			if err != nil {
				logger.Log.Errorf("Error occured while committing on topic %s message %s: %s", c.topicName, string(m.Value), err)
			}

			logger.Log.Debugf("Successfully read message and completed action on topic %s for: %s", c.topicName, string(m.Value))
			return m.Value, nil
		}

		logger.Log.Debugf("Successfully read message on topic %s but could not complete action for: %s", c.topicName, string(m.Value))
		return nil, fmt.Errorf("Unable to complete action on topic %s required by: %s", c.topicName, string(m.Value))
	}

	c.reader.CommitMessages(ctx, m)

	logger.Log.Debugf("Successfully read message on topic %s: %s", c.topicName, string(m.Value))
	return m.Value, nil
}
