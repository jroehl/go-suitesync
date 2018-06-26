package sdf

import (
	"os"
	"testing"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	lib.CliCache = ".clicache"
	res, err := GenerateToken(e, "123456")
	assert.Equal(t, "", res)
	assert.Nil(t, err)
	os.Remove(".clicache")
}

func TestCommand(t *testing.T) {
	res := Command(e, "commandstring", []Flag{Flag{F: "flag"}}, true)
	assert.Equal(t, " commandstring  -flag  -url  -email  -account  -role  \n", res)
}

func TestExecute(t *testing.T) {
	res := execute(e, "bin", "cmd", "flags", "prompt", "./", true)
	assert.Equal(t, "bin cmd flags\n", res)
}

func TestCreateAccountCustomizationProject(t *testing.T) {
	res := CreateAccountCustomizationProject("name", "")
	assert.Contains(t, res.Dir, "name")
	assert.Contains(t, res.FileCabinet, "FileCabinet")
	assert.Equal(t, "name", res.Name)
	assert.Equal(t, "name", res.Params.Name)
	err := lib.CheckDir(res.Dir)
	assert.Nil(t, err)
	lib.Remove(res.Dir)
}

func TestCreateSuiteAppProject(t *testing.T) {
	res := CreateSuiteAppProject("name", "", "id", "version", "publisherid")
	assert.Contains(t, res.Dir, "publisherid.id")
	assert.Contains(t, res.FileCabinet, "publisherid.id/FileCabinet")
	assert.Equal(t, "publisherid.id", res.Name)
	assert.Equal(t, "name", res.Params.Name)
	assert.Equal(t, "version", res.Params.Version)
	assert.Equal(t, "publisherid", res.Params.PublisherID)
	err := lib.CheckDir(res.Dir)
	assert.Nil(t, err)
	lib.Remove(res.Dir)
}
