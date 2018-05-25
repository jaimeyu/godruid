package models

import (
	"strings"
	"testing"

	"encoding/json"
	"fmt"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestPrettyPrint(t *testing.T) {
	prettyUser := `{
		"_id": "someID",
		"_rev": "someREV,
		"data": {
		  "createdTimestamp": 123,
		  "datatype": "user",
		  "lastModifiedTimestamp": 456,
		  "onboardingToken": "token",
		  "password": "admin",
		  "state": 2,
		  "tenantId": "tenant123",
		  "userVerified": true,
		  "username": "admin@nopers.com"
		}
		}`
	adminUserData := pb.TenantUserData{
		CreatedTimestamp:      123,
		Datatype:              "user",
		LastModifiedTimestamp: 456,
		OnboardingToken:       "token",
		Password:              "admin",
		State:                 2,
		TenantId:              "tenant123",
		UserVerified:          true,
		Username:              "admin@nopers.com"}
	adminUser := &pb.TenantUser{
		XId:  "someID",
		XRev: "someREV",
		Data: &adminUserData}

	result := AsJSONString(adminUser)
	assert.NotEmpty(t, result)
	assert.NotEqual(t, prettyUser, result)
	assert.True(t, strings.Contains(result, LogRedactStr))

}

/*TestMergeMap -- Tests the function MergeMap() */
func TestMergeMap(t *testing.T) {
	// Create some test data
	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})

	expected := make(map[string]interface{})
	expected["Hello"] = "world"
	expected["Boom"] = "shakalaka"
	expected["Dummy"] = "value"
	expected["IntKey"] = 42

	m2["Hello"] = "world"
	m2["Boom"] = "shakalaka"

	// This merge should pass because m1 is empty and m2 has data. So m1 should == m2
	MergeMaps(m1, m2)
	if !reflect.DeepEqual(m1, m2) {
		t.Log("Failed simple merge test")
		t.Fail()
	}

	//t.Log("Complex merge test")
	m1["Dummy"] = "value"
	m2["IntKey"] = 42

	// m1 has some data and m2 now has new data to overwrite
	MergeMaps(m1, m2)

	if reflect.DeepEqual(m1, m2) {
		t.Log("Failed complex merge test")
		t.Fail()
	} else {
		if !reflect.DeepEqual(m1, expected) {
			fmt.Printf("Expected dst:\t%v\n", expected)
			fmt.Printf("Received dst:\t%v\n", m1)
			t.Log("Complex merge did not validate")
			t.Fail()
		}
	}
}

type testStruct struct {
	Hello  string
	Boom   string
	Dummy  string
	IntKey int
}

//setupTest() -- Setups the variables to to feed the functions
func setupTest() (map[string]interface{}, interface{}) {
	// Create the attribute map
	expected := make(map[string]interface{})
	expected["Hello"] = "world"
	expected["Boom"] = "shakalaka"
	expected["Dummy"] = "value"
	expected["IntKey"] = 42

	e1 := make(map[string]interface{})
	e2 := make(map[string]interface{})

	// Now manually create a JSON-API like map
	e2["attributes"] = expected
	e1["data"] = e2

	// Fill in the expected output
	expectedSt := testStruct{Hello: "world",
		Boom:   "shakalaka",
		Dummy:  "value",
		IntKey: 42}
	return e1, expectedSt
}

//TestMergeMapIntoOBject -- Tests MergeMapIntoOBject with valid data
func TestMergeMapIntoObject(t *testing.T) {
	e1, expectedSt := setupTest()
	var newSt = testStruct{}

	js, errMarsh := json.Marshal(e1)
	if errMarsh != nil {
		t.Fail()
	}

	// Test empty values
	var errEmpty error
	errEmpty = MergeObjWithMap(&newSt, []byte("hello"))
	if errEmpty == nil {
		t.Fail()
	}

	//	fmt.Printf("Doing merge on struct\n\n")
	err := MergeObjWithMap(&newSt, js)
	if err != nil {
		fmt.Println("Failed merging with error: ", err)
		t.Fail()
	}

	//	fmt.Printf("Expectedst:%+v\n", expectedSt)
	//	fmt.Printf("newSt     :%+v\n", newSt)
	if !reflect.DeepEqual(newSt, expectedSt) {
		t.Fail()
	}
}

//TestMergeMapBadJSON -- Tests MergeObjWithMap with bad json data
func TestMergeMapBadJSON(t *testing.T) {
	newSt := &testStruct{}
	var errEmpty error
	errEmpty = MergeObjWithMap(&newSt, []byte("hello"))
	if errEmpty == nil {
		t.Fail()
	}
}

//TestMergeMapBadObj -- Tests MergeObjWithMap with bad object data
func TestMergeMapBadObj(t *testing.T) {
	var errEmpty error
	errEmpty = MergeObjWithMap(nil, []byte("hello"))
	if errEmpty == nil {
		t.Fail()
	}
}
