package logger_test

import (
	"strings"
	"testing"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/stretchr/testify/assert"
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

	result := models.AsJSONString(adminUser)
	assert.NotEmpty(t, result)
	assert.NotEqual(t, prettyUser, result)
	assert.True(t, strings.Contains(result, logger.LogRedactStr))

}
