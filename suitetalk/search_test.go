package suitetalk

import (
	"io/ioutil"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestSearchRequestFolder(t *testing.T) {
	reset([]string{"../tests/xml/searchFolder.xml"})
	res := SearchRequest(c, searchFolder)

	assert.Equal(t, "0", res[0].InternalID)
	assert.Equal(t, "1", res[1].InternalID)
	assert.Equal(t, "2", res[2].InternalID)
	assert.Equal(t, "3", res[3].InternalID)
	assert.Equal(t, "", res[0].Parent)
	assert.Equal(t, "0", res[1].Parent)
	assert.Equal(t, "1", res[2].Parent)
	assert.Equal(t, "1", res[3].Parent)
	assert.Equal(t, "FOOBAR", res[0].Name)
	assert.Equal(t, "FOO", res[1].Name)
	assert.Equal(t, "BAR", res[2].Name)
	assert.Equal(t, "BAZ", res[3].Name)
	assert.True(t, res[0].IsDir)
	assert.True(t, res[1].IsDir)
	assert.True(t, res[2].IsDir)
	assert.True(t, res[3].IsDir)
	assert.Zero(t, res[0].Children)
	assert.Zero(t, res[1].Children)
	assert.Zero(t, res[2].Children)
	assert.Zero(t, res[3].Children)

	assert.NotNil(t, res)
}

func TestSearchRequestFile(t *testing.T) {
	reset([]string{"../tests/xml/searchFile.xml"})
	res := SearchRequest(c, searchFile)

	assert.Equal(t, "00", res[0].InternalID)
	assert.Equal(t, "01", res[1].InternalID)
	assert.Equal(t, "10", res[2].InternalID)
	assert.Equal(t, "20", res[3].InternalID)
	assert.Equal(t, "0", res[0].Parent)
	assert.Equal(t, "0", res[1].Parent)
	assert.Equal(t, "1", res[2].Parent)
	assert.Equal(t, "2", res[3].Parent)
	assert.Equal(t, "foobar.js", res[0].Name)
	assert.Equal(t, "foo.js", res[1].Name)
	assert.Equal(t, "bar.js", res[2].Name)
	assert.Equal(t, "hashes.json", res[3].Name)
	assert.False(t, res[0].IsDir)
	assert.False(t, res[1].IsDir)
	assert.False(t, res[2].IsDir)
	assert.False(t, res[3].IsDir)

	assert.NotNil(t, res)
}

func TestSearchRequestFolderMore(t *testing.T) {
	reset([]string{"../tests/xml/searchFolderMore.xml", "../tests/xml/searchMoreWithIdResult.xml"})
	res := SearchRequest(c, searchFolder)

	assert.Equal(t, "666", res[4].InternalID)
	assert.Equal(t, "", res[4].Parent)
	assert.Equal(t, "SMORE", res[4].Name)
	assert.True(t, res[4].IsDir)
	assert.Zero(t, res[4].Children)

	assert.NotNil(t, res)
}

func TestSoapSearch(t *testing.T) {
	doc := etree.NewDocument()
	envelope := doc.CreateElement("soap:Envelope")
	soapSearch(envelope, folderSearchAdvanced)

	s, _ := doc.WriteToString()
	assert.Contains(t, s, "<searchRecord xmlns:q1=")
	assert.Contains(t, s, "<q1:columns>")
	assert.Contains(t, s, "<q1:basic>")
	assert.Contains(t, s, "xsi:type=\"q1:FolderSearchAdvanced\">")
	assert.Contains(t, s, "<parent xmlns=")
	assert.Contains(t, s, "<internalId xmlns=")
	assert.Contains(t, s, "<name xmlns=")
	assert.NotNil(t, doc)
}

func TestSoapSearchMore(t *testing.T) {
	res, doc := soapSearchMore(1, "SearchId123")

	s, _ := doc.WriteToString()
	assert.Contains(t, s, "<searchMoreWithId>")
	assert.Contains(t, s, "<searchId>SearchId123</searchId>")
	assert.Contains(t, s, "<pageIndex>1</pageIndex>")
	assert.Contains(t, s, "<ns1:account>account</ns1:account>")
	assert.Contains(t, s, "<ns1:consumerKey>consumer_key</ns1:consumerKey>")
	assert.Contains(t, s, "<ns1:token>tokenid</ns1:token>")
	assert.Contains(t, s, "<ns1:signature algorithm=")
	assert.Contains(t, s, "<ns1:nonce>")
	assert.Contains(t, s, "<ns1:timestamp>")
	assert.Contains(t, s, "<ns1:signature algorithm=")
	assert.NotNil(t, doc)
	assert.NotNil(t, res)
}

func TestParseSoapSearchFile(t *testing.T) {
	con, _ := ioutil.ReadFile("../tests/xml/searchFile.xml")
	res, meta, err := parseSoapSearch(con, false, false)

	assert.Nil(t, err)

	assert.Equal(t, "00", res[0].InternalID)
	assert.Equal(t, "01", res[1].InternalID)
	assert.Equal(t, "10", res[2].InternalID)
	assert.Equal(t, "20", res[3].InternalID)
	assert.Equal(t, "0", res[0].Parent)
	assert.Equal(t, "0", res[1].Parent)
	assert.Equal(t, "1", res[2].Parent)
	assert.Equal(t, "2", res[3].Parent)
	assert.Equal(t, "foobar.js", res[0].Name)
	assert.Equal(t, "foo.js", res[1].Name)
	assert.Equal(t, "bar.js", res[2].Name)
	assert.Equal(t, "hashes.json", res[3].Name)
	assert.False(t, res[0].IsDir)
	assert.False(t, res[1].IsDir)
	assert.False(t, res[2].IsDir)
	assert.False(t, res[3].IsDir)

	assert.Equal(t, true, meta.Successful)
	assert.Equal(t, 1, meta.TotalPages)
	assert.Equal(t, 4, meta.TotalRecords)
	assert.Equal(t, "WEBSERVICES_TSTDRV411742_06212018601356341139979712_a13e9488", meta.SearchID)

	assert.True(t, len(res) == 4)

	assert.NotNil(t, res)
	assert.NotNil(t, meta)
}

func TestParseSoapSearchFolder(t *testing.T) {
	con, _ := ioutil.ReadFile("../tests/xml/searchFolder.xml")
	res, meta, err := parseSoapSearch(con, true, false)

	assert.Nil(t, err)

	assert.Equal(t, "0", res[0].InternalID)
	assert.Equal(t, "1", res[1].InternalID)
	assert.Equal(t, "2", res[2].InternalID)
	assert.Equal(t, "3", res[3].InternalID)
	assert.Equal(t, "", res[0].Parent)
	assert.Equal(t, "0", res[1].Parent)
	assert.Equal(t, "1", res[2].Parent)
	assert.Equal(t, "1", res[3].Parent)
	assert.Equal(t, "FOOBAR", res[0].Name)
	assert.Equal(t, "FOO", res[1].Name)
	assert.Equal(t, "BAR", res[2].Name)
	assert.Equal(t, "BAZ", res[3].Name)
	assert.True(t, res[0].IsDir)
	assert.True(t, res[1].IsDir)
	assert.True(t, res[2].IsDir)
	assert.True(t, res[3].IsDir)
	assert.Zero(t, res[0].Children)
	assert.Zero(t, res[1].Children)
	assert.Zero(t, res[2].Children)
	assert.Zero(t, res[3].Children)

	assert.Equal(t, true, meta.Successful)
	assert.Equal(t, 1, meta.TotalPages)
	assert.Equal(t, 4, meta.TotalRecords)
	assert.Equal(t, "WEBSERVICES_TSTDRV411742_06212018601356341139979712_a13e9488", meta.SearchID)

	assert.True(t, len(res) == 4)

	assert.NotNil(t, res)
	assert.NotNil(t, meta)
}

func TestParseSoapSearchMore(t *testing.T) {
	con, _ := ioutil.ReadFile("../tests/xml/searchMoreWithId.xml")
	res, meta, err := parseSoapSearch(con, true, true)
	assert.Nil(t, err)
	assert.Zero(t, res)
	assert.NotNil(t, meta)
}

func TestParseSoapSearchFail(t *testing.T) {
	res, meta, err := parseSoapSearch([]byte{}, false, false)
	assert.Zero(t, res)
	assert.Zero(t, meta)
	assert.EqualError(t, err, "REQUEST_ERROR")
}
