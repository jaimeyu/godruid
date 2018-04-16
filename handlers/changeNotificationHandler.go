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

var changeNotifH ChangeNotificationHandler

type ChangeNotificationHandler struct {
	monitoredObjectChanges chan []*tenmod.MonitoredObject
	brokers                []string
	topic                  string
	adminDB                *datastore.AdminServiceDatastore
	tenantDB               *datastore.TenantServiceDatastore
}

func getChangeNotificationHandler() *ChangeNotificationHandler {
	return &changeNotifH
}

func CreateChangeNotificationHandler() *ChangeNotificationHandler {
	changeNotifH = ChangeNotificationHandler{}

	cfg := gather.GetConfig()
	broker := cfg.GetString(gather.CK_kafka_broker.String())
	if len(broker) < 1 {
		logger.Log.Warning("No Kafka broker configured for notifications")
		return nil
	}
	changeNotifH.brokers = []string{broker}
	changeNotifH.topic = defaultKafkaTopic

	tenantDB, err := getTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TenantDB: %s", err.Error())
		return nil
	}
	changeNotifH.tenantDB = &tenantDB

	adminDB, err := getAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminDB: %s", err.Error())
		return nil
	}
	changeNotifH.adminDB = &adminDB

	changeNotifH.monitoredObjectChanges = make(chan []*tenmod.MonitoredObject)

	return &changeNotifH
}

type changeNotifier struct {
	brokers              []string
	topic                string
	adminDB              *datastore.AdminServiceDatastore
	tenantDB             *datastore.TenantServiceDatastore
	lastSuccessTimestamp int64
	hasErrors            bool
	fullRefresh          bool
}

func (c *ChangeNotificationHandler) SendChangeNotifications() {

	lastFullRefresh := int64(0)
	lastSuccess := int64(0)

	// Run an auditer to do a refresh at regular intervals
	ticker := time.NewTicker(pollingFrequencySecs * time.Second)
	quit := make(chan struct{})

	for {
		var mo []*tenmod.MonitoredObject
		select {
		case <-ticker.C:

			// Time to run the audit to push changes we may have missed through the channel.
			startTime := time.Now().UnixNano() / int64(time.Millisecond)
			needsRefresh := lastFullRefresh <= (startTime - refreshFrequencyMillis)

			notifier := c.createNotifier(lastSuccess, needsRefresh)
			notifier.pollChanges()
			if !notifier.hasErrors {
				lastSuccess = startTime
				if needsRefresh {
					lastFullRefresh = startTime
				}
			}

		case mo = <-c.monitoredObjectChanges:
			// A monitoredObject was changed. Push to kafka
			logger.Log.Debugf("Received a changed notification for %v", mo)
			c.sendToKafkaAsync(mo)

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (c *ChangeNotificationHandler) sendToKafkaAsync(monitoredObjects []*tenmod.MonitoredObject) {
	producer := newAsyncKafkaProducer(c.brokers)
	// Create a callback for handling errors
	go func() {
		for err := range producer.Errors() {
			logger.Log.Errorf("Failed to write monitored object: %s", err)
		}
	}()

	sendMonitoredObjects(producer, c.topic, monitoredObjects)

	producer.AsyncClose()
}

func (c *ChangeNotificationHandler) createNotifier(lastSuccess int64, needsRefresh bool) *changeNotifier {
	notifier := new(changeNotifier)
	notifier.adminDB = c.adminDB
	notifier.tenantDB = c.tenantDB
	notifier.topic = c.topic
	notifier.brokers = c.brokers
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

	kafkaProducer := newAsyncKafkaProducer(cn.brokers)
	// Create a callback for handling errors
	go func() {
		for err := range kafkaProducer.Errors() {
			logger.Log.Errorf("Failed to write monitored object: %s", err)
			cn.hasErrors = true
		}
	}()

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

		sendMonitoredObject(kafkaProducer, cn.topic, mo)
		sentCount++
	}
	logger.Log.Infof("Sent %d monitored object notifications for tenant %s", sentCount, tenantID)

}

func sendMonitoredObjects(kafkaProducer sarama.AsyncProducer, topic string, monitoredObjects []*tenmod.MonitoredObject) {
	for _, mo := range monitoredObjects {
		if mo == nil {
			continue
		}
		sendMonitoredObject(kafkaProducer, topic, mo)
	}

}

func sendMonitoredObject(kafkaProducer sarama.AsyncProducer, topic string, monitoredObject *tenmod.MonitoredObject) {
	// Generate a json payload and send it.
	// Later we can serialized object but right now we don't guarantee the the receiver knows how
	// to deserialize objects.
	b, err := json.Marshal(monitoredObject)

	if err != nil {
		logger.Log.Error("Failed to marshal monitored object", err.Error())
		return
	}

	logger.Log.Debugf("sending %s", monitoredObject.ObjectName)

	kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(monitoredObject.ID),
		Value: sarama.StringEncoder(b),
	}
}

func newAsyncKafkaProducer(brokers []string) sarama.AsyncProducer {

	logger.Log.Debug("Starting Kafka Producer")

	config := sarama.NewConfig()

	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 100 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		logger.Log.Fatalf("Failed to start Kafka producer. %s", err)
	}

	return producer
}
