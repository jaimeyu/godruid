package handlers

import (
	"errors"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

func validateAdminUserRequest(request *pb.AdminUserRequest, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid AdminUserRequest: no Admin User data provided")
	}

	if len(request.GetXId()) == 0 {
		return errors.New("Invalid AdminUserRequest: no Admin User ID provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantUserRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantDescriptorRequest(request *pb.TenantDescriptorRequest, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantDescriptorRequest: no Tenant Descriptor data provided")
	}

	if len(request.GetXId()) == 0 {
		return errors.New("Invalid TenantDescriptorRequest: no Tenant ID provided")
	}

	if isUpdate && (len(request.GetXRev()) == 0 || request.GetData().GetCreatedTimestamp() == 0) {
		return errors.New("Invalid TenantUserRequest: must provide a createdTimestamp and revision for an update")
	}

	return nil
}

func validateTenantUserRequest(request *pb.TenantUserRequest, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantUserRequest: no Tenant User data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantUserRequest: no Tenant ID provided")
	}

	if len(request.GetXId()) == 0 {
		return errors.New("Invalid TenantUserRequest: no Tenant User ID provided")
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

func validateTenantDomainRequest(request *pb.TenantDomainRequest, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantDomainRequest: no Tenant Domain data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantDomainRequest: no Tenant Id provided")
	}

	if len(request.GetXId()) == 0 {
		return errors.New("Invalid TenantDomainRequest: no Tenant Domain ID provided")
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

func validateTenantIngPrfRequest(request *pb.TenantIngestionProfileRequest, isUpdate bool) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantIngestionProfileRequest: no Tenant Ingestion Profile data provided")
	}

	if len(request.GetData().GetTenantId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileRequest: no Tenant Id provided")
	}

	if len(request.GetXId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileRequest: no Ingestion Profile ID provided")
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

func validateMonitoredObjectRequest(request *pb.MonitoredObjectRequest, isUpdate bool) error {
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
