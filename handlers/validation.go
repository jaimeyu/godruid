package handlers

import (
	"errors"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

func validateAdminUserRequest(request *pb.AdminUser, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid AdminUserRequest: no Admin User data provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantUserRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantDescriptorRequest(request *pb.TenantDescriptor, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantDescriptorRequest: no Tenant Descriptor data provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantUserRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantUserRequest(request *pb.TenantUser, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantUserRequest: no Tenant User data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantUserRequest: no Tenant ID provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantUserRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantUserIDRequest(request *pb.TenantUserIdRequest) error {
	if request == nil || len(request.GetUserId()) == 0 {
		return errors.New("Invalid TenantUserIdRequest: no Tenant User ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantUserIdRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantDomainRequest(request *pb.TenantDomain, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantDomainRequest: no Tenant Domain data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantDomainRequest: no Tenant Id provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantDomainRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantDomainIDRequest(request *pb.TenantDomainIdRequest) error {
	if request == nil || len(request.GetDomainId()) == 0 {
		return errors.New("Invalid TenantDomainIdRequest: no Tenant Domain ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantDomainIdRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantIngPrfRequest(request *pb.TenantIngestionProfile, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantIngestionProfileRequest: no Tenant Ingestion Profile data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileRequest: no Tenant Id provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantIngestionProfileRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantIngPrfIDRequest(request *pb.TenantIngestionProfileIdRequest) error {
	if request == nil || len(request.GetIngestionProfileId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileIdRequest: no Ingestion Profile ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileIdRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantThreshPrfRequest(request *pb.TenantThresholdProfile, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantThresholdProfileRequest: no Tenant Threshold Profile data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantThresholdProfileRequest: no Tenant Id provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantThresholdProfileRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantThreshPrfIDRequest(request *pb.TenantThresholdProfileIdRequest) error {
	if request == nil || len(request.GetThresholdProfileId()) == 0 {
		return errors.New("Invalid TenantThresholdProfileIdRequest: no Threshold Profile ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantThresholdProfileIdRequest: no Tenant Id provided")
	}

	return nil
}

func validateMonitoredObjectRequest(request *pb.MonitoredObject, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid MonitoredObjectRequest: no Tenant Monitored Object data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid MonitoredObjectRequest: no Tenant Id provided")
	}

	if len(request.GetData().GetId()) == 0 {
		return errors.New("Invalid MonitoredObjectRequest: no Monitored Object Id provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid MonitoredObjectRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateMonitoredObjectIDRequest(request *pb.MonitoredObjectIdRequest) error {
	if request == nil || len(request.GetMonitoredObjectId()) == 0 {
		return errors.New("Invalid MonitoredObjectIdRequest: no Monitored Object ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid MonitoredObjectIdRequest: no Tenant Id provided")
	}

	return nil
}

func validateMonitoredObjectToDomainMapRequest(request *pb.MonitoredObjectCountByDomainRequest) error {
	if request == nil {
		return errors.New("Invalid MonitoredObjectCountByDomainRequest: no request data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid MonitoredObjectCountByDomainRequest: no Tenant Id provided")
	}

	return nil
}

func validateIngestionDictionary(request *pb.IngestionDictionary, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid IngestionDictionary: no IngestionDictionary data provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid IngestionDictionary: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantMetaRequest(request *pb.TenantMetadata, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantMeta: no Tenant Meta data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantMeta: no Tenant Id provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantMeta: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateValidTypes(request *pb.ValidTypes, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid ValidTypes: no ValidTypes data provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid ValidTypes: must provide a createdTimestamp and revision for an update")
	}

	return nil
}
