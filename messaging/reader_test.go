package messaging_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/messaging"
	kafka "github.com/segmentio/kafka-go"
)

func TestReader(t *testing.T) {
	kafkaConsumer := messaging.CreateKafkaReader("colt-mef")

	go func() {
		for {
			kafkaConsumer.ReadMessage(func(stuff []byte) bool {
				logger.Log.Debugf("JUST HERE IN THE ACTION: %s", string(stuff))
				return true
			})
		}
	}()

	// make a writer that produces to topic-A, using the least-bytes distribution
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "colt-mef",
		Balancer: &kafka.LeastBytes{},
	})

	for i := 0; i < 5; i++ {
		index := strconv.FormatInt(int64(i), 10)
		w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("Key-A" + index),
				Value: []byte(fmt.Sprintf("Hello World!: %s%s%s%s", index, index, index, index)),
			},
		)
	}

	w.Close()

	time.Sleep(time.Second * 5)

	kafkaConsumer.Destroy()
}
