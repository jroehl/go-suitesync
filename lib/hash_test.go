package lib

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirContentFail(t *testing.T) {
	_, err := DirContent("./foobar", "", true, true)
	assert.EqualError(t, err, "\"./foobar\" does not exist")
}

func TestDirContent(t *testing.T) {
	d, _ := os.Getwd()
	res, err := DirContent(path.Join(d, "..", "tests", "fs"), "/Prefix", true, true)

	expected := []Hash{
		Hash{
			Path: "/Prefix/Bar Foo",
			Hash: "d41d8cd98f00b204e9800998ecf8427e",
			Name: "Bar Foo"},
		Hash{
			Path: "/Prefix/hashes.json",
			Hash: "51ac1788ae90fa01a3f57e3af4b1252a",
			Name: "hashes.json"},
		Hash{
			Path: "/Prefix/foo bar.js",
			Hash: "d41d8cd98f00b204e9800998ecf8427e",
			Name: "foo bar.js"},
		Hash{
			Path: "/Prefix/subdir/subdircontent.xml",
			Hash: "d41d8cd98f00b204e9800998ecf8427e",
			Name: "subdircontent.xml"},
	}
	assert.Nil(t, err)
	assert.ElementsMatch(t, expected, res)
}

func TestNormalizeRootPath(t *testing.T) {
	res := NormalizeRootPath("/RemotePath/foo/foobar.js", "/RootPath", "/RemotePath")
	assert.Equal(t, "/RootPath/foo/foobar.js", res)
}
