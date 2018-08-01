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
	metmod "github.com/accedian/adh-gather/models/metrics"
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

func (runner *TenantServiceDatastoreTestRunner) RunTenantConnectorConfigCRUD(t *testing.T) {

	const COMPANY1 = "ConnectorCompany"
	const SUBDOMAIN1 = "subdom1"
	const CONN1 = "connector1"
	const CONN2 = "connector2"
	const EXPORTGROUP1 = "group1"
	const EXPORTGROUP2 = "group2"

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
	recList, err := runner.tenantDB.GetAllTenantConnectorConfigs(TENANT, "")
	assert.Nil(t, err)
	assert.NotNil(t, recList)

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetTenantConnectorConfig(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateTenantConnectorConfig(&tenmod.ConnectorConfig{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	TenantConnectorConfig := tenmod.ConnectorConfig{
		Name:        CONN1,
		TenantID:    TENANT,
		ExportGroup: EXPORTGROUP1,
	}

	created, err := runner.tenantDB.CreateTenantConnectorConfig(&TenantConnectorConfig)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantConnectorConfigType), created.Datatype)
	assert.Equal(t, created.Name, CONN1, "Name not the same")
	assert.Equal(t, created.ExportGroup, EXPORTGROUP1, "Export group not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantConnectorConfig(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.ConnectorConfig{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.ExportGroup = EXPORTGROUP2
	updated, err := runner.tenantDB.UpdateTenantConnectorConfig(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), updated.Datatype)
	assert.Equal(t, updated.Name, CONN1, "Name not the same")
	assert.Equal(t, updated.ExportGroup, EXPORTGROUP2, "Export Group was not updated")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	TenantConnectorConfig2 := tenmod.ConnectorConfig{
		Name:        CONN2,
		TenantID:    TENANT,
		ExportGroup: EXPORTGROUP1}
	created2, err := runner.tenantDB.CreateTenantConnectorConfig(&TenantConnectorConfig2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), created2.Datatype)
	assert.Equal(t, created2.Name, CONN2, "Name not the same")
	assert.Equal(t, created2.ExportGroup, EXPORTGROUP1, "Export group not the same")
	assert.Equal(t, created2.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantConnectorConfigs(TENANT, "")
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantConnectorConfig(TENANT, fetched.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, fetched.Name, "Deleted name not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantConnectorConfigs(TENANT, "")
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantConnectorConfig(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantConnectorConfig(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantConnectorConfig(TENANT, created2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, created2.Name, "Deleted name not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllTenantConnectorConfigs(TENANT, "")
	assert.Nil(t, err)
	assert.NotNil(t, recList)
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
		Name:     DOM1,
		TenantID: TENANT,
		Color:    COLOR1}

	created, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), created.Datatype)
	assert.Equal(t, created.Name, DOM1, "Name not the same")
	assert.Equal(t, created.Color, COLOR1, "Color not the same")
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
	const OBJNAME3 = "obj3"
	const OBJID3 = "object3"
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
	bulkReq := []*tenmod.MonitoredObject{
		&tenmod.MonitoredObject{
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
		},
		&tenmod.MonitoredObject{
			MonitoredObjectID: OBJID3,
			ObjectName:        OBJNAME3,
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
	assert.Equal(t, 3, len(bulkResult))

	fetchedList, err = runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
	assert.Equal(t, 3, len(fetchedList))
	// Make sure everything has an ID and revision
	for _, v := range fetchedList {
		assert.NotNil(t, v)
		assert.NotEmpty(t, v.ID)
		assert.NotEmpty(t, v.REV)
	}

	// Modify only 2 of the monitored objects
	bulkUpdateRequest := make([]*tenmod.MonitoredObject, 2)
	var objId1, objId2, objId3 string
	for i, v := range fetchedList {
		if fetchedList[i].MonitoredObjectID == OBJID1 {
			objId1 = fetchedList[i].ID
			// change the domain set
			obj := &tenmod.MonitoredObject{}
			deepcopy.Copy(obj, v)
			obj.DomainSet = DOMAINSET2
			bulkUpdateRequest[0] = obj
		} else if fetchedList[i].MonitoredObjectID == OBJID2 {
			objId2 = fetchedList[i].ID
			// clear the domains
			obj := &tenmod.MonitoredObject{}
			deepcopy.Copy(obj, v)
			obj.DomainSet = nil
			bulkUpdateRequest[1] = obj
		} else if fetchedList[i].MonitoredObjectID == OBJID3 {
			objId3 = fetchedList[i].ID
		}
		assert.ElementsMatch(t, DOMAINSET1, fetchedList[i].DomainSet)
	}

	bulkResult, err = runner.tenantDB.BulkUpdateMonitoredObjects(TENANT, bulkUpdateRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, bulkResult)
	assert.Equal(t, 2, len(bulkResult))

	fetched, err = runner.tenantDB.GetMonitoredObject(TENANT, objId1)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.ElementsMatch(t, DOMAINSET2, fetched.DomainSet)

	fetched, err = runner.tenantDB.GetMonitoredObject(TENANT, objId2)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.Nil(t, fetched.DomainSet)

	fetched, err = runner.tenantDB.GetMonitoredObject(TENANT, objId3)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.ElementsMatch(t, DOMAINSET1, fetched.DomainSet)

	// Delete the remaining records
	fetchedList, err = runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
	assert.Equal(t, 3, len(fetchedList))

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
		TenantName: COMPANY1,
		TenantID:   TENANT}
	created, err := runner.tenantDB.CreateTenantMeta(&meta)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantMetaType), created.Datatype)
	assert.Equal(t, created.TenantName, COMPANY1, "Name not the same")
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
	updated, err := runner.tenantDB.UpdateTenantMeta(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantMetaType), updated.Datatype)
	assert.Equal(t, updated.TenantName, COMPANY1, "Name not the same")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	meta2 := tenmod.Metadata{
		TenantName: COMPANY1,
		TenantID:   TENANT}
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
	updateRecord.Metrics.VendorMap["badStuff"] = tenmod.IngPrfMonitoredObjectTypeMap{}
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

func (runner *TenantServiceDatastoreTestRunner) RunTenantThresholdProfileCRUD(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const NAME1 = "name1"
	const NAME2 = "name2"

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
	recList, err := runner.tenantDB.GetAllTenantThresholdProfile(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetTenantThresholdProfile(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateTenantThresholdProfile(&tenmod.ThresholdProfile{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	defaultThresholdProfileShell := &tenmod.ThresholdProfile{}
	if err := json.Unmarshal(defaultThresholdsBytes, &defaultThresholdProfileShell); err != nil {
		logger.Log.Debugf("Unable to umarshal threshold profile: %s", err.Error())
	}

	original := tenmod.ThresholdProfile{
		Datatype:   string(tenmod.TenantThresholdProfileType),
		TenantID:   TENANT,
		Name:       NAME1,
		Thresholds: defaultThresholdProfileShell.Thresholds,
	}
	created, err := runner.tenantDB.CreateTenantThresholdProfile(&original)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantThresholdProfileType), created.Datatype)
	assert.Equal(t, created.Name, NAME1, "Names not the same")
	assert.Equal(t, created.Thresholds, defaultThresholdProfileShell.Thresholds, "Thresholds not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantThresholdProfile(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.ThresholdProfile{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Thresholds.VendorMap["badstuff"] = tenmod.ThrPrfMetric{}
	updated, err := runner.tenantDB.UpdateTenantThresholdProfile(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantThresholdProfileType), updated.Datatype)
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.NotEqual(t, updated.Thresholds, fetched.Thresholds, "Thresholds were not updated")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	second := tenmod.ThresholdProfile{
		Datatype:   string(tenmod.TenantThresholdProfileType),
		TenantID:   TENANT,
		Name:       NAME2,
		Thresholds: defaultThresholdProfileShell.Thresholds}
	created2, err := runner.tenantDB.CreateTenantThresholdProfile(&second)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(tenmod.TenantThresholdProfileType), created2.Datatype)
	assert.Equal(t, created2.Name, NAME2, "Name not the same")
	assert.Equal(t, created2.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, created2.Thresholds, defaultThresholdProfileShell.Thresholds, "Thresholds not the same")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantThresholdProfile(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantThresholdProfile(TENANT, fetched.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, fetched.Name, "Deleted name not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantThresholdProfile(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantThresholdProfile(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantThresholdProfile(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantThresholdProfile(TENANT, created2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, created2.Name, "Deleted name not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllTenantThresholdProfile(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
}

func (runner *TenantServiceDatastoreTestRunner) RunGetMonitoredObjectByDomainMapTest(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const NAME1 = "name1"
	const NAME2 = "name2"
	const DOM1 = "domain1"
	const DOM2 = "domain2"
	const DOM3 = "domain3"
	const COLOR1 = "color1"
	const COLOR2 = "color2"
	const THRPRF = "ThresholdPrf"

	const OBJNAME1 = "obj1"
	const OBJID1 = "object1"
	const OBJNAME2 = "obj2"
	const OBJID2 = "object2"
	const OBJNAME3 = "obj3"
	const OBJID3 = "object3"
	const ACTNAME1 = "actName1"
	const ACTTYPE1 = string(tenmod.AccedianVNID)
	const ACTNAME2 = "actName2"
	const ACTTYPE2 = string(tenmod.AccedianNID)
	const REFNAME1 = "refname1"
	const REFTYPE1 = string(tenmod.AccedianNID)

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

	// Create a couple Domains
	// Create a record
	tenantDomain := tenmod.Domain{
		Name:     DOM1,
		TenantID: TENANT,
		Color:    COLOR1}
	created, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	tenantDomain = tenmod.Domain{
		Name:     DOM2,
		TenantID: TENANT,
		Color:    COLOR2}
	created2, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, created2)

	// Validate they were created
	recList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, recList)
	assert.Equal(t, 2, len(recList))

	// Now create some MonitoredObjects
	DOMAINSET1 := []string{created.ID}
	DOMAINSET2 := []string{created2.ID}
	DOMAINSET3 := []string{created2.ID}
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
			DomainSet:         DOMAINSET2,
		},
		&tenmod.MonitoredObject{
			MonitoredObjectID: OBJID3,
			ObjectName:        OBJNAME3,
			TenantID:          TENANT,
			ReflectorName:     REFNAME1,
			ReflectorType:     REFTYPE1,
			DomainSet:         DOMAINSET3,
		}}
	bulkResult, err := runner.tenantDB.BulkInsertMonitoredObjects(TENANT, bulkReq)
	assert.Nil(t, err)
	assert.NotEmpty(t, bulkResult)
	assert.Equal(t, 3, len(bulkResult))

	moByDomReq := tenmod.MonitoredObjectCountByDomainRequest{}

	// Fail first due to no tenant ID
	resp, err := runner.tenantDB.GetMonitoredObjectToDomainMap(&moByDomReq)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Add tenant ID
	moByDomReq.TenantID = TENANT
	resp, err = runner.tenantDB.GetMonitoredObjectToDomainMap(&moByDomReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.DomainToMonitoredObjectCountMap)
	assert.Equal(t, int64(1), int64(len(resp.DomainToMonitoredObjectSetMap[created.ID])))
	assert.Equal(t, int64(2), int64(len(resp.DomainToMonitoredObjectSetMap[created2.ID])))

	// By Count
	moByDomReq.ByCount = true
	resp, err = runner.tenantDB.GetMonitoredObjectToDomainMap(&moByDomReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.DomainToMonitoredObjectSetMap)
	assert.Equal(t, int64(1), resp.DomainToMonitoredObjectCountMap[created.ID])
	assert.Equal(t, int64(2), resp.DomainToMonitoredObjectCountMap[created2.ID])

	// Filter the list
	moByDomReq.DomainSet = []string{created.ID}
	resp, err = runner.tenantDB.GetMonitoredObjectToDomainMap(&moByDomReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.DomainToMonitoredObjectSetMap)
	assert.Equal(t, int64(1), int64(resp.DomainToMonitoredObjectCountMap[created.ID]))
	assert.Equal(t, int64(0), resp.DomainToMonitoredObjectCountMap[created2.ID])

	// Cleanup
	domList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, domList)
	for _, val := range domList {
		del, err := runner.tenantDB.DeleteTenantDomain(TENANT, val.ID)
		assert.Nil(t, err)
		assert.NotNil(t, del)
	}

	domList, err = runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.Empty(t, domList)

	moList, err := runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, moList)
	for _, val := range moList {
		del, err := runner.tenantDB.DeleteMonitoredObject(TENANT, val.ID)
		assert.Nil(t, err)
		assert.NotNil(t, del)
	}

	moList, err = runner.tenantDB.GetAllMonitoredObjects(TENANT)
	assert.Nil(t, err)
	assert.Empty(t, moList)
}

