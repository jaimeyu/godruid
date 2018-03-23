package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

const (
	VendorMap              = "vendorMap"
	DirectionMap           = "directionMap"
	MonitoredObjectTypeMap = "monitoredObjectTypeMap"
	MetricMap              = "metricMap"
	EventAttrMap           = "eventAttrMap"
	EventMap               = "eventMap"
	Enabled                = "enabled"
	Critical               = "critical"
	Major                  = "major"
	Minor                  = "minor"
)

var (
	// ValidMonitoredObjectTypes - known Monitored Object types in the system.
	ValidMonitoredObjectTypes = map[string]tenmod.MonitoredObjectType{
		"pe": tenmod.TwampPE,
		"sf": tenmod.TwampSF,
		"sl": tenmod.TwampSL,
		string(tenmod.TwampPE): tenmod.TwampPE,
		string(tenmod.TwampSF): tenmod.TwampSF,
		string(tenmod.TwampSL): tenmod.TwampSL}

	// ValidMonitoredObjectDeviceTypes - known Monitored Object Device types in the system.
	ValidMonitoredObjectDeviceTypes = map[string]tenmod.MonitoredObjectDeviceType{
		string(tenmod.AccedianNID):  tenmod.AccedianNID,
		string(tenmod.AccedianVNID): tenmod.AccedianVNID}
)

// GRPCServiceHandler - implementer of all gRPC Services. Offloads
// implementation details to each unique service handler. When new
// gRPC services are added, a new Service Handler should be created,
// and a pointer to that object should be added to this wrapper.
type GRPCServiceHandler struct {
	ash               *AdminServiceHandler
	tsh               *TenantServiceHandler
	DefaultValidTypes *admmod.ValidTypes
}

// CreateCoordinator - used to create a gRPC service handler wrapper
// that coordinates the logic to satisfy all gRPC service
// interfaces.
func CreateCoordinator() *GRPCServiceHandler {
	result := new(GRPCServiceHandler)

	result.ash = CreateAdminServiceHandler()
	result.tsh = CreateTenantServiceHandler()

	// Setup the known values of the Valid Types for the system
	// by using the enumerated protobuf values
	validMonObjTypes := make(map[string]string, 0)
	validMonObjDevTypes := make(map[string]string, 0)

	for key, val := range ValidMonitoredObjectTypes {
		validMonObjTypes[key] = string(val)
	}
	for key, val := range ValidMonitoredObjectDeviceTypes {
		validMonObjDevTypes[key] = string(val)
	}

	result.DefaultValidTypes = &admmod.ValidTypes{
		MonitoredObjectTypes:       validMonObjTypes,
		MonitoredObjectDeviceTypes: validMonObjDevTypes}

	return result
}

func trackAPIMetrics(startTime time.Time, code string, objType string) {
	mon.TrackAPITimeMetricInSeconds(startTime, code, objType)
}

// CreateAdminUser - Create an Administrative User.
func (gsh *GRPCServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	startTime := time.Now()

	res, err := gsh.ash.CreateAdminUser(ctx, user)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateAdminUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateAdminUserStr)
	return res, nil
}

// UpdateAdminUser - Update an Administrative User.
func (gsh *GRPCServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	startTime := time.Now()

	res, err := gsh.ash.UpdateAdminUser(ctx, user)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateAdminUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateAdminUserStr)
	return res, nil
}

