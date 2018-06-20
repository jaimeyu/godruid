package druid

import (
	"net/http"
	pkgurl "net/url"
	"strings"

	"github.com/accedian/adh-gather/config"
)

type Auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

// These are the dev credentials to log into debug deployment
const (
	username   = "admin@datahub.com"
	password   = "AccedianPass"
	cookieName = "skylight-aaa"
)

// GetAuthCode - Only used for dev builds running on local deployment to access a remote druid (because it has real usuable datasets).
func GetAuthCode(cfg config.Provider) string {
	url := cfg.GetString("druid.auth")

	// If druid auth isn't setup (usually only applicable on dev machines)
	if url != "" {
		data := pkgurl.Values{}
		data.Set("username", username)
		data.Add("password", password)

		// Send the request to get a oauth cookie
		req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))

		req.Header.Add("Accept", "application/vnd.api+json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, _ := http.DefaultClient.Do(req)

		// Make sure to close the body/clean up
		// Wait, do we need to? The GC should take care of this.
		defer res.Body.Close()
		//body, _ := ioutil.ReadAll(res.Body)

		var token string
		// Find the token, if we can't return empty token
		for _, cookie := range res.Cookies() {
			if cookie.Name == cookieName {
				token = cookie.Value
			}
		}

		return token
	}
	return ""
}
