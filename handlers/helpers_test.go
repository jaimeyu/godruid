package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/handlers"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/scheduler"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/icrowley/fake"
	"github.com/spf13/viper"
)

const (
	adminDBName = "adh-admin"

	linksPrev  = "prev"
	linksFirst = "first"
	linksSelf  = "self"
	linksNext  = "next"
)

var (
	adminDB  datastore.AdminServiceDatastore
	tenantDB datastore.TenantServiceDatastore

	objectTypes = []string{string(tenmod.TwampPE), string(tenmod.TwampSF), string(tenmod.TwampSL), string(tenmod.Flowmeter)}
	deviceTypes = []string{string(tenmod.AccedianVNID), string(tenmod.AccedianNID)}
)

func setupTestDatastore() error {
	// Setup Test Env
	cfg := gather.LoadConfig("../config/adh-gather-test.yml", viper.New())
	monitoring.InitMetrics()

	cfg.Set("ingDict", "../files/defaultIngestionDictionary.json")
	cfg.Set("changeNotifications", "false")
	cfg.Set(gather.CK_args_authorizationAAA.String(), "true")

	handlers.InitializeAuthHelper()

	scheduler.Initialize(handlers.CreateMetricServiceHandler(), nil, nil, 5)

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
	if err != nil {
		return fmt.Errorf("Unable to create Admin DB: %s", err.Error())
	}

	err = adminDB.AddAdminViews()
	if err != nil {
		return fmt.Errorf("Unable to add Admin Views to Admin DB: %s", err.Error())
	}

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
	return createHttpRequestWithParams(tenantID, roles, "", "")
}

func createHttpRequestWithParams(tenantID string, roles string, url string, method string) *http.Request {
	if len(url) == 0 {
		url = "/i/am/made/up"
	}

	if len(method) == 0 {
		method = "GET"
	}

	req, _ := http.NewRequest(method, url, bytes.NewBufferString("whatever cause it was already read"))

	req.Header.Add(handlers.XFwdTenantId, tenantID)
	req.Header.Add(handlers.XFwdUserRoles, roles)
	return req
}