func (runner *TenantServiceDatastoreTestRunner) RunHasDashboardWithDomainTest(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const NAME1 = "name1"
	const NAME2 = "name2"
	const DOM1 = "domain1"
	const DOM2 = "domain2"
	const DOM3 = "domain3"
	const COLOR1 = "color1"
	const COLOR2 = "color2"
	const THRPRF = "ThresholdPrf"
	const DASHBOARD_NODOMAINS = "dash1"
	const DASHBOARD_DOM1 = "dash2"
	const DASHBOARD_DOM1_DOM2 = "dash3"

	const OBJNAME1 = "obj1"
	const OBJID1 = "object1"
	const OBJNAME2 = "obj2"
	const OBJID2 = "object2"
	const OBJNAME3 = "obj3"
	const OBJID3 = "object3"
	const ACTNAME1 = "actName1"
	const ACTTYPE1 = string(tenmod.AccedianVNID)
	const ACTNAME2 = "actName2"
	const ACTTYPE2 = string(tenmod.AccedianNID)
	const REFNAME1 = "refname1"
	const REFTYPE1 = string(tenmod.AccedianNID)

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

	// Create a couple Domains
	// Create a record
	tenantDomain := tenmod.Domain{
		Name:     DOM1,
		TenantID: TENANT,
		Color:    COLOR1}
	res, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	tenantDomain = tenmod.Domain{
		Name:     DOM2,
		TenantID: TENANT,
		Color:    COLOR2}
	res, err = runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	tenantDomain = tenmod.Domain{
		Name:     DOM3,
		TenantID: TENANT,
		Color:    COLOR2}
	res, err = runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	// Validate they were created
	recList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, recList)
	assert.Equal(t, 3, len(recList))

	res2, err := runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM1)
	assert.Nil(t, err)
	assert.False(t, res2)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM2)
	assert.Nil(t, err)
	assert.False(t, res2)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM3)
	assert.Nil(t, err)
	assert.False(t, res2)

	_, err = runner.tenantDB.CreateDashboard(&tenmod.Dashboard{Name: DASHBOARD_NODOMAINS, TenantID: TENANT})
	assert.Nil(t, err)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM1)
	assert.Nil(t, err)
	assert.False(t, res2)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM2)
	assert.Nil(t, err)
	assert.False(t, res2)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM3)
	assert.Nil(t, err)

	_, err = runner.tenantDB.CreateDashboard(&tenmod.Dashboard{Name: DASHBOARD_DOM1, TenantID: TENANT, DomainSet: []string{DOM1}})
	assert.Nil(t, err)

	_, err = runner.tenantDB.CreateDashboard(&tenmod.Dashboard{Name: DASHBOARD_DOM1_DOM2, TenantID: TENANT, DomainSet: []string{DOM2, DOM1}})
	assert.Nil(t, err)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM1)
	assert.Nil(t, err)
	assert.True(t, res2)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM2)
	assert.Nil(t, err)
	assert.True(t, res2)

	res2, err = runner.tenantDB.HasDashboardsWithDomain(TENANT, DOM3)
	assert.Nil(t, err)
	assert.False(t, res2)

}

