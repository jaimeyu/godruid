package messaging_test

// NOTE: this test was just used for debugging purposes and can't be used in a test suite until we change the
// CI build to include a running kafka instance as well as changing the test to not depend on long sleep sessions.
// The test is being left in for debugging purposes only if issues are encoutered witht eh Colt MEF demo
// func TestReader(t *testing.T) {
// 	gather.LoadConfig("../config/adh-gather-test.yml", viper.New())

// 	kafkaConsumer := messaging.CreateKafkaReader("colt-mef", "0")
// 	kafkaProducer := messaging.CreateKafkaWriter("colt-mef")

// 	go func() {
// 		for {
// 			kafkaConsumer.ReadMessage(func(stuff []byte) bool {
// 				logger.Log.Debugf("JUST HERE IN THE ACTION: %s", string(stuff))
// 				return true
// 			})
// 		}
// 	}()

// 	for i := 0; i < 2; i++ {
// 		index := strconv.FormatInt(int64(i), 10)
// 		kafkaProducer.WriteMessage("Key-A"+index, []byte(`{"service_id": "80033646","action": "INCREASE_BANDWIDTH","bandwidth_change": 200}`))
// 	}

// 	time.Sleep(time.Second * 20)

// 	kafkaProducer.Destroy()
// 	kafkaConsumer.Destroy()
// }
