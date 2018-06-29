package inMemory

import (
	"fmt"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/getlantern/deepcopy"
	uuid "github.com/satori/go.uuid"

	"github.com/accedian/adh-gather/models/common"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// TenantServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type TenantServiceDatastoreInMemory struct {
	tenantToIDtoTenantUserMap                 map[string]map[string]*tenmod.User
	tenantToIDtoTenantDomainMap               map[string]map[string]*tenmod.Domain
	tenantToIDtoTenantConnectorConfigMap      map[string]map[string]*tenmod.ConnectorConfig
	tenantToIDtoTenantConnectorInstanceMap    map[string]map[string]*tenmod.ConnectorInstance
	tenantToIDtoTenantMonitoredObjectMap      map[string]map[string]*tenmod.MonitoredObject
	tenantToIDtoTenantThrPrfMap               map[string]map[string]*tenmod.ThresholdProfile
	tenantToIDtoTenantReportScheduleConfigMap map[string]map[string]*metmod.ReportScheduleConfig
	tenantToIDtoTenantSLAReportMap            map[string]map[string]*metmod.SLAReport
	tenantToIDtoDashboardMap                  map[string]map[string]*tenmod.Dashboard

	tenantIDtoMetaSlice   map[string][]*tenmod.Metadata
	tenantIDtoIngPrfSlice map[string][]*tenmod.IngestionProfile
}

// CreateTenantServiceDAO - returns an in-memory implementation of the Tenant Service
// datastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreInMemory, error) {
	res := new(TenantServiceDatastoreInMemory)

	res.tenantToIDtoTenantUserMap = map[string]map[string]*tenmod.User{}
	res.tenantToIDtoTenantDomainMap = map[string]map[string]*tenmod.Domain{}
	res.tenantToIDtoTenantConnectorConfigMap = map[string]map[string]*tenmod.ConnectorConfig{}
	res.tenantToIDtoTenantMonitoredObjectMap = map[string]map[string]*tenmod.MonitoredObject{}
	res.tenantToIDtoTenantThrPrfMap = map[string]map[string]*tenmod.ThresholdProfile{}
	res.tenantToIDtoTenantSLAReportMap = map[string]map[string]*metmod.SLAReport{}
	res.tenantToIDtoDashboardMap = map[string]map[string]*tenmod.Dashboard{}

	res.tenantIDtoMetaSlice = map[string][]*tenmod.Metadata{}
	res.tenantIDtoIngPrfSlice = map[string][]*tenmod.IngestionProfile{}

	return res, nil
}

func (tsd *TenantServiceDatastoreInMemory) GetConnectorConfigUpdateChan() chan *tenmod.ConnectorConfig {
	return nil
}

