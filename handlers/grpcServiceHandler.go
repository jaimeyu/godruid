package handlers

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/accedian/adh-gather/gathergrpc"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

var (
	// ValidMonitoredObjectTypes - known Monitored Object types in the system.
	ValidMonitoredObjectTypes = map[pb.MonitoredObjectData_MonitoredObjectType]string{
		pb.MonitoredObjectData_MO_UNKNOWN: "unknown",
		pb.MonitoredObjectData_TWAMP:      "pe"}

	// ValidMonitoredObjectDeviceTypes - known Monitored Object Device types in the system.
	ValidMonitoredObjectDeviceTypes = map[pb.MonitoredObjectData_DeviceType]string{
		pb.MonitoredObjectData_DT_UNKNOWN:    "unknown",
		pb.MonitoredObjectData_ACCEDIAN_NID:  "accedianNID",
		pb.MonitoredObjectData_ACCEDIAN_VNID: "accedianVNID"}
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
	validMonObjTypes := make([]string, 0)
	validMonObjDevTypes := make([]string, 0)

	for _, val := range ValidMonitoredObjectTypes {
		validMonObjTypes = append(validMonObjTypes, val)
	}
	for _, val := range ValidMonitoredObjectDeviceTypes {
		validMonObjDevTypes = append(validMonObjDevTypes, val)
	}

	result.DefaultValidTypes = &pb.ValidTypesData{
		MonitoredObjectTypes:       validMonObjTypes,
		MonitoredObjectDeviceTypes: validMonObjDevTypes}

	return result
}

// CreateAdminUser - Create an Administrative User.
func (gsh *GRPCServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	return gsh.ash.CreateAdminUser(ctx, user)
}

// UpdateAdminUser - Update an Administrative User.
func (gsh *GRPCServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	return gsh.ash.UpdateAdminUser(ctx, user)
}

// DeleteAdminUser - Delete an Administrative User.
func (gsh *GRPCServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	return gsh.ash.DeleteAdminUser(ctx, userID)
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (gsh *GRPCServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	return gsh.ash.GetAdminUser(ctx, userID)
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (gsh *GRPCServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	return gsh.ash.GetAllAdminUsers(ctx, noValue)
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (gsh *GRPCServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Check if a tenant already exists with this name.
	existingTenantByName, _ := gsh.ash.GetTenantIDByAlias(ctx, &wr.StringValue{Value: strings.ToLower(tenantMeta.GetData().GetName())})
	if len(existingTenantByName.GetValue()) != 0 {
		return nil, fmt.Errorf("Unable to create Tenant %s. A Tenant with this name already exists", tenantMeta.GetData().GetName())
	}

	// Create the Tenant metadata record and reserve space to store isolated Tenant data
	result, err := gsh.ash.CreateTenant(ctx, tenantMeta)
	if err != nil {
		return nil, err
	}

	// Create a default Ingestion Profile for the Tenant.
	idForTenant := result.GetXId()
	ingPrfData := createDefaultTenantIngPrf(idForTenant)
	ingPrfReq := pb.TenantIngestionProfile{Data: ingPrfData}
	_, err = gsh.tsh.CreateTenantIngestionProfile(ctx, &ingPrfReq)
	if err != nil {
		return nil, err
	}

	// Create a default Threshold Profile for the Tenant
	threshPrfData := createDefaultTenantThresholdPrf(idForTenant)
	threshPrfReq := pb.TenantThresholdProfile{Data: threshPrfData}
	threshProfileResponse, err := gsh.tsh.CreateTenantThresholdProfile(ctx, &threshPrfReq)
	if err != nil {
		return nil, err
	}

	// Create the tenant metadata
	// For the IDs used as references inside other objects, need to strip off the 'thresholdProfile_2_'
	// as this is just relational pouch adaption:
	meta := createDefaultTenantMeta(idForTenant, threshProfileResponse.GetXId(), result.GetData().GetName())
	metaReq := pb.TenantMetadata{Data: meta}
	_, err = gsh.tsh.CreateTenantMeta(ctx, &metaReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	return gsh.ash.UpdateTenantDescriptor(ctx, tenantMeta)
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (gsh *GRPCServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// TODO: Add calls here to Tenant Service to delete any related
	// tenant data.

	return gsh.ash.DeleteTenant(ctx, tenantID)
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (gsh *GRPCServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	return gsh.ash.GetTenantDescriptor(ctx, tenantID)
}

// GetAllTenantDescriptors -  Retrieve all Tenant Descriptors.
func (gsh *GRPCServiceHandler) GetAllTenantDescriptors(ctx context.Context, noValue *emp.Empty) (*pb.TenantDescriptorList, error) {
	return gsh.ash.GetAllTenantDescriptors(ctx, noValue)
}

// CreateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) CreateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	return gsh.ash.CreateIngestionDictionary(ctx, ingDictionary)
}

// UpdateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) UpdateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	return gsh.ash.UpdateIngestionDictionary(ctx, ingDictionary)
}

// DeleteIngestionDictionary - Delete an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) DeleteIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	return gsh.ash.DeleteIngestionDictionary(ctx, noValue)
}

// GetIngestionDictionary - Retrieve an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) GetIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	return gsh.ash.GetIngestionDictionary(ctx, noValue)
}

// CreateTenantUser - creates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	return gsh.tsh.CreateTenantUser(ctx, tenantUserReq)
}

// UpdateTenantUser - updates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	return gsh.tsh.UpdateTenantUser(ctx, tenantUserReq)
}

