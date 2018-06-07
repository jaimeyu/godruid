package tenant

import (
	"errors"

	"github.com/manyminds/api2go/jsonapi"
)

// TenantDataType - enumeration of the types of data stored in the Tenant Datastore
type TenantDataType string

const (
	// TenantUserType - datatype string used to identify a Tenant User in the datastore record
	TenantUserType TenantDataType = "user"

	// TenantConnectorConfigType - datatype string used to identify a Tenant ConnectorConfig in the datastore record
	TenantConnectorConfigType TenantDataType = "connectorConfig"

	// TenantConnectorInstanceType - datatype string used to identify a Tenant ConnectorInstance in the datastore record
	TenantConnectorInstanceType TenantDataType = "connectorInstance"

	// TenantDomainType - datatype string used to identify a Tenant Domain in the datastore record
	TenantDomainType TenantDataType = "domain"

	// TenantIngestionProfileType - datatype string used to identify a Tenant Ingestion Profile in the datastore record
	TenantIngestionProfileType TenantDataType = "ingestionProfile"

	// TenantMonitoredObjectType - datatype string used to identify a Tenant MonitoredObject in the datastore record
	TenantMonitoredObjectType TenantDataType = "monitoredObject"

	// TenantThresholdProfileType - datatype string used to identify a Tenant Ingestion Profile in the datastore record
	TenantThresholdProfileType TenantDataType = "thresholdProfile"

	// TenantMetaType - datatype string used to identify a Tenant Meta in the datastore record
	TenantMetaType TenantDataType = "tenantMetadata"

	// TenantReportScheduleConfigType - datatype string used to identify a Tenant SLA Report in the datastore record
	TenantReportScheduleConfigType TenantDataType = "tenantReportScheduleConfig"

	// TenantReportType - datatype string used to identify a Tenant Report in the datastore record
	TenantReportType TenantDataType = "report"

	// DashboardType - datatype string used to identify a Tenant Meta in the datastore record
	TenantDashboardType TenantDataType = "dashboard"
)

// MonitoredObjectType - defines the known types of Monitored Objects for Skylight Datahub
type MonitoredObjectType string

const (
	// MonitoredObjectUnknown - value for Unnkown monitored objects
	MonitoredObjectUnknown MonitoredObjectType = "unknown"

	// TwampPE - value for TWAMP PE monitored objects
	TwampPE MonitoredObjectType = "twamp-pe"

	// TwampSF - value for TWAMP Stateful monitored objects
	TwampSF MonitoredObjectType = "twamp-sf"

	// TwampSL - value for TWAMP Stateless monitored objects
	TwampSL MonitoredObjectType = "twamp-sl"

	// Flowmeter - value for Flowmeter monitored objects
	Flowmeter MonitoredObjectType = "flowmeter"
)

// VendorMetricType - defines the known types of Vendor metric categories.
type VendorMetricType string

const (
	// AccedianTwamp - represents Accedian TWAMP vendor metrics.
	AccedianTwamp VendorMetricType = "accedian-twamp"

	// AccedianFlowmeter - represents Accedian Flowmeter vendor metrics.
	AccedianFlowmeter VendorMetricType = "accedian-flowmeter"
)

// MonitoredObjectDeviceType - defines the known types of devices (actuators / reflectors) for
// Skylight Datahub
type MonitoredObjectDeviceType string

const (
	// MonitoredObjectDeviceUnknown - value for TWAMP Light monitored objects
	MonitoredObjectDeviceUnknown MonitoredObjectDeviceType = "unknown"

	// AccedianNID - value for Accedian NID monitored objects device type
	AccedianNID MonitoredObjectDeviceType = "accedian-nid"

	// AccedianVNID - value for Accedian VNID monitored objects device type
	AccedianVNID MonitoredObjectDeviceType = "accedian-vnid"
)

