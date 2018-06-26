package sdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/stretchr/testify/assert"
)

var hashes = []lib.Hash{
	lib.Hash{
		Path: "/FOOBAR/FOO/BAR/Bar Foo",
		Hash: "d41d8cd98f00b204e9800998ecf8427e",
		Name: "Bar Foo"},
	lib.Hash{
		Path: "/FOOBAR/FOO/BAR/hashes.json",
		Hash: "51ac1788ae90fa01a3f57e3af4b1252a",
		Name: "hashes.json"},
	lib.Hash{
		Path: "/FOOBAR/FOO/BAR/foo bar.js",
		Hash: "d41d8cd98f00b204e9800998ecf8427e",
		Name: "foo bar.js"},
	lib.Hash{
		Path: "/FOOBAR/FOO/BAR/subdir/subdircontent.xml",
		Hash: "d41d8cd98f00b204e9800998ecf8427e",
		Name: "subdircontent.xml"},
}

func TestSyncInitial(t *testing.T) {
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"})
	lib.Credentials = map[string]string{}
	lib.Credentials[lib.HashFile] = "hashfile.json"
	wd, _ := os.Getwd()
	al, ar, dl, dr, err := Sync(e, c, filepath.Join(wd, "..", "tests", "fs"), "/SuiteScripts", false, nil)
	assert.Nil(t, al)
	assert.Nil(t, ar)
	assert.Nil(t, dl)
	assert.Nil(t, dr)
	assert.Nil(t, err)
}

func TestSyncErrSrc(t *testing.T) {
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"})
	lib.Credentials = map[string]string{}
	lib.Credentials[lib.HashFile] = "hashfile.json"
	al, ar, dl, dr, err := Sync(e, c, "foo", "/SuiteScripts", false, nil)
	assert.Nil(t, al)
	assert.Nil(t, ar)
	assert.Nil(t, dl)
	assert.Nil(t, dr)
	assert.Contains(t, err.Error(), "foo\" does not exist")
}

func TestSyncErrDest(t *testing.T) {
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"})
	lib.Credentials = map[string]string{}
	lib.Credentials[lib.HashFile] = "hashes.json"
	wd, _ := os.Getwd()
	al, ar, dl, dr, err := Sync(e, c, filepath.Join(wd, "..", "tests", "fs"), "FOOBAR/FOO/BAR", false, nil)
	assert.Nil(t, al)
	assert.Nil(t, ar)
	assert.Nil(t, dl)
	assert.Nil(t, dr)
	assert.Error(t, err, "destination has to be absolute")
}

func TestSync(t *testing.T) {
	lib.IsVerbose = true
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"})
	lib.Credentials = map[string]string{}
	lib.Credentials[lib.HashFile] = "hashes.json"
	lib.Credentials[lib.Realm] = "system.netsuite.com"
	wd, _ := os.Getwd()
	al, ar, dl, dr, err := Sync(e, c, filepath.Join(wd, "..", "tests", "fs"), "/FOOBAR/FOO/BAR", false, hashes)
	assert.Nil(t, al)
	assert.Nil(t, ar)
	assert.Nil(t, dl)
	assert.Nil(t, dr)
	assert.Nil(t, err)
}

func TestSyncAddedLocal(t *testing.T) {
	lib.IsVerbose = true
	reset([]string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml", "../tests/xml/deleteResult.xml"})
	lib.Credentials = map[string]string{}
	lib.Credentials[lib.HashFile] = "hashes.json"
	wd, _ := os.Getwd()
	al, ar, dl, dr, err := Sync(e, c, filepath.Join(wd, "..", "tests", "xml"), "/FOOBAR/FOO/BAR", false, hashes)
	assert.Contains(t, al[0], "xml/deleteResult.xml")
	assert.Nil(t, ar)
	assert.Nil(t, dl)
	assert.Nil(t, dr)
	assert.Nil(t, err)
}