// DoesTenantExist - helper function to determine if a Tenant does have data stored for a particular type of data.
func (tsd *TenantServiceDatastoreInMemory) DoesTenantExist(tenantID string, ctx tenmod.TenantDataType) error {
	if len(tenantID) == 0 {
		return fmt.Errorf("%s does not exist", tenantID)
	}

	tenantDNE := fmt.Errorf("%s does not exist", tenantID)
	switch ctx {
	case tenmod.TenantUserType:
		if tsd.tenantToIDtoTenantUserMap[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantDomainType:
		if tsd.tenantToIDtoTenantDomainMap[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantConnectorConfigType:
		if tsd.tenantToIDtoTenantConnectorConfigMap[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantMonitoredObjectType:
		if tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantMetaType:
		if tsd.tenantIDtoMetaSlice[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantIngestionProfileType:
		if tsd.tenantIDtoIngPrfSlice[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantThresholdProfileType:
		if tsd.tenantToIDtoTenantThrPrfMap[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantReportScheduleConfigType:
		if tsd.tenantToIDtoTenantReportScheduleConfigMap[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantSLAReportStr:
		if tsd.tenantToIDtoTenantSLAReportMap[tenantID] == nil {
			return tenantDNE
		}
	case tenmod.TenantDashboardType:
		if tsd.tenantToIDtoDashboardMap[tenantID] == nil {
			return tenantDNE
		}
	default:
		return fmt.Errorf("Invalid data type %s provided", string(ctx))
	}

	return nil
}

// CreateTenantUser - InMemory implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreInMemory) CreateTenantUser(tenantUserRequest *tenmod.User) (*tenmod.User, error) {
	if err := tsd.DoesTenantExist(tenantUserRequest.TenantID, tenmod.TenantUserType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantUserMap[tenantUserRequest.TenantID] = map[string]*tenmod.User{}
	}

	userCopy := tenmod.User{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.ID = uuid.NewV4().String()
	userCopy.REV = uuid.NewV4().String()
	userCopy.Datatype = string(tenmod.TenantUserType)
	userCopy.CreatedTimestamp = ds.MakeTimestamp()
	userCopy.LastModifiedTimestamp = userCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.TenantID][userCopy.ID] = &userCopy

	return &userCopy, nil
}

// UpdateTenantUser - InMemory implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantUser(tenantUserRequest *tenmod.User) (*tenmod.User, error) {
	if len(tenantUserRequest.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantUserStr)
	}
	if len(tenantUserRequest.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantUserStr)
	}
	if err := tsd.DoesTenantExist(tenantUserRequest.TenantID, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	userCopy := tenmod.User{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.REV = uuid.NewV4().String()
	userCopy.Datatype = string(tenmod.TenantUserType)
	userCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.TenantID][userCopy.ID] = &userCopy

	return &userCopy, nil
}

// DeleteTenantUser - InMemory implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantUser(tenantID string, userID string) (*tenmod.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", tenmod.TenantUserStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantUserStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	user, ok := tsd.tenantToIDtoTenantUserMap[tenantID][userID]
	if ok {
		delete(tsd.tenantToIDtoTenantUserMap[tenantID], userID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantUserMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantUserMap, tenantID)
		}
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantUserStr)
}

// GetTenantUser - InMemory implementation of GetTenantUser
func (tsd *TenantServiceDatastoreInMemory) GetTenantUser(tenantID string, userID string) (*tenmod.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", tenmod.TenantUserStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantUserStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	user, ok := tsd.tenantToIDtoTenantUserMap[tenantID][userID]
	if ok {
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantUserStr)
}

// GetAllTenantUsers - InMemory implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantUsers(tenantID string) ([]*tenmod.User, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantUserType)
	if err != nil {
		return []*tenmod.User{}, nil
	}

	tenantUserList := make([]*tenmod.User, 0)

	for _, user := range tsd.tenantToIDtoTenantUserMap[tenantID] {
		tenantUserList = append(tenantUserList, user)
	}

	return tenantUserList, nil
}

// CreateTenantConnectorInstance - InMemory implementation of CreateTenantConnectorInstance
func (tsd *TenantServiceDatastoreInMemory) CreateTenantConnectorInstance(tenantConnectorInstanceRequest *tenmod.ConnectorInstance) (*tenmod.ConnectorInstance, error) {
	if err := tsd.DoesTenantExist(tenantConnectorInstanceRequest.TenantID, tenmod.TenantConnectorInstanceType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantConnectorInstanceMap[tenantConnectorInstanceRequest.TenantID] = map[string]*tenmod.ConnectorInstance{}
	}

	recCopy := tenmod.ConnectorInstance{}
	deepcopy.Copy(&recCopy, tenantConnectorInstanceRequest)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantConnectorInstanceType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantConnectorInstanceMap[tenantConnectorInstanceRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// UpdateTenantConnectorInstance - InMemory implementation of UpdateTenantConnectorInstance
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantConnectorInstance(tenantConnectorInstanceRequest *tenmod.ConnectorInstance) (*tenmod.ConnectorInstance, error) {
	if len(tenantConnectorInstanceRequest.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantConnectorInstanceStr)
	}
	if len(tenantConnectorInstanceRequest.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantConnectorInstanceStr)
	}
	if err := tsd.DoesTenantExist(tenantConnectorInstanceRequest.TenantID, tenmod.TenantConnectorInstanceType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantConnectorInstanceStr)
	}

	recCopy := tenmod.ConnectorInstance{}
	deepcopy.Copy(&recCopy, tenantConnectorInstanceRequest)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantConnectorInstanceType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantConnectorInstanceMap[tenantConnectorInstanceRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// DeleteTenantConnectorInstance - InMemory implementation of DeleteTenantConnectorInstance
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantConnectorInstance(tenantID string, dataID string) (*tenmod.ConnectorInstance, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a ConnectorInstance ID", tenmod.TenantConnectorInstanceStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantConnectorInstanceStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantConnectorInstanceType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantConnectorInstanceStr)
	}

	rec, ok := tsd.tenantToIDtoTenantConnectorInstanceMap[tenantID][dataID]
	logger.Log.Debugf(models.AsJSONString(tsd.tenantToIDtoTenantConnectorInstanceMap))
	if ok {
		delete(tsd.tenantToIDtoTenantConnectorInstanceMap[tenantID], dataID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantConnectorInstanceMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantConnectorInstanceMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantConnectorInstanceStr)
}

// GetTenantConnectorInstance - InMemory implementation of GetTenantConnectorInstance
func (tsd *TenantServiceDatastoreInMemory) GetTenantConnectorInstance(tenantID string, dataID string) (*tenmod.ConnectorInstance, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a ConnectorInstance ID", tenmod.TenantConnectorInstanceStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantConnectorInstanceStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantConnectorInstanceType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantConnectorInstanceStr)
	}

	rec, ok := tsd.tenantToIDtoTenantConnectorInstanceMap[tenantID][dataID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantConnectorInstanceStr)
}

// GetAllTenantConnectorInstances - InMemory implementation of GetAllTenantConnectorInstances
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantConnectorInstances(tenantID string) ([]*tenmod.ConnectorInstance, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantConnectorInstanceType)
	if err != nil {
		return []*tenmod.ConnectorInstance{}, nil
	}

	recList := make([]*tenmod.ConnectorInstance, 0)

	for _, rec := range tsd.tenantToIDtoTenantConnectorInstanceMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

// CreateTenantConnector - InMemory implementation of CreateTenantConnectorConfig
func (tsd *TenantServiceDatastoreInMemory) CreateTenantConnectorConfig(TenantConnectorConfigRequest *tenmod.ConnectorConfig) (*tenmod.ConnectorConfig, error) {
	if err := tsd.DoesTenantExist(TenantConnectorConfigRequest.TenantID, tenmod.TenantConnectorConfigType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantConnectorConfigMap[TenantConnectorConfigRequest.TenantID] = map[string]*tenmod.ConnectorConfig{}
	}

	recCopy := tenmod.ConnectorConfig{}
	deepcopy.Copy(&recCopy, TenantConnectorConfigRequest)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantConnectorConfigType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantConnectorConfigMap[TenantConnectorConfigRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// UpdateTenantConnectorConfig - InMemory implementation of UpdateTenantConnectorConfig
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantConnectorConfig(TenantConnectorConfigRequest *tenmod.ConnectorConfig) (*tenmod.ConnectorConfig, error) {
	if len(TenantConnectorConfigRequest.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantConnectorConfigStr)
	}
	if len(TenantConnectorConfigRequest.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantConnectorConfigStr)
	}
	if err := tsd.DoesTenantExist(TenantConnectorConfigRequest.TenantID, tenmod.TenantConnectorConfigType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantConnectorConfigStr)
	}

	recCopy := tenmod.ConnectorConfig{}
	deepcopy.Copy(&recCopy, TenantConnectorConfigRequest)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantConnectorConfigType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantConnectorConfigMap[TenantConnectorConfigRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// DeleteTenantConnectorConfig - InMemory implementation of DeleteTenantConnectorConfig
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantConnectorConfig(tenantID string, dataID string) (*tenmod.ConnectorConfig, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Connector ID", tenmod.TenantConnectorConfigStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantConnectorConfigStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantConnectorConfigType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantConnectorConfigStr)
	}

	rec, ok := tsd.tenantToIDtoTenantConnectorConfigMap[tenantID][dataID]
	logger.Log.Debugf(models.AsJSONString(tsd.tenantToIDtoTenantConnectorConfigMap))
	if ok {
		delete(tsd.tenantToIDtoTenantConnectorConfigMap[tenantID], dataID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantConnectorConfigMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantConnectorConfigMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantConnectorConfigStr)
}

// GetTenantConnectorConfig - InMemory implementation of GetTenantConnectorConfig
func (tsd *TenantServiceDatastoreInMemory) GetTenantConnectorConfig(tenantID string, dataID string) (*tenmod.ConnectorConfig, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Connector ID", tenmod.TenantConnectorConfigStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantConnectorConfigStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantConnectorConfigType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantConnectorConfigStr)
	}

	rec, ok := tsd.tenantToIDtoTenantConnectorConfigMap[tenantID][dataID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantConnectorConfigStr)
}

// GetAllTenantConnectorConfigs - InMemory implementation of GetAllTenantConnectorConfigs
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantConnectorConfigs(tenantID, zone string) ([]*tenmod.ConnectorConfig, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantConnectorConfigType)
	if err != nil {
		return []*tenmod.ConnectorConfig{}, nil
	}

	recList := make([]*tenmod.ConnectorConfig, 0)

	for _, rec := range tsd.tenantToIDtoTenantConnectorConfigMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

// GetAllAvailableTenantConnectorConfigs - Returns all tenant connectors matching tenantID, zone, that aren't already being used
func (tsd *TenantServiceDatastoreInMemory) GetAllAvailableTenantConnectorConfigs(tenantID, zone string) ([]*tenmod.ConnectorConfig, error) {
	return nil, nil
}

// GetAllTenantConnectorConfigsByInstanceID - Returns the TenantConnectorConfigConfigs with the given instance ID
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantConnectorConfigsByInstanceID(tenantID, instanceID string) ([]*tenmod.ConnectorConfig, error) {
	return nil, nil
}

// CreateTenantDomain - InMemory implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) CreateTenantDomain(tenantDomainRequest *tenmod.Domain) (*tenmod.Domain, error) {
	if err := tsd.DoesTenantExist(tenantDomainRequest.TenantID, tenmod.TenantDomainType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.TenantID] = map[string]*tenmod.Domain{}
	}

	recCopy := tenmod.Domain{}
	deepcopy.Copy(&recCopy, tenantDomainRequest)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantDomainType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// UpdateTenantDomain - InMemory implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantDomain(tenantDomainRequest *tenmod.Domain) (*tenmod.Domain, error) {
	if len(tenantDomainRequest.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantDomainStr)
	}
	if len(tenantDomainRequest.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantDomainStr)
	}
	if err := tsd.DoesTenantExist(tenantDomainRequest.TenantID, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	recCopy := tenmod.Domain{}
	deepcopy.Copy(&recCopy, tenantDomainRequest)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantDomainType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// DeleteTenantDomain - InMemory implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantDomainStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDomainStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	rec, ok := tsd.tenantToIDtoTenantDomainMap[tenantID][dataID]
	logger.Log.Debugf(models.AsJSONString(tsd.tenantToIDtoTenantDomainMap))
	if ok {
		delete(tsd.tenantToIDtoTenantDomainMap[tenantID], dataID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantDomainMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantDomainMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantDomainStr)
}

// GetTenantDomain - InMemory implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreInMemory) GetTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantDomainStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDomainStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	rec, ok := tsd.tenantToIDtoTenantDomainMap[tenantID][dataID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantDomainStr)
}

