package restlet

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
)

type Body struct {
	Action string   `json:"action"`
	Items  []string `json:"items"`
}

func request(body []byte) []byte {

	uri := strings.Join([]string{
		"https://rest.",
		strings.Replace(lib.Credentials[lib.Realm], "system.", "", 1),
		"/app/site/hosting/restlet.nl",
	}, "")

	var (
		req    *http.Request
		err    error
		method string
	)
	if body != nil {
		method = "POST"
		req, err = http.NewRequest(method, uri, bytes.NewBuffer(body))
	} else {
		method = "GET"
		req, err = http.NewRequest(method, uri, nil)
	}

	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("script", lib.RestletScriptID)
	q.Add("deploy", lib.RestletScriptDeployment)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/json")

	a, _ := lib.GetAuthRestlet(method, uri, lib.HmacSha256)

	h := strings.Join([]string{
		"OAuth realm=\"", lib.PercentEncode(a.Account), "\",",
		"oauth_consumer_key=\"", lib.PercentEncode(a.ConsumerKey), "\",",
		"oauth_token=\"", lib.PercentEncode(a.Token), "\",",
		"oauth_signature_method=\"", lib.PercentEncode(a.Algorithm), "\",",
		"oauth_timestamp=\"", lib.PercentEncode(a.Timestamp), "\",",
		"oauth_nonce=\"", lib.PercentEncode(a.Nonce), "\",",
		"oauth_version=\"", lib.PercentEncode("1.0"), "\",",
		"oauth_signature=\"", lib.PercentEncode(a.Signature), "\"",
	}, "")

	req.Header.Add("Authorization", h)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bits, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bits
}

// Healthcheck check if restlet is set up
func Healthcheck() (*lib.Response, error) {
	if lib.IsVerbose {
		lib.PrNoticeF("Healthcheck request issued\n")
	}

	a := request(nil)

	decoded := new(lib.Response)
	err := json.Unmarshal(a, decoded)
	if err != nil {
		return nil, err
	}
	if decoded.Code != 200 {
		return nil, errors.New("healthcheck error")
	}
	return decoded, nil
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
	if decoded.Successful == nil || decoded.Unsuccessful == nil {
		log.Fatal(string(a))
	}
	lib.PrintResponse("Deleted items", decoded.Successful)
	lib.PrintResponse("Failed deleting items", decoded.Unsuccessful)
	return decoded
}