const (
	// TenantUserStr - common name of the TenantUser data type for use in logs.
	TenantUserStr = "Tenant User"

	// TenantConnectorConfigStr - common name of the Tenant Connector data type for use in logs.
	TenantConnectorConfigStr = "Tenant Connector"

	// TenantConnectorInstanceStr - common name of the Tenant ConnectorInstance data type for use in logs.
	TenantConnectorInstanceStr = "Tenant ConnectorInstance"

	// TenantDomainStr - common name of the Tenant Domain data type for use in logs.
	TenantDomainStr = "Tenant Domain"

	// TenantIngestionProfileStr - common name of the Tenant Ingestion Profile data type for use in logs.
	TenantIngestionProfileStr = "Tenant Ingestion Profile"

	// TenantMonitoredObjectStr - common name of the Tenant Monitored Object data type for use in logs.
	TenantMonitoredObjectStr = "Tenant Monitored Object"

	// TenantThresholdProfileStr - common name of the Tenant Ingestion Profile data type for use in logs.
	TenantThresholdProfileStr = "Tenant Threshold Profile"

	// MonitoredObjectToDomainMapStr - common name for the Monitored Object to Domain Map for use in logs.
	MonitoredObjectToDomainMapStr = "Monitored Object to Domain Map"

	// TenantMetaStr - common name for the Meta for use in logs.
	TenantMetaStr = "Tenant Meta"

	// TenantReportScheduleConfigStr - common name for the report schedule configuration for use in logs.
	TenantReportScheduleConfigStr = "Tenant Report Schedule Configuration"

	// TenantSLAReportStr - common name for the sla report for use in logs.
	TenantSLAReportStr = "Tenant SLA Report"

	// TenantDashboardStr - common name for the Dashboard for use in logs.
	TenantDashboardStr = "Tenant Dashboard"
)