// GetAllTenantDomains - InMemory implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantDomains(tenantID string) ([]*tenmod.Domain, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantDomainType)
	if err != nil {
		return []*tenmod.Domain{}, nil
	}

	recList := make([]*tenmod.Domain, 0)

	for _, rec := range tsd.tenantToIDtoTenantDomainMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

// CreateTenantIngestionProfile - InMemory implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantIngestionProfile(tenantIngPrfReq *tenmod.IngestionProfile) (*tenmod.IngestionProfile, error) {
	if err := tsd.DoesTenantExist(tenantIngPrfReq.TenantID, tenmod.TenantIngestionProfileType); err != nil {
		// Make a place for the tenant
		tsd.tenantIDtoIngPrfSlice[tenantIngPrfReq.TenantID] = make([]*tenmod.IngestionProfile, 1)
	}

	existing, _ := tsd.GetActiveTenantIngestionProfile(tenantIngPrfReq.TenantID)
	if existing != nil {
		return nil, fmt.Errorf("Unable to create %s, it already exists", tenmod.TenantIngestionProfileStr)
	}

	recCopy := tenmod.IngestionProfile{}
	deepcopy.Copy(&recCopy, tenantIngPrfReq)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantIngestionProfileType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantIDtoIngPrfSlice[tenantIngPrfReq.TenantID][0] = &recCopy

	return &recCopy, nil
}