// DeleteAdminUser - Delete an Administrative User.
func (gsh *GRPCServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	startTime := time.Now()

	res, err := gsh.ash.DeleteAdminUser(ctx, userID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteAdminUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteAdminUserStr)
	return res, nil
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (gsh *GRPCServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetAdminUser(ctx, userID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetAdminUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetAdminUserStr)
	return res, nil
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (gsh *GRPCServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetAllAdminUsers(ctx, noValue)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetAllAdminUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetAllAdminUserStr)
	return res, nil
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (gsh *GRPCServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	startTime := time.Now()

	// Check if a tenant already exists with this name.
	existingTenantByName, _ := gsh.ash.GetTenantIDByAlias(ctx, &wr.StringValue{Value: strings.ToLower(tenantMeta.GetData().GetName())})
	if len(existingTenantByName.GetValue()) != 0 {
		trackAPIMetrics(startTime, "500", mon.CreateTenantStr)
		msg := fmt.Sprintf("Unable to create Tenant %s. A Tenant with this name already exists", tenantMeta.GetData().GetName())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Create the Tenant metadata record and reserve space to store isolated Tenant data
	result, err := gsh.ash.CreateTenant(ctx, tenantMeta)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateTenantStr)
		msg := fmt.Sprintf("Unable to create Tenant %s", err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Create a default Ingestion Profile for the Tenant.
	idForTenant := result.GetXId()
	ingPrfData := createDefaultTenantIngPrf(idForTenant)

	// Convert to PB object
	convertedIP := pb.TenantIngestionProfile{}
	if err := pb.ConvertToPBObject(ingPrfData, &convertedIP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	_, err = gsh.tsh.CreateTenantIngestionProfile(ctx, &convertedIP)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateTenantStr)
		msg := fmt.Sprintf("Unable to create default Ingestion Profile %s", err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Create a default Threshold Profile for the Tenant
	threshPrfData := createDefaultTenantThresholdPrf(idForTenant)

	// Convert to PB object
	convertedTP := pb.TenantThresholdProfile{}
	if err := pb.ConvertToPBObject(threshPrfData, &convertedTP); err != nil {
		msg := fmt.Sprintf("Unable to convert request to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	threshProfileResponse, err := gsh.tsh.CreateTenantThresholdProfile(ctx, &convertedTP)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateTenantStr)
		msg := fmt.Sprintf("Unable to create default Threshold Profile %s", err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Create the tenant metadata
	// For the IDs used as references inside other objects, need to strip off the 'thresholdProfile_2_'
	// as this is just relational pouch adaption:
	meta := createDefaultTenantMeta(idForTenant, threshProfileResponse.GetXId(), result.GetData().GetName())

	// Convert to PB object
	convertedMD := pb.TenantMetadata{}
	if err := pb.ConvertToPBObject(meta, &convertedMD); err != nil {
		msg := fmt.Sprintf("Unable to convert request to store %s: %s", tenmod.TenantMetaStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	_, err = gsh.tsh.CreateTenantMeta(ctx, &convertedMD)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateTenantStr)
		msg := fmt.Sprintf("Unable to create Tenant metadata %s", err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	trackAPIMetrics(startTime, "200", mon.CreateTenantStr)
	return result, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// startTime := time.Now()

	// res, err := gsh.ash.UpdateTenantDescriptor(ctx, tenantMeta)
	// if err != nil {
	// 	trackAPIMetrics(startTime, "500", mon.UpdateTenantStr)
	// 	return nil, err
	// }

	// trackAPIMetrics(startTime, "200", mon.UpdateTenantStr)
	// return res, nil
	return nil, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (gsh *GRPCServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// startTime := time.Now()

	// res, err := gsh.ash.DeleteTenant(ctx, tenantID)
	// if err != nil {
	// 	trackAPIMetrics(startTime, "500", mon.DeleteTenantStr)
	// 	return nil, err
	// }

	// trackAPIMetrics(startTime, "200", mon.DeleteTenantStr)
	// return res, nil
	return nil, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (gsh *GRPCServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// startTime := time.Now()

	// res, err := gsh.ash.GetTenantDescriptor(ctx, tenantID)
	// if err != nil {
	// 	trackAPIMetrics(startTime, "500", mon.GetTenantStr)
	// 	return nil, err
	// }

	// trackAPIMetrics(startTime, "200", mon.GetTenantStr)
	// return res, nil
	return nil, nil
}

// GetAllTenantDescriptors -  Retrieve all Tenant Descriptors.
func (gsh *GRPCServiceHandler) GetAllTenantDescriptors(ctx context.Context, noValue *emp.Empty) (*pb.TenantDescriptorList, error) {
	// startTime := time.Now()

	// res, err := gsh.ash.GetAllTenantDescriptors(ctx, noValue)
	// if err != nil {
	// 	trackAPIMetrics(startTime, "500", mon.GetTenantStr)
	// 	return nil, err
	// }

	// trackAPIMetrics(startTime, "200", mon.GetTenantStr)
	// return res, nil
	return nil, nil
}

// CreateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) CreateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	startTime := time.Now()

	res, err := gsh.ash.CreateIngestionDictionary(ctx, ingDictionary)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateIngDictStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetTenantStr)
	return res, nil
}

// UpdateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) UpdateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	startTime := time.Now()

	res, err := gsh.ash.UpdateIngestionDictionary(ctx, ingDictionary)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateIngDictStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateIngDictStr)
	return res, nil
}

// DeleteIngestionDictionary - Delete an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) DeleteIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	startTime := time.Now()

	res, err := gsh.ash.DeleteIngestionDictionary(ctx, noValue)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteIngDictStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteIngDictStr)
	return res, nil
}

// GetIngestionDictionary - Retrieve an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) GetIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetIngestionDictionary(ctx, noValue)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetIngDictStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetIngDictStr)
	return res, nil
}

// CreateTenantUser - creates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	startTime := time.Now()

	res, err := gsh.tsh.CreateTenantUser(ctx, tenantUserReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateTenantUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateTenantUserStr)
	return res, nil
}

// UpdateTenantUser - updates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	startTime := time.Now()

	res, err := gsh.tsh.UpdateTenantUser(ctx, tenantUserReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateTenantUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateTenantUserStr)
	return res, nil
}

// DeleteTenantUser - deletes a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	startTime := time.Now()

	res, err := gsh.tsh.DeleteTenantUser(ctx, tenantUserIDReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteTenantUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteTenantUserStr)
	return res, nil
}

// GetTenantUser - retrieves a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetTenantUser(ctx, tenantUserIDReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetTenantUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetTenantUserStr)
	return res, nil
}

