package sdf

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/jroehl/go-suitesync/suitetalk"
	"github.com/stretchr/testify/mock"
)

type Ex struct {
}

var (
	i        = 0
	reqFiles = []string{}
	c        = new(ClientMock)
	e        = new(Ex)
)

func (e *Ex) Command(name string, arg ...string) *exec.Cmd {
	cmd := append([]string{name}, arg...)
	return exec.Command("echo", strings.Join(cmd, " "))
}

// ClientMock HTTPClient mock
type ClientMock struct {
	mock.Mock
}

func reset(requests []string) {
	i = 0
	reqFiles = requests
	suitetalk.Cache = nil
	suitetalk.Pathlookup = make(map[string]*lib.SearchResult)
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
