package restlet

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gomodule/oauth1/oauth"
	"github.com/jroehl/go-suitesync/lib"
)

type Body struct {
	Action string   `json:"action"`
	Items  []string `json:"items"`
}

type client struct {
	client oauth.Client
	token  oauth.Credentials
}

func request(body []byte) []byte {

	r := lib.Creds

	var c client
	c.client.Credentials.Token = r.ConsumerKey
	c.client.Credentials.Secret = r.ConsumerSecret
	c.token.Token = r.TokenID
	c.token.Secret = r.TokenSecret

	uri := strings.Join([]string{
		"https://rest.",
		strings.Replace(lib.Creds.Realm, "system.", "", 1),
		"/app/site/hosting/restlet.nl",
	}, "")

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))

	if r.Script != "" || r.Deployment != "" {
		q := req.URL.Query()
		if r.Script != "" {
			q.Add("script", r.Script)
		}
		if r.Deployment != "" {
			q.Add("deploy", r.Deployment)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Get the authorization header (deprecated but needed in this case)
	a := c.client.AuthorizationHeader(&c.token, "POST", req.URL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", strings.Join([]string{a, ", realm=\"", r.Account, "\""}, ""))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bits, err := ioutil.ReadAll(res.Body)

	return bits
}

// Delete items from filecabinet
func Delete(items []string) *lib.Responses {

	if lib.IsVerbose {
		lib.PrettyList("DELETE request issued for", items)
	}

	b := Body{Action: "delete", Items: items}

	j := lib.ToJson(b)

	a := request(j)

	decoded := new(lib.Responses)
	err := json.Unmarshal(a, decoded)
	if err != nil {
		log.Fatal(err)
	}
	lib.PrintResponse("Deleted items", decoded.Successful)
	lib.PrintResponse("Failed deleting items", decoded.Unsuccessful)
	return decoded
}