// UpdateTenantIngestionProfile - InMemory implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantIngestionProfile(tenantIngPrfReq *tenmod.IngestionProfile) (*tenmod.IngestionProfile, error) {
	if len(tenantIngPrfReq.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantIngestionProfileStr)
	}
	if len(tenantIngPrfReq.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantIngestionProfileStr)
	}
	if err := tsd.DoesTenantExist(tenantIngPrfReq.TenantID, tenmod.TenantIngestionProfileType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantIngestionProfileStr)
	}

	recCopy := tenmod.IngestionProfile{}
	deepcopy.Copy(&recCopy, tenantIngPrfReq)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantIngestionProfileType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantIDtoIngPrfSlice[tenantIngPrfReq.TenantID][0] = &recCopy

	return &recCopy, nil
}

// GetTenantIngestionProfile - InMemory implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantIngestionProfileStr)
	}
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Ingestion Proile ID", tenmod.TenantIngestionProfileStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantIngestionProfileType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantIngestionProfileStr)
	}
	existing := tsd.tenantIDtoIngPrfSlice[tenantID][0]
	if dataID != existing.ID {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantIngestionProfileStr)
	}

	return existing, nil
}

// DeleteTenantIngestionProfile - InMemory implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantIngestionProfileStr)
	}
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Ingestion Proile ID", tenmod.TenantIngestionProfileStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantIngestionProfileType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantIngestionProfileStr)
	}
	existing, err := tsd.GetActiveTenantIngestionProfile(tenantID)
	if err != nil {
		return nil, err
	}
	if dataID != existing.ID {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantIngestionProfileStr)
	}

	tsd.tenantIDtoIngPrfSlice[tenantID][0] = nil

	return existing, nil
}

