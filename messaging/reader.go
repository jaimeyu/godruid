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
	return CreateKafkaReaderWithSyncTime(topicName, groupTag, time.Second)
}

func CreateKafkaReaderWithSyncTime(topicName string, groupTag string, syncTimeInterval time.Duration) *KafkaConsumer {
	cfg := gather.GetConfig()
	result := &KafkaConsumer{}

	result.topicName = topicName

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.GetString(gather.CK_kafka_broker.String())},
		Topic:          result.topicName,
		GroupID:        result.topicName + "-" + groupTag,
		CommitInterval: syncTimeInterval,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
	})

	logger.Log.Debugf("Created kafka reader for topic %s", topicName)

	result.reader = r

	return result
}

func (c *KafkaConsumer) Destroy() {
	if err := c.reader.Close(); err != nil {
		logger.Log.Errorf("Unable to close Kafka Consumer for topic: %s", c.topicName, err.Error())
	}
}

// ReadMessage - used to read a message when an explicit action must be completed successfully before committing the offset for the message read.
func (c *KafkaConsumer) ReadMessage(action func([]byte) bool) ([]byte, error) {
	ctx := context.Background()
	m, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("message for topic %s at offset %d: %s = %s\n", c.topicName, m.Offset, string(m.Key), string(m.Value))

	if action != nil {
		if action(m.Value) {

			err = c.reader.CommitMessages(ctx, m)

			if err != nil {
				logger.Log.Errorf("Error occured while committing on topic %s message %s: %s", c.topicName, string(m.Value), err)
			}

			logger.Log.Debugf("Successfully read message and completed action on topic %s for: %s", c.topicName, string(m.Value))
			return m.Value, nil
		}

		logger.Log.Debugf("Successfully read message on topic %s but could not complete action for: %s", c.topicName, string(m.Value))
		return nil, fmt.Errorf("Unable to complete action on topic %s required by: %s", c.topicName, string(m.Value))
	}

	err = c.reader.CommitMessages(ctx, m)

	if err != nil {
		logger.Log.Errorf("Error occured while committing on topic %s message %s: %s", c.topicName, string(m.Value), err)
	}

	logger.Log.Debugf("Successfully read message on topic %s: %s", c.topicName, string(m.Value))
	return m.Value, nil
}

// ReadMessageWithoutExplicitOffsetManagement - used when you just want to read a message from a kafka topic without concern for managing the offset
func (c *KafkaConsumer) ReadMessageWithoutExplicitOffsetManagement() ([]byte, error) {
	ctx := context.Background()
	m, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("Successfully read message on topic %s: %s", c.topicName, string(m.Value))
	return m.Value, nil
}