func (runner *TenantServiceDatastoreTestRunner) RunMonitoredObjectGetAllInList(t *testing.T) {
	const COMPANY1 = "GetAllInListCo"
	const SUBDOMAIN1 = "adom1"
	const NAME1 = "aname1"
	const NAME2 = "aname2"
	const DOM1 = "adomain1"
	const DOM2 = "adomain2"
	const DOM3 = "adomain3"
	const COLOR1 = "acolor1"
	const COLOR2 = "acolor2"
	const THRPRF = "aThresholdPrf"

	const OBJNAME1 = "aobj1"
	const OBJID1 = "aobject1"
	const OBJNAME2 = "aobj2"
	const OBJID2 = "aobject2"
	const OBJNAME3 = "aobj3"
	const OBJID3 = "aobject3"
	const ACTNAME1 = "aactName1"
	const ACTTYPE1 = string(tenmod.AccedianVNID)
	const ACTNAME2 = "aactName2"
	const ACTTYPE2 = string(tenmod.AccedianNID)
	const REFNAME1 = "arefname1"
	const REFTYPE1 = string(tenmod.AccedianNID)

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

	// Create a couple Domains
	// Create a record
	tenantDomain := tenmod.Domain{
		Name:     DOM1,
		TenantID: TENANT,
		Color:    COLOR1}
	created, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	tenantDomain = tenmod.Domain{
		Name:     DOM2,
		TenantID: TENANT,
		Color:    COLOR2}
	created2, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, created2)

	// Validate they were created
	recList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, recList)
	assert.Equal(t, 2, len(recList))

	// Now create some MonitoredObjects
	DOMAINSET1 := []string{created.ID}
	DOMAINSET2 := []string{created2.ID}
	DOMAINSET3 := []string{created2.ID}
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
			DomainSet:         DOMAINSET2,
		},
		&tenmod.MonitoredObject{
			MonitoredObjectID: OBJID3,
			ObjectName:        OBJNAME3,
			TenantID:          TENANT,
			ReflectorName:     REFNAME1,
			ReflectorType:     REFTYPE1,
			DomainSet:         DOMAINSET3,
		}}
	bulkResult, err := runner.tenantDB.BulkInsertMonitoredObjects(TENANT, bulkReq)
	assert.Nil(t, err)
	assert.NotEmpty(t, bulkResult)
	assert.Equal(t, 3, len(bulkResult))

	// Now get a couple of the MOs using their IDs
	byIDResult, err := runner.tenantDB.GetAllMonitoredObjectsInIDList(TENANT, []string{bulkResult[0].ID, bulkResult[2].ID})
	assert.Nil(t, err)
	assert.NotEmpty(t, byIDResult)
	assert.Equal(t, 2, len(byIDResult))

	byIDResult, err = runner.tenantDB.GetAllMonitoredObjectsInIDList(TENANT, []string{})
	assert.Nil(t, err)
	assert.Empty(t, byIDResult)

}

