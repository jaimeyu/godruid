package tenant

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models/common"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	testUtil "github.com/accedian/adh-gather/models/test"
	uuid "github.com/satori/go.uuid"
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

func TestTenantUserSerialization(t *testing.T) {
	original := &User{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantUserType),
		TenantID:              fake.CharactersN(12),
		Username:              fake.EmailAddress(),
		Password:              fake.Password(6, 8, true, true, false),
		State:                 string(common.UserActive),
		SendOnboardingEmail:   true,
		OnboardingToken:       fake.CharactersN(14),
		UserVerified:          false,
		Domains:               []string{fake.CharactersN(5), fake.CharactersN(7)},
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "username", "password", "sendOnboardingEmail", "state", "onboardingToken", "userVerified", "domains", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &User{}, original.ID, attrKeys)
}

func TestTenantUserValidation(t *testing.T) {
	original := &User{}

	// Must have TenantID
	err := original.Validate(false)
	assert.NotNil(t, err)

	original.TenantID = "something"
	err = original.Validate(false)
	assert.Nil(t, err)

	// Must not provide REV if it is not an Update
	original.REV = "oops"
	err = original.Validate(false)
	assert.NotNil(t, err)

	original.REV = ""
	err = original.Validate(false)
	assert.Nil(t, err)

	// If it is an update, make sure REV and CreatedTimestamp are there
	original.CreatedTimestamp = 0
	err = original.Validate(true)
	assert.NotNil(t, err)

	original.REV = "update"
	original.CreatedTimestamp = 123
	err = original.Validate(true)
	assert.Nil(t, err)
}

func TestTenantDomainSerialization(t *testing.T) {
	original := &Domain{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantDomainType),
		TenantID:              fake.CharactersN(12),
		Name:                  fake.State(),
		Color:                 fake.HexColor(),
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "name", "color", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &Domain{}, original.ID, attrKeys)
}

func TestTenantDomainValidation(t *testing.T) {
	original := &Domain{}

	// Must have TenantID
	err := original.Validate(false)
	assert.NotNil(t, err)

	original.TenantID = "something"
	err = original.Validate(false)
	assert.Nil(t, err)

	// Must not provide REV if it is not an Update
	original.REV = "oops"
	err = original.Validate(false)
	assert.NotNil(t, err)

	original.REV = ""
	err = original.Validate(false)
	assert.Nil(t, err)

	// If it is an update, make sure REV and CreatedTimestamp are there
	original.CreatedTimestamp = 0
	err = original.Validate(true)
	assert.NotNil(t, err)

	original.REV = "update"
	original.CreatedTimestamp = 123
	err = original.Validate(true)
	assert.Nil(t, err)
}

func TestTenantIngestionProfileSerialization(t *testing.T) {

	defaultIngestionProfileShell := &IngestionProfile{}
	if err := json.Unmarshal(defaultIngestionProfileBytes, &defaultIngestionProfileShell); err != nil {
		logger.Log.Debugf("Unable to umarshal ingestion profile: %s", err.Error())
	}

	original := &IngestionProfile{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantIngestionProfileType),
		TenantID:              fake.CharactersN(12),
		Metrics:               defaultIngestionProfileShell.Metrics,
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "metrics", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &IngestionProfile{}, original.ID, attrKeys)
}

func TestTenantIngestionProfileValidation(t *testing.T) {
	original := &IngestionProfile{}

	// Must have TenantID
	err := original.Validate(false)
	assert.NotNil(t, err)

	original.TenantID = "something"
	err = original.Validate(false)
	assert.Nil(t, err)

	// Must not provide REV if it is not an Update
	original.REV = "oops"
	err = original.Validate(false)
	assert.NotNil(t, err)

	original.REV = ""
	err = original.Validate(false)
	assert.Nil(t, err)

	// If it is an update, make sure REV and CreatedTimestamp are there
	original.CreatedTimestamp = 0
	err = original.Validate(true)
	assert.NotNil(t, err)

	original.REV = "update"
	original.CreatedTimestamp = 123
	err = original.Validate(true)
	assert.Nil(t, err)
}

func TestTenantThresholdProfileSerialization(t *testing.T) {

	defaultThresholdProfileShell := &ThresholdProfile{}
	if err := json.Unmarshal(defaultThresholdsBytes, &defaultThresholdProfileShell); err != nil {
		logger.Log.Debugf("Unable to umarshal threshold profile: %s", err.Error())
	}

	original := &ThresholdProfile{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantThresholdProfileType),
		TenantID:              fake.CharactersN(12),
		Thresholds:            defaultThresholdProfileShell.Thresholds,
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "thresholds", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &ThresholdProfile{}, original.ID, attrKeys)
}

