package sdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jroehl/go-suitesync/lib"

	"github.com/stretchr/testify/assert"
)

var (
	dest, _ = os.Getwd()
	fsReqs  = []string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"}
	dls     = []FileTransfer{
		FileTransfer{Root: "", Path: "/bar.js", Dest: dest, Src: "/FOOBAR/FOO/bar.js"},
		FileTransfer{Root: "", Path: "/BAR/hashes.json", Dest: dest, Src: "/FOOBAR/FOO/BAR/hashes.json"},
		FileTransfer{Root: "", Path: "foo.js", Dest: dest, Src: "/FOOBAR/foo.js"}}
	dirs = []FileTransfer{
		FileTransfer{Root: "", Path: "/BAR", Dest: dest, Src: "/FOOBAR/FOO/BAR"},
		FileTransfer{Root: "", Path: "/BAZ", Dest: dest, Src: "/FOOBAR/FOO/BAZ"},
	}
)

func TestGetDirDownloads(t *testing.T) {
	reset(fsReqs)
	downloads, dirs := getDirDownloads(c, "/FOOBAR/FOO", "./Downloads")
	expectedDls := []FileTransfer{
		FileTransfer{Root: "", Path: "/bar.js", Dest: "./Downloads", Src: "/FOOBAR/FOO/bar.js"},
		FileTransfer{Root: "", Path: "/BAR/hashes.json", Dest: "./Downloads", Src: "/FOOBAR/FOO/BAR/hashes.json"},
	}
	expectedDirs := []FileTransfer{
		FileTransfer{Root: "", Path: "/BAR", Dest: "./Downloads", Src: "/FOOBAR/FOO/BAR"},
		FileTransfer{Root: "", Path: "/BAZ", Dest: "./Downloads", Src: "/FOOBAR/FOO/BAZ"},
	}
	assert.ElementsMatch(t, expectedDls, downloads)
	assert.ElementsMatch(t, expectedDirs, dirs)
}

func TestProcessDownloads(t *testing.T) {
	reset(fsReqs)
	lib.IsVerbose = true
	downloaded, created := processDownloads(e, dls, dirs)
	assert.Contains(t, downloaded[0].Dest, dest)
	assert.Contains(t, downloaded[0].Src, "/FileCabinet/")
	assert.Contains(t, downloaded[0].Src, "suitesync_")
	assert.ElementsMatch(t, dirs, created)
}

func TestDownload(t *testing.T) {
	reset(fsReqs)
	downloads, err := Download(e, c, []string{"/FOOBAR/FOO", "/FOOBAR/foo.js"}, "./Downloads")
	expected := []FileTransfer{
		FileTransfer{Root: "", Path: "/bar.js", Dest: filepath.Join(dest, "Downloads"), Src: "/FOOBAR/FOO/bar.js"},
		FileTransfer{Root: "", Path: "/BAR/hashes.json", Dest: filepath.Join(dest, "Downloads"), Src: "/FOOBAR/FOO/BAR/hashes.json"},
		FileTransfer{Root: "", Path: "foo.js", Dest: filepath.Join(dest, "Downloads"), Src: "/FOOBAR/foo.js"},
	}
	assert.Nil(t, err)
	assert.ElementsMatch(t, expected, downloads)
}

func TestDownloadFail(t *testing.T) {
	reset(fsReqs)
	downloads, err := Download(e, c, []string{"/snoos"}, "./Downloads")
	assert.Nil(t, downloads)
	assert.EqualError(t, err, "\nNo result for \"/snoos\"\n\n")
}

func TestDownloadWD(t *testing.T) {
	reset(fsReqs)
	downloads, err := Download(e, c, []string{"/FOOBAR/FOO", "/FOOBAR/foo.js"}, "./")
	assert.Nil(t, err)
	assert.ElementsMatch(t, dls, downloads)
}

func TestCleanup(t *testing.T) {
	lib.Remove(filepath.Join(dest, "FOOBAR"))
	lib.Remove(filepath.Join(dest, "Downloads"))
	lib.Remove(filepath.Join(dest, "BAR"))
	lib.Remove(filepath.Join(dest, "BAZ"))
}
