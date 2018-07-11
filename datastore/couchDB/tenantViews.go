package couchDB

import (
	"encoding/json"

	"github.com/accedian/adh-gather/logger"
)

var (
	monitoredObjectCountIndexBytes = []byte(`{
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
  }`)

	monitoredObjectMetaIndexBytes = []byte(`{
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
  }`)

	uniqueMetaIndexBytes = []byte(`{
	"_id": "_design/uniqueMeta",
	"language": "javascript",
	"views": {
	  "uniqueKeys": {
		"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n          emit(key, doc.data.meta[key]);\n      }\n    }\n}",
		"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    // All the reduce functions have completed and returned an aggregate of collections\n    if (rereduce) {\n        var resp = [];\n        resp = values;\n        // We need to flatten the collection because each reduce function returns an array.\n        // If we left the structure as is, each row is filtered as a whole, rather than filtering on each individual item in the row.\n        resp = flatten(resp);\n        \n        // Filter the collection for only the unique items.\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        // Now get the count of unique values in this collections\n        var cnt = 0;\n        for( var v in filteredArray) {\n            cnt++;\n            }\n        return cnt;\n\n    } else {\n     return values;\n    }\n\n}\n"
	  },
	  "allValuesPerKey": {
		"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        emit(key, doc.data.meta[key]);\n        \n      }\n    }\n}",
		"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    if (rereduce) {\n        var resp = [];\n        resp = values;\n        \n        // Flatten the collection since we receives an array of arrays\n        resp = flatten(resp);\n        \n        // Filter for only the unique values\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        return filteredArray;\n    } else {\n     return values;\n    }\n\n}\n"
	  },
	  "allValuesReduce": {
		"map": "function(doc) {\n    if(doc.data.meta) {\n      for (var key in doc.data.meta) {\n      \n        // Reverse the order so we can search per existing value, rather than by key\n        emit(doc.data.meta[key], key);\n        \n      }\n    }\n}",
		"reduce": "function(keys, values, rereduce) {\n    function flatten(arr) {\n      return arr.reduce(function (flat, toFlatten) {\n        return flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten);\n      }, []);\n    }\n\n    if (rereduce) {\n        var resp = [];\n        resp = values;//.flat();\n        resp = flatten(resp);\n        var filteredArray = resp.filter(function(item, pos){\n            return resp.indexOf(item)== pos; \n        });\n\n        return filteredArray.slice(0,1000);\n    } else {\n     return values;\n    }\n\n}\n"
	  }
	}
  }`)

	uniqueMetaIndexObject           = map[string]interface{}{}
	monitoredObjectMetaIndexObject  = map[string]interface{}{}
	monitoredObjectCountIndexObject = map[string]interface{}{}

	tenantMonitoredObjectViews = []map[string]interface{}{}
)

func init() {
	if err := json.Unmarshal(uniqueMetaIndexBytes, &uniqueMetaIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate Unique Meta Index: %s", err.Error())
	}
	if err := json.Unmarshal(monitoredObjectMetaIndexBytes, &monitoredObjectMetaIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate Unique Meta Index: %s", err.Error())
	}
	if err := json.Unmarshal(monitoredObjectCountIndexBytes, &monitoredObjectCountIndexObject); err != nil {
		logger.Log.Errorf("Unable to generate Unique Meta Index: %s", err.Error())
	}

	tenantMonitoredObjectViews = append(tenantMonitoredObjectViews, uniqueMetaIndexObject, monitoredObjectMetaIndexObject, monitoredObjectCountIndexObject)
}
