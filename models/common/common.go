package common

import "fmt"

// RESTValidator - Interface for validation of models used in REST requests.
type RESTValidator interface {
	Validate(isUpdate bool) error
}

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

// BulkOperationResult - response for a bulk update.
type BulkOperationResult struct {
	OK     bool   `json:"ok"`
	ID     string `json:"id"`
	REV    string `json:"rev"`
	ERROR  string `json:"error"`
	REASON string `json:"reason"`
}

// PaginationOffsets - struct to hold offsets for pagination related queries
type PaginationOffsets struct {
	Self string
	Prev string
	Next string
}