// CreateTenantThresholdProfile - InMemory implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantThresholdProfile(tenantThreshPrfReq *tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error) {
	if err := tsd.DoesTenantExist(tenantThreshPrfReq.TenantID, tenmod.TenantThresholdProfileType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantThrPrfMap[tenantThreshPrfReq.TenantID] = map[string]*tenmod.ThresholdProfile{}
	}

	recCopy := tenmod.ThresholdProfile{}
	deepcopy.Copy(&recCopy, tenantThreshPrfReq)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantThresholdProfileType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantThrPrfMap[tenantThreshPrfReq.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// UpdateTenantThresholdProfile - InMemory implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantThresholdProfile(tenantThreshPrfReq *tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error) {
	if len(tenantThreshPrfReq.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantThresholdProfileStr)
	}
	if len(tenantThreshPrfReq.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantThresholdProfileStr)
	}
	if err := tsd.DoesTenantExist(tenantThreshPrfReq.TenantID, tenmod.TenantThresholdProfileType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantThresholdProfileStr)
	}

	recCopy := tenmod.ThresholdProfile{}
	deepcopy.Copy(&recCopy, tenantThreshPrfReq)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantThresholdProfileType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantThrPrfMap[tenantThreshPrfReq.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// GetTenantThresholdProfile - InMemory implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantDomainStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDomainStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantThresholdProfileType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantThresholdProfileStr)
	}

	rec, ok := tsd.tenantToIDtoTenantThrPrfMap[tenantID][dataID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantThresholdProfileStr)
}

// DeleteTenantThresholdProfile - InMemory implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantThresholdProfileStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantThresholdProfileStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantThresholdProfileType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantThresholdProfileStr)
	}

	rec, ok := tsd.tenantToIDtoTenantThrPrfMap[tenantID][dataID]
	if ok {
		delete(tsd.tenantToIDtoTenantThrPrfMap[tenantID], dataID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantThrPrfMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantThrPrfMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantDomainStr)
}

// CreateMonitoredObject - InMemory implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) CreateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error) {
	if err := tsd.DoesTenantExist(monitoredObjectReq.TenantID, tenmod.TenantMonitoredObjectType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantMonitoredObjectMap[monitoredObjectReq.TenantID] = map[string]*tenmod.MonitoredObject{}
	}

	recCopy := tenmod.MonitoredObject{}
	deepcopy.Copy(&recCopy, monitoredObjectReq)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantMonitoredObjectType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantMonitoredObjectMap[monitoredObjectReq.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// UpdateMonitoredObject - InMemory implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) UpdateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error) {
	if len(monitoredObjectReq.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantMonitoredObjectStr)
	}
	if len(monitoredObjectReq.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantMonitoredObjectStr)
	}
	if err := tsd.DoesTenantExist(monitoredObjectReq.TenantID, tenmod.TenantMonitoredObjectType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantMonitoredObjectStr)
	}

	recCopy := tenmod.MonitoredObject{}
	deepcopy.Copy(&recCopy, monitoredObjectReq)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantMonitoredObjectType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantMonitoredObjectMap[monitoredObjectReq.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// GetMonitoredObject - InMemory implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantMonitoredObjectStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantMonitoredObjectStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantMonitoredObjectType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantMonitoredObjectStr)
	}

	rec, ok := tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID][dataID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantMonitoredObjectStr)
}

// DeleteMonitoredObject - InMemory implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) DeleteMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantMonitoredObjectStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantMonitoredObjectStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantMonitoredObjectType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantMonitoredObjectStr)
	}

	rec, ok := tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID][dataID]
	if ok {
		delete(tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID], dataID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantMonitoredObjectMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantMonitoredObjectStr)
}

// GetAllMonitoredObjects - InMemory implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) GetAllMonitoredObjects(tenantID string) ([]*tenmod.MonitoredObject, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantMonitoredObjectType)
	if err != nil {
		return []*tenmod.MonitoredObject{}, nil
	}

	recList := make([]*tenmod.MonitoredObject, 0)

	for _, rec := range tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

// GetAllMonitoredObjectsInIDList - InMemory implementation of GetAllMonitoredObjectsInIDList
func (tsd *TenantServiceDatastoreInMemory) GetAllMonitoredObjectsInIDList(tenantID string, idList []string) ([]*tenmod.MonitoredObject, error) {
	if idList == nil || len(idList) == 0 {
		return []*tenmod.MonitoredObject{}, nil
	}

	recList := make([]*tenmod.MonitoredObject, 0)

	for _, rec := range tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID] {
		if doesSliceContainString(idList, rec.ID) {
			recList = append(recList, rec)
		}
	}

	return recList, nil
}

