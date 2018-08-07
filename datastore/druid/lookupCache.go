package druid

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/accedian/adh-gather/logger"
)

// A simple cache to mirror what is stored on Druid.
// It helps when building queries and need to supply lookup names that should be fully active on druid.
// The lookup cache is updated when pushing lookups to Druid but there is some delay before they are fully
// propagated and active Druid so each lookup is timestamped in the cache in order to identify those
// that should be active. Note, by 'active' we mean created and known to druid nodes.  It may have some pending
// updates but that's ok we are just trying to avoid getting 'lookup not found' errors. We allow queries
// on out-of-date lookups as we accept eventual consistency in Druid queries
var lookups = lookupCache{
	lookupNames: nil,
}

const (
	druidLookupWriteDelay time.Duration = 15 * time.Second // This should really be coordinated with Druid's config
	druidLookupSeparator                = "*"
)

// Lookup cache with lazy initialization.  The cache is not guaranteed to be initialized immediately because
// on server startup the DB may not be ready.
type lookupCache struct {
	lookupNames map[string]*time.Time
}

// Construct a lookup name
// Looks up have the format: type*tenantID*key*value*partition
func buildLookupName(dimType, tenantID, dimKey string, dimValue string, dimPartition int) string {
	name := strings.ToLower(druidLookupSeparator + dimType + druidLookupSeparator + tenantID + druidLookupSeparator + dimKey + druidLookupSeparator + dimValue + druidLookupSeparator + fmt.Sprintf("%d", dimPartition))
	hash := sha256.Sum256([]byte(name))
	return fmt.Sprintf("%X%s", hash, name)
}

func buildLookupNamePrefix(dimType, tenantID string) string {
	return strings.ToLower(dimType + druidLookupSeparator + tenantID)
}

/* Heads up. The following is so when a getLookup is executed, it only allows it
 * after about 15 seconds since the monitored object was updated.
 * This is so Druid has some time to process and update the LookUp view.
 * The logic has changed a bit, I found a condition where a caller was sending empty
 * lookup maps which wrecks havok to the lookups global and nulls it out, hence the new
 * length checks and lazy memory assignment.
 * The added debugging is important because it will tell you what the logic thinks its doing
 * when getLookUpNames is called and returns an error and you need to work backwards.
 */
func updateLookupCache(lookupMap map[string]*lookup) {

	// Don't work on map if its empty
	if lookupMap == nil {
		logger.Log.Infof("Lookup map is empty, skipping cache update")
		return
	}

	if len(lookupMap) == 0 {
		logger.Log.Infof("Lookup map is empty, skipping cache update")
		return
	}

	// Update the lookup cache.  Map each lookup name to the earliest timestamp it was created.
	curTs := time.Now()
	epoch := time.Unix(0, 0)
	var earliestTs *time.Time
	// Lazy build the look up cache
	if lookups.lookupNames == nil {
		lookups.lookupNames = make(map[string]*time.Time)
	}
	for k, v := range lookupMap {
		earliestTs = &curTs
		if lookups.lookupNames == nil {
			// This is the initial loading of the cache so just use epoch for active lookups
			if v.active {
				earliestTs = &epoch
			}
		} else if prevTs, ok := lookups.lookupNames[k]; ok {
			// The lookup already exists so use its existing timestamp.
			earliestTs = prevTs
		}
		lookups.lookupNames[k] = earliestTs
		logger.Log.Debugf("Updated lookup [%s] -> %+v", k, earliestTs)
	}

	logger.Log.Debugf("Overwrote lookup: %+v", lookups.lookupNames)

}
