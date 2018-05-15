package druid

import (
	"strings"
	"time"
)

// A simple cache to mirror what is stored on Druid.
// It helps when building queries and need to supply lookup names that should be fully active on druid.
// The lookup cache is updated when pushing lookups to Druid but there is some delay before they are fully
// propagated and active Druid so each lookup is timestamped in the cache in order to identify those
// that should be active.
var lookups = lookupCache{
	lookupNames: nil,
}

const (
	druidLookupWriteDelay time.Duration = 15 * time.Second // This should really be coordinated with Druid's config
)

type lookupCache struct {
	lookupNames map[string]*time.Time
}

// Generates a lookup name from the parameters and returns it only if it is considered to be valid and active
func getLookupName(dimType, tenantID, dimValue string) (string, bool) {
	lookup := buildLookupName(dimType, tenantID, dimValue)

	if lookups.lookupNames == nil {
		// The cache hasn't been initialized. Assume the lookup name is
		// valid.
		return lookup, true
	}
	// Check the cache.  The lookup is valid to use if it's present in
	// the cache and it has been committed for at least the write delay.
	if ts, ok := lookups.lookupNames[lookup]; ok {
		if ts.Add(druidLookupWriteDelay).Before(time.Now()) {
			// Enough time has passed to commit the lookup
			return lookup, true
		}
	}

	return "", false
}

// Construct a lookup name
func buildLookupName(dimType, tenantID, dimValue string) string {
	return strings.ToLower(dimType + "|" + tenantID + "|" + dimValue)
}

func buildLookupNamePrefix(dimType, tenantID string) string {
	return strings.ToLower(dimType + "|" + tenantID)
}

func getLookupNamePrefix(dimType, tenantID string) string {
	return strings.ToLower(dimType + "|" + tenantID)
}

func refreshLookups(lookupMap map[string]lookup) {

	// Update the lookup cache.  Map each lookup name to the earliest timestamp it was created.
	curTs := time.Now()
	var earliestTs *time.Time
	newLookupNames := make(map[string]*time.Time, len(lookupMap))
	for k := range lookupMap {
		earliestTs = &curTs
		if lookups.lookupNames == nil {
			// This is the initial loading of the cache so just use epoch and assume these
			// lookups have existed for some time and are ready to use.
			earliestTs = &time.Time{}
		} else if prevTs, ok := lookups.lookupNames[k]; ok {
			// The lookup already exists so use its existing timestamp.
			earliestTs = prevTs
		}
		newLookupNames[k] = earliestTs
	}

	lookups.lookupNames = newLookupNames

}
