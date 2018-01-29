package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/accedian/adh-gather/gathergrpc"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
	mon "github.com/accedian/adh-gather/monitoring"
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

var (
	// ValidMonitoredObjectTypes - known Monitored Object types in the system.
	ValidMonitoredObjectTypes = map[string]MonitoredObjectType{
		"pe":            TwampPE,
		"sf":            TwampSF,
		"sl":            TwampSL,
		string(TwampPE): TwampPE,
		string(TwampSF): TwampSF,
		string(TwampSL): TwampSL}

	// ValidMonitoredObjectDeviceTypes - known Monitored Object Device types in the system.
	ValidMonitoredObjectDeviceTypes = map[string]MonitoredObjectDeviceType{
		string(AccedianNID):  AccedianNID,
		string(AccedianVNID): AccedianVNID}
)

// GRPCServiceHandler - implementer of all gRPC Services. Offloads
// implementation details to each unique service handler. When new
// gRPC services are added, a new Service Handler should be created,
// and a pointer to that object should be added to this wrapper.
type GRPCServiceHandler struct {
	ash               *AdminServiceHandler
	tsh               *TenantServiceHandler
	msh               *MetricServiceHandler
	DefaultValidTypes *pb.ValidTypesData
}

// CreateCoordinator - used to create a gRPC service handler wrapper
// that coordinates the logic to satisfy all gRPC service
// interfaces.
func CreateCoordinator() *GRPCServiceHandler {
	result := new(GRPCServiceHandler)

	result.ash = CreateAdminServiceHandler()
	result.tsh = CreateTenantServiceHandler()
	result.msh = CreateMetricServiceHandler()

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

	result.DefaultValidTypes = &pb.ValidTypesData{
		MonitoredObjectTypes:       validMonObjTypes,
		MonitoredObjectDeviceTypes: validMonObjDevTypes}

	return result
}

// CreateAdminUser - Create an Administrative User.
func (gsh *GRPCServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	startTime := time.Now()

	res, err := gsh.ash.CreateAdminUser(ctx, user)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.CreateAdminUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateAdminUserStr)
	return res, nil
}

// UpdateAdminUser - Update an Administrative User.
func (gsh *GRPCServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	startTime := time.Now()
	res, err := gsh.ash.UpdateAdminUser(ctx, user)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.UpdateAdminUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateAdminUserStr)
	return res, nil
}

// DeleteAdminUser - Delete an Administrative User.
func (gsh *GRPCServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	startTime := time.Now()
	res, err := gsh.ash.DeleteAdminUser(ctx, userID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.DeleteAdminUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteAdminUserStr)
	return res, nil
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (gsh *GRPCServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	startTime := time.Now()
	res, err := gsh.ash.GetAdminUser(ctx, userID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.GetAdminUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetAdminUserStr)
	return res, nil
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (gsh *GRPCServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	startTime := time.Now()
	res, err := gsh.ash.GetAllAdminUsers(ctx, noValue)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.GetAllAdminUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetAllAdminUserStr)
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
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.CreateTenantStr)
		return nil, fmt.Errorf("Unable to create Tenant %s. A Tenant with this name already exists", tenantMeta.GetData().GetName())
	}

	// Create the Tenant metadata record and reserve space to store isolated Tenant data
	result, err := gsh.ash.CreateTenant(ctx, tenantMeta)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.CreateTenantStr)
		return nil, err
	}

	// Create a default Ingestion Profile for the Tenant.
	idForTenant := result.GetXId()
	ingPrfData := createDefaultTenantIngPrf(idForTenant)
	ingPrfReq := pb.TenantIngestionProfile{Data: ingPrfData}
	_, err = gsh.tsh.CreateTenantIngestionProfile(ctx, &ingPrfReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.CreateTenantStr)
		return nil, err
	}

	// Create a default Threshold Profile for the Tenant
	threshPrfData := createDefaultTenantThresholdPrf(idForTenant)
	threshPrfReq := pb.TenantThresholdProfile{Data: threshPrfData}
	threshProfileResponse, err := gsh.tsh.CreateTenantThresholdProfile(ctx, &threshPrfReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.CreateTenantStr)
		return nil, err
	}

	// Create the tenant metadata
	// For the IDs used as references inside other objects, need to strip off the 'thresholdProfile_2_'
	// as this is just relational pouch adaption:
	meta := createDefaultTenantMeta(idForTenant, threshProfileResponse.GetXId(), result.GetData().GetName())
	metaReq := pb.TenantMetadata{Data: meta}
	_, err = gsh.tsh.CreateTenantMeta(ctx, &metaReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.CreateTenantStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateTenantStr)
	return result, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	startTime := time.Now()
	
	res, err := gsh.ash.UpdateTenantDescriptor(ctx, tenantMeta)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.UpdateTenantStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateTenantStr)
	return res, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (gsh *GRPCServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	startTime := time.Now()

	res, err := gsh.ash.DeleteTenant(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime, "500", mon.DeleteTenantStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteTenantStr)
	return res, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (gsh *GRPCServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	startTime := time.Now()

	res, err := gsh.ash.GetTenantDescriptor(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetTenantStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetTenantStr)
	return res, nil
}

// GetAllTenantDescriptors -  Retrieve all Tenant Descriptors.
func (gsh *GRPCServiceHandler) GetAllTenantDescriptors(ctx context.Context, noValue *emp.Empty) (*pb.TenantDescriptorList, error) {
	startTime := time.Now()
	res, err := gsh.ash.GetAllTenantDescriptors(ctx, noValue)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetTenantStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetTenantStr)
	return res, nil
}

// CreateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) CreateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	startTime := time.Now()
	res, err := gsh.ash.CreateIngestionDictionary(ctx, ingDictionary)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateIngDictStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetTenantStr)
	return res, nil
}

// UpdateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) UpdateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	startTime := time.Now()
	res, err := gsh.ash.UpdateIngestionDictionary(ctx, ingDictionary)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateIngDictStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateIngDictStr)
	return res, nil
}

// DeleteIngestionDictionary - Delete an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) DeleteIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	startTime := time.Now()
	res, err := gsh.ash.DeleteIngestionDictionary(ctx, noValue)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.DeleteIngDictStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteIngDictStr)
	return res, nil
}

// GetIngestionDictionary - Retrieve an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) GetIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	startTime := time.Now()
	res, err := gsh.ash.GetIngestionDictionary(ctx, noValue)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetIngDictStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetIngDictStr)
	return res, nil
}

// CreateTenantUser - creates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	startTime := time.Now()
	res, err := gsh.tsh.CreateTenantUser(ctx, tenantUserReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateTenantUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateTenantUserStr)
	return res, nil
}

// UpdateTenantUser - updates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	startTime := time.Now()
	res, err := gsh.tsh.UpdateTenantUser(ctx, tenantUserReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateTenantUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateTenantUserStr)
	return res, nil
}

// DeleteTenantUser - deletes a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	startTime := time.Now()
	res, err := gsh.tsh.DeleteTenantUser(ctx, tenantUserIDReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.DeleteTenantUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteTenantUserStr)
	return res, nil
}

// GetTenantUser - retrieves a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetTenantUser(ctx, tenantUserIDReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetTenantUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetTenantUserStr)
	return res, nil
}

// GetAllTenantUsers - retrieves all users scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantUsers(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantUserList, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetAllTenantUsers(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetAllTenantUserStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetAllTenantUserStr)
	return res, nil
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	startTime := time.Now()
	res, err := gsh.tsh.CreateTenantDomain(ctx, tenantDomainRequest)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateTenantDomainStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateTenantDomainStr)
	return res, nil
}

// UpdateTenantDomain - updates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	startTime := time.Now()
	res, err := gsh.tsh.UpdateTenantDomain(ctx, tenantDomainRequest)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateTenantDomainStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateTenantDomainStr)
	return res, nil
}

// DeleteTenantDomain - deletes a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	startTime := time.Now()
	res, err := gsh.tsh.DeleteTenantDomain(ctx, tenantDomainIDRequest)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.DeleteTenantDomainStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteTenantDomainStr)
	return res, nil
}

// GetTenantDomain - retrieves a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetTenantDomain(ctx, tenantDomainIDRequest)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetTenantDomainStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetTenantDomainStr)
	return res, nil
}

// GetAllTenantDomains - retrieves all Domains scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantDomains(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDomainList, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetAllTenantDomains(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetAllTenantDomainStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetAllTenantDomainStr)
	return res, nil
}

// CreateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.CreateTenantIngestionProfile(ctx, tenantIngPrfReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateIngPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateIngPrfStr)
	return res, nil
}

// UpdateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.UpdateTenantIngestionProfile(ctx, tenantIngPrfReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateIngPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateIngPrfStr)
	return res, nil
}

// GetTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetTenantIngestionProfile(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetIngPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetIngPrfStr)
	return res, nil
}

// GetActiveTenantIngestionProfile - retrieves the active Ingestion Profile for a single Tenant.
func (gsh *GRPCServiceHandler) GetActiveTenantIngestionProfile(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetActiveTenantIngestionProfile(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetActiveIngPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetActiveIngPrfStr)
	return res, nil
}

// DeleteTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.DeleteTenantIngestionProfile(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.DeleteIngPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteIngPrfStr)
	return res, nil
}

// CreateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.CreateTenantThresholdProfile(ctx, tenantThreshPrfReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateThrPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateThrPrfStr)
	return res, nil
}

// UpdateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.UpdateTenantThresholdProfile(ctx, tenantThreshPrfReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateThrPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateThrPrfStr)
	return res, nil
}

// GetTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetTenantThresholdProfile(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetThrPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetThrPrfStr)
	return res, nil
}

// DeleteTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()
	res, err := gsh.tsh.DeleteTenantThresholdProfile(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.DeleteThrPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteThrPrfStr)
	return res, nil
}

// GetAllTenantThresholdProfiles - retieve all Tenant Thresholds.
func (gsh *GRPCServiceHandler) GetAllTenantThresholdProfiles(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantThresholdProfileList, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetAllTenantThresholdProfiles(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetAllThrPrfStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetAllThrPrfStr)
	return res, nil
}

// CreateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	startTime := time.Now()
	res, err := gsh.tsh.CreateMonitoredObject(ctx, monitoredObjectReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateMonObjStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateMonObjStr)
	return res, nil
}

// UpdateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	startTime := time.Now()
	res, err := gsh.tsh.UpdateMonitoredObject(ctx, monitoredObjectReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateMonObjStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateMonObjStr)
	return res, nil
}

// GetMonitoredObject - retrieves the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) GetMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetMonitoredObject(ctx, monitoredObjectIDReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetMonObjStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetMonObjStr)
	return res, nil
}

// DeleteMonitoredObject - deletes the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	startTime := time.Now()
	res, err := gsh.tsh.DeleteMonitoredObject(ctx, monitoredObjectIDReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.DeleteMonObjStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteMonObjStr)
	return res, nil
}

// GetAllMonitoredObjects - retrieves all MonitoredObjects scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllMonitoredObjects(ctx context.Context, tenantID *wr.StringValue) (*pb.MonitoredObjectList, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetAllMonitoredObjects(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetAllMonObjStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetAllMonObjStr)
	return res, nil
}

// GetMonitoredObjectToDomainMap - retrieves a mapping of MonitoredObjects to each Domain. Will retrieve the mapping either as a count, or as a set of all
// MonitoredObjects that use each Domain.
func (gsh *GRPCServiceHandler) GetMonitoredObjectToDomainMap(ctx context.Context, moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetMonitoredObjectToDomainMap(ctx, moByDomReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetMonObjToDomMapStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetMonObjToDomMapStr)
	return res, nil
}

// GetThresholdCrossing - Retrieves the Threshold crossings for a given threshold profile,
// interval, tenant, domain
func (gsh *GRPCServiceHandler) GetThresholdCrossing(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest) (*pb.JSONAPIObject, error) {
	startTime := time.Now()
	tenantID := thresholdCrossingReq.Tenant

	thresholdProfile, err := gsh.GetTenantThresholdProfile(ctx, &pb.TenantThresholdProfileIdRequest{
		TenantId:           tenantID,
		ThresholdProfileId: thresholdCrossingReq.ThresholdProfileId,
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err)
	}

	res, err := gsh.msh.GetThresholdCrossing(ctx, thresholdCrossingReq, thresholdProfile)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetThrCrossStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetThrCrossStr)
	return res, nil
}

// GetThresholdCrossingByMonitoredObject - Retrieves the Threshold crossings for a given threshold profile,
// interval, tenant, domain, and groups by monitoredObjectID
func (gsh *GRPCServiceHandler) GetThresholdCrossingByMonitoredObject(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest) (*pb.JSONAPIObject, error) {
	startTime := time.Now()

	tenantID := thresholdCrossingReq.Tenant

	thresholdProfile, err := gsh.GetTenantThresholdProfile(ctx, &pb.TenantThresholdProfileIdRequest{
		TenantId:           tenantID,
		ThresholdProfileId: thresholdCrossingReq.ThresholdProfileId,
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err)
	}

	res, err := gsh.msh.GetThresholdCrossingByMonitoredObject(ctx, thresholdCrossingReq, thresholdProfile)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetThrCrossByMonObjStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetThrCrossByMonObjStr)
	return res, nil
}