// GetMonitoredObjectToDomainMap - InMemory implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObjectToDomainMap(moByDomReq *tenmod.MonitoredObjectCountByDomainRequest) (*tenmod.MonitoredObjectCountByDomainResponse, error) {
	err := tsd.DoesTenantExist(moByDomReq.TenantID, tenmod.TenantMonitoredObjectType)
	if err != nil {
		return nil, fmt.Errorf("Tenant does not have any Monitored Objects")
	}
	err = tsd.DoesTenantExist(moByDomReq.TenantID, tenmod.TenantDomainType)
	if err != nil {
		return nil, fmt.Errorf("Tenant does not have any Domains")
	}

	// Get response data either by subset, or for all domains
	response := tenmod.MonitoredObjectCountByDomainResponse{}
	domainSet := moByDomReq.DomainSet
	if domainSet == nil || len(domainSet) == 0 {
		// Retrieve values for all domains
		response.DomainToMonitoredObjectSetMap = map[string][]string{}
		for _, mo := range tsd.tenantToIDtoTenantMonitoredObjectMap[moByDomReq.TenantID] {
			for _, dom := range mo.DomainSet {
				if response.DomainToMonitoredObjectSetMap[dom] == nil {
					response.DomainToMonitoredObjectSetMap[dom] = []string{mo.ID}
					continue
				}
				if !contains(response.DomainToMonitoredObjectSetMap[dom], mo.ID) {
					response.DomainToMonitoredObjectSetMap[dom] = append(response.DomainToMonitoredObjectSetMap[dom], mo.ID)
				}
			}

		}
	} else {
		// Retrieve just the subset of values.
		response.DomainToMonitoredObjectSetMap = map[string][]string{}
		for _, mo := range tsd.tenantToIDtoTenantMonitoredObjectMap[moByDomReq.TenantID] {
			for _, dom := range mo.DomainSet {
				if !contains(domainSet, dom) {
					continue
				}
				if response.DomainToMonitoredObjectSetMap[dom] == nil {
					response.DomainToMonitoredObjectSetMap[dom] = []string{mo.ID}
					continue
				}
				if !contains(response.DomainToMonitoredObjectSetMap[dom], mo.ID) {
					response.DomainToMonitoredObjectSetMap[dom] = append(response.DomainToMonitoredObjectSetMap[dom], mo.ID)
				}
			}

		}
	}

	// Filtering is done, do the count
	if moByDomReq.ByCount {
		response.DomainToMonitoredObjectCountMap = map[string]int64{}
		for key, val := range response.DomainToMonitoredObjectSetMap {
			response.DomainToMonitoredObjectCountMap[key] = int64(len(val))
		}

		// Null out the other response
		response.DomainToMonitoredObjectSetMap = nil
	}

	return &response, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// CreateTenantMeta - InMemory implementation of CreateTenantMeta
func (tsd *TenantServiceDatastoreInMemory) CreateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error) {
	if err := tsd.DoesTenantExist(meta.TenantID, tenmod.TenantMetaType); err != nil {
		// Make a place for the tenant
		tsd.tenantIDtoMetaSlice[meta.TenantID] = make([]*tenmod.Metadata, 1)
	}

	existing, _ := tsd.GetTenantMeta(meta.TenantID)
	if existing != nil {
		return nil, fmt.Errorf("Unable to create %s, it already exists", tenmod.TenantMetaStr)
	}

	recCopy := tenmod.Metadata{}
	deepcopy.Copy(&recCopy, meta)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantMetaType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantIDtoMetaSlice[meta.TenantID][0] = &recCopy

	return &recCopy, nil
}

// UpdateTenantMeta - InMemory implementation of UpdateTenantMeta
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error) {
	if len(meta.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantMetaStr)
	}
	if len(meta.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantMetaStr)
	}
	if err := tsd.DoesTenantExist(meta.TenantID, tenmod.TenantMetaType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantMetaStr)
	}

	recCopy := tenmod.Metadata{}
	deepcopy.Copy(&recCopy, meta)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantMetaType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantIDtoMetaSlice[meta.TenantID][0] = &recCopy

	return &recCopy, nil
}

// DeleteTenantMeta - InMemory implementation of DeleteTenantMeta
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantMeta(tenantID string) (*tenmod.Metadata, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantMetaStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantMetaType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantMetaStr)
	}
	existing, err := tsd.GetTenantMeta(tenantID)
	if err != nil {
		return nil, err
	}

	tsd.tenantIDtoMetaSlice[tenantID][0] = nil

	return existing, nil
}

// GetTenantMeta - InMemory implementation of GetTenantMeta
func (tsd *TenantServiceDatastoreInMemory) GetTenantMeta(tenantID string) (*tenmod.Metadata, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantMonitoredObjectStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantMetaType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantMetaStr)
	}

	if tsd.tenantIDtoMetaSlice[tenantID][0] == nil {
		return nil, fmt.Errorf("%s not found", tenmod.TenantMetaStr)
	}

	return tsd.tenantIDtoMetaSlice[tenantID][0], nil
}

