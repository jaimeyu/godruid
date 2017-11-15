package models

type CouchChangeFeedResponse struct {
	LastSeq string                  `json:"last_seq"`
	Pending int32                   `json:"pending"`
	Results []CouchChangeFeedResult `json:"results"`
}

type CouchChangeFeedResult struct {
	Changes []CouchRevWrapper `json:"changes"`
	ID      string            `json:"id"`
	Seq     string            `json:"seq"`
	Deleted bool              `json:"deleted"`
}

type CouchRevWrapper struct {
	Rev string `json:"rev"`
}
