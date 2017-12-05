package druid

import (
	"encoding/json"

	"github.com/accedian/godruid"
)

func queryToString(query godruid.Query, debug bool) string {
	var reqJson []byte
	var err error

	if debug {
		reqJson, err = json.MarshalIndent(query, "", "  ")
	} else {
		reqJson, err = json.Marshal(query)
	}

	if err != nil {
		return ""
	}

	return string(reqJson)
}
