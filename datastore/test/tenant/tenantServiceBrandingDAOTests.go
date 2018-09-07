package tenant

import (
	"testing"

	ds "github.com/accedian/adh-gather/datastore"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/getlantern/deepcopy"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func (runner *TenantServiceDatastoreTestRunner) RunBrandingCRUD(t *testing.T) {
	const COMPANY1 = "BrandingCompany"
	const SUBDOMAIN1 = "subdom1"
	const NAME1 = "name1"
	const NAME2 = "name2"
	const THRPRF1 = "ThresholdPrf1"

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

	// Make sure there are not records to start
	fetchList, err := runner.tenantDB.GetAllTenantBrandings(TENANT)
	assert.NotNil(t, err)
	assert.Empty(t, fetchList)

	// Create a record
	rec := tenmod.Branding{
		Color: fake.CharactersN(6),
		Logo: &tenmod.BrandingLogo{
			File: &tenmod.BrandingLogoFile{
				ContentType: fake.CharactersN(12),
				Data:        fake.CharactersN(100),
			},
		},
		TenantID: TENANT,
	}
	created, err := runner.tenantDB.CreateTenantBranding(&rec)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, TENANT, created.TenantID)
	assert.Equal(t, rec.Color, created.Color)
	assert.Equal(t, rec.Logo, created.Logo)
	assert.Equal(t, string(tenmod.TenantBrandingType), created.Datatype)
	assert.True(t, created.CreatedTimestamp > 0)
	assert.True(t, created.CreatedTimestamp == created.LastModifiedTimestamp)

	// Fetch by ID - unknown ID should fail
	fetched, err := runner.tenantDB.GetTenantBranding(TENANT, "notreal")
	assert.NotNil(t, err)

	// Fetch by ID - success
	fetched, err = runner.tenantDB.GetTenantBranding(TENANT, created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, created.REV, fetched.REV)
	assert.Equal(t, created, fetched)

	// Update a record that does not exist
	dne := tenmod.Branding{
		ID:    fake.CharactersN(12),
		REV:   fake.CharactersN(15),
		Color: fake.CharactersN(6),
		Logo: &tenmod.BrandingLogo{
			File: &tenmod.BrandingLogoFile{
				ContentType: fake.CharactersN(12),
				Data:        fake.CharactersN(100),
			},
		},
		TenantID: TENANT,
	}
	_, err = runner.tenantDB.UpdateTenantBranding(&dne)
	assert.NotNil(t, err)

	// Update that does work:
	updateRecord := tenmod.Branding{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Color = fake.CharactersN(8)
	updated, err := runner.tenantDB.UpdateTenantBranding(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.NotEqual(t, fetched.REV, updated.REV)
	assert.NotEqual(t, fetched.Color, updated.Color)
	assert.Equal(t, updateRecord.Color, updated.Color)
	assert.Equal(t, fetched.CreatedTimestamp, updated.CreatedTimestamp)

	// Create additional records
	numAdditionalRecords := 4
	for i := 0; i < numAdditionalRecords; i++ {
		newOne := tenmod.Branding{
			Color: fake.CharactersN(6),
			Logo: &tenmod.BrandingLogo{
				File: &tenmod.BrandingLogoFile{
					ContentType: fake.CharactersN(12),
					Data:        fake.CharactersN(100),
				},
			},
			TenantID: TENANT,
		}
		createdNew, err := runner.tenantDB.CreateTenantBranding(&newOne)
		assert.Nil(t, err)
		assert.NotNil(t, createdNew)
		assert.NotEmpty(t, createdNew.ID)
		assert.NotEmpty(t, createdNew.REV)
		assert.Equal(t, newOne.Color, createdNew.Color)
		assert.Equal(t, TENANT, createdNew.TenantID)
		assert.Equal(t, newOne.Logo, createdNew.Logo)
		assert.Equal(t, string(tenmod.TenantBrandingType), createdNew.Datatype)
		assert.True(t, createdNew.CreatedTimestamp > 0)
		assert.True(t, createdNew.CreatedTimestamp == createdNew.LastModifiedTimestamp)
	}

	// Get All records
	fetchList, err = runner.tenantDB.GetAllTenantBrandings(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchList)
	assert.Equal(t, numAdditionalRecords+1, len(fetchList))

	// Delete a record that DNE
	deleted, err := runner.tenantDB.DeleteTenantBranding(TENANT, "notReal")
	assert.NotNil(t, err)
	assert.Nil(t, deleted)

	// Delete a record that does exist
	deleted, err = runner.tenantDB.DeleteTenantBranding(TENANT, updated.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.Equal(t, updated, deleted)

	// Get All records - should have 1 less
	fetchList, err = runner.tenantDB.GetAllTenantBrandings(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchList)
	assert.Equal(t, numAdditionalRecords, len(fetchList))

	// Delete the remaining records
	for _, rec := range fetchList {
		gone, err := runner.tenantDB.DeleteTenantBranding(TENANT, rec.ID)
		assert.Nil(t, err)
		assert.NotNil(t, gone)
	}

	// Make sure there are not records to end
	fetchList, err = runner.tenantDB.GetAllTenantBrandings(TENANT)
	assert.NotNil(t, err)
	assert.Empty(t, fetchList)
}

func (runner *TenantServiceDatastoreTestRunner) RunLocaleCRUD(t *testing.T) {
	const COMPANY1 = "LocaleCompany"
	const SUBDOMAIN1 = "subdom1"
	const NAME1 = "name1"
	const NAME2 = "name2"
	const THRPRF1 = "ThresholdPrf1"

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

	// Make sure there are not records to start
	fetchList, err := runner.tenantDB.GetAllTenantLocales(TENANT)
	assert.NotNil(t, err)
	assert.Empty(t, fetchList)

	// Create a record
	rec := tenmod.Locale{
		Intl:     fake.CharactersN(6),
		Moment:   fake.CharactersN(2),
		Timezone: fake.CharactersN(15),
		TenantID: TENANT,
	}
	created, err := runner.tenantDB.CreateTenantLocale(&rec)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, TENANT, created.TenantID)
	assert.Equal(t, rec.Intl, created.Intl)
	assert.Equal(t, rec.Moment, created.Moment)
	assert.Equal(t, rec.Timezone, created.Timezone)
	assert.Equal(t, string(tenmod.TenantLocaleType), created.Datatype)
	assert.True(t, created.CreatedTimestamp > 0)
	assert.True(t, created.CreatedTimestamp == created.LastModifiedTimestamp)

	// Fetch by ID - unknown ID should fail
	fetched, err := runner.tenantDB.GetTenantLocale(TENANT, "notreal")
	assert.NotNil(t, err)

	// Fetch by ID - success
	fetched, err = runner.tenantDB.GetTenantLocale(TENANT, created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, created.REV, fetched.REV)
	assert.Equal(t, created, fetched)

	// Update a record that does not exist
	dne := tenmod.Locale{
		ID:       fake.CharactersN(12),
		REV:      fake.CharactersN(15),
		Intl:     fake.CharactersN(6),
		Moment:   fake.CharactersN(2),
		Timezone: fake.CharactersN(15),
		TenantID: TENANT,
	}
	_, err = runner.tenantDB.UpdateTenantLocale(&dne)
	assert.NotNil(t, err)

	// Update that does work:
	updateRecord := tenmod.Locale{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Moment = fake.CharactersN(8)
	updated, err := runner.tenantDB.UpdateTenantLocale(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.NotEqual(t, fetched.REV, updated.REV)
	assert.NotEqual(t, fetched.Moment, updated.Moment)
	assert.Equal(t, updateRecord.Moment, updated.Moment)
	assert.Equal(t, fetched.CreatedTimestamp, updated.CreatedTimestamp)

	// Create additional records
	numAdditionalRecords := 4
	for i := 0; i < numAdditionalRecords; i++ {
		newOne := tenmod.Locale{
			Intl:     fake.CharactersN(6),
			Moment:   fake.CharactersN(2),
			Timezone: fake.CharactersN(15),
			TenantID: TENANT,
		}
		createdNew, err := runner.tenantDB.CreateTenantLocale(&newOne)
		assert.Nil(t, err)
		assert.NotNil(t, createdNew)
		assert.NotEmpty(t, createdNew.ID)
		assert.NotEmpty(t, createdNew.REV)
		assert.Equal(t, TENANT, createdNew.TenantID)
		assert.Equal(t, rec.Intl, created.Intl)
		assert.Equal(t, rec.Moment, created.Moment)
		assert.Equal(t, rec.Timezone, created.Timezone)
		assert.Equal(t, string(tenmod.TenantLocaleType), createdNew.Datatype)
		assert.True(t, createdNew.CreatedTimestamp > 0)
		assert.True(t, createdNew.CreatedTimestamp == createdNew.LastModifiedTimestamp)
	}

	// Get All records
	fetchList, err = runner.tenantDB.GetAllTenantLocales(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchList)
	assert.Equal(t, numAdditionalRecords+1, len(fetchList))

	// Delete a record that DNE
	deleted, err := runner.tenantDB.DeleteTenantLocale(TENANT, "notReal")
	assert.NotNil(t, err)
	assert.Nil(t, deleted)

	// Delete a record that does exist
	deleted, err = runner.tenantDB.DeleteTenantLocale(TENANT, updated.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.Equal(t, updated, deleted)

	// Get All records - should have 1 less
	fetchList, err = runner.tenantDB.GetAllTenantLocales(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchList)
	assert.Equal(t, numAdditionalRecords, len(fetchList))

	// Delete the remaining records
	for _, rec := range fetchList {
		gone, err := runner.tenantDB.DeleteTenantLocale(TENANT, rec.ID)
		assert.Nil(t, err)
		assert.NotNil(t, gone)
	}

	// Make sure there are not records to end
	fetchList, err = runner.tenantDB.GetAllTenantLocales(TENANT)
	assert.NotNil(t, err)
	assert.Empty(t, fetchList)
}
