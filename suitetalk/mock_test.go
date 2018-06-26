package suitetalk

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/stretchr/testify/mock"
)

// ClientMock HTTPClient mock
type ClientMock struct {
	mock.Mock
}

var (
	i        = 0
	reqFiles = []string{}
	c        = new(ClientMock)
)

func reset(requests []string) {
	i = 0
	reqFiles = requests
	Cache = nil
	Pathlookup = make(map[string]*lib.SearchResult)
	c = new(ClientMock)
	c.On("Do").Return(mock.AnythingOfType("string"), nil)
}

// Do mock
func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	con, _ := ioutil.ReadFile(reqFiles[i])
	t := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(string(con))),
	}
	c.Called()
	i = i + 1
	return &t, nil
}
