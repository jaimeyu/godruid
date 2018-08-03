package handlers

<<<<<<< HEAD
import (
	//"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"testing"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestExtractHeader(t *testing.T) {

	logger.Log.Debug("Starting Header Text")
	h := http.Header{}

	h.Add(XFwdUserRoles, UserRoleSkylight)
	h.Add(XFwdUserName, "user")
	h.Add(XFwdUserId, "0")
	h.Add(XFwdTenantId, "0")

	user, err := ExtractHeaderToUserAuthRequest(h)
	assert.Nil(t, err)
	assert.Equal(t, len(user.UserRoles), 1)
	assert.Equal(t, user.UserRoles[0], UserRoleSkylight)
	assert.Equal(t, user.UserName, "user")
	assert.Equal(t, user.UserID, "0")
	assert.Equal(t, user.TenantID, "0")

	rolesstr := fmt.Sprintf("%s,%s", UserRoleSkylight, UserRoleTenantUser)
	h = http.Header{}
	user = &RequestUserAuth{}
	h.Add(XFwdUserRoles, rolesstr)

	user, err = ExtractHeaderToUserAuthRequest(h)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(user.UserRoles))
	assert.Equal(t, UserRoleSkylight, user.UserRoles[0])
	assert.Equal(t, UserRoleTenantUser, user.UserRoles[1])

}

type mockProvider struct {
	Active bool
	val    string
}

func (m *mockProvider) GetBool(x string) bool {
	return m.Active
}
func (m *mockProvider) Get(x string) interface{} {
	return nil
}

func (m *mockProvider) GetInt(x string) int {
	return 0
}
func (m *mockProvider) GetString(x string) string {
	return ""
}
func (m *mockProvider) GetStringSlice(x string) []string {
	return nil
}
func (m *mockProvider) GetStringMap(x string) map[string]interface{} {
	return nil
}

func (m *mockProvider) GetStringMapString(key string) map[string]string {
	return nil
}

func (m *mockProvider) Set(key string, value interface{}) {
	return
}

func (m *mockProvider) IsSet(key string) bool {
	return false
}

type mockGather struct {
	Cfg mockProvider
}

func (m *mockGather) GetConfig() mockProvider {
	return m.Cfg
}

