package sdf

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	pwd, _ = os.Getwd()
	uls    = []FileTransfer{
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/xml/deleteResult.xml"),
			Root: path.Join(pwd, "..", "/tests/xml"),
			Path: "/deleteResult.xml",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/xml/searchFolder.xml"),
			Root: path.Join(pwd, "..", "/tests/xml"),
			Path: "/searchFolder.xml",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/Bar Foo"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/Bar Foo",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/hashes.json"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/hashes.json",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/foo bar.js"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/foo bar.js",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/subdir/subdircontent.xml"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/subdir/subdircontent.xml",
			Dest: "/SuiteScripts"},
	}
)

func TestGetDirUploads(t *testing.T) {
	d, _ := os.Getwd()
	uploads := getDirUploads(path.Join(d, "..", "tests", "fs"), "/SuiteScripts")
	expected := []FileTransfer{
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/Bar Foo"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/Bar Foo",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/hashes.json"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/hashes.json",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/foo bar.js"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/foo bar.js",
			Dest: "/SuiteScripts"},
		FileTransfer{
			Src:  path.Join(pwd, "..", "/tests/fs/subdir/subdircontent.xml"),
			Root: path.Join(pwd, "..", "/tests/fs"),
			Path: "/subdir/subdircontent.xml",
			Dest: "/SuiteScripts"},
	}
	assert.ElementsMatch(t, expected, uploads)
}

func TestProcessUploads(t *testing.T) {
	copied, err := processUploads(e, uls, "/SuiteScripts")
	expected := []string{
		"/SuiteScripts/deleteResult.xml",
		"/SuiteScripts/searchFolder.xml",
		"/SuiteScripts/Bar Foo",
		"/SuiteScripts/hashes.json",
		"/SuiteScripts/foo bar.js",
		"/SuiteScripts/subdir/subdircontent.xml",
	}
	assert.ElementsMatch(t, expected, copied)
	assert.Nil(t, err)
}

func TestProcessUploadsFail(t *testing.T) {
	copied, err := processUploads(e, []FileTransfer{FileTransfer{Root: "/foobar", Dest: "/SuiteScripts", Path: "test", Src: "/foobar/test"}}, "/SuiteScripts")
	assert.Nil(t, copied)
	assert.Contains(t, err.Error(), "\" does not exist")
}

func TestUpload(t *testing.T) {
	d, _ := os.Getwd()
	fs := path.Join(d, "..", "tests", "fs")
	dr := path.Join(d, "..", "tests", "xml", "deleteResult.xml")
	sf := path.Join(d, "..", "tests", "xml", "searchFolder.xml")
	uploads, err := Upload(e, []string{fs, dr, sf}, "/SuiteScripts")
	expected := uls
	assert.Nil(t, err)
	assert.ElementsMatch(t, expected, uploads)
}

func TestUploadFail(t *testing.T) {
	uploads, err := Upload(e, []string{"../123456"}, "/SuiteScripts")
	assert.Nil(t, uploads)
	assert.Contains(t, err.Error(), "\" does not exist")
}
