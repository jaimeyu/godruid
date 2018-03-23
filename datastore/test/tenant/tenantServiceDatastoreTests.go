package tenant

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/getlantern/deepcopy"

	"github.com/stretchr/testify/assert"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

var (
	defaultThresholdsBytes = []byte(`{
		"thresholds": {
			"vendorMap": {
				"accedian-flowmeter": {
					"metricMap": {
						"throughputAvg": {
							"eventAttrMap": {
								"critical": "25000000",
								"enabled": "true",
								"major": "20000000",
								"minor": "18000000"
							}
						}
					},
					"monitoredObjectTypeMap": {
						"flowmeter": {
							"metricMap": {
								"throughputAvg": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "25000000",
														"lowerStrict": "true",
														"unit": "bps"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "20000000",
														"lowerStrict": "true",
														"unit": "bps",
														"upperLimit": "25000000",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "18000000",
														"lowerStrict": "true",
														"unit": "bps",
														"upperLimit": "20000000"
													}
												}
											}
										}
									}
								}
							}
						}
					}
				},
				"accedian-twamp": {
					"metricMap": {
						"delayP95": {
							"eventAttrMap": {
								"critical": "100000",
								"enabled": "true",
								"major": "95000",
								"minor": "92500"
							}
						},
						"jitterP95": {
							"eventAttrMap": {
								"critical": "30000",
								"enabled": "true",
								"major": "20000",
								"minor": "15000"
							}
						},
						"packetsLostPct": {
							"eventAttrMap": {
								"critical": "0.8",
								"enabled": "true",
								"major": "0.3",
								"minor": "0.1"
							}
						}
					},
					"monitoredObjectTypeMap": {
						"twamp-pe": {
							"metricMap": {
								"delayP95": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "100000",
														"lowerStrict": "true",
														"unit": "ms"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "95000",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "100000",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "92500",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "95000"
													}
												}
											}
										}
									}
								},
								"jitterP95": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "30000",
														"lowerStrict": "true",
														"unit": "ms"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "20000",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "30000",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "15000",
														"lowerStrict": "true",
														"unit": "ms",
														"upperLimit": "20000"
													}
												}
											}
										}
									}
								},
								"packetsLostPct": {
									"directionMap": {
										"0": {
											"eventMap": {
												"critical": {
													"eventAttrMap": {
														"lowerLimit": "0.8",
														"lowerStrict": "true",
														"unit": "pct"
													}
												},
												"major": {
													"eventAttrMap": {
														"lowerLimit": "0.3",
														"lowerStrict": "true",
														"unit": "pct",
														"upperLimit": "0.8",
														"upperStrict": "false"
													}
												},
												"minor": {
													"eventAttrMap": {
														"lowerLimit": "0.1",
														"lowerStrict": "true",
														"unit": "pct",
														"upperLimit": "0.3"
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`)

	defaultIngestionProfileBytes = []byte(`{"metrics": {
		"vendorMap": {
		  "accedian-flowmeter": {
			"monitoredObjectTypeMap": {
			  "flowmeter": {
				"metricMap": {
				  "bytesReceived": true,
				  "packetsReceived": true,
				  "throughputAvg": true,
				  "throughputMax": true,
				  "throughputMin": true
				}
			  }
			}
		  },
		  "accedian-twamp": {
			"monitoredObjectTypeMap": {
			  "twamp-pe": {
				"metricMap": {
				  "delayMax": true,
				  "delayP95": true,
				  "delayPHi": true,
				  "delayVarP95": true,
				  "delayVarPHi": true,
				  "jitterMax": true,
				  "jitterP95": true,
				  "jitterPHi": true,
				  "lostBurstMax": true,
				  "packetsLost": true,
				  "packetsLostPct": true,
				  "packetsReceived": true
				}
			  },
			  "twamp-sf": {
				"metricMap": {
				  "delayMax": true,
				  "delayP95": true,
				  "delayPHi": true,
				  "delayVarP95": true,
				  "delayVarPHi": true,
				  "jitterMax": true,
				  "jitterP95": true,
				  "jitterPHi": true,
				  "lostBurstMax": true,
				  "packetsLost": true,
				  "packetsLostPct": true,
				  "packetsReceived": true
				}
			  },
			  "twamp-sl": {
				"metricMap": {
				  "delayMax": true,
				  "delayP95": true,
				  "delayPHi": true,
				  "delayVarP95": true,
				  "delayVarPHi": true,
				  "jitterMax": true,
				  "jitterP95": true,
				  "jitterPHi": true,
				  "lostBurstMax": true,
				  "packetsLost": true,
				  "packetsLostPct": true,
				  "packetsReceived": true
				}
			  }
			}
		  }
		}
	  }}`)
)

