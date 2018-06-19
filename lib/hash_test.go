package lib

import (
	"os"
	"path"
	"testing"
)

func TestDirContentFail(t *testing.T) {
	_, err := DirContent("./foobar", "", true, true)
	if err == nil {
		t.Errorf("DirContent failed, got: %s, want: %s.", "nil", "\"./foobar\" does not exist")
	}
}

func TestDirContent(t *testing.T) {
	d, _ := os.Getwd()
	res, err := DirContent(path.Join(d, "..", "tests", "fs"), "/Prefix", true, true)
	if err != nil {
		t.Errorf("DirContent failed, got: %s, want: %s.", err.Error(), "nil")
	}
	if len(res) != 4 {
		t.Errorf("DirContent failed, got: %d, want: %d entries.", len(res), 4)
	}
	if res[0].Path != "/Prefix/Bar Foo" {
		t.Errorf("DirContent failed, got: %s, want: %s.", res[0].Path, "/Prefix/Bar Foo")
	}
	if res[1].Name != "Foo.png" {
		t.Errorf("DirContent failed, got: %s, want: %s.", res[1].Name, "Foo.png")
	}
	if res[3].Hash != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("DirContent failed, got: %s, want: %s.", res[3].Hash, "d41d8cd98f00b204e9800998ecf8427e")
	}
}

func TestNormalizeRootPath(t *testing.T) {
	res := NormalizeRootPath("/RemotePath/foo/foobar.js", "/RootPath", "/RemotePath")
	if res != "/RootPath/foo/foobar.js" {
		t.Errorf("NormalizeRootPath failed, got: %s, want: %s entries.", res, "/RootPath/foo/foobar.js")
	}
}
