package couchDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

const (
	tenantIDByNameIndex               = "_design/tenant/_view/byAlias"
	monitoredObjectCountByDomainIndex = "_design/monitoredObjectCount"

	monitoredObjectDBSuffix = "_monitored-objects"
	reportObjectDBSuffix    = "_reports"

	// Indexer for monitored objects by name
	monitoredObjectsByObjectNameIndex = "byObjectName"
	monitoredObjectsByObjectNameKey   = "objectName"
	monitoredObjectIndex              = "indexOfObjectName"

	monitoredObjectCountIndexBytes = `{
	"_id": "_design/monitoredObjectCount",
	"language": "javascript",
	"views": {
	  "byDomain": {
		"map": "function(doc) {\n    if (doc.data && doc.data.datatype && doc.data.datatype === 'monitoredObject' && doc.data.domainSet) {\n      for (var i in doc.data.domainSet) {\n        emit(doc.data.domainSet[i], doc._id);\n      }\n    }\n}"
	  },
	  "count": {
		"map": "function(doc) {  if (doc.data && doc.data.datatype && doc.data.datatype === 'monitoredObject') { emit(doc.id, 1) } }",
		"reduce": "_count"
	  }
	}
  }`

	monitoredObjectMetaIndexBytes = `{
	"_id": "_design/monitoredObject-meta",
	"language": "query",
	"views": {
	  "objectName": {
		"map": {
		  "fields": {
			"data.objectName": "asc"
		  },
		  "partial_filter_selector": {}
		},
		"reduce": "_count",
		"options": {
		  "def": {
			"fields": [
			  "data.objectName"
			]
		  }
		}
	  }
	}
  }`

	keyViewName = "%sView"
	keyViewFn   = `function (doc) {
			if (doc.data.meta["%s"]) {
				emit(doc.data.datatype, doc.id)
			}
		}`

	mapFnName = "map"

	metaFieldPrefix              = "meta"
	metakeysViewDdocName         = "uniqueMeta"
	metakeysViewUniqueKeysURI    = "uniqueMeta/uniquesKeys"
	metakeysViewUniqueValuessURI = "uniqueMeta/uniqueValues"
	metaKeyName                  = "{{KeyName}}"
	metaKeyField                 = "{{KeyField}}"
	uniqueMetaIndexBytes         = `{
		"_id": "_design/uniqueMeta",
		"language": "javascript",
		"views": {
			"uniqueKeys": {
				"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n          emit(key, doc.data.meta[key]);\n      }\n    }\n}",
				"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    // All the reduce functions have completed and returned an aggregate of collections\n    if (rereduce) {\n        var resp = [];\n        resp = values;\n        // We need to flatten the collection because each reduce function returns an array.\n        // If we left the structure as is, each row is filtered as a whole, rather than filtering on each individual item in the row.\n        resp = flatten(resp);\n        \n        // Filter the collection for only the unique items.\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        // Now get the count of unique values in this collections\n        var cnt = 0;\n        for( var v in filteredArray) {\n            cnt++;\n            }\n        return cnt;\n\n    } else {\n     return values;\n    }\n\n}\n"
			},
			"allValuesPerKey": {
				"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        emit(key, doc.data.meta[key]);\n        \n      }\n    }\n}",
				"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    if (rereduce) {\n        var resp = [];\n        resp = values;//.flat();\n        resp = flatten(resp);\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        return filteredArray.slice(0,1000);\n    } else {\n     return values;\n    }\n\n}\n"
			},
			"allValuesReduce": {
				"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        emit(doc.data.meta[key], key);\n        \n      }\n    }\n}",
				"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    if (rereduce) {\n        var resp = [];\n        resp = values;//.flat();\n        resp = flatten(resp);\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        return filteredArray.slice(0,1000);\n    } else {\n     return values;\n    }\n\n}\n"
			},
			"searchLookup": {
				"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        emit([doc.data.meta[key], key], key);\n        \n      }\n    }\n}",
				"reduce": "function (keys, values, rereduce) {\n  if (rereduce) {\n    return sum(values);\n  } else {\n    return values.length;\n  }\n}"
			},
			"wordSearch": {
				"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n        var sentence = doc.data.meta[key];\n        \n        // split on space\n        var words = sentence.split(\" \");\n        \n        for (var word in words) {\n            emit(words[word], {\"category\":key, \"sentence\":sentence});\n        }\n      }\n    }\n}"
			}
		}
	}`
	metaUniqueValuesViewsURI          = "{{KeyName}}/{{KeyField}}"
	metaUniqueValuesViewsDdocTemplate = `{"_id": "_design/viewOf{{KeyName}}","views": {"by{{KeyName}}": {"reduce": "function(keys, values) {return sum(values);}","map": "function(doc) {if (doc.data.{{KeyField}}) {emit(doc.data.{{KeyField}}, 1);}}"}},"language": "javascript"}`

	metaDdocTemplate  = "%s"
	metaIndexTemplate = `{
		"_id": "_design/indexOf{{KeyName}}",
		"language": "query",
		"views": {
			"by{{KeyName}}": {
				"map": {
					"fields": {
						"data.{{KeyField}}": "asc"
					},
					"partial_filter_selector": {}
				},
				"options": {
					"def": {
						"fields": [
							"data.{{KeyField}}"
						]
					}
				}
			}
		}
	}`
)