// User - defines a Tenant user.
type User struct {
	ID                    string   `json:"_id"`
	REV                   string   `json:"_rev"`
	Datatype              string   `json:"datatype"`
	TenantID              string   `json:"tenantId"`
	Username              string   `json:"username"`
	Password              string   `json:"password"`
	SendOnboardingEmail   bool     `json:"sendOnboardingEmail"`
	OnboardingToken       string   `json:"onboardingToken"`
	UserVerified          bool     `json:"userVerified"`
	State                 string   `json:"state"`
	Domains               []string `json:"domains"`
	CreatedTimestamp      int64    `json:"createdTimestamp"`
	LastModifiedTimestamp int64    `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (u *User) GetID() string {
	return u.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (u *User) SetID(s string) error {
	u.ID = s
	return nil
}

// Validate - used during validation of incoming REST requests for this object
func (u *User) Validate(isUpdate bool) error {
	if len(u.TenantID) == 0 {
		return errors.New("Invalid Tenant User request: must provide a Tenant ID")
	}
	if !isUpdate && len(u.REV) != 0 {
		return errors.New("Invalid Tenant User request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(u.REV) == 0) {
		return errors.New("Invalid Tenant User request: must provide a revision (_rev) for an update")
	}

	return nil
}

// ConnectorInstance - defines a Tenant ConnectorInstnace
type ConnectorInstance struct {
	ID                    string `json:"_id"`
	REV                   string `json:"_rev"`
	Datatype              string `json:"datatype"`
	TenantID              string `json:"tenantId"`
	Status                string `json:"Status"`
	Hostname              string `json:"hostname"`
	CreatedTimestamp      int64  `json:"createdTimestamp"`
	LastModifiedTimestamp int64  `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (d *ConnectorInstance) GetID() string {
	return d.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (d *ConnectorInstance) SetID(s string) error {
	d.ID = s
	return nil
}

// Validate - used during validation of incoming REST requests for this object
func (d *ConnectorInstance) Validate(isUpdate bool) error {
	if len(d.TenantID) == 0 {
		return errors.New("Invalid Tenant ConnectorInstance request: must provide a Tenant ID")
	}
	if !isUpdate && len(d.REV) != 0 {
		return errors.New("Invalid Tenant ConnectorInstance request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(d.REV) == 0 || d.CreatedTimestamp == 0) {
		return errors.New("Invalid Tenant ConnectorInstance request: must provide a revision (_rev) for an update")
	}

	return nil
}

// ConnectorConfig - defines a Tenant ConnectorConfig
type ConnectorConfig struct {
	ID                              string `json:"_id"`
	REV                             string `json:"_rev"`
	Datatype                        string `json:"datatype"`
	URL                             string `json:"url"`
	Port                            int    `json:"port"`
	PollingFrequency                int    `json:"pollingFrequency"`
	Username                        string `json:"username"`
	Password                        string `json:"password"`
	ExportGroup                     string `json:"exportGroup"`
	DatahubHearbeatFrequency        int    `json:"datahubHeartbeatFrequency"`
	DatahubConnectionRetryFrequency int    `json:"datahubConnectionRetryFrequency"`
	ConnectorInstanceID             string `json:"connectorInstanceId"`
	TenantID                        string `json:"tenantId"`
	Name                            string `json:"name"`
	Zone                            string `json:"zone"`
	Type                            string `json:"type"`
	CreatedTimestamp                int64  `json:"createdTimestamp"`
	LastModifiedTimestamp           int64  `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (d *ConnectorConfig) GetID() string {
	return d.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (d *ConnectorConfig) SetID(s string) error {
	d.ID = s
	return nil
}

// Validate - used during validation of incoming REST requests for this object
func (d *ConnectorConfig) Validate(isUpdate bool) error {
	if len(d.TenantID) == 0 {
		return errors.New("Invalid Tenant ConnectorConfig request: must provide a Tenant ID")
	}
	if !isUpdate && len(d.REV) != 0 {
		return errors.New("Invalid Tenant ConnectorConfig request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(d.REV) == 0) {
		return errors.New("Invalid Tenant ConnectorConfig request: must provide a revision for an update")
	}

	return nil
}

// Domain - defines a Tenant Domain.
type Domain struct {
	ID                    string `json:"_id"`
	REV                   string `json:"_rev"`
	Datatype              string `json:"datatype"`
	TenantID              string `json:"tenantId"`
	Name                  string `json:"name"`
	Color                 string `json:"color"`
	CreatedTimestamp      int64  `json:"createdTimestamp"`
	LastModifiedTimestamp int64  `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (d *Domain) GetID() string {
	return d.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (d *Domain) SetID(s string) error {
	d.ID = s
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (d *Domain) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "thresholdProfiles",
			Name: "thresholdProfiles",
		},
	}
}

// Validate - used during validation of incoming REST requests for this object
func (d *Domain) Validate(isUpdate bool) error {
	if len(d.TenantID) == 0 {
		return errors.New("Invalid Tenant Domain request: must provide a Tenant ID")
	}
	if !isUpdate && len(d.REV) != 0 {
		return errors.New("Invalid Tenant Domain request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(d.REV) == 0) {
		return errors.New("Invalid Tenant Domain request: must provide a revision (_rev) for an update")
	}

	return nil
}

// IngestionProfile - defines a Tenant Ingestion Profile.
type IngestionProfile struct {
	ID                    string          `json:"_id"`
	REV                   string          `json:"_rev"`
	Datatype              string          `json:"datatype"`
	TenantID              string          `json:"tenantId"`
	Metrics               IngPrfVendorMap `json:"metrics"`
	CreatedTimestamp      int64           `json:"createdTimestamp"`
	LastModifiedTimestamp int64           `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (prf *IngestionProfile) GetID() string {
	return prf.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (prf *IngestionProfile) SetID(s string) error {
	prf.ID = s
	return nil
}

// Validate - used during validation of incoming REST requests for this object
func (prf *IngestionProfile) Validate(isUpdate bool) error {
	if len(prf.TenantID) == 0 {
		return errors.New("Invalid Tenant Ingestion Profile request: must provide a Tenant ID")
	}
	if !isUpdate && len(prf.REV) != 0 {
		return errors.New("Invalid Tenant Ingestion Profile request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(prf.REV) == 0) {
		return errors.New("Invalid Tenant Ingestion Profile request: must provide a revision (_rev) for an update")
	}

	return nil
}

type IngPrfVendorMap struct {
	VendorMap map[string]IngPrfMonitoredObjectTypeMap `json:"vendorMap"`
}

type IngPrfMonitoredObjectTypeMap struct {
	MonitoredObjectTypeMap map[string]IngPrfMetricMap `json:"monitoredObjectTypeMap"`
}

type IngPrfMetricMap struct {
	MetricMap map[string]bool `json:"metricMap"`
}

// ThresholdProfile - defines a Tenant Threshold Profile.
type ThresholdProfile struct {
	ID                    string          `json:"_id"`
	REV                   string          `json:"_rev"`
	Datatype              string          `json:"datatype"`
	TenantID              string          `json:"tenantId"`
	Name                  string          `json:"name"`
	Thresholds            ThrPrfVendorMap `json:"thresholds"`
	CreatedTimestamp      int64           `json:"createdTimestamp"`
	LastModifiedTimestamp int64           `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (prf *ThresholdProfile) GetID() string {
	return prf.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (prf *ThresholdProfile) SetID(s string) error {
	prf.ID = s
	return nil
}

type ThrPrfVendorMap struct {
	VendorMap map[string]ThrPrfMetric `json:"vendorMap"`
}

type ThrPrfMetric struct {
	MetricMap              map[string]ThrPrfUIEvtAttrMap `json:"metricMap"`
	MonitoredObjectTypeMap map[string]ThrPrfMetricMap    `json:"monitoredObjectTypeMap"`
}

type ThrPrfUIEvtAttrMap struct {
	EventAttrMap map[string]string `json:"eventAttrMap"`
}

type ThrPrfMetricMap struct {
	MetricMap map[string]ThrPrfDirectionMap `json:"metricMap"`
}

type ThrPrfDirectionMap struct {
	DirectionMap map[string]ThrPrfEventMap `json:"directionMap"`
}

type ThrPrfEventMap struct {
	EventMap map[string]ThrPrfEventAttrMap `json:"eventMap"`
}

type ThrPrfEventAttrMap struct {
	EventAttrMap map[string]string `json:"eventAttrMap"`
}

// Validate - used during validation of incoming REST requests for this object
func (prf *ThresholdProfile) Validate(isUpdate bool) error {
	if len(prf.TenantID) == 0 {
		return errors.New("Invalid Tenant Threshold Profile request: must provide a Tenant ID")
	}
	if !isUpdate && len(prf.REV) != 0 {
		return errors.New("Invalid Tenant Threshold Profile request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(prf.REV) == 0) {
		return errors.New("Invalid Tenant Threshold Profile request: must provide a revision (_rev) for an update")
	}

	return nil
}

type MonitoredObjectGroup struct {
	MonitoredObjectTypeMap map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]string `json:"monitoredObjectTypeMap"`
	MetricMap              map[string]map[string]map[string]string                                                                   `json:"metricMap"`
}

// MonitoredObject - defines a Tenant Monitored Object.
type MonitoredObject struct {
	ID                    string   `json:"_id"`
	REV                   string   `json:"_rev"`
	Datatype              string   `json:"datatype"`
	TenantID              string   `json:"tenantId"`
	MonitoredObjectID     string   `json:"id"`
	ActuatorType          string   `json:"actuatorType"`
	ActuatorName          string   `json:"actuatorName"`
	ReflectorType         string   `json:"reflectorType"`
	ReflectorName         string   `json:"reflectorName"`
	ObjectType            string   `json:"objectType"`
	ObjectName            string   `json:"objectName"`
	DomainSet             []string `json:"domainSet"`
	CreatedTimestamp      int64    `json:"createdTimestamp"`
	LastModifiedTimestamp int64    `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (mo *MonitoredObject) GetID() string {
	return mo.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (mo *MonitoredObject) SetID(s string) error {
	mo.ID = s
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (mo *MonitoredObject) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "domains",
			Name: "domains",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (mo *MonitoredObject) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, domID := range mo.DomainSet {
		result = append(result, jsonapi.ReferenceID{
			ID:   domID,
			Type: "domains",
			Name: "domains",
		})
	}

	return result
}

// SetToManyReferenceIDs sets domain reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (mo *MonitoredObject) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "domains" {
		mo.DomainSet = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds new domains to the reference list
func (mo *MonitoredObject) AddToManyIDs(name string, IDs []string) error {
	if name == "thresholdProfiles" {
		mo.DomainSet = append(mo.DomainSet, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes domains from the reference list
func (mo *MonitoredObject) DeleteToManyIDs(name string, IDs []string) error {
	if name == "thresholdProfiles" {
		for _, ID := range IDs {
			for pos, oldID := range mo.DomainSet {
				if ID == oldID {
					// match, this ID must be removed
					mo.DomainSet = append(mo.DomainSet[:pos], mo.DomainSet[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// Validate - used during validation of incoming REST requests for this object
func (mo *MonitoredObject) Validate(isUpdate bool) error {
	if len(mo.TenantID) == 0 {
		return errors.New("Invalid Tenant Monitored Object request: must provide a Tenant ID")
	}
	if !isUpdate && len(mo.REV) != 0 {
		return errors.New("Invalid Tenant Monitored Object request: must not provide a revision value in a creation request")
	}
	if len(mo.MonitoredObjectID) == 0 {
		return errors.New("Invalid Tenant Monitored Object request: must provide a Monitored Object ID")

	}

	if isUpdate && (len(mo.REV) == 0) {
		return errors.New("Invalid Tenant Monitored object request: must provide a revision (_rev) for an update")
	}

	return nil
}

// Metadata - defines a Tenant Metadata.
type Metadata struct {
	ID                      string `json:"_id"`
	REV                     string `json:"_rev"`
	Datatype                string `json:"datatype"`
	TenantID                string `json:"tenantId"`
	TenantName              string `json:"tenantName"`
	DefaultThresholdProfile string `json:"defaultThresholdProfile"`
	CreatedTimestamp        int64  `json:"createdTimestamp"`
	LastModifiedTimestamp   int64  `json:"lastModifiedTimestamp"`
}

// GetID - required implementation for jsonapi marshalling
func (meta *Metadata) GetID() string {
	return meta.ID
}

// SetID - required implementation for jsonapi unmarshalling
func (meta *Metadata) SetID(s string) error {
	meta.ID = s
	return nil
}

// GetName - required implementation for renaming the type in jsonapi payload
func (meta *Metadata) GetName() string {
	return string(TenantMetaType)
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (meta *Metadata) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "defaultThresholdProfile",
			Name: "defaultThresholdProfile",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (meta *Metadata) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	result = append(result, jsonapi.ReferenceID{
		ID:   meta.DefaultThresholdProfile,
		Type: "defaultThresholdProfile",
		Name: "defaultThresholdProfile",
	})

	return result
}

// SetToOneReferenceID sets domain reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (meta *Metadata) SetToOneReferenceID(name string, ID string) error {
	if name == "defaultThresholdProfile" {
		meta.DefaultThresholdProfile = ID
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// Validate - used during validation of incoming REST requests for this object
func (meta *Metadata) Validate(isUpdate bool) error {
	if len(meta.TenantID) == 0 {
		return errors.New("Invalid Tenant Metadata request: must provide a Tenant ID")
	}
	if !isUpdate && len(meta.REV) != 0 {
		return errors.New("Invalid Tenant Metadata request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(meta.REV) == 0) {
		return errors.New("Invalid Tenant Metadata request: must provide a revision (_rev) for an update")
	}

	return nil
}

// MonitoredObjectCountByDomainRequest - request type for retrieving a MonitoredObject Count by Domain
type MonitoredObjectCountByDomainRequest struct {
	TenantID  string   `json:"tenantId"`
	ByCount   bool     `json:"byCount"`
	DomainSet []string `json:"domainSet"`
}

// Validate - used during validation of incoming REST requests for this object
func (req *MonitoredObjectCountByDomainRequest) Validate(isUpdate bool) error {
	if len(req.TenantID) == 0 {
		return errors.New("Invalid Tenant Metadata request: must provide a Tenant ID")
	}

	return nil
}

// MonitoredObjectCountByDomainResponse response for a request for MonitoredObject Count by Domain
type MonitoredObjectCountByDomainResponse struct {
	DomainToMonitoredObjectCountMap map[string]int64    `json:"domainToMonitoredObjectCountMap"`
	DomainToMonitoredObjectSetMap   map[string][]string `json:"domainToMonitoredObjectSetMap"`
}

// BulkMonitoredObjectRequest - used for requests that pass in a set of Monitored Objects
type BulkMonitoredObjectRequest struct {
	MonitoredObjectSet []*MonitoredObject `json:"monitoredObjectSet"`
}

type Dashboard struct {
	ID        string   `json:"_id"`
	REV       string   `json:"_rev"`
	TenantID  string   `json:"-"` // UI does not write this property
	Name      string   `json:"name"`
	DomainSet []string `json:"domainSet"`
}