func TestTenantThresholdProfileValidation(t *testing.T) {
	original := &ThresholdProfile{}

	// Must have TenantID
	err := original.Validate(false)
	assert.NotNil(t, err)

	original.TenantID = "something"
	err = original.Validate(false)
	assert.Nil(t, err)

	// Must not provide REV if it is not an Update
	original.REV = "oops"
	err = original.Validate(false)
	assert.NotNil(t, err)

	original.REV = ""
	err = original.Validate(false)
	assert.Nil(t, err)

	// If it is an update, make sure REV and CreatedTimestamp are there
	original.CreatedTimestamp = 0
	err = original.Validate(true)
	assert.NotNil(t, err)

	original.REV = "update"
	original.CreatedTimestamp = 123
	err = original.Validate(true)
	assert.Nil(t, err)
}

func TestTenantMonitoredObjectSerialization(t *testing.T) {
	actName := fake.City()
	refName := fake.City()
	original := &MonitoredObject{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantMonitoredObjectType),
		TenantID:              fake.CharactersN(12),
		ActuatorType:          fake.Company(),
		ActuatorName:          actName,
		ReflectorType:         fake.Company(),
		ReflectorName:         refName,
		ObjectType:            fake.Brand(),
		ObjectName:            fake.City(),
		MonitoredObjectID:     strings.Join([]string{actName, refName}, "-"),
		DomainSet:             []string{uuid.NewV4().String(), uuid.NewV4().String()},
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "actuatorName", "actuatorType",
		"reflectorName", "reflectorType", "objectName", "objectType", "domainSet",
		"objectId", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &MonitoredObject{}, original.ID, attrKeys)
}

func TestTenantMonitoredObjectValidation(t *testing.T) {
	original := &MonitoredObject{MonitoredObjectID: "something"}

	// Must have TenantID
	err := original.Validate(false)
	assert.NotNil(t, err)

	original.TenantID = "something"
	err = original.Validate(false)
	assert.Nil(t, err)

	// Must have MonitoredObjectID
	original.MonitoredObjectID = ""
	err = original.Validate(false)
	assert.NotNil(t, err)

	original.MonitoredObjectID = "newone"
	err = original.Validate(false)
	assert.Nil(t, err)

	// Must not provide REV if it is not an Update
	original.REV = "oops"
	err = original.Validate(false)
	assert.NotNil(t, err)

	original.REV = ""
	err = original.Validate(false)
	assert.Nil(t, err)

	// If it is an update, make sure REV and CreatedTimestamp are there
	original.CreatedTimestamp = 0
	err = original.Validate(true)
	assert.NotNil(t, err)

	original.REV = "update"
	original.CreatedTimestamp = 123
	err = original.Validate(true)
	assert.Nil(t, err)
}

func TestTenantMetadataSerialization(t *testing.T) {
	original := &Metadata{
		ID:                    uuid.NewV4().String(),
		REV:                   uuid.NewV4().String(),
		Datatype:              string(TenantMetaType),
		TenantID:              fake.CharactersN(12),
		TenantName:            fake.Company(),
		CreatedTimestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		LastModifiedTimestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}

	attrKeys := []string{"_rev", "datatype", "tenantId", "tenantName", "createdTimestamp", "lastModifiedTimestamp"}

	testUtil.RunSerializationTest(t, original, &Metadata{}, original.ID, attrKeys)
}

func TestTenantMetadataValidation(t *testing.T) {
	original := &Metadata{}

	// Must have TenantID
	err := original.Validate(false)
	assert.NotNil(t, err)

	original.TenantID = "something"
	err = original.Validate(false)
	assert.Nil(t, err)

	// Must not provide REV if it is not an Update
	original.REV = "oops"
	err = original.Validate(false)
	assert.NotNil(t, err)

	original.REV = ""
	err = original.Validate(false)
	assert.Nil(t, err)

	// If it is an update, make sure REV and CreatedTimestamp are there
	original.CreatedTimestamp = 0
	err = original.Validate(true)
	assert.NotNil(t, err)

	original.REV = "update"
	original.CreatedTimestamp = 123
	err = original.Validate(true)
	assert.Nil(t, err)
}
