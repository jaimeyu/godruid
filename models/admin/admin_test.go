package admin

import (
	"testing"
	"time"

	"github.com/accedian/adh-gather/models/common"
	test "github.com/accedian/adh-gather/models/test"
	"github.com/icrowley/fake"

	uuid "github.com/satori/go.uuid"
)

func TestTenantSerialization(t *testing.T) {
	original := &Tenant{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantType),
		Name:                  fake.Company(),
		URLSubdomain:          fake.DomainName(),
		State:                 string(common.UserActive),
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "name", "urlSubdomain", "state", "createdTimestamp", "lastModifiedTimestamp"}

	test.RunSerializationTest(t, original, &Tenant{}, original.ID, attrKeys)
}

func TestAdminUserSerialization(t *testing.T) {
	original := &User{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(AdminUserType),
		Username:              fake.EmailAddress(),
		Password:              fake.Password(6, 8, true, true, false),
		State:                 string(common.UserActive),
		SendOnboardingEmail:   true,
		OnboardingToken:       fake.CharactersN(14),
		UserVerified:          false,
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "username", "password", "sendOnboardingEmail", "state", "onboardingToken", "userVerified", "createdTimestamp", "lastModifiedTimestamp"}

	test.RunSerializationTest(t, original, &User{}, original.ID, attrKeys)
}

func TestIngestionDictionarySerialization(t *testing.T) {
	uiData := &UIData{
		Group:    fake.CharactersN(6),
		Position: fake.CharactersN(8),
	}
	monObj := &MonitoredObjectType{
		Key:         fake.CharactersN(8),
		RawMetricID: fake.CharactersN(8),
		Units:       []string{fake.CharactersN(3), fake.CharactersN(3)},
		Directions:  []string{fake.Character(), fake.Character()},
	}
	monObj2 := &MonitoredObjectType{
		Key:         fake.CharactersN(8),
		RawMetricID: fake.CharactersN(8),
		Units:       []string{fake.CharactersN(3), fake.CharactersN(3)},
		Directions:  []string{fake.Character(), fake.Character()},
	}
	metricDefinition := &MetricDefinition{
		UIData:               uiData,
		MonitoredObjectTypes: []*MonitoredObjectType{monObj, monObj2},
	}
	uiGroupData := &UIGroupData{
		MetricGroups: []string{fake.CharactersN(8)},
	}
	metricMap := &MetricMap{
		MetricMap: map[string]*MetricDefinition{fake.CharactersN(6): metricDefinition},
		UI:        uiGroupData,
	}
	original := &IngestionDictionary{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(IngestionDictionaryType),
		Metrics:               map[string]*MetricMap{fake.CharactersN(8): metricMap},
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "metrics", "createdTimestamp", "lastModifiedTimestamp"}

	test.RunSerializationTest(t, original, &IngestionDictionary{}, original.ID, attrKeys)
}

func TestValidTypesSerialization(t *testing.T) {
	monitoredObjectTypesMap := map[string]string{
		fake.CharactersN(8): fake.CharactersN(10),
		fake.CharactersN(8): fake.CharactersN(10),
	}
	monitoredObjectDeviceTypesMap := map[string]string{
		fake.CharactersN(8): fake.CharactersN(10),
		fake.CharactersN(8): fake.CharactersN(10),
	}
	original := &ValidTypes{
		ID:                         uuid.NewV4().String(),
		REV:                        uuid.NewV4().String(),
		Datatype:                   string(ValidTypesType),
		MonitoredObjectTypes:       monitoredObjectTypesMap,
		MonitoredObjectDeviceTypes: monitoredObjectDeviceTypesMap,
		CreatedTimestamp:           time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "monitoredObjectTypes", "monitoredObjectDeviceTypes", "createdTimestamp", "lastModifiedTimestamp"}

	test.RunSerializationTest(t, original, &ValidTypes{}, original.ID, attrKeys)
}