func getTenantViews() []map[string]interface{} {
	uniqueMetaIndexObject := map[string]interface{}{}
	monitoredObjectMetaIndexObject := map[string]interface{}{}
	monitoredObjectCountIndexObject := map[string]interface{}{}

	if err := json.Unmarshal([]byte(uniqueMetaIndexBytes), &uniqueMetaIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate Unique Meta Index: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(monitoredObjectMetaIndexBytes), &monitoredObjectMetaIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate Unique Meta Index: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(monitoredObjectCountIndexBytes), &monitoredObjectCountIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate Unique Meta Index: %s", err.Error())
	}

	return []map[string]interface{}{uniqueMetaIndexObject, monitoredObjectMetaIndexObject, monitoredObjectCountIndexObject}
}

/*
 This function takes a key and then creates an index for it and then start the indexer.
 We currently only support generating an index based on a singular key.
*/
func createCouchDBViewIndex(dbName string, template string, ddocName string, keyNames []string, prefix string) error {

	if len(keyNames) == 0 {
		return errors.New("keyNames cannot be 0")
	}

	var item string
	if len(prefix) == 0 {
		item = keyNames[0]
	} else {
		item = fmt.Sprintf("%s.%s", prefix, keyNames[0])
	}
	ckey := keyNames[0]

	var docret tenmod.MonitoredObjectMetaDesignDocument
	//var document = fmt.Sprintf(metaIndexTemplate, ckey, ckey, item, item)
	document := strings.Replace(template, metaKeyName, ckey, -1)
	document = strings.Replace(document, metaKeyField, item, -1)

	//ddocName := fmt.Sprintf(metaIndexDdocTemplate, ckey)
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Creating new Index for key '%s' file with payload:%s", keyNames[0], models.AsJSONString(document))
	}
	err := updateCouchDBDocWithStringDoc(dbName, document, string(tenmod.TenantMetaType), tenmod.TenantMetaStr, docret)

	if err != nil {
		logger.Log.Errorf("Error creating index design document %s: %s :%s\n", tenmod.TenantMetaStr, models.AsJSONString(document), err.Error())

		return err
	}
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Successfully created Indexer -> %s", models.AsJSONString(document))
	}

	return nil
}

func indexViewTriggerBuild(dbName string, ddoc string, key string) {

	// When we do a bulk update on monitored objects, we'll be issuing
	// a lot of view queries so instead. So now we check if there is already
	// generating a view and if so, then just quit. There's no point in hammering
	// couch to update the views.
	_, stored := couchdbViewBuilderBusyMap.LoadOrStore(ddoc, true)
	if stored == true {
		// We're already building the index, don't interrupt it.
		return
	}

	db, err := getDatabase(dbName)
	if err != nil {
		logger.Log.Errorf("Could not load db %s", dbName)
	}
	uri := fmt.Sprintf("_design/%s/_view/by%s", ddoc, key)
	logger.Log.Debugf("Starting to Index %s%s", dbName, uri)
	// Now go get the view (we don't actually look at it, we just want couch to start the indexer)
	_, err = db.Get(uri, nil)
	if err != nil {
		logger.Log.Errorf("Unsuccessfully Indexed %sbecause %s", uri, err.Error())
		return
	}
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Successfully Indexed %s -> %s", uri, "") //models.AsJSONString(v))
	}

	couchdbViewBuilderBusyMap.Delete(ddoc)
}