// GetActiveTenantIngestionProfile - InMemory implementation of GetActiveTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) GetActiveTenantIngestionProfile(tenantID string) (*tenmod.IngestionProfile, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantIngestionProfileStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantIngestionProfileType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantIngestionProfileStr)
	}

	if tsd.tenantIDtoIngPrfSlice[tenantID][0] == nil {
		return nil, fmt.Errorf("%s not found", tenmod.TenantIngestionProfileStr)
	}

	return tsd.tenantIDtoIngPrfSlice[tenantID][0], nil
}

// GetAllTenantThresholdProfile - InMemory implementation of GetAllTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantThresholdProfile(tenantID string) ([]*tenmod.ThresholdProfile, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantThresholdProfileType)
	if err != nil {
		return []*tenmod.ThresholdProfile{}, nil
	}

	recList := make([]*tenmod.ThresholdProfile, 0)

	for _, rec := range tsd.tenantToIDtoTenantThrPrfMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

// BulkInsertMonitoredObjects - InMemory implementation of BulkInsertMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) BulkInsertMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error) {
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantMonitoredObjectType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantMonitoredObjectMap[tenantID] = map[string]*tenmod.MonitoredObject{}
	}

	result := make([]*common.BulkOperationResult, 0)
	for _, val := range value {
		created, err := tsd.CreateMonitoredObject(val)
		if err != nil {
			entry := common.BulkOperationResult{
				OK:     false,
				REASON: err.Error(),
			}
			result = append(result, &entry)
		} else {
			entry := common.BulkOperationResult{
				OK: true,
				ID: created.ID,
			}
			result = append(result, &entry)
		}
	}

	return result, nil
}

// BulkUpdatetMonitoredObjects - InMemory implementation of BulkUpdatetMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) BulkUpdateMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error) {

	result := make([]*common.BulkOperationResult, 0)
	for _, val := range value {
		created, err := tsd.UpdateMonitoredObject(val)
		if err != nil {
			entry := common.BulkOperationResult{
				OK:     false,
				REASON: err.Error(),
			}
			result = append(result, &entry)
		} else {
			entry := common.BulkOperationResult{
				OK: true,
				ID: created.ID,
			}
			result = append(result, &entry)
		}
	}

	return result, nil
}

func (tsd *TenantServiceDatastoreInMemory) CreateReportScheduleConfig(config *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error) {
	if err := tsd.DoesTenantExist(config.TenantID, tenmod.TenantReportScheduleConfigType); err != nil {
		tsd.tenantToIDtoTenantReportScheduleConfigMap[config.TenantID] = map[string]*metmod.ReportScheduleConfig{}
	}

	recCopy := metmod.ReportScheduleConfig{}
	deepcopy.Copy(&recCopy, config)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantReportScheduleConfigMap[config.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}
func (tsd *TenantServiceDatastoreInMemory) UpdateReportScheduleConfig(config *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error) {
	if len(config.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantReportScheduleConfigStr)
	}
	if len(config.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantReportScheduleConfigStr)
	}
	if err := tsd.DoesTenantExist(config.TenantID, tenmod.TenantReportScheduleConfigType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantReportScheduleConfigStr)
	}

	recCopy := metmod.ReportScheduleConfig{}
	deepcopy.Copy(&recCopy, config)
	recCopy.REV = uuid.NewV4().String()
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantReportScheduleConfigMap[config.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}
func (tsd *TenantServiceDatastoreInMemory) DeleteReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error) {
	if len(configID) == 0 {
		return nil, fmt.Errorf("%s must provide a Report schedule config ID", tenmod.TenantReportScheduleConfigStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantReportScheduleConfigStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantReportScheduleConfigType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantReportScheduleConfigStr)
	}

	rec, ok := tsd.tenantToIDtoTenantReportScheduleConfigMap[tenantID][configID]
	logger.Log.Debugf(models.AsJSONString(tsd.tenantToIDtoTenantReportScheduleConfigMap))
	if ok {
		delete(tsd.tenantToIDtoTenantReportScheduleConfigMap[tenantID], configID)

		if len(tsd.tenantToIDtoTenantReportScheduleConfigMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantReportScheduleConfigMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantReportScheduleConfigStr)
}
func (tsd *TenantServiceDatastoreInMemory) GetReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error) {
	if len(configID) == 0 {
		return nil, fmt.Errorf("%s must provide a Report schedule config ID", tenmod.TenantReportScheduleConfigStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantReportScheduleConfigStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantReportScheduleConfigType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantReportScheduleConfigStr)
	}

	rec, ok := tsd.tenantToIDtoTenantReportScheduleConfigMap[tenantID][configID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantReportScheduleConfigStr)
}
func (tsd *TenantServiceDatastoreInMemory) GetAllReportScheduleConfigs(tenantID string) ([]*metmod.ReportScheduleConfig, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantReportScheduleConfigType)
	if err != nil {
		return []*metmod.ReportScheduleConfig{}, nil
	}

	recList := make([]*metmod.ReportScheduleConfig, 0)

	for _, rec := range tsd.tenantToIDtoTenantReportScheduleConfigMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