func TestRAC(t *testing.T) {
=======
//"encoding/json"

// func TestExtractHeader(t *testing.T) {

// 	logger.Log.Debug("Starting Header Text")
// 	h := http.Header{}

// 	h.Add(XFwdUserRoles, UserRoleSkylight)
// 	h.Add(XFwdUserName, "user")
// 	h.Add(XFwdUserId, "0")
// 	h.Add(XFwdTenantId, "0")

// 	user, err := ExtractHeaderToUserAuthRequest(h)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(user.UserRoles), 1)
// 	assert.Equal(t, user.UserRoles[0], UserRoleSkylight)
// 	assert.Equal(t, user.UserName, "user")
// 	assert.Equal(t, user.UserID, "0")
// 	assert.Equal(t, user.TenantID, "0")

// 	rolesstr := fmt.Sprintf("%s,%s", UserRoleSkylight, UserRoleTenantUser)
// 	h = http.Header{}
// 	user = &RequestUserAuth{}
// 	h.Add(XFwdUserRoles, rolesstr)

// 	user, err = ExtractHeaderToUserAuthRequest(h)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 2, len(user.UserRoles))
// 	assert.Equal(t, UserRoleSkylight, user.UserRoles[0])
// 	assert.Equal(t, UserRoleTenantUser, user.UserRoles[1])

// }

// type mockProvider struct {
// 	Active bool
// 	val    string
// }

// func (m *mockProvider) GetBool(x string) bool {
// 	return m.Active
// }
// func (m *mockProvider) Get(x string) interface{} {
// 	return nil
// }

// func (m *mockProvider) GetInt(x string) int {
// 	return 0
// }
// func (m *mockProvider) GetString(x string) string {
// 	return ""
// }
// func (m *mockProvider) GetStringSlice(x string) []string {
// 	return nil
// }
// func (m *mockProvider) GetStringMap(x string) map[string]interface{} {
// 	return nil
// }

// func (m *mockProvider) GetStringMapString(key string) map[string]string {
// 	return nil
// }

// func (m *mockProvider) Set(key string, value interface{}) {
// 	return
// }
>>>>>>> upstream/master

// func (m *mockProvider) IsSet(key string) bool {
// 	return false
// }

// type mockGather struct {
// 	Cfg mockProvider
// }

// func (m *mockGather) GetConfig() mockProvider {
// 	return m.Cfg
// }

// func TestRAC(t *testing.T) {

// 	h := http.Header{}

// 	v := viper.New()
// 	gather.LoadDefaults(v)

// 	cfg := gather.LoadConfig("../config/adh-gather-debug.yml", v)
// 	cfg.Set(gather.CK_args_authorizationAAA.String(), true)

// 	// Deny access
// 	res := RoleAccessControl(h, nil)
// 	assert.Equal(t, false, res)

// 	res = RoleAccessControl(h, []string{})
// 	assert.Equal(t, false, res)

// 	h.Add(XFwdUserRoles, UserRoleSkylight)
// 	h.Add(XFwdUserName, "user")
// 	h.Add(XFwdUserId, "0")
// 	h.Add(XFwdTenantId, "0")

// 	// Deny access
// 	res = RoleAccessControl(h, []string{UserRoleTenantUser})
// 	assert.Equal(t, false, res)

// 	res = RoleAccessControl(h, []string{UserRoleTenantAdmin})
// 	assert.Equal(t, false, res)

// 	res = RoleAccessControl(h, []string{UserRoleUnknown})
// 	assert.Equal(t, false, res)

// 	// Allow access
// 	res = RoleAccessControl(h, []string{UserRoleSkylight})
// 	assert.Equal(t, true, res)

// 	// Allow access because RAC is disabled
// 	cfg.Set(gather.CK_args_authorizationAAA.String(), false)
// 	res = RoleAccessControl(h, []string{UserRoleTenantUser})
// 	assert.Equal(t, true, res)

// }

// type mockResp struct{}

// func (m *mockResp) Header() http.Header {
// 	return http.Header{}
// }
// func (m *mockResp) Write([]byte) (int, error) {
// 	return 0, nil
// }
// func (m *mockResp) WriteHeader(statusCode int) {
// }

// func TestBuildFunctor(t *testing.T) {

// 	w := mockResp{}
// 	r := http.Request{}

// 	r.URL = &url.URL{Path: "test/path"}
// 	r.Header = http.Header{}
// 	h := http.Header{}
// 	h.Add(XFwdUserRoles, UserRoleSkylight)
// 	h.Add(XFwdUserName, "user")
// 	h.Add(XFwdUserId, "0")
// 	h.Add(XFwdTenantId, "0")
// 	r.Header = h

// 	passed := false

// 	v := viper.New()
// 	gather.LoadDefaults(v)

// 	cfg := gather.LoadConfig("../config/adh-gather-debug.yml", v)
// 	cfg.Set(gather.CK_args_authorizationAAA.String(), true)

// 	functorHttpHandler := func(w http.ResponseWriter, r *http.Request) {
// 		passed = true

// 		logger.Log.Debug("Test passed!")
// 	}

// 	// Should allow access (AAA is disabled)
// 	passed = false

// 	cfg.Set(gather.CK_args_authorizationAAA.String(), false)

// 	functor2test3 := BuildRouteHandlerWithRAC([]string{UserRoleTenantUser}, functorHttpHandler)

// 	functor2test3(&w, &r)

// 	if passed == false {
// 		t.Fail()
// 	}

// 	// The following test works but depends on a panic to occur because the ReportError
// 	// expects working http objects which I do not have (can work if we use a http mock module).
// 	// Disabling the test for now since we shouldn't allow the test to panic even if we recover.

// 	// Should deny access
// 	//functor2test := BuildRouteHandlerWithRAC([]string{UserRoleSkylight}, functorHttpHandler)

// 	//functor2test(&w, &r)

<<<<<<< HEAD
	//if passed == true {
	//	t.Fail()
	//}
}
func TestUnion(t *testing.T) {

	listA := []string{"a", "b"}
	listB := []string{"a", "c"}

	unionList := listUnion(listA, listB)
	sort.Strings(unionList)

	if unionList[0] != "a" && unionList[1] != "b" && unionList[2] != "c" {
		t.Errorf("Incorrect union between %v and %v: %v", listA, listB, unionList)
	}

}

func TestIntersection(t *testing.T) {
	listA := []string{"a", "b"}
	listB := []string{"a", "c"}

	intersectionList := listIntersection(listA, listB)

	if len(intersectionList) != 1 {
		t.Errorf("Was expecting intersection list size of %d but got %d instead", 1, len(intersectionList))
	}

	if intersectionList[0] != "a" {
		t.Errorf("Incorrect union between %v and %v: %v", listA, listB, intersectionList)
	}
}
=======
// 	//if passed == false {
// 	//	t.Fail()
// 	//}

// 	//cfg.Set(gather.CK_args_authorizationAAA.String(), true)
// 	//passed = false

// 	//cfg.Set(gather.CK_args_authorizationAAA.String(), true)

// 	//functor2test2 := BuildRouteHandlerWithRAC([]string{UserRoleTenantUser}, functorHttpHandler)

// 	//defer func() {
// 	//	if r := recover(); r != nil {
// 	//		fmt.Println("Recovered expected runtime error in test", r)
// 	//		passed = false
// 	//	}
// 	//}()
// 	//functor2test2(&w, &r)

// 	//if passed == true {
// 	//	t.Fail()
// 	//}
// }
>>>>>>> upstream/master