// TenantServiceDatastoreTestRunner - object used to run tests for any iplementation
// of the TenantServiceDatastore interface
type TenantServiceDatastoreTestRunner struct {
	tenantDB ds.TenantServiceDatastore
	adminDB  ds.AdminServiceDatastore
}

func InitTestRunner(tdb ds.TenantServiceDatastore, adb ds.AdminServiceDatastore) *TenantServiceDatastoreTestRunner {
	return &TenantServiceDatastoreTestRunner{
		tenantDB: tdb,
		adminDB:  adb,
	}
}

func (runner *TenantServiceDatastoreTestRunner) RunTenantUserCRUD(t *testing.T) {
	const COMPANY1 = "UserCompany"
	const SUBDOMAIN1 = "subdom1"
	const USER1 = "test1"
	const USER2 = "test2"
	const PASS1 = "pass1"
	const PASS2 = "pass2"
	const PASS3 = "pass3"
	const TOKEN1 = "token1"
	const TOKEN2 = "token2"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	tenantUserList, err := runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, tenantUserList)

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetTenantUser(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateTenantUser(&tenmod.User{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	tenantUser := tenmod.User{
		Username:        USER1,
		Password:        PASS1,
		OnboardingToken: TOKEN1,
		TenantID:        TENANT,
		State:           string(common.UserActive)}
	created, err := runner.tenantDB.CreateTenantUser(&tenantUser)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantUserType), created.Datatype)
	assert.Equal(t, created.Username, USER1, "Username not the same")
	assert.Equal(t, created.Password, PASS1, "Password not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, created.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantUser(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.User{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Password = PASS2
	updated, err := runner.tenantDB.UpdateTenantUser(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantUserType), updated.Datatype)
	assert.Equal(t, updated.Username, USER1, "Username not the same")
	assert.Equal(t, updated.Password, PASS2, "Password was not updated")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenantUser2 := tenmod.User{
		Username:        USER2,
		Password:        PASS3,
		OnboardingToken: TOKEN2,
		TenantID:        TENANT,
		State:           string(common.UserInvited)}
	created2, err := runner.tenantDB.CreateTenantUser(&tenantUser2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(tenmod.TenantUserType), created2.Datatype)
	assert.Equal(t, created2.Username, USER2, "Username not the same")
	assert.Equal(t, created2.Password, PASS3, "Password not the same")
	assert.Equal(t, created2.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, created2.OnboardingToken, TOKEN2, "OnboardingToken not the same")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantUser(TENANT, fetched.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Username, fetched.Username, "Deleted Username not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantUser(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantUser(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantUser(TENANT, created2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Username, created2.Username, "Deleted Username not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, tenantUserList)
}