func (tsd *TenantServiceDatastoreInMemory) CreateSLAReport(slaReport *metmod.SLAReport) (*metmod.SLAReport, error) {
	if err := tsd.DoesTenantExist(slaReport.TenantID, tenmod.TenantReportType); err != nil {
		tsd.tenantToIDtoTenantSLAReportMap[slaReport.TenantID] = map[string]*metmod.SLAReport{}
	}

	recCopy := metmod.SLAReport{}
	deepcopy.Copy(&recCopy, slaReport)
	recCopy.ID = uuid.NewV4().String()

	tsd.tenantToIDtoTenantSLAReportMap[slaReport.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

func (tsd *TenantServiceDatastoreInMemory) DeleteSLAReport(tenantID string, slaReportID string) (*metmod.SLAReport, error) {
	if len(slaReportID) == 0 {
		return nil, fmt.Errorf("%s must provide a SLA Report ID", tenmod.TenantSLAReportStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantSLAReportStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantReportType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantSLAReportStr)
	}

	rec, ok := tsd.tenantToIDtoTenantSLAReportMap[tenantID][slaReportID]
	logger.Log.Debugf(models.AsJSONString(tsd.tenantToIDtoTenantSLAReportMap))
	if ok {
		delete(tsd.tenantToIDtoTenantSLAReportMap[tenantID], slaReportID)

		if len(tsd.tenantToIDtoTenantSLAReportMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantSLAReportMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantSLAReportStr)
}

func (tsd *TenantServiceDatastoreInMemory) GetSLAReport(tenantID string, slaReportID string) (*metmod.SLAReport, error) {
	if len(slaReportID) == 0 {
		return nil, fmt.Errorf("%s must provide a SLA report ID", tenmod.TenantSLAReportStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantSLAReportStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantReportType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantSLAReportStr)
	}

	rec, ok := tsd.tenantToIDtoTenantSLAReportMap[tenantID][slaReportID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantSLAReportStr)
}

func (tsd *TenantServiceDatastoreInMemory) GetAllSLAReports(tenantID string) ([]*metmod.SLAReport, error) {
	err := tsd.DoesTenantExist(tenantID, tenmod.TenantReportType)
	if err != nil {
		return []*metmod.SLAReport{}, nil
	}

	recList := make([]*metmod.SLAReport, 0)

	for _, rec := range tsd.tenantToIDtoTenantSLAReportMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

func (tsd *TenantServiceDatastoreInMemory) CreateDashboard(dashboard *tenmod.Dashboard) (*tenmod.Dashboard, error) {
	if err := tsd.DoesTenantExist(dashboard.TenantID, tenmod.TenantDashboardType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoDashboardMap[dashboard.TenantID] = map[string]*tenmod.Dashboard{}
	}

	recCopy := tenmod.Dashboard{}
	deepcopy.Copy(&recCopy, dashboard)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()

	tsd.tenantToIDtoDashboardMap[dashboard.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

func (tsd *TenantServiceDatastoreInMemory) DeleteDashboard(tenantID string, dataID string) (*tenmod.Dashboard, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Dashboard ID", tenmod.TenantDashboardStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDashboardStr)
	}
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantDashboardType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	rec, ok := tsd.tenantToIDtoDashboardMap[tenantID][dataID]
	logger.Log.Debugf(models.AsJSONString(tsd.tenantToIDtoDashboardMap))
	if ok {
		delete(tsd.tenantToIDtoDashboardMap[tenantID], dataID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoDashboardMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoDashboardMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantReportScheduleConfigStr)
}

func (tsd *TenantServiceDatastoreInMemory) HasDashboardsWithDomain(tenantID string, domainID string) (bool, error) {
	if err := tsd.DoesTenantExist(tenantID, tenmod.TenantDashboardType); err != nil {
		return false, nil
	}

	for _, dashboard := range tsd.tenantToIDtoDashboardMap[tenantID] {
		if dashboard == nil {
			continue
		}
		logger.Log.Debugf("Processing tenantID %s dashboard %s", tenantID, models.AsJSONString(dashboard))
		for _, v := range dashboard.DomainSet {
			if v == domainID {
				return true, nil
			}
		}
	}
	return false, nil

}

func doesSliceContainString(container []string, value string) bool {
	for _, s := range container {
		if s == value {
			return true
		}
	}
	return false
}

func (tsd *TenantServiceDatastoreInMemory) MonitoredObjectKeysUpdate(tenantID string, meta map[string]string) error {
	return nil
}

func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObjectByObjectName(name string, tenantID string) (*tenmod.MonitoredObject, error) {
	return nil, nil
}
