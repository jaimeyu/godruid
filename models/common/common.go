package common

import "fmt"

// UserState - enum describing the state of a user.
type UserState string

const (
	UserUnknown       UserState = "USER_UNKNOWN"
	UserInvited       UserState = "INVITED"
	UserActive        UserState = "ACTIVE"
	UserSuspended     UserState = "SUSPENDED"
	UserPendingDelete UserState = "PENDING_DELETE"
)

func FailToSetID(s string) error {
	return fmt.Errorf("Can't set id for %s as no object was provided", s)
}
