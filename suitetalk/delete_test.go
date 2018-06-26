package suitetalk

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRequest(t *testing.T) {
	lib.IsVerbose = true
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml", "../tests/xml/deleteResult.xml"})
	res, sr := DeleteRequest(c, []string{"/FOOBAR/FOO/BAR", "/FOOBAR/foo.js"})

	assert.Equal(t, "foo.js", sr[0].Name)
	assert.Equal(t, "hashes.json", sr[1].Name)
	assert.Equal(t, "BAR", sr[2].Name)
	assert.Equal(t, "01", sr[0].InternalID)
	assert.Equal(t, "20", sr[1].InternalID)
	assert.Equal(t, "2", sr[2].InternalID)
	assert.Equal(t, "/FOOBAR/foo.js", sr[0].Path)
	assert.Equal(t, "/FOOBAR/FOO/BAR/hashes.json", sr[1].Path)
	assert.Equal(t, "/FOOBAR/FOO/BAR", sr[2].Path)

	assert.Equal(t, "1221", res[2].ID)
	assert.Equal(t, "999999", res[1].ID)
	assert.Equal(t, "123456", res[0].ID)
	assert.Equal(t, "DELETED", res[2].Code)
	assert.Equal(t, "MEDIA_NOT_FOUND", res[1].Code)
	assert.Equal(t, "RCRD_DSNT_EXIST", res[0].Code)
	assert.Equal(t, "Record was successfully deleted", res[2].Message)
	assert.Equal(t, "Media item not found 999999", res[1].Message)
	assert.Equal(t, "That record does not exist.", res[0].Message)

	lib.PrintResponse("foo", res)

	assert.NotNil(t, res)
	assert.NotNil(t, sr)
}

func TestSoapDelete(t *testing.T) {
	a := []*lib.SearchResult{
		&lib.SearchResult{InternalID: "1", IsDir: true},
		&lib.SearchResult{InternalID: "2", IsDir: true},
		&lib.SearchResult{InternalID: "3", IsDir: false},
		&lib.SearchResult{InternalID: "4", IsDir: true},
		&lib.SearchResult{InternalID: "5", IsDir: true},
		&lib.SearchResult{InternalID: "6", IsDir: false},
		&lib.SearchResult{InternalID: "7", IsDir: true},
		&lib.SearchResult{InternalID: "8", IsDir: true},
		&lib.SearchResult{InternalID: "9", IsDir: false},
	} // len 10

	chunks, docs := soapDelete(a, 5)
	assert.NotZero(t, chunks)
	assert.NotZero(t, docs)
	assert.True(t, len(chunks) == 2 && len(docs) == 2)
	for i, d := range docs {
		assert.NotZero(t, d)
		s, err := d.WriteToString()
		assert.Nil(t, err)
		if i == 0 {
			assert.Contains(t, s, "<baseRef type=\"folder\" internalId=\"1\"")
			assert.Contains(t, s, "<baseRef type=\"folder\" internalId=\"2\"")
			assert.Contains(t, s, "<baseRef type=\"file\" internalId=\"3\"")
			assert.Contains(t, s, "<baseRef type=\"folder\" internalId=\"4\"")
			assert.Contains(t, s, "<baseRef type=\"folder\" internalId=\"5\"")
		} else {
			assert.Contains(t, s, "<baseRef type=\"file\" internalId=\"6\"")
			assert.Contains(t, s, "<baseRef type=\"folder\" internalId=\"7\"")
			assert.Contains(t, s, "<baseRef type=\"folder\" internalId=\"8\"")
			assert.Contains(t, s, "<baseRef type=\"file\" internalId=\"9\"")
		}
	}
}

func TestParsesoapDelete(t *testing.T) {
	con, _ := ioutil.ReadFile("../tests/xml/deleteResult.xml")
	res, err := parseSoapDelete(con)
	assert.Nil(t, err)
	assert.Equal(t, "1221", res[2].ID)
	assert.Equal(t, "999999", res[1].ID)
	assert.Equal(t, "123456", res[0].ID)
	assert.Equal(t, "DELETED", res[2].Code)
	assert.Equal(t, "MEDIA_NOT_FOUND", res[1].Code)
	assert.Equal(t, "RCRD_DSNT_EXIST", res[0].Code)
	assert.Equal(t, "Record was successfully deleted", res[2].Message)
	assert.Equal(t, "Media item not found 999999", res[1].Message)
	assert.Equal(t, "That record does not exist.", res[0].Message)
	assert.True(t, len(res) == 3)
}

func TestParseFail(t *testing.T) {
	res, err := parseSoapDelete([]byte{})
	assert.Nil(t, res)
	assert.EqualError(t, err, "REQUEST_ERROR")
}

func TestSplit(t *testing.T) {
	a := []*lib.SearchResult{
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
		&lib.SearchResult{},
	} // len 11

	chunks := split(a, 5)

	assert.NotZero(t, chunks)
	assert.True(t, len(chunks) == 3)
	assert.True(t, len(chunks[0]) == 5)
	assert.True(t, len(chunks[2]) == 1)
}

func TestPrintResponse(t *testing.T) {
	con, _ := ioutil.ReadFile("../tests/xml/deleteResult.xml")
	res, _ := parseSoapDelete(con)

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lib.PrintResponse("HEADER", res)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	assert.Contains(t, out, "| # |   ID   |  TYPE  |      CODE       |             MESSAGE             |")
	assert.Contains(t, out, "|---|--------|--------|-----------------|---------------------------------|")
	assert.Contains(t, out, "| 1 |   1221 | file   | DELETED         | Record was successfully deleted |")
	assert.Contains(t, out, "| 2 | 999999 | file   | MEDIA_NOT_FOUND | Media item not found 999999     |")
	assert.Contains(t, out, "| 3 | 123456 | folder | RCRD_DSNT_EXIST | That record does not exist.     |")
	assert.Contains(t, out, "HEADER")
	assert.NotNil(t, out)
}
