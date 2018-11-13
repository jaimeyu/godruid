package datastore

import (
	"fmt"
	"strings"
	"time"

	pb "github.com/accedian/adh-gather/gathergrpc"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	uuid "github.com/satori/go.uuid"
)

var (
	NotFoundStr      = "status 404 - not found"
	ConflictStr      = "already exists"
	ConflictErrorStr = "status 409 - conflict"
)

const (

	// PouchDBIdBridgeStr - value required for pouchDB to properly identify a an item. Used in the
	// following way to generate an ID:
	//      id = <dataType> + PouchDBIdBridgeStr + <generatedUUID>
	PouchDBIdBridgeStr = "_2_"
)

// GenerateID - generates an ID for an object based on the type of the
// provided object. Returns 2 versions of the ID, the first is the full ID used
// to refer to the object in Couch, the second is the hash portion of the Couch ID
// which is used to refer to the data object.
func GenerateID(obj interface{}, dataType string) string {

	switch obj.(type) {
	case *pb.MonitoredObjectData:
		cast := obj.(*pb.MonitoredObjectData)
		return PrependToDataID(strings.TrimSpace(cast.GetId()), dataType)
	case *tenmod.MonitoredObject:
		cast := obj.(*tenmod.MonitoredObject)
		return PrependToDataID(strings.TrimSpace(cast.MonitoredObjectID), dataType)
	// case *tenmod.MetricBaseline:
	// 	cast := obj.(*tenmod.MetricBaseline)
	// 	return PrependToDataID(strings.TrimSpace(cast.MonitoredObjectID), dataType)
	// case tenmod.MetricBaseline:
	// 	cast := obj.(tenmod.MetricBaseline)
	// 	return PrependToDataID(strings.TrimSpace(cast.MonitoredObjectID), dataType)
	default:
		uuid := uuid.NewV4()
		return PrependToDataID(uuid.String(), dataType)
	}
}

