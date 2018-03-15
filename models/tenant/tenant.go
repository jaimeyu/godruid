package tenant

import (
	"errors"
)

// TenantDataType - enumeration of the types of data stored in the Tenant Datastore
type TenantDataType string

const (
	// TenantUserType - datatype string used to identify a Tenant User in the datastore record
	TenantUserType TenantDataType = "user"

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
)

const (
	// TenantUserStr - common name of the TenantUser data type for use in logs.
	TenantUserStr = "Tenant User"

	// TenantDomainStr - common name of the Tenant Domain data type for use in logs.
	TenantDomainStr = "Tenant Domain"

	// TenantIngestionProfileStr - common name of the Tenant Ingestion Profile data type for use in logs.
	TenantIngestionProfileStr = "Tenant Ingestion Profile"

	// TenantMonitoredObjectStr - common name of the Tenant Monitored Object data type for use in logs.
	TenantMonitoredObjectStr = "Tenant Monitored Object"

	// TenantThresholdProfileStr - common name of the Tenant Ingestion Profile data type for use in logs.
	TenantThresholdProfileStr = "Tenant Threshold Profile"

	// MonitoredObjectToDomainMapStr - common name for the Monitored Object to Doamin Map for use in logs.
	MonitoredObjectToDomainMapStr = "Monitored Object to Doamin Map"

	// TenantMetaStr - common name for the Meta for use in logs.
	TenantMetaStr = "Tenant Meta"
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
	if isUpdate && (len(u.REV) == 0 || u.CreatedTimestamp == 0) {
		return errors.New("Invalid Tenant User request: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

// Domain - defines a Tenant Domain.
type Domain struct {
	ID                    string   `json:"_id"`
	REV                   string   `json:"_rev"`
	Datatype              string   `json:"datatype"`
	TenantID              string   `json:"tenantId"`
	Name                  string   `json:"name"`
	Color                 string   `json:"color"`
	ThresholdProfileSet   []string `json:"thresholdProfileSet"`
	CreatedTimestamp      int64    `json:"createdTimestamp"`
	LastModifiedTimestamp int64    `json:"lastModifiedTimestamp"`
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

// Validate - used during validation of incoming REST requests for this object
func (d *Domain) Validate(isUpdate bool) error {
	if len(d.TenantID) == 0 {
		return errors.New("Invalid Tenant Domain request: must provide a Tenant ID")
	}
	if !isUpdate && len(d.REV) != 0 {
		return errors.New("Invalid Tenant Domain request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(d.REV) == 0 || d.CreatedTimestamp == 0) {
		return errors.New("Invalid Tenant Domain request: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

// IngestionProfile - defines a Tenant Ingestion Profile.
type IngestionProfile struct {
	ID                    string                                                                 `json:"_id"`
	REV                   string                                                                 `json:"_rev"`
	Datatype              string                                                                 `json:"datatype"`
	TenantID              string                                                                 `json:"tenantId"`
	Metrics               map[string]map[string]map[string]map[string]map[string]map[string]bool `json:"metrics"`
	CreatedTimestamp      int64                                                                  `json:"createdTimestamp"`
	LastModifiedTimestamp int64                                                                  `json:"lastModifiedTimestamp"`
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
	if isUpdate && (len(prf.REV) == 0 || prf.CreatedTimestamp == 0) {
		return errors.New("Invalid Tenant Ingestion Profile request: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

// ThresholdProfile - defines a Tenant Threshold Profile.
type ThresholdProfile struct {
	ID                    string                                     `json:"_id"`
	REV                   string                                     `json:"_rev"`
	Datatype              string                                     `json:"datatype"`
	TenantID              string                                     `json:"tenantId"`
	Name                  string                                     `json:"name"`
	Thresholds            map[string]map[string]MonitoredObjectGroup `json:"thresholds"`
	CreatedTimestamp      int64                                      `json:"createdTimestamp"`
	LastModifiedTimestamp int64                                      `json:"lastModifiedTimestamp"`
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

// Validate - used during validation of incoming REST requests for this object
func (prf *ThresholdProfile) Validate(isUpdate bool) error {
	if len(prf.TenantID) == 0 {
		return errors.New("Invalid Tenant Threshold Profile request: must provide a Tenant ID")
	}
	if !isUpdate && len(prf.REV) != 0 {
		return errors.New("Invalid Tenant Threshold Profile request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(prf.REV) == 0 || prf.CreatedTimestamp == 0) {
		return errors.New("Invalid Tenant Threshold Profile request: must provide a createdTimestamp and revision for an update")
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
	Color                 string   `json:"color"`
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

	if isUpdate && (len(mo.REV) == 0 || mo.CreatedTimestamp == 0) {
		return errors.New("Invalid Tenant Monitored object request: must provide a createdTimestamp and revision for an update")
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

// Validate - used during validation of incoming REST requests for this object
func (meta *Metadata) Validate(isUpdate bool) error {
	if len(meta.TenantID) == 0 {
		return errors.New("Invalid Tenant Metadata request: must provide a Tenant ID")
	}
	if !isUpdate && len(meta.REV) != 0 {
		return errors.New("Invalid Tenant Metadata request: must not provide a revision value in a creation request")
	}
	if isUpdate && (len(meta.REV) == 0 || meta.CreatedTimestamp == 0) {
		return errors.New("Invalid Tenant Metadata request: must provide a createdTimestamp and revision for an update")
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
