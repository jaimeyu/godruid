package couchDB

// CouchdbViewResultItem - Model for Extracting Couchdb data from Views
type ViewResultItem struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

// CouchdbViewResults - View results
type ViewResults struct {
	Rows []ViewResultItem `json:"rows"`
}

// CouchdbIndexResults -Model for Index results
type IndexResults struct {
	Docs     map[string]interface{} `json:"docs"`
	Bookmark string                 `json:"bookmark"`
}
