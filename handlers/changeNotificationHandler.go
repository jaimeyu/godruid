package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"

	"github.com/segmentio/kafka-go"
)

type EventType int16

const (
	MonitoredObjectCreated = 0
	MonitoredObjectDeleted = 1
	MonitoredObjectUpdated = 2
	DomainCreated          = 10
	DomainUpdated          = 11
	DomainDeleted          = 12
)

type ChangeEvent struct {
	eventType EventType
	tenantID  string
	payload   interface{}
}

const defaultPollingFrequency = 45 * time.Second // How often to poll tenantDB for recent changes
//const refreshFrequencyMillis = int64(gather. * time.Second / time.Millisecond) // How often to push a full refresh of tenantDB
const defaultKafkaTopic = "monitored-object" // The topic where changes are pushed.

/*
 The ChangeNotificationHandler is the entry point for handling changes to provisioning resources.
 Provisioning changes are either pushed to the ChangeNotificationHandler from provisioning workflows (i.e. API calls) or
 they are polled from the tenantDB.
 The polling mechanism serves more of a backup mechanism in situations where pushed events are lost and gives the opportunity
 for subscribers of events to perform a full resync of their systems with the tenantDB.
*/
type ChangeNotificationHandler struct {
	provisioningEvents chan *ChangeEvent
	brokers            []string
	topic              string
	adminDB            *datastore.AdminServiceDatastore
	tenantDB           *datastore.TenantServiceDatastore
	metricsDB          datastore.DruidDatastore
	batchSize          int64
	// To block pollChanges from overlapping
	locker uint32
}

// ChangeNotificationHandler singleton
var changeNotifH ChangeNotificationHandler

func getChangeNotificationHandler() *ChangeNotificationHandler {
	return &changeNotifH
}

func CreateChangeNotificationHandler() *ChangeNotificationHandler {

	cfg := gather.GetConfig()
	broker := cfg.GetString(gather.CK_kafka_broker.String())
	if len(broker) < 1 {
		logger.Log.Warning("No Kafka broker configured for notifications")
		return nil
	}

	tenantDB, err := GetTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TenantDB: %s", err.Error())
		return nil
	}
	adminDB, err := GetAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminDB: %s", err.Error())
		return nil
	}

	batchSize := int64(1000)
	cfgBatchSize := cfg.GetInt(gather.CK_server_datastore_batchsize.String())
	if cfgBatchSize > 0 {
		batchSize = int64(cfgBatchSize)
	}

	changeNotifH = ChangeNotificationHandler{
		brokers:            []string{broker},
		topic:              defaultKafkaTopic,
		tenantDB:           &tenantDB,
		adminDB:            &adminDB,
		provisioningEvents: make(chan *ChangeEvent, 20),
		metricsDB:          druid.NewDruidDatasctoreClient(),
		batchSize:          batchSize,
		locker:             0,
	}

	//	go changeNotifH.readFromKafka(broker, defaultKafkaTopic)

	return &changeNotifH
}

/*
The main loop
*/
func (c *ChangeNotificationHandler) SendChangeNotifications() {

	lastFullRefresh := time.Time{}
	lastSuccess := time.Time{}
	refreshFrequency := (time.Duration)(gather.GetConfig().GetInt(gather.CK_server_changenotif_refreshFreqSeconds.String())) * time.Second
	pollingFrequency := defaultPollingFrequency
	if refreshFrequency < pollingFrequency {
		pollingFrequency = refreshFrequency
	}
	// Run an auditer to do a refresh at regular intervals
	ticker := time.NewTicker(pollingFrequency)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:

			// Time to run the audit to push changes we may have missed through the channel.
			// If needsRefresh is false, just push changes detected since last push; otherwise
			// push all provisioning data that others are interested in.
			// Note: right now, this is a synchronous operation. If needed it could be handled in
			// a separate dedicated thread.
			startTime := time.Now().Truncate(time.Second)
			needsRefresh := !lastFullRefresh.Add(refreshFrequency).After(startTime)
			if err := c.pollChanges(lastSuccess.UnixNano()/int64(1000), needsRefresh); err == nil {
				lastSuccess = startTime
				if needsRefresh {
					lastFullRefresh = startTime
				}
			}

		case event := <-c.provisioningEvents:
			// Something changed, lets batch the events if we can.  This helps to
			// reduce the number of updates in the metrics DB to update metadata.
			logger.Log.Debugf("Received a changed notification %v", event)
			bufferedEvents := []*ChangeEvent{}
			bufferedEvents = append(bufferedEvents, event)
			t := time.After(5 * time.Second)
			buffering := true
			for buffering {
				select {
				case event := <-c.provisioningEvents:
					bufferedEvents = append(bufferedEvents, event)
				case <-t:
					buffering = false
				}
			}
			c.processEvents(bufferedEvents)

		case <-quit:
			ticker.Stop()
			return

		}
	}
}

