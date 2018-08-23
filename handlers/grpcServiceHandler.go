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

// GRPCServiceHandler - implementer of all gRPC Services. Offloads
// implementation details to each unique service handler. When new
// gRPC services are added, a new Service Handler should be created,
// and a pointer to that object should be added to this wrapper.
type GRPCServiceHandler struct {
	ash *AdminServiceHandler
	Tsh *TenantServiceHandler
}

// CreateCoordinator - used to create a gRPC service handler wrapper
// that coordinates the logic to satisfy all gRPC service
// interfaces.
func CreateCoordinator() *GRPCServiceHandler {
	result := new(GRPCServiceHandler)

	result.ash = CreateAdminServiceHandler()
	result.Tsh = CreateTenantServiceHandler()

	return result
}

func trackAPIMetrics(startTime time.Time, code string, objType string) {
	mon.TrackAPITimeMetricInSeconds(startTime, code, objType)
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

	_, err = gsh.Tsh.CreateTenantIngestionProfile(ctx, &convertedIP)
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

	threshProfileResponse, err := gsh.Tsh.CreateTenantThresholdProfile(ctx, &convertedTP)
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
	_, err = gsh.Tsh.CreateTenantMeta(ctx, &convertedMD)
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

// GetIngestionDictionary - Retrieve an IngestionDictionary used for the entire deployment.
func (gsh *GRPCServiceHandler) GetIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	startTime := time.Now()

	ingDict := admmod.GetIngestionDictionaryFromFile()

	// Convert to PB object
	converted := pb.IngestionDictionary{}
	if err := pb.ConvertToPBObject(ingDict, &converted); err != nil {
		msg := fmt.Sprintf("Unable to convert request to store %s: %s", "Ingestion Dictionary", err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	trackAPIMetrics(startTime, "200", mon.GetIngDictStr)
	return &converted, nil
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.CreateTenantDomain(ctx, tenantDomainRequest)
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

	res, err := gsh.Tsh.UpdateTenantDomain(ctx, tenantDomainRequest)
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

	res, err := gsh.Tsh.DeleteTenantDomain(ctx, tenantDomainIDRequest)
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

	res, err := gsh.Tsh.GetTenantDomain(ctx, tenantDomainIDRequest)
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

	res, err := gsh.Tsh.GetAllTenantDomains(ctx, tenantID)
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

	res, err := gsh.Tsh.CreateTenantIngestionProfile(ctx, tenantIngPrfReq)
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

	res, err := gsh.Tsh.UpdateTenantIngestionProfile(ctx, tenantIngPrfReq)
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

	res, err := gsh.Tsh.GetTenantIngestionProfile(ctx, tenantID)
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

	res, err := gsh.Tsh.GetActiveTenantIngestionProfile(ctx, tenantID)
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

	res, err := gsh.Tsh.DeleteTenantIngestionProfile(ctx, tenantID)
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

	res, err := gsh.Tsh.CreateTenantThresholdProfile(ctx, tenantThreshPrfReq)
	if err != nil {
		logger.Log.Errorf("Could not create Tenant ThresholdProfile for Tenant %s: %s", tenantThreshPrfReq.Data.GetTenantId(), err.Error())
		trackAPIMetrics(startTime, "500", mon.CreateThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateThrPrfStr)
	return res, nil
}

// UpdateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.UpdateTenantThresholdProfile(ctx, tenantThreshPrfReq)
	if err != nil {
		logger.Log.Errorf("Could not update Tenant ThresholdProfile for Tenant %s: %s", tenantThreshPrfReq.Data.GetTenantId(), err.Error())
		trackAPIMetrics(startTime, "500", mon.UpdateThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateThrPrfStr)
	return res, nil
}

// GetTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.GetTenantThresholdProfile(ctx, tenantID)
	if err != nil {
		logger.Log.Errorf("Could not retrieve Tenant ThresholdProfile %s for Tenant %s: %s", tenantID.GetTenantId(), tenantID.GetThresholdProfileId(), err.Error())
		trackAPIMetrics(startTime, "500", mon.GetThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetThrPrfStr)
	return res, nil
}

// DeleteTenantThresholdProfile - retrieves the Threshold Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantThresholdProfile(ctx context.Context, tenantID *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.DeleteTenantThresholdProfile(ctx, tenantID)
	if err != nil {
		logger.Log.Errorf("Could not delete Tenant ThresholdProfile %s for Tenant %s: %s", tenantID.GetTenantId(), tenantID.GetThresholdProfileId(), err.Error())
		trackAPIMetrics(startTime, "500", mon.DeleteThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteThrPrfStr)
	return res, nil
}

// GetAllTenantThresholdProfiles - retieve all Tenant Thresholds.
func (gsh *GRPCServiceHandler) GetAllTenantThresholdProfiles(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantThresholdProfileList, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.GetAllTenantThresholdProfiles(ctx, tenantID)
	if err != nil {
		logger.Log.Errorf("Could not retrieve all Tenant ThresholdProfiles for Tenant %s: %s", tenantID, err.Error())
		trackAPIMetrics(startTime, "500", mon.GetAllThrPrfStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetAllThrPrfStr)
	return res, nil
}

// CreateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.CreateMonitoredObject(ctx, monitoredObjectReq)
	if err != nil {
		logger.Log.Errorf("Could not create Monitored Object for Tenant %s: %s", monitoredObjectReq.Data.GetTenantId(), err.Error())
		trackAPIMetrics(startTime, "500", mon.CreateMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.CreateMonObjStr)
	return res, nil
}

// UpdateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.UpdateMonitoredObject(ctx, monitoredObjectReq)
	if err != nil {
		logger.Log.Errorf("Could not update Monitored Object for Tenant %s: %s", monitoredObjectReq.Data.GetTenantId(), err.Error())
		trackAPIMetrics(startTime, "500", mon.UpdateMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.UpdateMonObjStr)
	return res, nil
}

// GetMonitoredObject - retrieves the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) GetMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.GetMonitoredObject(ctx, monitoredObjectIDReq)
	if err != nil {
		logger.Log.Errorf("Could not Get Monitored Object %s for Tenant %s: %s", monitoredObjectIDReq.MonitoredObjectId, monitoredObjectIDReq.TenantId, err.Error())
		trackAPIMetrics(startTime, "500", mon.GetMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.GetMonObjStr)
	return res, nil
}

// DeleteMonitoredObject - deletes the MonitoredObject for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.DeleteMonitoredObject(ctx, monitoredObjectIDReq)
	if err != nil {
		logger.Log.Errorf("Could not delete Monitored Object %s for Tenant %s: %s", monitoredObjectIDReq.MonitoredObjectId, monitoredObjectIDReq.TenantId, err.Error())
		trackAPIMetrics(startTime, "500", mon.DeleteMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.DeleteMonObjStr)
	return res, nil
}

// GetAllMonitoredObjects - retrieves all MonitoredObjects scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllMonitoredObjects(ctx context.Context, tenantID *wr.StringValue) (*pb.MonitoredObjectList, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.GetAllMonitoredObjects(ctx, tenantID)
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

	res, err := gsh.Tsh.GetMonitoredObjectToDomainMap(ctx, moByDomReq)
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

	res, err := gsh.Tsh.CreateTenantMeta(ctx, meta)
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

	res, err := gsh.Tsh.UpdateTenantMeta(ctx, meta)
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

	res, err := gsh.Tsh.DeleteTenantMeta(ctx, tenantID)
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

	res, err := gsh.Tsh.GetTenantMeta(ctx, tenantID)
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

// GetValidTypes - retrieve the enire list of ValidTypes in the system.
func (gsh *GRPCServiceHandler) GetValidTypes(ctx context.Context, value *emp.Empty) (*pb.ValidTypes, error) {
	startTime := time.Now()

	validTypes := admmod.GetValidTypes()

	// Convert to PB object
	converted := pb.ValidTypes{}
	if err := pb.ConvertToPBObject(validTypes, &converted); err != nil {
		msg := fmt.Sprintf("Unable to convert request to store %s: %s", "Valid Types", err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	trackAPIMetrics(startTime, "200", mon.GetValidTypesStr)
	return &converted, nil
}

// BulkInsertMonitoredObjects - perform a bulk operation on a set of Monitored Objects.
func (gsh *GRPCServiceHandler) BulkInsertMonitoredObjects(ctx context.Context, value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error) {
	startTime := time.Now()

	res, err := gsh.Tsh.BulkInsertMonitoredObjects(ctx, value)
	if err != nil {
		trackAPIMetrics(startTime, "500", mon.BulkUpdateMonObjStr)
		return nil, err
	}

	trackAPIMetrics(startTime, "200", mon.BulkUpdateMonObjStr)
	return res, nil
}