// GetHistogram - Retrieve bucket data from druid
func (gsh *GRPCServiceHandler) GetHistogram(ctx context.Context, histogramReq *pb.HistogramRequest) (*pb.JSONAPIObject, error) {
	startTime := time.Now()
	res, err := gsh.msh.GetHistogram(ctx, histogramReq)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetHistogramObjStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetHistogramObjStr)
	return res, nil
}

// GetRawMetrics - Retrieve raw metric data from druid
func (gsh *GRPCServiceHandler) GetRawMetrics(ctx context.Context, rawMetricReq *pb.RawMetricsRequest) (*pb.JSONAPIObject, error) {

	return gsh.msh.GetRawMetrics(ctx, rawMetricReq)
}

// CreateTenantMeta - Create TenantMeta scoped to a Single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	startTime := time.Now()
	res, err := gsh.tsh.CreateTenantMeta(ctx, meta)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateTenantMetaStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateTenantMetaStr)
	return res, nil
}

// UpdateTenantMeta - Update TenantMeta scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	startTime := time.Now()
	res, err := gsh.tsh.UpdateTenantMeta(ctx, meta)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateTenantMetaStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateTenantMetaStr)
	return res, nil
}

// DeleteTenantMeta - Delete TenantMeta scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {
	startTime := time.Now()
	res, err := gsh.tsh.DeleteTenantMeta(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.DeleteTenantMetaStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.DeleteTenantMetaStr)
	return res, nil
}

// GetTenantMeta - Retrieve a User scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {
	startTime := time.Now()
	res, err := gsh.tsh.GetTenantMeta(ctx, tenantID)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetTenantMetaStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetTenantMetaStr)
	return res, nil
}

// GetTenantIDByAlias - retrieve a Tenant ID by the common name of the Tenant
func (gsh *GRPCServiceHandler) GetTenantIDByAlias(ctx context.Context, value *wr.StringValue) (*wr.StringValue, error) {
	startTime := time.Now()
	res, err := gsh.ash.GetTenantIDByAlias(ctx, value)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetTenantIDByAliasStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetTenantIDByAliasStr)
	return res, nil
}

// AddAdminViews - add views to admin db
func (gsh *GRPCServiceHandler) AddAdminViews() error {
	startTime := time.Now()
	err := gsh.ash.AddAdminViews()
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.AddAdminViewsStr)
		return err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.AddAdminViewsStr)
	return nil
}

// CreateValidTypes - Create the valid type definition in the system.
func (gsh *GRPCServiceHandler) CreateValidTypes(ctx context.Context, value *pb.ValidTypes) (*pb.ValidTypes, error) {
	startTime := time.Now()
	res, err := gsh.ash.CreateValidTypes(ctx, value)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.CreateValidTypesStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.CreateValidTypesStr)
	return res, nil
}

// UpdateValidTypes - Update the valid type definition in the system.
func (gsh *GRPCServiceHandler) UpdateValidTypes(ctx context.Context, value *pb.ValidTypes) (*pb.ValidTypes, error) {
	startTime := time.Now()
	res, err := gsh.ash.UpdateValidTypes(ctx, value)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.UpdateValidTypesStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.UpdateValidTypesStr)
	return res, nil
}

// GetValidTypes - retrieve the enire list of ValidTypes in the system.
func (gsh *GRPCServiceHandler) GetValidTypes(ctx context.Context, value *emp.Empty) (*pb.ValidTypes, error) {
	startTime := time.Now()
	res, err := gsh.ash.GetValidTypes(ctx, value)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetValidTypesStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetValidTypesStr)
	return res, nil
}

// GetSpecificValidTypes - retrieve a subset of the known ValidTypes in the system.
func (gsh *GRPCServiceHandler) GetSpecificValidTypes(ctx context.Context, value *pb.ValidTypesRequest) (*pb.ValidTypesData, error) {
	startTime := time.Now()
	res, err := gsh.ash.GetSpecificValidTypes(ctx, value)
	if err != nil {
		mon.TrackAPITimeMetricInSeconds(startTime,"500", mon.GetSpecificValidTypesStr)
		return nil, err
	}

	mon.TrackAPITimeMetricInSeconds(startTime, "200", mon.GetSpecificValidTypesStr)
	return res, nil
}