// GetAllTenantUsers - retrieves all users scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantUsers(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantUserList, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetAllTenantUsers(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetAllTenantUserStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetAllTenantUserStr)
	return res, nil
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	startTime := time.Now()

	res, err := gsh.tsh.CreateTenantDomain(ctx, tenantDomainRequest)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateTenantDomainStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateTenantDomainStr)
	return res, nil
}

// UpdateTenantDomain - updates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	startTime := time.Now()

	res, err := gsh.tsh.UpdateTenantDomain(ctx, tenantDomainRequest)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateTenantDomainStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateTenantDomainStr)
	return res, nil
}

// DeleteTenantDomain - deletes a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	startTime := time.Now()

	res, err := gsh.tsh.DeleteTenantDomain(ctx, tenantDomainIDRequest)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteTenantDomainStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteTenantDomainStr)
	return res, nil
}

// GetTenantDomain - retrieves a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetTenantDomain(ctx, tenantDomainIDRequest)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetTenantDomainStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetTenantDomainStr)
	return res, nil
}

// GetAllTenantDomains - retrieves all Domains scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantDomains(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDomainList, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetAllTenantDomains(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetAllTenantDomainStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetAllTenantDomainStr)
	return res, nil
}

// CreateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.CreateTenantIngestionProfile(ctx, tenantIngPrfReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateIngPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateIngPrfStr)
	return res, nil
}

// UpdateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.UpdateTenantIngestionProfile(ctx, tenantIngPrfReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateIngPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateIngPrfStr)
	return res, nil
}

// GetTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetTenantIngestionProfile(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetIngPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetIngPrfStr)
	return res, nil
}

// GetActiveTenantIngestionProfile - retrieves the active Ingestion Profile for a single Tenant.
func (gsh *GRPCServiceHandler) GetActiveTenantIngestionProfile(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetActiveTenantIngestionProfile(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetActiveIngPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetActiveIngPrfStr)
	return res, nil
}

// DeleteTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.DeleteTenantIngestionProfile(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteIngPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteIngPrfStr)
	return res, nil
}

// CreateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.CreateTenantThresholdProfile(ctx, tenantThreshPrfReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateThrPrfStr)
	return res, nil
}

// UpdateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.UpdateTenantThresholdProfile(ctx, tenantThreshPrfReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateThrPrfStr)
	return res, nil
}

// GetTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetTenantThresholdProfile(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetThrPrfStr)
	return res, nil
}

// DeleteTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()

	res, err := gsh.tsh.DeleteTenantThresholdProfile(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteThrPrfStr)
	return res, nil
}

// GetAllTenantThresholdProfiles - retieve all Tenant Thresholds.
func (gsh *GRPCServiceHandler) GetAllTenantThresholdProfiles(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantThresholdProfileList, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetAllTenantThresholdProfiles(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetAllThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetAllThrPrfStr)
	return res, nil
}

// CreateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.tsh.CreateMonitoredObject(ctx, monitoredObjectReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateMonObjStr)
	return res, nil
}

// UpdateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.tsh.UpdateMonitoredObject(ctx, monitoredObjectReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateMonObjStr)
	return res, nil
}

// GetMonitoredObject - retrieves the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) GetMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetMonitoredObject(ctx, monitoredObjectIDReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetMonObjStr)
	return res, nil
}

// DeleteMonitoredObject - deletes the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.tsh.DeleteMonitoredObject(ctx, monitoredObjectIDReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteMonObjStr)
	return res, nil
}

// GetAllMonitoredObjects - retrieves all MonitoredObjects scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllMonitoredObjects(ctx context.Context, tenantID *wr.StringValue) (*pb.MonitoredObjectList, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetAllMonitoredObjects(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetAllMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetAllMonObjStr)
	return res, nil
}

