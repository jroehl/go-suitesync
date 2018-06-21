package suitetalk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/beevik/etree"
	"github.com/jroehl/go-suitesync/lib"
)

func SOAPRequest(command string, searchString string) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)
	envelope := doc.CreateElement("soap:Envelope")
	envelope.CreateAttr("xmlns:soap", "http://schemas.xmlsoap.org/soap/envelope/")
	envelope.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	envelope.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	addTokenHeader(envelope.CreateElement("soap:Header"))

	body := envelope.CreateElement("soap:Body")
	action := ""
	switch command {
	case "searchFolder":
		action = "search"
		search(body, "q1:FolderSearchAdvanced", searchString)
	case "searchFile":
		action = "search"
		search(body, "q1:FileSearchAdvanced", searchString)
	default:
		log.Fatalf("Command \"%s\" not implemented", command)
	}

	bytes, err := doc.WriteToBytes()
	if err != nil {
		log.Fatal(err)
	}

	doc.Indent(2)
	doc.WriteTo(os.Stdout)

	res := request(bytes, action)
	arr := parseSearch(res)
	fmt.Println(arr)
}

type SearchResult struct {
	InternalID string
	Parent     string
	Name       string
}

func request(body []byte, action string) []byte {

	req, err := http.NewRequest("POST", fmt.Sprintf("https://webservices.%s/services/NetSuitePort_2018_1", strings.Replace(lib.Credentials[lib.Realm], "system.", "", 1)), bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("SOAPAction", action)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bits, err := ioutil.ReadAll(res.Body)

	return bits
}
