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
	monitoredObjectIndex   = "indexOfObjectName"
	mapFnName              = "map"
	legacyIndexTemplateStr = "_design/%s/_view/%s"
	indexTemplateStr       = "_design/indexOf%s/_view/by%s"
	viewTemplateStr        = "_design/viewOf%s/_view/by%s"
	legacyViewTemplateStr  = "_design/%s/_view/%s"

	metaFieldPrefix              = "meta"
	metakeysViewDdocName         = "metaViews"
	metaViewUniqueKeys           = "uniqueKeys"
	MetakeysViewUniqueKeysURI    = "_design/metaViews/_view/uniqueKeys"
	metakeysViewUniqueValuessURI = "uniqueValues"
	metaViewAllValuesPerKey      = "allValuesByKeyWithCounts"
	metaViewLookupWords          = "lookupWords"
	metaViewSearchLookup         = "searchLookup"

	MetaKeyIndexOf = "indexOf"
	MetaKeyViewOf  = "indexOf"

	metaKeyName  = "{{KeyName}}"
	metaKeyField = "{{KeyField}}"

	moIndexDdoc = "moIndex"
	moIndexView = "byName"

	objectCountDdoc       = "monitoredObjectCount"
	objectCountByNameView = "byName"
	objectCountView       = "byCount"

	monitoredObjectCountIndexBytes = `{
	"_id": "_design/monitoredObjectCount",
	"language": "javascript",
	"views": {
	  "byDomain": {
		"map": "function(doc) {\n    if (doc.data && doc.data.datatype && doc.data.datatype === 'monitoredObject' && doc.data.domainSet) {\n      for (var i in doc.data.domainSet) {\n        emit(doc.data.domainSet[i], doc._id);\n      }\n    }\n}"
	  },
    "byName": {
      "map": "function(doc) {\n        emit(doc.data.objectName, null);\n\n}"
     },
	  "count": {
		"map": "function(doc) {  if (doc.data && doc.data.datatype && doc.data.datatype === 'monitoredObject') { emit(doc.id, 1) } }",
		"reduce": "_count"
	  }, "byCount": {
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

	metaViews = `{
		"_id": "_design/metaViews",
		"language": "javascript",
		"views": {
		  "allValuesByKeyWithCounts": {
			"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        emit(key, doc.data.meta[key].toLowerCase());\n        \n      }\n    }\n}",
			"reduce": "function (keys, values, rereduce) {\n  if (rereduce) {\n    var valuesByKeySubtrees = values\n    return valuesByKeySubtrees.reduce(function(mergedValuesByKey, valuesByKeySubtree) {\n      if (!valuesByKeySubtree) {\n        return mergedValuesByKey\n      }\n      \n      Object.keys(valuesByKeySubtree).forEach(function (key) {\n        if (mergedValuesByKey[key]) {\n          if (mergedValuesByKey[key] === 'extended' || valuesByKeySubtree[key] === 'extended') {\n            mergedValuesByKey[key] = 'extended'\n            return\n          }\n\n          Object.keys(valuesByKeySubtree[key]).forEach(function (valueKey) {\n            var count = mergedValuesByKey[key][valueKey]\n            if (count) {\n              mergedValuesByKey[key][valueKey] = count + valuesByKeySubtree[key][valueKey]\n            } else {\n              mergedValuesByKey[key][valueKey] = valuesByKeySubtree[key][valueKey]\n            }\n          })\n          \n          \n          if (Object.keys(mergedValuesByKey[key]).length > 7) {\n            mergedValuesByKey[key] = 'extended'\n          }\n        } else {\n          mergedValuesByKey[key] = valuesByKeySubtree[key]\n        }\n      })\n      \n      return mergedValuesByKey\n    }, {})\n  }\n  \n  return keys.reduce(function (valuesByKey, key, index) {\n    var keyName = key[0]\n    \n    if (!valuesByKey[keyName]) {\n      valuesByKey[keyName] = {}\n    }\n      \n    var count = valuesByKey[keyName][values[index]]\n    if (count) {\n      valuesByKey[keyName][values[index]] = count + 1\n    } else {\n      valuesByKey[keyName][values[index]] = 1\n    }\n    return valuesByKey\n  }, {})\n}"
		  },
		  "allValuesPerKey": {
			"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        emit(key, doc.data.meta[key]);\n        \n      }\n    }\n}",
			"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    if (rereduce) {\n        var resp = [];\n        resp = values;//.flat();\n        resp = flatten(resp);\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n        \n        if (filteredArray > 1000) {\n        return [];\n        } else {\n            return filteredArray;\n        }\n    } else {\n     if (values.length > 1000) {\n        return [];\n     } else {\n        return values;\n     }\n    }\n\n}\n"
		  },
		  "allValuesReduce": {
			"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        emit(doc.data.meta[key], key);\n        \n      }\n    }\n}",
			"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    if (rereduce) {\n        var resp = [];\n        resp = values;//.flat();\n        resp = flatten(resp);\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        return filteredArray.slice(0,1000);\n    } else {\n     return values;\n    }\n\n}\n"
		  },
		  "lookupWords": {
			"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n        var sentence = doc.data.meta[key];\n        \n        // split on space\n        var words = sentence.split(\" \");\n        \n        for (var word in words) {\n            emit(key + \"__\" + words[word].toLowerCase(), {\"sentence\":sentence.toLowerCase()});\n        }\n      }\n    }\n}",
			"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n    \n    var sentences = []\n    \n    if (!rereduce) {\n      log(\"reducing:\" + JSON.stringify(keys) + \" || val:\" + JSON.stringify(values).toLowerCase())\n      \n      return flatten(values)\n    }\n\n    if (rereduce) {\n      if (!keys) {\n        keys = \"unknown\"\n      }\n      log(\"rereducing:\" + JSON.stringify(keys) + \" || val:\" + JSON.stringify(values).toLowerCase())\n\n      var flat = flatten(values)\n      log(\"rereducing:\" + JSON.stringify(keys) + \" || flattened:\" + JSON.stringify(flat).toLowerCase())\n    \n      for (var v in flat) {\n        sentences.push(flat[v][\"sentence\"])\n      }      \n\n      var filtered = sentences.filter(function(item, pos){\n            return sentences.indexOf(item)== pos; \n        });\n      \n      var ls = []\n      for (var d in filtered) {\n        var item = {\n          \"sentence\": filtered[d]\n        }\n        ls.push(item)\n      }\n      \n      \n      return ls\n    }\n\n}\n"
		  },
		  "uniqueKeys": {
			"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n          emit(key, doc.data.meta[key].toLowerCase());\n      }\n    }\n}",
			"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    // All the reduce functions have completed and returned an aggregate of collections\n    if (rereduce) {\n        var resp = [];\n        resp = values;\n        // We need to flatten the collection because each reduce function returns an array.\n        // If we left the structure as is, each row is filtered as a whole, rather than filtering on each individual item in the row.\n        resp = flatten(resp);\n        \n        // Filter the collection for only the unique items.\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        // Now get the count of unique values in this collections\n        var cnt = 0;\n        for( var v in filteredArray) {\n            cnt++;\n            }\n        return cnt;\n\n    } else {\n     return values;\n    }\n\n}\n"
		  }
		}
	  }`
	metaUniqueValuesViewsURI          = "{{KeyName}}/{{KeyField}}"
	metaUniqueValuesViewsDdocTemplate = `{
		"_id": "_design/viewOf{{KeyName}}",
		"language": "javascript",
		"views": {
			"by{{KeyName}}": {
				"reduce": "function(keys, values) {\n    return sum(values);\n}",
				"map": "function(doc) {\n    if (doc.data.meta) {\n        if (doc.data.{{KeyField}}) {\n            emit(doc.data.{{KeyField}}, 1);\n        }\n    }\n}"
			}
		}
	}`

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
	moIndexBytes = `{
		"_id": "_design/moIndex",
		"views": {
			"byName": {
				"map": "function (doc) {\n  if (doc.data && doc.data.datatype && doc.data.datatype === 'monitoredObject') {\n    emit(doc.data.objectName, doc.data.objectId);\n  }\n}"
			}
		},
		"language": "javascript"
	}`

	indexMonObjectNames = `{
		"_id": "_design/indexOfObjectName",
		"language": "query",
		"views": {
			"byObjectName": {
				"map": {
					"fields": {
						"data.objectName": "asc"
					},
					"partial_filter_selector": {}
				},
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
)

func getTenantViews() []map[string]interface{} {
	indexMonObjectNamesObject := map[string]interface{}{}
	metaViewObject := map[string]interface{}{}
	monitoredObjectMetaIndexObject := map[string]interface{}{}
	monitoredObjectCountIndexObject := map[string]interface{}{}
	moIndexObject := map[string]interface{}{}

	if err := json.Unmarshal([]byte(indexMonObjectNames), &indexMonObjectNamesObject); err != nil {
		logger.Log.Errorf("Unable to generate Meta View Index: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(metaViews), &metaViewObject); err != nil {
		logger.Log.Errorf("Unable to generate Meta View Index: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(monitoredObjectMetaIndexBytes), &monitoredObjectMetaIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate Unique Meta Index: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(moIndexBytes), &moIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate MO Index: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(monitoredObjectCountIndexBytes), &monitoredObjectCountIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate MO Index: %s", err.Error())
	}

	return []map[string]interface{}{metaViewObject, monitoredObjectMetaIndexObject, monitoredObjectCountIndexObject, moIndexObject, indexMonObjectNamesObject}
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
		if logger.IsDebugEnabled() {
			logger.Log.Debugf("Error creating index design document %s: %s :%s\n", tenmod.TenantMetaStr, models.AsJSONString(document), err.Error())
		}
		return err
	}
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Successfully created Indexer -> %s", models.AsJSONString(document))
	}

	return nil
}

// TriggerBuildCouchView - Couchviews are not auto generated so this function makes couchdb start the build process
func TriggerBuildCouchView(dbName string, ddoc string, key string, legacy bool) {
	//http://docs.couchdb.org/en/latest/api/ddoc/views.html#querying-views-and-indexes
	// When we do a bulk update on monitored objects, we'll be issuing
	// a lot of view queries so instead. So now we check if there is already
	// generating a view and if so, then just quit. There's no point in hammering
	// couch to update the views.
	_, stored := couchdbViewBuilderBusyMap.LoadOrStore(ddoc, true)
	// Should we defer for 5 seconds to let the system settle?
	if stored == true {
		// We're already building the index, don't interrupt it.
		return
	}

	db, err := getDatabase(dbName)
	if err != nil {
		logger.Log.Errorf("Could not load db %s", dbName)
	}
	var uri string
	uri = fmt.Sprintf(viewTemplateStr, ddoc, key)
	if legacy {
		uri = fmt.Sprintf(legacyViewTemplateStr, ddoc, key)
	}
	logger.Log.Debugf("Starting to build view %s%s", dbName, uri)
	// Now go get the view (we don't actually look at it, we just want couch to start the indexer)
	_, err = db.Get(uri, nil)
	if err != nil {
		logger.Log.Errorf("Unsuccessfully built view: %s because %s", uri, err.Error())
		return
	}
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Successfully built view %s -> %s", uri, "") //models.AsJSONString(v))
	}

	couchdbViewBuilderBusyMap.Delete(ddoc)
}

