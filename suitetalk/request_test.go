package suitetalk

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSoap(t *testing.T) {
	doc, body := soap()

	s, _ := doc.WriteToString()
	assert.Contains(t, s, "<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	assert.Contains(t, s, "<soap:Envelope xmlns:soap=")
	assert.Contains(t, s, "<soap:Header>")
	assert.Contains(t, s, "<soap:Body/>")

	assert.NotNil(t, doc)
	assert.NotNil(t, body)
}

func TestParseByteFail(t *testing.T) {
	con, _ := ioutil.ReadFile("../tests/xml/faultyResponse.xml")
	doc, err := parseByte(con)

	assert.Nil(t, doc)
	assert.EqualError(t, err, "\nsoapenv:Server.userException\nThe request could not be understood by the server due to malformed syntax.\n")
}

func TestParseByteSuccess(t *testing.T) {
	doc, err := parseByte([]byte{})

	assert.Nil(t, err)
	assert.NotNil(t, doc)
}
