package handlers_test

import (
	"fmt"
	"net/http"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/handlers"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	"github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/icrowley/fake"
	"github.com/spf13/viper"
)

const (
	adminDBName = "adh-admin"
)

var (
	adminDB  datastore.AdminServiceDatastore
	tenantDB datastore.TenantServiceDatastore
)

func setupTestDatastore() error {
	// Setup Test Env
	gather.LoadConfig("../config/adh-gather-test.yml", viper.New())
	monitoring.InitMetrics()

	var err error
	adminDB, err = handlers.GetAdminServiceDatastore()
	if err != nil {
		return fmt.Errorf("Unable to instantiate Admin Service DAO: %s", err.Error())
	}

	tenantDB, err = handlers.GetTenantServiceDatastore()
	if err != nil {
		return fmt.Errorf("Unable to instantiate Tenant Service DAO: %s", err.Error())
	}

	adminDB.DeleteDatabase(adminDBName)

	_, err = adminDB.CreateDatabase(adminDBName)
	return err
}

func destroyTestDatastore() error {
	return adminDB.DeleteDatabase(adminDBName)
}

func getRandomTenantDescriptor() *admmod.Tenant {
	return &admmod.Tenant{
		Name:         fake.CharactersN(12),
		URLSubdomain: fake.DomainName(),
		State:        string(common.UserActive),
	}
}

func createRandomDataCleaningProfileCreateRequest() *swagmodels.DataCleaningProfileCreateRequestData {
	dcpType := "dataCleaningProfiles"
	return &swagmodels.DataCleaningProfileCreateRequestData{
		Type: &dcpType,
		Attributes: &swagmodels.DataCleaningProfileCreateRequestDataAttributes{
			Rules: []*swagmodels.DataCleaningRule{createRandomDataCleaningRule()},
		},
	}
}

func createRandomDataCleaningRule() *swagmodels.DataCleaningRule {
	vend := fake.CharactersN(12)
	label := fake.CharactersN(8)
	return &swagmodels.DataCleaningRule{
		MetricVendor:     &vend,
		MetricLabel:      &label,
		TriggerCondition: createRandomDataCleaningProfileRuleCondition(),
		ClearCondition:   createRandomDataCleaningProfileRuleCondition(),
	}
}

func createRandomDataCleaningProfileRuleCondition() *swagmodels.DataCleaningCondition {
	comp := fake.CharactersN(3)
	dur := fake.CharactersN(6)
	val := fake.CharactersN(6)
	agg := fake.CharactersN(3)
	return &swagmodels.DataCleaningCondition{
		Comparator:     &comp,
		Duration:       &dur,
		Value:          &val,
		ValueAggregate: &agg,
	}
}

func createHttpRequest(tenantID string, roles string) *http.Request {
	req := http.Request{
		Header: make(http.Header),
	}
	req.Header.Add(handlers.XFwdTenantId, tenantID)
	req.Header.Add(handlers.XFwdUserRoles, roles)
	return &req
}