// DeleteTenantUser - deletes a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantUser(ctx context.Context, tenantUserIdReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	return gsh.tsh.DeleteTenantUser(ctx, tenantUserIdReq)
}

// GetTenantUser - retrieves a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantUser(ctx context.Context, tenantUserIdReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	return gsh.tsh.GetTenantUser(ctx, tenantUserIdReq)
}

// GetAllTenantUsers - retrieves all users scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantUsers(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantUserList, error) {
	return gsh.tsh.GetAllTenantUsers(ctx, tenantID)
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	return gsh.tsh.CreateTenantDomain(ctx, tenantDomainRequest)
}

// UpdateTenantDomain - updates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	return gsh.tsh.UpdateTenantDomain(ctx, tenantDomainRequest)
}

// DeleteTenantDomain - deletes a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	return gsh.tsh.DeleteTenantDomain(ctx, tenantDomainIDRequest)
}

// GetTenantDomain - retrieves a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	return gsh.tsh.GetTenantDomain(ctx, tenantDomainIDRequest)
}

// GetAllTenantDomains - retrieves all Domains scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantDomains(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDomainList, error) {
	return gsh.tsh.GetAllTenantDomains(ctx, tenantID)
}

// CreateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	return gsh.tsh.CreateTenantIngestionProfile(ctx, tenantIngPrfReq)
}

// UpdateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	return gsh.tsh.UpdateTenantIngestionProfile(ctx, tenantIngPrfReq)
}

// GetTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	return gsh.tsh.GetTenantIngestionProfile(ctx, tenantID)
}

// GetActiveTenantIngestionProfile - retrieves the active Ingestion Profile for a single Tenant.
func (gsh *GRPCServiceHandler) GetActiveTenantIngestionProfile(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantIngestionProfile, error) {
	return gsh.tsh.GetActiveTenantIngestionProfile(ctx, tenantID)
}

// DeleteTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	return gsh.tsh.DeleteTenantIngestionProfile(ctx, tenantID)
}

// CreateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	return gsh.tsh.CreateTenantThresholdProfile(ctx, tenantThreshPrfReq)
}

// UpdateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	return gsh.tsh.UpdateTenantThresholdProfile(ctx, tenantThreshPrfReq)
}

// GetTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	return gsh.tsh.GetTenantThresholdProfile(ctx, tenantID)
}

// DeleteTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	return gsh.tsh.DeleteTenantThresholdProfile(ctx, tenantID)
}

// GetAllTenantThresholdProfiles - retieve all Tenant Thresholds.
func (gsh *GRPCServiceHandler) GetAllTenantThresholdProfiles(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantThresholdProfileList, error) {
	return gsh.tsh.GetAllTenantThresholdProfiles(ctx, tenantID)
}

// CreateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	return gsh.tsh.CreateMonitoredObject(ctx, monitoredObjectReq)
}

// UpdateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	return gsh.tsh.UpdateMonitoredObject(ctx, monitoredObjectReq)
}

// GetMonitoredObject - retrieves the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) GetMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	return gsh.tsh.GetMonitoredObject(ctx, monitoredObjectIDReq)
}

// DeleteMonitoredObject - deletes the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	return gsh.tsh.DeleteMonitoredObject(ctx, monitoredObjectIDReq)
}

// GetAllMonitoredObjects - retrieves all MonitoredObjects scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllMonitoredObjects(ctx context.Context, tenantID *wr.StringValue) (*pb.MonitoredObjectList, error) {
	return gsh.tsh.GetAllMonitoredObjects(ctx, tenantID)
}