func (runner *TenantServiceDatastoreTestRunner) RunTenantReportScheduleCRUD(t *testing.T) {
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

	reports, err := runner.tenantDB.GetAllSLAReports(TENANT)
	// Should be 0 reports in db
	assert.Equal(t, len(reports), 0)
	assert.Nil(t, err)

	request := metmod.SLAReportRequest{
		TenantID:          TENANT,
		SlaScheduleConfig: "1000",
	}

	report := metmod.SLAReport{
		TenantID:      TENANT,
		ReportSummary: metmod.ReportSummary{},
		TimeSeriesResult: []metmod.TimeSeriesEntry{
			metmod.TimeSeriesEntry{
				Timestamp: "1000",
				Result: metmod.TimeSeriesResult{
					TotalDuration:          1000,
					TotalViolationCount:    42,
					TotalViolationDuration: 9001,
				},
			},
		},
		ReportScheduleConfig: request.SlaScheduleConfig,
	}
	tdb := runner.tenantDB

	res, err := tdb.CreateSLAReport(&report)
	assert.Nil(t, err)
	report.ID = res.ID
	report.REV = res.REV
	assert.Equal(t, *res, report)

	res, err = tdb.GetSLAReport(TENANT, res.ID)
	assert.Nil(t, err)
	assert.Equal(t, *res, report)

	reports, err = runner.tenantDB.GetAllSLAReports(TENANT)
	// Should be 1 reports in db
	assert.Equal(t, len(reports), 1)
	assert.Nil(t, err)

	res, err = tdb.DeleteSLAReport(TENANT, res.ID)
	assert.Nil(t, err)

	reports, err = runner.tenantDB.GetAllSLAReports(TENANT)
	// Should be 0 reports in db
	assert.Equal(t, len(reports), 0)
	assert.Nil(t, err)

	ReportConfig := metmod.ReportScheduleConfig{
		TenantID:          TENANT,
		Name:              "Test",
		Timeout:           5000,
		Hour:              "0",
		Minute:            "0",
		DayOfWeek:         "*",
		Month:             "*",
		DayOfMonth:        "*",
		Active:            true,
		TimeRangeDuration: "P1Y",
		Granularity:       "PT1H",
	}

	configs, err := tdb.GetAllReportScheduleConfigs(TENANT)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(configs))

	cfg, err := tdb.CreateReportScheduleConfig(&ReportConfig)
	assert.Nil(t, err)
	ReportConfig.ID = cfg.ID
	ReportConfig.REV = cfg.REV
	ReportConfig.CreatedTimestamp = cfg.CreatedTimestamp
	ReportConfig.LastModifiedTimestamp = cfg.LastModifiedTimestamp
	assert.Equal(t, *cfg, ReportConfig)

	configs, err = tdb.GetAllReportScheduleConfigs(TENANT)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(configs))

	c, err := tdb.GetReportScheduleConfig(TENANT, cfg.ID)
	assert.Nil(t, err)
	assert.Equal(t, *c, ReportConfig)
	logger.Log.Debugf("ReportConfig from couch: %+v vs real %+v", c, ReportConfig)

	c, err = tdb.DeleteReportScheduleConfig(TENANT, cfg.ID)
	assert.Nil(t, err)
	assert.Equal(t, *c, ReportConfig)

	configs, err = tdb.GetAllReportScheduleConfigs(TENANT)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(configs))

	// The following test string should pass the json unmarshaller but
	// when we check for validity, it should fail. POC not to trust the
	// json.Unmarshal error code
	invalidTestString := `{
		  "data": {
		    "byDayOfWeekResult": null,
		    "byHourOfDayResult": null,
		    "createdTimestamp": 1528132980052,
		    "datatype": "tenantSLAReport",
		    "lastModifiedTimestamp": 1528132980052,
		    "reportCompletionTime": "",
		    "reportTimeRange": "",
		    "slaReportRequest": {
		      "domain": [
		        "f1f0aa6c-294c-4f50-96e6-3b6fbb7fd5f1"
		      ],
		      "granularity": "PT1H",
		      "interval": "P1D/2018-06-04T00:23:00Z",
		      "slaScheduleConfigId": "d0281f08-d64d-4425-93b5-3aa1bbd2b36f",
		      "tenantId": "ade3010a-a70a-4444-8cc7-c12c57a9ada5",
		      "thresholdProfileId": "30eec8cd-d742-4af9-9431-c944db1ce6a5",
		      "timeout": 5000,
		      "timezone": "UTC"
		    },
		    "reportSummary": {
		      "objectCount": 0,
		      "perMetricSummary": null,
		      "slaCompliancePercent": 0,
		      "totalDuration": 0,
		      "totalViolationCount": 0,
		      "totalViolationDuration": 0
		    },
		    "tenantId": "ade3010a-a70a-4444-8cc7-c12c57a9ada5",
		    "timeSeriesResult": null
		  }
		}`

	var jsReport metmod.SLAReport
	err = json.Unmarshal([]byte(invalidTestString), &jsReport)
	assert.Nil(t, err)

	stored, err := tdb.CreateSLAReport(&jsReport)
	assert.NotNil(t, err)

	assert.NotEqual(t, "ade3010a-a70a-4444-8cc7-c12c5", jsReport.TenantID)

	// The following string is a valid response after generating a new SLA Report.
	validTestString := `{
                "_id": "e29e0871-1f40-4ddb-94c8-0df9a1bdbbd2",
                "_rev": "",
                "reportCompletionTime": "2018-06-04T19:17:32Z",
                "tenantId": "ade3010a-a70a-4444-8cc7-c12c57a9ada5",
                "reportTimeRange": "P90Y/2018-06-04T00:23:00Z",
                "reportSummary": {
                    "totalDuration": 114030012,
                    "totalViolationCount": 9402,
                    "totalViolationDuration": 114030012,
                    "slaCompliancePercent": 0,
                    "objectCount": 4,
                    "perMetricSummary": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "totalDuration": 114030012,
                                            "violationCount": 4701,
                                            "violationDuration": 114030012
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "totalDuration": 114030012,
                                            "violationCount": 4701,
                                            "violationDuration": 114030012
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "timeSeriesResult": [
                    {
                        "timestamp": "2018-05-31T17:00:00.000Z",
                        "result": {
                            "totalDuration": 2280008,
                            "totalViolationCount": 144,
                            "totalViolationDuration": 2280008,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 2280008,
                                                    "violationCount": 72,
                                                    "violationDuration": 2280008
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 2280008,
                                                    "violationCount": 72,
                                                    "violationDuration": 2280008
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-05-31T18:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-05-31T19:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-05-31T20:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-05-31T21:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-05-31T22:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-05-31T23:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T00:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T01:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T02:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T03:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T04:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T05:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T06:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T07:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T08:00:00.000Z",
                        "result": {
                            "totalDuration": 7200000,
                            "totalViolationCount": 600,
                            "totalViolationDuration": 7200000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 7200000,
                                                    "violationCount": 300,
                                                    "violationDuration": 7200000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T09:00:00.000Z",
                        "result": {
                            "totalDuration": 3150000,
                            "totalViolationCount": 254,
                            "totalViolationDuration": 3150000,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 3150000,
                                                    "violationCount": 127,
                                                    "violationDuration": 3150000
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 3150000,
                                                    "violationCount": 127,
                                                    "violationDuration": 3150000
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    {
                        "timestamp": "2018-06-01T14:00:00.000Z",
                        "result": {
                            "totalDuration": 600004,
                            "totalViolationCount": 4,
                            "totalViolationDuration": 600004,
                            "perMetricResult": {
                                "accedian-twamp": {
                                    "twamp-pe": {
                                        "delayP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 600004,
                                                    "violationCount": 2,
                                                    "violationDuration": 600004
                                                }
                                            }
                                        },
                                        "jitterP95": {
                                            "sla": {
                                                "0": {
                                                    "totalDuration": 600004,
                                                    "violationCount": 2,
                                                    "violationDuration": 600004
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                ],
                "byHourOfDayResult": {
                    "14": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 2
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 2
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "17": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 72
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 72
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "18": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "19": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "20": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "21": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "22": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "23": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "00": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "01": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "02": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "03": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "04": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "05": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "06": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "07": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "08": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 300
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "09": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 127
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 127
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "byDayOfWeekResult": {
                    "4": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 1872
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 1872
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "5": {
                        "accedian-twamp": {
                            "twamp-pe": {
                                "delayP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 2829
                                        }
                                    }
                                },
                                "jitterP95": {
                                    "sla": {
                                        "0": {
                                            "violationCount": 2829
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "reportScheduleConfig": "1000"
			}`

	err = json.Unmarshal([]byte(validTestString), &jsReport)
	assert.Nil(t, err)
	jsReport.ID = ""
	jsReport.REV = ""
	assert.Equal(t, "ade3010a-a70a-4444-8cc7-c12c57a9ada5", jsReport.TenantID)
	assert.Equal(t, "P90Y/2018-06-04T00:23:00Z", jsReport.ReportTimeRange)
	assert.Equal(t, int64(114030012), jsReport.ReportSummary.TotalDuration)
	assert.Equal(t, int32(9402), jsReport.ReportSummary.TotalViolationCount)
	assert.Equal(t, int64(114030012), jsReport.ReportSummary.TotalViolationDuration)
	assert.Equal(t, float32(0), jsReport.ReportSummary.SLACompliancePercent)
	jsReport.TenantID = TENANT

	stored, err = tdb.CreateSLAReport(&jsReport)
	assert.Nil(t, err)
	logger.Log.Debugf("Stored: %+v", stored)

	sRep, err := tdb.GetSLAReport(stored.TenantID, stored.ID)
	assert.Nil(t, err)
	assert.NotNil(t, sRep)
	jsReport.ID = stored.ID
	jsReport.REV = stored.REV

	assert.Equal(t, jsReport, *sRep)
	assert.Equal(t, TENANT, sRep.TenantID)
	assert.Equal(t, "P90Y/2018-06-04T00:23:00Z", sRep.ReportTimeRange)
	assert.Equal(t, int64(114030012), sRep.ReportSummary.TotalDuration)
	assert.Equal(t, int32(9402), sRep.ReportSummary.TotalViolationCount)
	assert.Equal(t, int64(114030012), sRep.ReportSummary.TotalViolationDuration)
	assert.Equal(t, float32(0), sRep.ReportSummary.SLACompliancePercent)

}