// TriggerBuildCouchIndex - Couchviews are not auto generated so this function makes couchdb start the build process
func TriggerBuildCouchIndex(dbName string, ddoc string, key string, legacyName bool) {

	//http://docs.couchdb.org/en/latest/api/ddoc/views.html#querying-views-and-indexes
	// When we do a bulk update on monitored objects, we'll be issuing
	// a lot of view queries so instead. So now we check if there is already
	// generating a view and if so, then just quit. There's no point in hammering
	// couch to update the views.
	_, stored := couchdbViewBuilderBusyMap.LoadOrStore(ddoc, true)
	// Should we defer for 5 seconds to let the system settle?
	if stored == true {
		// We're already building the index, don't interrupt it.
		return
	}

	db, err := getDatabase(dbName)
	if err != nil {
		logger.Log.Errorf("Could not load db %s", dbName)
	}
	var uri string
	uri = fmt.Sprintf(indexTemplateStr, ddoc, key)
	if legacyName {
		uri = fmt.Sprintf(legacyIndexTemplateStr, ddoc, key)
	}
	logger.Log.Debugf("Starting to Index %s%s", dbName, uri)
	// Now go get the view (we don't actually look at it, we just want couch to start the indexer)
	_, err = db.Get(uri, nil)
	if err != nil {
		logger.Log.Errorf("Unsuccessfully Indexed view: %s because %s", uri, err.Error())
	} else {
		if logger.IsDebugEnabled() {
			logger.Log.Debugf("Successfully Indexed %s -> %s", uri, "") //models.AsJSONString(v))
		}
	}

	couchdbViewBuilderBusyMap.Delete(ddoc)
}

// createNewTenantMetadataViews - Create an index based on metadata keys
func createNewTenantMetadataViews(dbName string, key string) error {
	err := createCouchDBViewIndex(dbName, metaIndexTemplate, key, []string{key}, metaFieldPrefix)
	if err != nil {
		if !strings.Contains(err.Error(), "status 409 - conflict") {
			msg := fmt.Sprintf("Index failed db:%s key:%s Error: %s", dbName, key, err.Error())
			return errors.New(msg)
		}
	}
	// Create a view based on unique values per new value
	err = createCouchDBViewIndex(dbName, metaUniqueValuesViewsDdocTemplate, key, []string{key}, metaFieldPrefix)
	if err != nil {
		if !strings.Contains(err.Error(), "status 409 - conflict") {
			msg := fmt.Sprintf("View failed db:%s key:%s Error: %s", dbName, key, err.Error())
			return errors.New(msg)
		}
	}
	return nil
}