// GetMonitoredObjectToDomainMap - retrieves a mapping of MonitoredObjects to each Domain. Will retrieve the mapping either as a count, or as a set of all
// MonitoredObjects that use each Domain.
func (gsh *GRPCServiceHandler) GetMonitoredObjectToDomainMap(ctx context.Context, moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	return gsh.tsh.GetMonitoredObjectToDomainMap(ctx, moByDomReq)
}

// GetThresholdCrossing - Retrieves the Threshold crossings for a given threshold profile,
// interval, tenant, domain
func (gsh *GRPCServiceHandler) GetThresholdCrossing(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest) (*pb.JSONAPIObject, error) {
	tenantID := thresholdCrossingReq.Tenant

	thresholdProfile, err := gsh.GetTenantThresholdProfile(ctx, &pb.TenantThresholdProfileIdRequest{
		TenantId:           tenantID,
		ThresholdProfileId: thresholdCrossingReq.ThresholdProfileId,
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err)
	}

	return gsh.msh.GetThresholdCrossing(ctx, thresholdCrossingReq, thresholdProfile)
}

// GetThresholdCrossingByMonitoredObject - Retrieves the Threshold crossings for a given threshold profile,
// interval, tenant, domain, and groups by monitoredObjectID
func (gsh *GRPCServiceHandler) GetThresholdCrossingByMonitoredObject(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest) (*pb.JSONAPIObject, error) {

	tenantID := thresholdCrossingReq.Tenant

	thresholdProfile, err := gsh.GetTenantThresholdProfile(ctx, &pb.TenantThresholdProfileIdRequest{
		TenantId:           tenantID,
		ThresholdProfileId: thresholdCrossingReq.ThresholdProfileId,
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to find threshold profile for given query parameters: %s. Error: %s", thresholdCrossingReq, err)
	}

	return gsh.msh.GetThresholdCrossingByMonitoredObject(ctx, thresholdCrossingReq, thresholdProfile)
}

// GetHistogram -
func (gsh *GRPCServiceHandler) GetHistogram(ctx context.Context, histogramReq *pb.HistogramRequest) (*pb.JSONAPIObject, error) {

	return gsh.msh.GetHistogram(ctx, histogramReq)
}

// CreateTenantMeta - Create TenantMeta scoped to a Single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	return gsh.tsh.CreateTenantMeta(ctx, meta)
}

// UpdateTenantMeta - Update TenantMeta scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	return gsh.tsh.UpdateTenantMeta(ctx, meta)
}

// DeleteTenantMeta - Delete TenantMeta scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {
	return gsh.tsh.DeleteTenantMeta(ctx, tenantID)
}

// GetTenantMeta - Retrieve a User scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {
	return gsh.tsh.GetTenantMeta(ctx, tenantID)
}

// GetTenantIDByAlias - retrieve a Tenant ID by the common name of the Tenant
func (gsh *GRPCServiceHandler) GetTenantIDByAlias(ctx context.Context, value *wr.StringValue) (*wr.StringValue, error) {
	return gsh.ash.GetTenantIDByAlias(ctx, value)
}

// AddAdminViews - add views to admin db
func (gsh *GRPCServiceHandler) AddAdminViews() error {
	return gsh.ash.AddAdminViews()
}

// CreateValidTypes - Create the valid type definition in the system.
func (gsh *GRPCServiceHandler) CreateValidTypes(ctx context.Context, value *pb.ValidTypes) (*pb.ValidTypes, error) {
	return gsh.ash.CreateValidTypes(ctx, value)
}

// UpdateValidTypes - Update the valid type definition in the system.
func (gsh *GRPCServiceHandler) UpdateValidTypes(ctx context.Context, value *pb.ValidTypes) (*pb.ValidTypes, error) {
	return gsh.ash.UpdateValidTypes(ctx, value)
}

// GetValidTypes - retrieve the enire list of ValidTypes in the system.
func (gsh *GRPCServiceHandler) GetValidTypes(ctx context.Context, value *emp.Empty) (*pb.ValidTypes, error) {
	return gsh.ash.GetValidTypes(ctx, value)
}

// GetSpecificValidTypes - retrieve a subset of the known ValidTypes in the system.
func (gsh *GRPCServiceHandler) GetSpecificValidTypes(ctx context.Context, value *pb.ValidTypesRequest) (*pb.ValidTypesData, error) {
	return gsh.ash.GetSpecificValidTypes(ctx, value)
}