func (runner *TenantServiceDatastoreTestRunner) RunTenantDataCleaningProfileCRUD(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const THRPRF = "ThresholdPrf"
	const THRPRF2 = "ThresholdPrf2"

	RULE1 := &tenmod.DataCleaningRule{
		MetricLabel:  "SomeLabel1",
		MetricVendor: "SomeVendor1",
		TriggerCondition: &tenmod.DataCleaningRuleCondition{
			Comparator:     "trigcomp1",
			Duration:       "trigdur1",
			Value:          "trigvalue1",
			ValueAggregate: "trigagg1",
		},
		ClearCondition: &tenmod.DataCleaningRuleCondition{
			Comparator:     "clearcomp1",
			Duration:       "cleardur1",
			Value:          "clearvalue1",
			ValueAggregate: "clearagg1",
		},
	}
	RULE2 := &tenmod.DataCleaningRule{
		MetricLabel:  "SomeLabel2",
		MetricVendor: "SomeVendor2",
		TriggerCondition: &tenmod.DataCleaningRuleCondition{
			Comparator:     "trigcomp2",
			Duration:       "trigdur2",
			Value:          "trigvalue2",
			ValueAggregate: "trigagg2",
		},
		ClearCondition: &tenmod.DataCleaningRuleCondition{
			Comparator:     "clearcomp2",
			Duration:       "cleardur2",
			Value:          "clearvalue2",
			ValueAggregate: "clearagg2",
		},
	}
	RULES1 := []*tenmod.DataCleaningRule{RULE1, RULE2}

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
	record, err := runner.tenantDB.GetAllTenantDataCleaningProfiles(TENANT)
	assert.NotNil(t, err)
	assert.Empty(t, record)

	// Try to Update a record that does not exist:
	fail, err := runner.tenantDB.UpdateTenantDataCleaningProfile(&tenmod.DataCleaningProfile{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	dcp := tenmod.DataCleaningProfile{
		TenantID: TENANT,
		Rules:    RULES1,
	}
	created, err := runner.tenantDB.CreateTenantDataCleaningProfile(&dcp)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantDataCleaningProfileType), created.Datatype)
	assert.Equal(t, 2, len(created.Rules), "Not the correct number of rules")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantDataCleaningProfile(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, string(tenmod.TenantDataCleaningProfileType), fetched.Datatype)
	assert.Equal(t, 2, len(fetched.Rules), "Not the correct number of rules")
	assert.Equal(t, TENANT, fetched.TenantID, "Tenant ID not the same")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.DataCleaningProfile{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Rules = []*tenmod.DataCleaningRule{}
	updated, err := runner.tenantDB.UpdateTenantDataCleaningProfile(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantDataCleaningProfileType), updated.Datatype)
	assert.Equal(t, 0, len(updated.Rules), "Should not have any rules")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	dcp2 := tenmod.DataCleaningProfile{
		TenantID: TENANT}
	created2, err := runner.tenantDB.CreateTenantDataCleaningProfile(&dcp2)
	assert.NotNil(t, err)
	assert.Nil(t, created2)

	// Try the get all
	allRecords, err := runner.tenantDB.GetAllTenantDataCleaningProfiles(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, allRecords)
	assert.Equal(t, 1, len(allRecords), "Should have found 1 record")

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantDataCleaningProfile(TENANT, updated.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.Equal(t, 0, len(deleted.Rules), "Deleted should not have any rules")

	// Delete a record that does not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantDataCleaningProfile(TENANT, updated.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	allRecords, err = runner.tenantDB.GetAllTenantDataCleaningProfiles(TENANT)
	assert.NotNil(t, err)
	assert.Nil(t, allRecords)
}
