package sdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHashFile(t *testing.T) {
	lib.IsVerbose = true
	lib.Credentials = map[string]string{}
	lib.Credentials[lib.HashFile] = "hashfile.json"
	wd, _ := os.Getwd()
	res, hf := UpdateHashFile(e, filepath.Join(wd, "..", "tests", "fs"), "/SuiteScripts", true, []string{})
	assert.Equal(t, "/SuiteScripts/hashfile.json", hf)

	expected := []lib.Hash{
		lib.Hash{
			Path: "/SuiteScripts/Bar Foo",
			Hash: "d41d8cd98f00b204e9800998ecf8427e",
			Name: "Bar Foo"},
		lib.Hash{
			Path: "/SuiteScripts/hashes.json",
			Hash: "51ac1788ae90fa01a3f57e3af4b1252a",
			Name: "hashes.json"},
		lib.Hash{
			Path: "/SuiteScripts/foo bar.js",
			Hash: "d41d8cd98f00b204e9800998ecf8427e",
			Name: "foo bar.js"},
		lib.Hash{
			Path: "/SuiteScripts/subdir/subdircontent.xml",
			Hash: "d41d8cd98f00b204e9800998ecf8427e",
			Name: "subdircontent.xml"},
	}
	assert.ElementsMatch(t, expected, res)
}

func TestGetRemoteHash(t *testing.T) {
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"})
	lib.IsVerbose = true
	res := getRemoteHash(e, c, "/FOOBAR/FOO/BAR/hashes.json")
	assert.Nil(t, res)
}

func TestGetRemoteHashNil(t *testing.T) {
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"})
	res := getRemoteHash(e, c, "/werwer")
	assert.Nil(t, res)
}
