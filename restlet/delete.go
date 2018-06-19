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

	r := lib.Credentials

	var c client
	c.client.Credentials.Token = r[lib.ConsumerKey]
	c.client.Credentials.Secret = r[lib.ConsumerSecret]
	c.token.Token = r[lib.TokenID]
	c.token.Secret = r[lib.TokenSecret]

	uri := strings.Join([]string{
		"https://rest.",
		strings.Replace(r[lib.Realm], "system.", "", 1),
		"/app/site/hosting/restlet.nl",
	}, "")

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))

	q := req.URL.Query()
	q.Add("script", lib.RestletScriptID)
	q.Add("deploy", lib.RestletScriptDeployment)
	req.URL.RawQuery = q.Encode()

	// Get the authorization header (deprecated but needed in this case)
	a := c.client.AuthorizationHeader(&c.token, "POST", req.URL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", strings.Join([]string{a, ", realm=\"", r[lib.Account], "\""}, ""))

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
