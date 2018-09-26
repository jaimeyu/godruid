package messaging_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/messaging"
	"github.com/spf13/viper"
)

func TestReader(t *testing.T) {
	gather.LoadConfig("../config/adh-gather-test.yml", viper.New())

	kafkaConsumer := messaging.CreateKafkaReader("colt-mef")
	kafkaProducer := messaging.CreateKafkaWriter("colt-mef")

	go func() {
		for {
			kafkaConsumer.ReadMessage(func(stuff []byte) bool {
				logger.Log.Debugf("JUST HERE IN THE ACTION: %s", string(stuff))
				return false
			})
		}
	}()

	for i := 0; i < 5; i++ {
		index := strconv.FormatInt(int64(i), 10)
		kafkaProducer.WriteMessage("Key-A"+index, []byte(`{"service_id": "80033646","action": "INCREASE_BANDWIDTH","bandwidth_change": 200}`))
	}

	time.Sleep(time.Second * 5)

	kafkaProducer.Destroy()
	kafkaConsumer.Destroy()
}
