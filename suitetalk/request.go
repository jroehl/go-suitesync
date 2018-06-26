package suitetalk

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/beevik/etree"
	"github.com/jroehl/go-suitesync/lib"
)

const (
	chunkSize = 200 // web services limit size requests

	searchFolder         = "searchFolder"
	deleteFolder         = "deleteFolder"
	folderSearchAdvanced = "q1:FolderSearchAdvanced"
	searchFile           = "searchFile"
	deleteFile           = "deleteFile"
	fileSearchAdvanced   = "q1:FileSearchAdvanced"
	search               = "search"
	searchMoreWithID     = "searchMoreWithId"
	deleteList           = "deleteList"

	// types for soap requests
	common      = "urn:common_2018_1.platform.webservices.netsuite.com"
	messages    = "urn:messages_2018_1.platform.webservices.netsuite.com"
	filecabinet = "urn:filecabinet_2018_1.documents.webservices.netsuite.com"
	core        = "urn:core_2018_1.platform.webservices.netsuite.com"
)

// HTTPClient resource interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func soap() (*etree.Document, *etree.Element) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)
	envelope := doc.CreateElement("soap:Envelope")
	envelope.CreateAttr("xmlns:soap", "http://schemas.xmlsoap.org/soap/envelope/")
	envelope.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	envelope.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	addTokenHeader(envelope.CreateElement("soap:Header"))
	body := envelope.CreateElement("soap:Body")
	return doc, body
}

func doRequest(client HTTPClient, body []byte, action string) []byte {
	url := strings.Join([]string{"https://webservices.", strings.Replace(lib.Credentials[lib.Realm], "system.", "", 1), "/services/NetSuitePort_2018_1"}, "")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("SOAPAction", action)

	res, err := client.Do(req)
	if err != nil {
		lib.PrFatalf("\nRequest to \"%s\" failed - %s", url, err.Error())
	}
	defer res.Body.Close()

	bits, err := ioutil.ReadAll(res.Body)

	return bits
}

func parseByte(xml []byte) (*etree.Document, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xml); err != nil {
		panic(err)
	}
	if err := doc.FindElement("soapenv:Envelope/soapenv:Body/soapenv:Fault"); err != nil {
		fc := err.FindElement("faultcode")
		fs := err.FindElement("faultstring")
		e := strings.Join([]string{"\n", fc.Text(), "\n", fs.Text(), "\n"}, "")
		return nil, errors.New(e)
	}
	return doc, nil
}
