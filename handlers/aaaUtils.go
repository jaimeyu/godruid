package handlers

import (
	"net/http"

	"github.com/accedian/adh-gather/gather"
)

var (
	authEnabled = true

	SkylightAdminRoleOnly       = []string{userRoleSkylight}
	SkylightAndTenantAdminRoles = []string{userRoleSkylight, userRoleTenantAdmin}
)

func InitializeAuthHelper() {
	authEnabled = GetAuthorizationToggle()
}

func isRequestAuthorized(request *http.Request, allowedRoles []string) bool {
	if !authEnabled {
		return true
	}

	return gather.DoesSliceContainString(allowedRoles, request.Header.Get(xFwdUserRoles))
}