func NotifyMonitoredObjectCreated(tenantID string, obj ...*tenmod.MonitoredObject) {
	NotifyEvent(&ChangeEvent{
		eventType: MonitoredObjectCreated,
		tenantID:  tenantID,
		payload:   obj,
	})
}

func NotifyMonitoredObjectUpdated(tenantID string, obj ...*tenmod.MonitoredObject) {
	NotifyEvent(&ChangeEvent{
		eventType: MonitoredObjectUpdated,
		tenantID:  tenantID,
		payload:   obj,
	})
}

func NotifyMonitoredObjectDeleted(tenantID string, obj ...*tenmod.MonitoredObject) {
	NotifyEvent(&ChangeEvent{
		eventType: MonitoredObjectDeleted,
		tenantID:  tenantID,
		payload:   obj,
	})
}

func NotifyDomainCreated(tenantID string, obj ...*tenmod.Domain) {
	NotifyEvent(&ChangeEvent{
		eventType: DomainCreated,
		tenantID:  tenantID,
		payload:   obj,
	})
}

func NotifyDomainUpdated(tenantID string, obj ...*tenmod.Domain) {
	NotifyEvent(&ChangeEvent{
		eventType: DomainUpdated,
		tenantID:  tenantID,
		payload:   obj,
	})
}

func NotifyDomainDeleted(tenantID string, obj ...*tenmod.Domain) {
	NotifyEvent(&ChangeEvent{
		eventType: DomainDeleted,
		tenantID:  tenantID,
		payload:   obj,
	})
}

func NotifyEvent(event *ChangeEvent) {
	// Found a potential bug if change notification is disabled
	// and someone tries to send an event, the goroutine will lock and never exit.
	// This is problem for dev loads rather than production.
	if changeNotifH.provisioningEvents != nil {
		changeNotifH.provisioningEvents <- event
	}
}

func (c *ChangeNotificationHandler) processEvents(events []*ChangeEvent) {
	processedTenantIds := make(map[string]bool)
	for _, e := range events {

		metadataChange := false
		switch e.eventType {
		case MonitoredObjectCreated, MonitoredObjectUpdated:
			c.sendToKafka(e.tenantID, e.payload.([]*tenmod.MonitoredObject))
			metadataChange = true
		case MonitoredObjectDeleted, DomainCreated, DomainUpdated, DomainDeleted:
			metadataChange = true
		}

		if metadataChange {
			// Currently any metadataChange is handled by resynchronizing all metadata for the tenant so we
			// don't really care what the nature of the change was.
			// This approach is nice and simple but is also effectively similar to dropping a table
			// and re-populating it. If it becomes inefficient we'll have to update this to do
			// more of a CRUD approach to the metadata.
			if _, ok := processedTenantIds[e.tenantID]; !ok {
				c.updateMetricsDatastoreMetadata(e.tenantID)
				processedTenantIds[e.tenantID] = true
			}
		}
	}
}

func (c *ChangeNotificationHandler) sendToKafka(tenantID string, monitoredObjects []*tenmod.MonitoredObject) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      c.brokers,
		Topic:        c.topic,
		RequiredAcks: 0,
		Async:        true,
		Balancer:     &kafka.LeastBytes{},
	})
	defer func() {
		logger.Log.Info("closing kafka producer")
		w.Close()
	}()

	sendMonitoredObjects(w, tenantID, monitoredObjects)

}

func debugAddFakeMonitoredObjects() []*tenmod.MonitoredObject {
	var monitoredObjects []*tenmod.MonitoredObject
	// For metadata, we need to build a list of known qualifiers
	logger.Log.Infof("Dumping Updating metadata from poll change")
	//debugging
	//colors := []string{"black", "white", "orange", "blue", "green", "red", "purple", "gold", "yellow", "brown", "aqua"}

	testNodes := 50000 // change this for testing!
	for i := 0; i < testNodes; i++ {
		mo := tenmod.MonitoredObject{
			ID:                fmt.Sprintf("debug_%d", i),
			ObjectName:        fmt.Sprintf("debug_%d", i),
			MonitoredObjectID: fmt.Sprintf("debug_%d", i),
			// Meta:              map[string]string{"colors": colors[i%len(colors)]},
			Meta:             map[string]string{"colors": "paris"},
			CreatedTimestamp: 10,
		}

		monitoredObjects = append(monitoredObjects, &mo)
	}
	return monitoredObjects
}

// Obsolete!
func (c *ChangeNotificationHandler) updateMetricsDatastoreMetadata(tenantID string) {

}