func (runner *TenantServiceDatastoreTestRunner) RunTenantDomainCRUD(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const DOM1 = "domain1"
	const DOM2 = "domain2"
	const DOM3 = "domain3"
	const COLOR1 = "color1"
	const COLOR2 = "color2"
	const THRPRF = "ThresholdPrf"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	recList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetTenantDomain(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateTenantDomain(&tenmod.Domain{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	tenantDomain := tenmod.Domain{
		Name:                DOM1,
		TenantID:            TENANT,
		Color:               COLOR1,
		ThresholdProfileSet: []string{THRPRF}}
	created, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), created.Datatype)
	assert.Equal(t, created.Name, DOM1, "Name not the same")
	assert.Equal(t, created.Color, COLOR1, "Color not the same")
	assert.Equal(t, created.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantDomain(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.Domain{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Color = COLOR2
	updated, err := runner.tenantDB.UpdateTenantDomain(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), updated.Datatype)
	assert.Equal(t, updated.Name, DOM1, "Name not the same")
	assert.Equal(t, updated.Color, COLOR2, "Password was not updated")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenantDomain2 := tenmod.Domain{
		Name:     DOM2,
		TenantID: TENANT,
		Color:    COLOR1}
	created2, err := runner.tenantDB.CreateTenantDomain(&tenantDomain2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), created2.Datatype)
	assert.Equal(t, created2.Name, DOM2, "Name not the same")
	assert.Equal(t, created2.Color, COLOR1, "Password not the same")
	assert.Equal(t, created2.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, len(created2.ThresholdProfileSet) == 0, "Should not be a Threshold Profile ID")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantDomain(TENANT, fetched.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, fetched.Name, "Deleted name not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantDomain(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantDomain(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantDomain(TENANT, created2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, created2.Name, "Deleted name not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
}

func (runner *TenantServiceDatastoreTestRunner) RunTenantMonitoredObjectCRUD(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdomain.cool"
	const OBJNAME1 = "obj1"
	const OBJID1 = "object1"
	const OBJNAME2 = "obj2"
	const OBJID2 = "object2"
	const ACTNAME1 = "actName1"
	const ACTTYPE1 = string(tenmod.AccedianVNID)
	const ACTNAME2 = "actName2"
	const ACTTYPE2 = string(tenmod.AccedianNID)
	const REFNAME1 = "refname1"
	const REFTYPE1 = string(tenmod.AccedianNID)
	const COLOR2 = "color2"
	const DOM1 = "domain1"
	const DOM2 = "domain2"
	const DOM3 = "domain3"
	DOMAINSET1 := []string{DOM1, DOM2}
	DOMAINSET2 := []string{DOM1, DOM3}

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	recList, err := runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetMonitoredObject(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateMonitoredObject(&tenmod.MonitoredObject{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	tenantMonObj := tenmod.MonitoredObject{
		MonitoredObjectID: OBJID1,
		ObjectName:        OBJNAME1,
		TenantID:          TENANT,
		ActuatorName:      ACTNAME1,
		ActuatorType:      ACTTYPE1,
		ReflectorName:     REFNAME1,
		ReflectorType:     REFTYPE1,
		DomainSet:         DOMAINSET1,
	}
	created, err := runner.tenantDB.CreateMonitoredObject(&tenantMonObj)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantMonitoredObjectType), created.Datatype)
	assert.Equal(t, created.ObjectName, OBJNAME1, "Name not the same")
	assert.Equal(t, created.MonitoredObjectID, OBJID1, "ID not the same")
	assert.Equal(t, created.ActuatorName, ACTNAME1, "Actuator Name not the same")
	assert.Equal(t, created.ActuatorType, ACTTYPE1, "Actuator Type not the same")
	assert.Equal(t, created.ReflectorName, REFNAME1, "Reflector Name not the same")
	assert.Equal(t, created.ReflectorType, REFTYPE1, "Reflector Type not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, created.DomainSet, DOMAINSET1)
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetMonitoredObject(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.MonitoredObject{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.DomainSet = DOMAINSET2
	updated, err := runner.tenantDB.UpdateMonitoredObject(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantMonitoredObjectType), created.Datatype)
	assert.Equal(t, updated.ObjectName, OBJNAME1, "Name not the same")
	assert.Equal(t, updated.MonitoredObjectID, OBJID1, "ID not the same")
	assert.Equal(t, updated.ActuatorName, ACTNAME1, "Actuator Name not the same")
	assert.Equal(t, updated.ActuatorType, ACTTYPE1, "Actuator Type not the same")
	assert.Equal(t, updated.ReflectorName, REFNAME1, "Reflector Name not the same")
	assert.Equal(t, updated.ReflectorType, REFTYPE1, "Reflector Type not the same")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.DomainSet, DOMAINSET2, "The domain set was not updated")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenantMonObj2 := tenmod.MonitoredObject{
		MonitoredObjectID: OBJID2,
		ObjectName:        OBJNAME2,
		TenantID:          TENANT,
		ActuatorName:      ACTNAME2,
		ActuatorType:      ACTTYPE2,
		DomainSet:         DOMAINSET1}
	created2, err := runner.tenantDB.CreateMonitoredObject(&tenantMonObj2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, created2.ObjectName, OBJNAME2, "Name not the same")
	assert.Equal(t, created2.MonitoredObjectID, OBJID2, "ID not the same")
	assert.Equal(t, string(tenmod.TenantMonitoredObjectType), created2.Datatype)
	assert.Equal(t, created2.ActuatorName, ACTNAME2, "Actuator Name not the same")
	assert.Equal(t, created2.ActuatorType, ACTTYPE2, "Actuator Type not the same")
	assert.Empty(t, created2.ReflectorName, REFNAME1, "Reflector Name should not be set")
	assert.Empty(t, created2.ReflectorType, REFTYPE1, "Reflector Type should not be set")
	assert.Equal(t, created2.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, created2.DomainSet, DOMAINSET1, "The domain set was not updated")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteMonitoredObject(TENANT, fetched.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.MonitoredObjectID, fetched.MonitoredObjectID, "Deleted name not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantDomain(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that does not exist
	deleteDNE, err := runner.tenantDB.DeleteMonitoredObject(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteMonitoredObject(TENANT, created2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.MonitoredObjectID, created2.MonitoredObjectID, "Deleted name not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)

	// Create some records in Bulk:
	bulkReq := []*tenmod.MonitoredObject{&tenmod.MonitoredObject{
		MonitoredObjectID: OBJID2,
		ObjectName:        OBJNAME2,
		TenantID:          TENANT,
		ActuatorName:      ACTNAME2,
		ActuatorType:      ACTTYPE2,
		DomainSet:         DOMAINSET1},
		&tenmod.MonitoredObject{
			MonitoredObjectID: OBJID1,
			ObjectName:        OBJNAME1,
			TenantID:          TENANT,
			ActuatorName:      ACTNAME1,
			ActuatorType:      ACTTYPE1,
			ReflectorName:     REFNAME1,
			ReflectorType:     REFTYPE1,
			DomainSet:         DOMAINSET1,
		}}
	bulkResult, err := runner.tenantDB.BulkInsertMonitoredObjects(TENANT, bulkReq)
	assert.Nil(t, err)
	assert.NotEmpty(t, bulkResult)
	assert.Equal(t, 2, len(bulkResult))

	// Delete the remaining records
	fetchedList, err = runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
	assert.Equal(t, 2, len(fetchedList))

	for _, val := range fetchedList {
		del, err := runner.tenantDB.DeleteMonitoredObject(TENANT, val.ID)
		assert.Nil(t, err)
		assert.NotNil(t, del)
		assert.NotEmpty(t, del.ID)
		assert.NotEmpty(t, del.REV)
	}

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
	assert.Empty(t, fetchedList)
}

func (runner *TenantServiceDatastoreTestRunner) RunTenantMetadataCRUD(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const THRPRF = "ThresholdPrf"
	const THRPRF2 = "ThresholdPrf2"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	record, err := runner.tenantDB.GetTenantMeta(TENANT)
	assert.NotNil(t, err)
	assert.Nil(t, record)

	// Try to Update a record that does not exist:
	fail, err := runner.tenantDB.UpdateTenantMeta(&tenmod.Metadata{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	meta := tenmod.Metadata{
		TenantName:              COMPANY1,
		TenantID:                TENANT,
		DefaultThresholdProfile: THRPRF}
	created, err := runner.tenantDB.CreateTenantMeta(&meta)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantMetaType), created.Datatype)
	assert.Equal(t, created.TenantName, COMPANY1, "Name not the same")
	assert.Equal(t, created.DefaultThresholdProfile, THRPRF, "Threshold Profile not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantMeta(TENANT)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.Metadata{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.DefaultThresholdProfile = THRPRF2
	updated, err := runner.tenantDB.UpdateTenantMeta(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantMetaType), updated.Datatype)
	assert.Equal(t, updated.TenantName, COMPANY1, "Name not the same")
	assert.Equal(t, updated.DefaultThresholdProfile, THRPRF2, "Threshold Profile was not updated")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	meta2 := tenmod.Metadata{
		TenantName:              COMPANY1,
		TenantID:                TENANT,
		DefaultThresholdProfile: THRPRF}
	created2, err := runner.tenantDB.CreateTenantMeta(&meta2)
	assert.NotNil(t, err)
	assert.Nil(t, created2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantMeta(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.Equal(t, deleted.TenantName, fetched.TenantName, "Deleted name not the same")

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantMeta(TENANT)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)
}

func (runner *TenantServiceDatastoreTestRunner) RunTenantIngestionProfileCRUD(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const DOM1 = "domain1"
	const DOM2 = "domain2"
	const DOM3 = "domain3"
	const COLOR1 = "color1"
	const COLOR2 = "color2"
	const THRPRF = "ThresholdPrf"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	rec, err := runner.tenantDB.GetActiveTenantIngestionProfile(TENANT)
	assert.NotNil(t, err)
	assert.Nil(t, rec)

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetTenantIngestionProfile(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateTenantIngestionProfile(&tenmod.IngestionProfile{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	defaultIngestionProfileShell := &tenmod.IngestionProfile{}
	if err := json.Unmarshal(defaultIngestionProfileBytes, &defaultIngestionProfileShell); err != nil {
		logger.Log.Debugf("Unable to umarshal ingestion profile: %s", err.Error())
	}

	ingPrf := tenmod.IngestionProfile{
		Datatype: string(tenmod.TenantIngestionProfileType),
		TenantID: TENANT,
		Metrics:  defaultIngestionProfileShell.Metrics,
	}
	created, err := runner.tenantDB.CreateTenantIngestionProfile(&ingPrf)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantIngestionProfileType), created.Datatype)
	assert.Equal(t, created.Metrics, defaultIngestionProfileShell.Metrics, "Metrics not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantIngestionProfile(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.IngestionProfile{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Metrics["badStuff"] = nil
	updated, err := runner.tenantDB.UpdateTenantIngestionProfile(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantIngestionProfileType), updated.Datatype)
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.NotEqual(t, updated.Metrics, defaultIngestionProfileShell.Metrics, "Metrics were not updated")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record - should fail.
	tenantDomain2 := tenmod.IngestionProfile{
		Datatype: string(tenmod.TenantIngestionProfileType),
		TenantID: TENANT,
		Metrics:  defaultIngestionProfileShell.Metrics,
	}
	created2, err := runner.tenantDB.CreateTenantIngestionProfile(&tenantDomain2)
	assert.NotNil(t, err)
	assert.Nil(t, created2)

	// Get active records
	active, err := runner.tenantDB.GetActiveTenantIngestionProfile(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, active)
	assert.Equal(t, updated, active)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantIngestionProfile(TENANT, active.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, active, deleted, "Deleted not the same")

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantDomain(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)
}
