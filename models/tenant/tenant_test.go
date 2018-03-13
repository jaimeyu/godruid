package tenant

import (
	"strings"
	"testing"
	"time"

	"github.com/accedian/adh-gather/models/common"
	"github.com/icrowley/fake"

	testUtil "github.com/accedian/adh-gather/models/test"
	uuid "github.com/satori/go.uuid"
)

func TestTenantUserSerialization(t *testing.T) {
	original := &User{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantUserType),
		TenantID:              fake.CharactersN(12),
		Username:              fake.EmailAddress(),
		Password:              fake.Password(6, 8, true, true, false),
		State:                 string(common.UserActive),
		SendOnboardingEmail:   true,
		OnboardingToken:       fake.CharactersN(14),
		UserVerified:          false,
		Domains:               []string{fake.CharactersN(5), fake.CharactersN(7)},
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "username", "password", "sendOnboardingEmail", "state", "onboardingToken", "userVerified", "domains", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &User{}, original.ID, attrKeys)
}

func TestTenantDomainSerialization(t *testing.T) {
	original := &Domain{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantDomainType),
		TenantID:              fake.CharactersN(12),
		Name:                  fake.State(),
		Color:                 fake.HexColor(),
		ThresholdProfileSet:   []string{uuid.NewV4().String(), uuid.NewV4().String()},
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "name", "color", "thresholdProfileSet", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &Domain{}, original.ID, attrKeys)
}

func TestTenantIngestionProfileSerialization(t *testing.T) {
	metrics := make(map[string]map[string]map[string]bool)
	company := fake.Company()
	brand := fake.Brand()
	metrics[company] = make(map[string]map[string]bool)
	metrics[company][brand] = make(map[string]bool)
	metrics[company][brand][fake.JobTitle()] = true
	metrics[company][brand][fake.JobTitle()] = true
	metrics[company][brand][fake.JobTitle()] = false
	company2 := fake.Company()
	brand2 := fake.Brand()
	metrics[company2] = make(map[string]map[string]bool)
	metrics[company2][brand2] = make(map[string]bool)
	metrics[company2][brand2][fake.JobTitle()] = false
	original := &IngestionProfile{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantIngestionProfileType),
		TenantID:              fake.CharactersN(12),
		Metrics:               metrics,
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "metrics", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &IngestionProfile{}, original.ID, attrKeys)
}

func TestTenantThresholdProfileSerialization(t *testing.T) {
	monObjectTypeMap := make(map[string]map[string]map[string]map[string]string)
	monobj := fake.Brand()
	metric := fake.JobTitle()
	dir := "0"
	dir2 := "1"
	attr := fake.Industry()
	attr2 := fake.Industry()
	monObjectTypeMap[monobj] = make(map[string]map[string]map[string]string)
	monObjectTypeMap[monobj][metric] = make(map[string]map[string]string)
	monObjectTypeMap[monobj][metric][dir] = make(map[string]string)
	monObjectTypeMap[monobj][metric][dir][attr] = fake.Continent()
	monObjectTypeMap[monobj][metric][dir][attr2] = fake.Continent()
	monObjectTypeMap[monobj][metric][dir2] = make(map[string]string)
	monObjectTypeMap[monobj][metric][dir2][attr] = fake.Continent()
	monObjectTypeMap[monobj][metric][dir2][attr2] = fake.Continent()

	metricMap := map[string]string{fake.Color(): fake.Brand(), fake.Color(): fake.Brand()}
	thresholds := map[string]MonitoredObjectTypeMap{
		fake.Company(): MonitoredObjectTypeMap{
			MonitoredObjectTypeMap: monObjectTypeMap,
			MetricMap:              metricMap,
		},
	}
	original := &ThresholdProfile{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantThresholdProfileType),
		TenantID:              fake.CharactersN(12),
		Thresholds:            thresholds,
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "thresholds", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &ThresholdProfile{}, original.ID, attrKeys)
}

func TestTenantMonitoredObjectSerialization(t *testing.T) {
	actName := fake.City()
	refName := fake.City()
	original := &MonitoredObject{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantMonitoredObjectType),
		TenantID:              fake.CharactersN(12),
		ActuatorType:          fake.Company(),
		ActuatorName:          actName,
		ReflectorType:         fake.Company(),
		ReflectorName:         refName,
		ObjectType:            fake.Brand(),
		ObjectName:            fake.City(),
		MonitoredObjectID:     strings.Join([]string{actName, refName}, "-"),
		DomainSet:             []string{uuid.NewV4().String(), uuid.NewV4().String()},
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "actuatorName", "actuatorType",
		"reflectorName", "reflectorType", "objectName", "objectType", "domainSet",
		"id", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &MonitoredObject{}, original.ID, attrKeys)
}

func TestTenantMetadataSerialization(t *testing.T) {
	original := &Metadata{
		ID:                      uuid.NewV4().String(),
		REV:                     uuid.NewV4().String(),
		Datatype:                string(TenantMetaType),
		TenantID:                fake.CharactersN(12),
		TenantName:              fake.Company(),
		DefaultThresholdProfile: uuid.NewV4().String(),
		CreatedTimestamp:        time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp:   time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "tenantName", "defaultThresholdProfile", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &Metadata{}, original.ID, attrKeys)
}