func trimAndLowercase(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// GetDataIDFromFullID - returns just the hash portion of an ID from the full ID.
// Note that for monitoredObjects, a hash is not stored in the ID field.
// The datatype is appended by a name passed in through ingestion and may have
// underscores in its name. There was a bug in the previous logic that incorrectly
// stripped all substrings prepended with _.
// The fix here is to find the index of the _FIRST TWO_ underscores in the fullID string.
// Then only providing the slice AFTER the second underscore.
// The behaviour here is based on the fact we prepend the unique ID with the
// document's data type, eg: 'monitoredObject_2_'. See #GenerateID() on the actual behavior.
func GetDataIDFromFullID(fullID string) string {
	const OFFSET int = 2
	const UNDERSCORE string = "_"
	if len(fullID) == 0 {
		return ""
	}
	loc1 := strings.Index(fullID, UNDERSCORE)
	if loc1 < 0 {
		return fullID
	}
	loc2 := strings.Index(fullID[loc1+1:], UNDERSCORE)
	if loc2 < 0 {
		return fullID
	}

	stripIdx := loc1 + loc2 + OFFSET
	stripped := fullID[stripIdx:]
	return stripped
}

// PrependToDataID - generates a full ID from the dataID and the dataType
func PrependToDataID(dataID string, dataType string) string {
	return fmt.Sprintf("%s%s%s", dataType, PouchDBIdBridgeStr, dataID)
}

// MakeTimestamp - get a timestamp from epoch in milliseconds
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// EnsureIngestionProfileHasBothModels - Helper function to make sure that we include both flattened and hierarchical models
// until we can deprecate the hierarchical model.
// TODO: Remove this when we deprecate the hierarchical model
func EnsureIngestionProfileHasBothModels(profile *tenmod.IngestionProfile) {
	hasFlattenedModel := profile.MetricList != nil && len(profile.MetricList) != 0
	hasHierarchicalModel := profile.Metrics != nil

	if hasFlattenedModel {
		// using the flattened mode - add hierarchical model
		newMetrics := tenmod.IngPrfVendorMap{}
		for _, metric := range profile.MetricList {
			if newMetrics.VendorMap == nil {
				newMetrics.VendorMap = map[string]*tenmod.IngPrfMonitoredObjectTypeMap{}
			}
			if newMetrics.VendorMap[metric.Vendor] == nil {
				newMetrics.VendorMap[metric.Vendor] = &tenmod.IngPrfMonitoredObjectTypeMap{}
			}

			moTypeMap := newMetrics.VendorMap[metric.Vendor]
			if moTypeMap.MonitoredObjectTypeMap == nil {
				moTypeMap.MonitoredObjectTypeMap = map[string]*tenmod.IngPrfMetricMap{}
			}
			if moTypeMap.MonitoredObjectTypeMap[metric.MonitoredObjectType] == nil {
				moTypeMap.MonitoredObjectTypeMap[metric.MonitoredObjectType] = &tenmod.IngPrfMetricMap{}
			}

			metricMap := moTypeMap.MonitoredObjectTypeMap[metric.MonitoredObjectType]
			if metricMap.MetricMap == nil {
				metricMap.MetricMap = map[string]bool{}
			}

			metricMap.MetricMap[metric.Metric] = metric.Enabled

		}

		profile.Metrics = &newMetrics

		return
	}

	if !hasFlattenedModel && hasHierarchicalModel {
		ingestionDictionary := admmod.GetIngestionDictionaryFromFile()

		// using the hierarchical model - add flattened model
		flattened := []*tenmod.IngestionProfileMetric{}

		for vk, v := range profile.Metrics.VendorMap {
			for moKey, mo := range v.MonitoredObjectTypeMap {
				for m, enabled := range mo.MetricMap {
					// Get the directions from the Ingestion dictionary:
					directions := getDirectionValuesForMetricFromIngestionDictionary(m, vk, moKey, ingestionDictionary)

					for _, d := range directions {
						dimensions := map[string][]string{}

						flattened = append(flattened, &tenmod.IngestionProfileMetric{
							Dimensions:          dimensions,
							Direction:           d,
							Enabled:             enabled,
							Metric:              m,
							MonitoredObjectType: moKey,
							Vendor:              vk,
						})
					}

				}
			}
		}

		profile.MetricList = flattened

		return
	}
}

func getDimensionValueForMetricFromIngestionDictionary(dimensionKey string, metricName string, vendorName string, moTypeName string, ingestionDictionary *admmod.IngestionDictionary) []string {
	for _, dictItem := range ingestionDictionary.MetricList {
		if dictItem.Metric == metricName && dictItem.MonitoredObjectType == moTypeName && dictItem.Vendor == vendorName {
			for dimKey, dim := range dictItem.Dimensions {
				if dimKey == dimensionKey {
					return dim
				}
			}
		}
	}

	return []string{}
}

func getDirectionValuesForMetricFromIngestionDictionary(metricName string, vendorName string, moTypeName string, ingestionDictionary *admmod.IngestionDictionary) []string {

	for _, dictItem := range ingestionDictionary.MetricList {
		if dictItem.Metric == metricName && dictItem.MonitoredObjectType == moTypeName && dictItem.Vendor == vendorName {
			return dictItem.Directions
		}
	}

	return []string{}
}

// EnsureThresholdProfileHasBothModels - Helper function to make sure that we include both flattened and hierarchical models
// until we can deprecate the hierarchical model.
// TODO: Remove this when we deprecate the hierarchical model
func EnsureThresholdProfileHasBothModels(profile *tenmod.ThresholdProfile) {
	hasFlattenedModel := profile.ThresholdList != nil && len(profile.ThresholdList) != 0
	hasHierarchicalModel := profile.Thresholds != nil

	if hasFlattenedModel {
		// using the flattened mode - add hierarchical model
		newThresholds := tenmod.ThrPrfVendorMap{}
		for _, thresh := range profile.ThresholdList {
			if newThresholds.VendorMap == nil {
				newThresholds.VendorMap = map[string]*tenmod.ThrPrfMetric{}
			}

			if newThresholds.VendorMap[thresh.Vendor] == nil {
				newThresholds.VendorMap[thresh.Vendor] = &tenmod.ThrPrfMetric{}
			}

			vendor := newThresholds.VendorMap[thresh.Vendor]
			if vendor.MonitoredObjectTypeMap == nil {
				vendor.MonitoredObjectTypeMap = map[string]*tenmod.ThrPrfMetricMap{}
			}

			moMap := vendor.MonitoredObjectTypeMap
			if moMap[thresh.MonitoredObjectType] == nil {
				moMap[thresh.MonitoredObjectType] = &tenmod.ThrPrfMetricMap{}
			}

			mo := moMap[thresh.MonitoredObjectType]
			if mo.MetricMap == nil {
				mo.MetricMap = map[string]*tenmod.ThrPrfDirectionMap{}
			}

			if mo.MetricMap[thresh.Metric] == nil {
				mo.MetricMap[thresh.Metric] = &tenmod.ThrPrfDirectionMap{}
			}

			met := mo.MetricMap[thresh.Metric]
			if met.DirectionMap == nil {
				met.DirectionMap = map[string]*tenmod.ThrPrfEventMap{}
			}

			if met.DirectionMap[thresh.Direction] == nil {
				met.DirectionMap[thresh.Direction] = &tenmod.ThrPrfEventMap{}
			}

			if met.DirectionMap[thresh.Direction].EventMap == nil {
				met.DirectionMap[thresh.Direction].EventMap = map[string]*tenmod.ThrPrfEventAttrMap{}
			}

			eventMap := met.DirectionMap[thresh.Direction].EventMap
			for _, event := range thresh.Events {
				eventName := event["eventName"]
				eventAttrMap := map[string]string{}
				for eventKey, eventValue := range event {
					if eventKey == "eventName" {
						continue
					}

					eventAttrMap[eventKey] = eventValue
				}
				eventMap[eventName] = &tenmod.ThrPrfEventAttrMap{
					EventAttrMap: eventAttrMap,
				}

			}

			metMap := vendor.MetricMap
			if metMap == nil {
				metMap = map[string]*tenmod.ThrPrfUIEvtAttrMap{}
			}
			if metMap[thresh.Metric] == nil {
				metMap[thresh.Metric] = &tenmod.ThrPrfUIEvtAttrMap{}
			}

			uiEventAttrMap := metMap[thresh.Metric]
			if uiEventAttrMap.EventAttrMap == nil {
				uiEventAttrMap.EventAttrMap = map[string]string{}
			}

			uiEventAttrMap.EventAttrMap["enabled"] = thresh.Enabled
		}

		profile.Thresholds = &newThresholds
		return
	}

	if !hasFlattenedModel && hasHierarchicalModel {
		// using the hierarchical model - add flattened model
		flattened := []*tenmod.ThresholdProfileThreshold{}

		for vk, v := range profile.Thresholds.VendorMap {
			for moKey, mo := range v.MonitoredObjectTypeMap {
				for metKey, met := range mo.MetricMap {
					enabled := "true"
					if v.MetricMap != nil {
						enabled = v.MetricMap[metKey].EventAttrMap["enabled"]
					}

					for dk, d := range met.DirectionMap {
						events := []map[string]string{}
						for ek, e := range d.EventMap {
							eventAttrs := map[string]string{}
							eventAttrs["eventName"] = ek
							for evAttrKey, evAttrVal := range e.EventAttrMap {
								eventAttrs[evAttrKey] = evAttrVal
							}
							events = append(events, eventAttrs)
						}
						flattened = append(flattened, &tenmod.ThresholdProfileThreshold{
							Enabled:             enabled,
							Metric:              metKey,
							MonitoredObjectType: moKey,
							Vendor:              vk,
							Direction:           dk,
							Events:              events,
							Dimensions:          map[string][]string{},
						})
					}
				}
			}
		}

		profile.ThresholdList = flattened

		return
	}
}