// getAllMonitoredObjects - uses the paginated DB call to acquire all monitored objects
func (c *ChangeNotificationHandler) getAllMonitoredObjects(tenantID string) ([]*tenmod.MonitoredObject, error) {

	result := make([]*tenmod.MonitoredObject, 0)
	startKey := ""

	for true {
		monitoredObjects, paginationOffsets, err := (*c.tenantDB).GetAllMonitoredObjectsByPage(tenantID, startKey, c.batchSize)
		if err != nil {
			return nil, err
		}

		result = append(result, monitoredObjects...)

		if len(paginationOffsets.Next) == 0 {
			break
		}

		startKey = paginationOffsets.Next
	}

	return result, nil
}

func (c *ChangeNotificationHandler) pollChanges(lastSyncTimestamp int64, fullRefresh bool) error {
	// Avoid running overlapping pollChanges
	if !atomic.CompareAndSwapUint32(&c.locker, 0, 1) {
		return nil
	}
	defer atomic.StoreUint32(&c.locker, 0)

	startTime := time.Now()

	logger.Log.Debugf("pollChanges fullRefresh=%v, lastSuccess=%d", fullRefresh, lastSyncTimestamp)
	tenants, err := (*c.adminDB).GetAllTenantDescriptors()
	if err != nil {
		logger.Log.Error("Unable to fetch list of tenants: %s", err.Error())
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, "500", mon.PollChanges)

		return err
	}

	if len(tenants) < 1 {
		logger.Log.Warning("No tenants found")
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, "500", mon.PollChanges)
		return nil
	}

	kafkaProducer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      c.brokers,
		Topic:        c.topic,
		RequiredAcks: 0,
		Async:        true,
		Balancer:     &kafka.LeastBytes{},
	})
	defer func() {
		kafkaProducer.Close()
	}()

	logger.Log.Debug("Started Kafka Producer")

	var lastError error
	for _, t := range tenants {

		// changeDetected := false

		monitoredObjects, err := c.getAllMonitoredObjects(t.ID)
		if err != nil {
			logger.Log.Warningf("Failed to fetch Monitored Objects for tenant %s: %s", t.ID, err.Error())
			continue
		}

		// Enable this to add arbitary number of items into the druid look ups
		//monitoredObjects = debugAddFakeMonitoredObjects()

		// Update counters
		setMonitoredObjectCount(len(monitoredObjects))

		if fullRefresh {
			sendMonitoredObjects(kafkaProducer, t.ID, monitoredObjects)
		} else {
			//TODO at a later time we could use change notification mechanism from DB rather than query all
			for _, mo := range monitoredObjects {
				if mo.CreatedTimestamp > lastSyncTimestamp || mo.LastModifiedTimestamp > lastSyncTimestamp {
					// changeDetected = true
					sendMonitoredObject(kafkaProducer, mo)
				}
			}

		}
		// Obsolete
		// if fullRefresh || changeDetected {

		// 	// For metadata, we need to build a list of known qualifiers
		// 	logger.Log.Infof("Dumping Updating metadata from poll change")

		// 	if err = c.metricsDB.AddMonitoredObjectToLookup(t.ID, monitoredObjects, "meta"); err != nil {
		// 		logger.Log.Errorf("Failed to update metrics metadata for tenant %s: %s", t.ID, err.Error())
		// 		lastError = err
		// 		continue
		// 	} else {
		// 		logger.Log.Infof("Updated metadata in metric DB for tenant %s", t.ID)
		// 	}

		// }

	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, "200", mon.PollChanges)

	return lastError
}

func sendMonitoredObjects(writer *kafka.Writer, tenantID string, monitoredObjects []*tenmod.MonitoredObject) {

	logger.Log.Debugf("Sending %d monitored objects to kafka for tenant %s", len(monitoredObjects), tenantID)
	sentCount := 0
	for _, mo := range monitoredObjects {

		// Workaround for bug where tenantId and id attributes were cleared by UI.
		mo.TenantID = tenantID
		if len(mo.ID) == 0 {
			mo.ID = mo.ObjectName
		}

		sendMonitoredObject(writer, mo)
		sentCount++
	}
	logger.Log.Infof("Sent %d monitored object to kafka for tenant %s ", sentCount, tenantID)

}

func sendMonitoredObject(writer *kafka.Writer, monitoredObject *tenmod.MonitoredObject) {
	// Generate a json payload and send it.
	// Later we can serialized object but right now we don't guarantee the the receiver knows how
	// to deserialize objects.
	b, err := json.Marshal(monitoredObject)

	if err != nil {
		logger.Log.Error("Failed to marshal monitored object", err.Error())
		return
	}

	writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(monitoredObject.ID),
		Value: []byte(b),
	})
}
