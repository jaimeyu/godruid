package handlers

import (
	"net/http"

	"github.com/accedian/adh-gather/gather"
)

var (
	authEnabled               = true
	changeNotificationEnabled = true

	SkylightAdminRoleOnly       = []string{userRoleSkylight}
	SkylightAndTenantAdminRoles = []string{userRoleSkylight, userRoleTenantAdmin}
	AllRoles                    = []string{userRoleSkylight, userRoleTenantAdmin, userRoleTenantUser}
)

// TODO: Make this better as I do not like how it is just free-floating vars on the package
func InitializeAuthHelper() {
	authEnabled = GetAuthorizationToggle()
	changeNotificationEnabled = GetChangeNotificationsToggle()
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
