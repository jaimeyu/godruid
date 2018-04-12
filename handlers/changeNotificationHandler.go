package handlers

import (
	"encoding/json"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	tenmod "github.com/accedian/adh-gather/models/tenant"

	"github.com/shopify/sarama"
)

/*
	Pushes provisioning changes to other system via Kafka.
	Currently this just sends all monitored object data to Kafka at a regular interval.
*/

const pollingFrequencySecs = 60
const refreshFrequencyMillis = int64(5 * time.Minute / time.Millisecond)
const defaultKafkaTopic = "monitored-object"

type changeNotifier struct {
	Brokers              []string
	Topic                string
	adminDB              *datastore.AdminServiceDatastore
	tenantDB             *datastore.TenantServiceDatastore
	lastSuccessTimestamp int64
	hasErrors            bool
	fullRefresh          bool
}

func SendChangeNotifications() {

	cfg := gather.GetConfig()
	broker := cfg.GetString(gather.CK_kafka_broker.String())
	if len(broker) < 1 {
		logger.Log.Warning("No Kafka broker configured for notifications")
		return
	}

	lastFullRefresh := int64(0)
	lastSuccess := int64(0)
	tenantDB, err := getTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TenantDB: %s", err.Error())
		return
	}
	adminDB, err := getAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminDB: %s", err.Error())
		return
	}

	// Monitor jobs at a regular interval
	ticker := time.NewTicker(pollingFrequencySecs * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			startTime := time.Now().UnixNano() / int64(time.Millisecond)
			needsRefresh := lastFullRefresh <= (startTime - refreshFrequencyMillis)

			notifier := createNotifier(&adminDB, &tenantDB, []string{broker}, defaultKafkaTopic, lastSuccess, needsRefresh)
			notifier.pollChanges()
			if !notifier.hasErrors {
				lastSuccess = startTime
				if needsRefresh {
					lastFullRefresh = startTime
				}
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func createNotifier(adminDB *datastore.AdminServiceDatastore, tenantDB *datastore.TenantServiceDatastore, brokers []string, topic string, lastSuccess int64, needsRefresh bool) *changeNotifier {
	notifier := new(changeNotifier)

	// We use the REST handler as a cheap way to get references to the DBs.
	notifier.adminDB = adminDB
	notifier.tenantDB = tenantDB
	notifier.Topic = topic
	notifier.Brokers = brokers
	notifier.lastSuccessTimestamp = lastSuccess
	notifier.fullRefresh = needsRefresh
	return notifier
}

func (cn *changeNotifier) pollChanges() {

	logger.Log.Infof("pollChanges fullRefresh=%v, lastSuccess=%d", cn.fullRefresh, cn.lastSuccessTimestamp)
	tenants, err := (*cn.adminDB).GetAllTenantDescriptors()
	if err != nil {
		logger.Log.Error("Unable to fetch list of tenants: %s", err.Error())
		cn.hasErrors = true
		return
	}

	if len(tenants) < 1 {
		logger.Log.Warning("No tenants found")
		return
	}

	if cn.fullRefresh {
		logger.Log.Debugf("Performing a full refresh")
	}

	kafkaProducer := cn.newKafkaProducer()
	logger.Log.Debug("Started Kafka Producer")

	for _, t := range tenants {

		logger.Log.Debugf("Fetching Monitored Objects for tenant %s", t.ID)

		monitoredObjects, err := (*cn.tenantDB).GetAllMonitoredObjects(t.ID)
		if err != nil {
			logger.Log.Errorf("Failed to fetch Monitored Objects for tenant %s: %s", t.ID, err.Error())
			cn.hasErrors = true
			continue
		}

		cn.sendMonitoredObjects(kafkaProducer, t.ID, monitoredObjects)
	}

	// Perform a synchronous close. Wait for remaining messages to be sent and close the producer.
	err = kafkaProducer.Close()
	if err != nil {
		cn.hasErrors = true
	}

}

func (cn *changeNotifier) sendMonitoredObjects(kafkaProducer sarama.AsyncProducer, tenantID string, monitoredObjects []*tenmod.MonitoredObject) {

	logger.Log.Debugf("Got %d monitored objects for tenant %s", len(monitoredObjects), tenantID)
	sentCount := 0
	for _, mo := range monitoredObjects {

		if !cn.fullRefresh && mo.CreatedTimestamp < cn.lastSuccessTimestamp && mo.LastModifiedTimestamp < cn.lastSuccessTimestamp {
			// This MO was already sent since it last changed and we aren't doing a full refresh
			continue
		}

		// Workaround for bug where tenantId and id attributes were cleared by UI.
		mo.TenantID = tenantID
		if len(mo.ID) == 0 {
			mo.ID = mo.ObjectName
		}

		// Generate a json payload and send it.
		// Later we can serialized object but right now we don't guarantee the the receiver knows how
		// to deserialize objects.
		b, err := json.Marshal(mo)

		if err != nil {
			logger.Log.Error("Failed to marshal monitored object", err.Error())
			continue
		}

		logger.Log.Debugf("sending %s", mo.ObjectName)

		kafkaProducer.Input() <- &sarama.ProducerMessage{
			Topic: cn.Topic,
			Key:   sarama.StringEncoder(mo.ID),
			Value: sarama.StringEncoder(b),
		}
		sentCount++
	}
	logger.Log.Infof("Sent %d monitored object notifications for tenant %s", sentCount, tenantID)

}

func (cn *changeNotifier) newKafkaProducer() sarama.AsyncProducer {

	logger.Log.Debug("Starting Kafka Producer")

	config := sarama.NewConfig()

	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 100 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(cn.Brokers, config)
	if err != nil {
		logger.Log.Fatalf("Failed to start Kafka producer. %s", err)
	}

	// Create a callback for handling errors
	go func() {
		for err := range producer.Errors() {
			logger.Log.Errorf("Failed to write access log entry: %s", err)
			cn.hasErrors = true
		}
	}()

	return producer
}
