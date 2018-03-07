package admin

import (
	"github.com/manyminds/api2go/jsonapi"
)

// DataType - data type descriptors for objects stored in the admin datastore
type DataType string

const (
	// AdminUserType - datatype string used to identify an Admin User in the datastore record
	AdminUserType DataType = "adminUser"

	// TenantType - datatype string used to identify a Tenant Descriptor in the datastore record
	TenantType DataType = "tenant"

	// IngestionDictionaryType - datatype string used to identify an IngestionDictionary in the datastore record
	IngestionDictionaryType DataType = "ingestionDictionary"

	// ValidTypesType - datatype string used to identify a ValidTypes object in the datastore record
	ValidTypesType DataType = "validTypes"
)

const (
	// AdminUserStr - common name of the AdminUser data type for use in logs.
	AdminUserStr = "Admin User"

	// TenantStr - common name of the TenantDescriptor data type for use in logs.
	TenantStr = "Tenant"

	// IngestionDictionaryStr - common name of the IngestionDictionary data type for use in logs.
	IngestionDictionaryStr = "Ingestion Dictionary"

	// ValidTypesStr - common name of the ValidTypes data type for use in logs.
	ValidTypesStr = "Valid Types object"
)

// Tenant - defines a tenant
type Tenant struct {
	ID                    string `json:"-"`
	REV                   string `json:"_rev"`
	Datatype              string `json:"datatype"`
	Name                  string `json:"name"`
	URLSubdomain          string `json:"urlSubdomain"`
	State                 string `json:"state"`
	CreatedTimestamp      int64  `json:"createdTimestamp"`
	LastModifiedTimestamp int64  `json:"lastModifiedTimestamp"`
}

func (t *Tenant) GetID() string {
	return t.ID
}

func (t *Tenant) SetID(s string) error {
	t.ID = s
	return nil
}

// User - defines an Admin user.
type User struct {
	ID                    string `json:"-"`
	REV                   string `json:"_rev"`
	Datatype              string `json:"datatype"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	SendOnboardingEmail   bool   `json:"sendOnboardingEmail"`
	OnboardingToken       string `json:"onboardingToken"`
	UserVerified          bool   `json:"userVerified"`
	State                 string `json:"state"`
	CreatedTimestamp      int64  `json:"createdTimestamp"`
	LastModifiedTimestamp int64  `json:"lastModifiedTimestamp"`
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) SetID(s string) error {
	u.ID = s
	return nil
}

func (u *User) GetName() string {
	return jsonapi.Pluralize(string(AdminUserType))
}

// IngestionDictionary - defines an IngestionDictionary.
type IngestionDictionary struct {
	ID                    string                `json:"-"`
	REV                   string                `json:"_rev"`
	Datatype              string                `json:"datatype"`
	Metrics               map[string]*MetricMap `json:"metrics"`
	CreatedTimestamp      int64                 `json:"createdTimestamp"`
	LastModifiedTimestamp int64                 `json:"lastModifiedTimestamp"`
}

func (dict *IngestionDictionary) GetID() string {
	return dict.ID
}

func (dict *IngestionDictionary) SetID(s string) error {
	dict.ID = s
	return nil
}

type MetricMap struct {
	MetricMap map[string]*MetricDefinition `json:"metricMap"`
	UI        *UIGroupData                 `json:"ui"`
}

type UIGroupData struct {
	MetricGroups []string `json:"metricGroups"`
}

type MetricDefinition struct {
	MonitoredObjectTypes []*MonitoredObjectType `json:"monitoredObjectTypes"`
	UIData               *UIData                `json:"uiData"`
}

type UIData struct {
	Group    string `json:"group"`
	Position string `json:"position"`
}

type MonitoredObjectType struct {
	Key         string   `json:"key"`
	RawMetricID string   `json:"rawMetricId"`
	Units       []string `json:"units"`
	Directions  []string `json:"directions"`
}