// GetMonitoredObjectToDomainMap - retrieves a mapping of MonitoredObjects to each Domain. Will retrieve the mapping either as a count, or as a set of all
// MonitoredObjects that use each Domain.
func (gsh *GRPCServiceHandler) GetMonitoredObjectToDomainMap(ctx context.Context, moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetMonitoredObjectToDomainMap(ctx, moByDomReq)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetMonObjToDomMapStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetMonObjToDomMapStr)
	return res, nil
}

// CreateTenantMeta - Create TenantMeta scoped to a Single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	startTime := time.Now()

	res, err := gsh.tsh.CreateTenantMeta(ctx, meta)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateTenantMetaStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateTenantMetaStr)
	return res, nil
}

// UpdateTenantMeta - Update TenantMeta scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	startTime := time.Now()

	res, err := gsh.tsh.UpdateTenantMeta(ctx, meta)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateTenantMetaStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateTenantMetaStr)
	return res, nil
}

// DeleteTenantMeta - Delete TenantMeta scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {
	startTime := time.Now()

	res, err := gsh.tsh.DeleteTenantMeta(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.DeleteTenantMetaStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteTenantMetaStr)
	return res, nil
}

// GetTenantMeta - Retrieve a User scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {
	startTime := time.Now()

	res, err := gsh.tsh.GetTenantMeta(ctx, tenantID)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetTenantMetaStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetTenantMetaStr)
	return res, nil
}

// GetTenantIDByAlias - retrieve a Tenant ID by the common name of the Tenant
func (gsh *GRPCServiceHandler) GetTenantIDByAlias(ctx context.Context, value *wr.StringValue) (*wr.StringValue, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetTenantIDByAlias(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetTenantIDByAliasStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetTenantIDByAliasStr)
	return res, nil
}

// GetTenantSummaryByAlias - retrieves Tenant summary by a known alias.
func (gsh *GRPCServiceHandler) GetTenantSummaryByAlias(ctx context.Context, value *wr.StringValue) (*pb.TenantSummary, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetTenantIDByAlias(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetTenantSummaryByAliasStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetTenantSummaryByAliasStr)
	return &pb.TenantSummary{Alias: value.Value, Id: res.Value}, nil
}

// AddAdminViews - add views to admin db
func (gsh *GRPCServiceHandler) AddAdminViews() error {
	// startTime := time.Now()

	// err := gsh.ash.AddAdminViews()
	// if err != nil {
	// 	trackAPIMetrics(startTime, "500", mon.AddAdminViewsStr)
	// 	return err
	// }

	// trackAPIMetrics(startTime, "200", mon.AddAdminViewsStr)
	// return nil
	return nil
}

// CreateValidTypes - Create the valid type definition in the system.
func (gsh *GRPCServiceHandler) CreateValidTypes(ctx context.Context, value *pb.ValidTypes) (*pb.ValidTypes, error) {
	startTime := time.Now()

	res, err := gsh.ash.CreateValidTypes(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.CreateValidTypesStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateValidTypesStr)
	return res, nil
}

// UpdateValidTypes - Update the valid type definition in the system.
func (gsh *GRPCServiceHandler) UpdateValidTypes(ctx context.Context, value *pb.ValidTypes) (*pb.ValidTypes, error) {
	startTime := time.Now()

	res, err := gsh.ash.UpdateValidTypes(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.UpdateValidTypesStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateValidTypesStr)
	return res, nil
}

// GetValidTypes - retrieve the enire list of ValidTypes in the system.
func (gsh *GRPCServiceHandler) GetValidTypes(ctx context.Context, value *emp.Empty) (*pb.ValidTypes, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetValidTypes(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetValidTypesStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetValidTypesStr)
	return res, nil
}

// GetSpecificValidTypes - retrieve a subset of the known ValidTypes in the system.
func (gsh *GRPCServiceHandler) GetSpecificValidTypes(ctx context.Context, value *pb.ValidTypesRequest) (*pb.ValidTypesData, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetSpecificValidTypes(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.GetSpecificValidTypesStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetSpecificValidTypesStr)
	return res, nil
}

// DeleteValidTypes - Delete valid types used for the entire deployment.
func (gsh *GRPCServiceHandler) DeleteValidTypes(ctx context.Context, noValue *emp.Empty) (*pb.ValidTypes, error) {
	startTime := time.Now()

	res, err := gsh.ash.DeleteValidTypes(ctx, noValue)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.ValidTypesStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.ValidTypesStr)
	return res, nil
}

// BulkInsertMonitoredObjects - perform a bulk operation on a set of Monitored Objects.
func (gsh *GRPCServiceHandler) BulkInsertMonitoredObjects(ctx context.Context, value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error) {
	startTime := time.Now()

	res, err := gsh.tsh.BulkInsertMonitoredObjects(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.BulkUpdateMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.BulkUpdateMonObjStr)
	return res, nil
}
