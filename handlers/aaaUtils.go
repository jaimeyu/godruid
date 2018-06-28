package handlers

import (
	"net/http"

	"github.com/accedian/adh-gather/gather"
)

var (
	authEnabled = true

	SkylightAdminRoleOnly       = []string{userRoleSkylight}
	SkylightAndTenantAdminRoles = []string{userRoleSkylight, userRoleTenantAdmin}
	AllRoles                    = []string{userRoleSkylight, userRoleTenantAdmin, userRoleTenantUser}
)

func InitializeAuthHelper() {
	authEnabled = GetAuthorizationToggle()
}

func isRequestAuthorized(request *http.Request, allowedRoles []string) bool {
	// No need for Authorization check if Authorization is disabled
	if !authEnabled {
		return true
	}

	requestRole := request.Header.Get(xFwdUserRoles)

	// Allow system elements to have access
	if requestRole == userRoleSystem {
		return true
	}

	return gather.DoesSliceContainString(allowedRoles, requestRole)
}
