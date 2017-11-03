package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/leesper/couchdb-golang"

	pb "github.com/accedian/adh-gather/gathergrpc"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

const dbName string = "adh-admin"

// AdminServiceHandler - implementation of the interface for the gRPC
// Admin service. Anytime the Admin service changes, the logic to handle the
// API will be modified here.
type AdminServiceHandler struct {
	provDB string
}

// CreateHandler - used to generate a handler for the Admin Service.
func CreateHandler(provDBURL string) *AdminServiceHandler {
	result := new(AdminServiceHandler)
	result.provDB = provDBURL
	return result
}

// CreateAdminUser - Create an Administrative User.
func (ash *AdminServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Connect to PROV DB
	dbName := getDbName(ash.provDB)
	db, err := couchdb.NewDatabase(dbName)
	if err != nil {
		log.Printf("Unable to connect to Prov DB %s: %v\n", ash.provDB, err)
		return nil, err
	}

	log.Printf("Using DB %s to Create useer %v \n", dbName, user)

	// Give the user a known id and timestamps:
	user.XId = user.Username
	user.CreatedTimestamp = time.Now().Unix()
	user.LastModifiedTimestamp = user.GetCreatedTimestamp()

	// Marshal the Admin and read the bytes as string.
	userToBytes, err := json.Marshal(user)
	if err != nil {
		log.Printf("Unable to convert user to format to persist: %v\n", err)
		return nil, err
	}
	var storeFormat map[string]interface{}
	err = json.Unmarshal(userToBytes, &storeFormat)
	if err != nil {
		log.Printf("Unable to convert user to format to persist: %v\n", err)
		return nil, err
	}

	log.Printf("Attempting to store user: %v", storeFormat)

	// Store the user in PROV DB
	options := new(url.Values)
	id, rev, err := db.Save(storeFormat, *options)
	if err != nil {
		log.Printf("Unable to store admin user: %v\n", err)
		return nil, err
	}
	log.Printf("Successfully stored user %s with rev %s", id, rev)

	// Return the provisioned user.
	log.Printf("Stored user: %v\n", user)
	return user, nil
}

// UpdateAdminUser - Update an Administrative User.
func (ash *AdminServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// DeleteAdminUser - Delete an Administrative User.
func (ash *AdminServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (ash *AdminServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	// Connect to PROV DB
	dbName := getDbName(ash.provDB)
	db, err := couchdb.NewDatabase(dbName)
	if err != nil {
		log.Printf("Unable to connect to Prov DB %s: %v\n", ash.provDB, err)
		return nil, err
	}

	log.Printf("Using db %s to GET Admin User %s\n", dbName, userID.Value)

	// Get the user from PROV DB
	options := new(url.Values)
	fetchedUser, err := db.Get(userID.Value, *options)
	if err != nil {
		log.Printf("Error retrieving user %s: %v\n", userID.Value, err)
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	fetchedUserInBytes, err := json.Marshal(fetchedUser)
	if err != nil {
		fmt.Printf("Error converting retrieved user to proper type: %v\n", err)
	}

	res := pb.AdminUser{}
	json.Unmarshal(fetchedUserInBytes, &res)
	log.Printf("Retrieved user: %v\n", res)
	return &res, nil
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (ash *AdminServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	// Stub to implement
	return nil, nil
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (ash *AdminServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (ash *AdminServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (ash *AdminServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (ash *AdminServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// Helper methods

// func generateAdminUser(username string,
// 	password string,
// 	onboardingToken string,
// 	sendEmail bool,
// 	verified bool,
// 	state pb.UserState,
// 	created int64,
// 	lastMod int64) *pb.AdminUser {

// 	return &pb.AdminUser{
// 		XId:                   generateAdminUserId(username),
// 		Username:              username,
// 		Password:              password,
// 		SendOnboardingEmail:   sendEmail,
// 		OnboardingToken:       onboardingToken,
// 		UserVerified:          verified,
// 		State:                 state,
// 		CreatedTimestamp:      created,
// 		LastModifiedTimestamp: lastMod}
// }

func getDbName(provDbUrl string) string {
	return provDbUrl + "/" + dbName
}
