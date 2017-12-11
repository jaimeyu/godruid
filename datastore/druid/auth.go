package druid

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/accedian/adh-gather/config"
)

type Auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

func GetAuthCode(cfg config.Provider) string {
	// include the URL from this step:
	// https://www.notion.so/accedian/Accessing-IAP-protected-resources-using-OAUTH2-bf19cf64e46a4fcbac0f49926c82ae39#2e70c2a6017d4c1288b8a447db5756de
	// inside the yaml config file.
	// this is only here for dev purposes and should go away once we
	// implement oauth
	url := cfg.GetString("druid.auth")

	if url != "" {
		req, _ := http.NewRequest("POST", url, nil)

		req.Header.Add("cache-control", "no-cache")
		req.Header.Add("postman-token", "1d5ed9ef-2e83-a103-c4e1-f68ec0e56134")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var auth Auth

		json.Unmarshal(body, &auth)

		return auth.IDToken
	}
	return ""
}
