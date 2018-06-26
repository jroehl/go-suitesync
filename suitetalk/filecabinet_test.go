package suitetalk

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fsReqs = []string{"../tests/xml/searchFolder.xml", "../tests/xml/searchFile.xml"}

func TestListFiles7(t *testing.T) {
	reset(fsReqs)

	res, err := ListFiles(c, "/FOOBAR")

	assert.Nil(t, err)
	assert.Len(t, res, 7)
	assert.True(t, len(c.Calls) == 2)
	assert.NotNil(t, Cache)
	assert.NotNil(t, Pathlookup)
}

func TestListFiles4(t *testing.T) {
	reset(fsReqs)

	res, err := ListFiles(c, "/FOOBAR/FOO")

	assert.Nil(t, err)
	assert.Len(t, res, 4)
	assert.True(t, len(c.Calls) == 2)
	assert.NotNil(t, Cache)
	assert.NotNil(t, Pathlookup)
}

func TestListFilesFailDir(t *testing.T) {
	reset(fsReqs)

	res, err := ListFiles(c, "/FOOBAR/foo")

	assert.Nil(t, res)
	assert.EqualError(t, err, "\nNo result for \"/FOOBAR/foo\"\n\n")
}

func TestListFilesFailNotFound(t *testing.T) {
	reset(fsReqs)

	res, err := ListFiles(c, "/FOOBAR/foobar.js")

	assert.Nil(t, res)
	assert.EqualError(t, err, "\"/FOOBAR/foobar.js\" is not a directory")
}

func TestGetFs(t *testing.T) {
	reset(fsReqs)

	getFs(c)

	printTree(Cache, 0)

	assert.Len(t, c.Calls, 2)
	assert.NotNil(t, Cache)
	assert.NotNil(t, Pathlookup)
}

func TestCache(t *testing.T) {
	reset(fsReqs)

	getFs(c)
	getFs(c)

	assert.Len(t, c.Calls, 2)
	assert.Len(t, Pathlookup, 8)

	assert.NotNil(t, Cache)
	assert.NotNil(t, Pathlookup)
}

func TestGetPath(t *testing.T) {
	reset(fsReqs)

	it, err := GetPath(c, "/FOOBAR/FOO")

	assert.Len(t, it.Children, 3)
	assert.NotNil(t, it)
	assert.Nil(t, err)
}

func TestGetPathNotFound(t *testing.T) {
	reset(fsReqs)

	it, err := GetPath(c, "/ewrwerwre")

	assert.Nil(t, it)
	assert.EqualError(t, err, "\nNo result for \"/ewrwerwre\"\n\n")
}

func TestFlattenChildren(t *testing.T) {
	reset(fsReqs)

	it, err := GetPath(c, "/FOOBAR")
	flattened := FlattenChildren(it.Children)
	assert.Len(t, flattened, 7)
	assert.NotNil(t, flattened)
	assert.Nil(t, err)
}

func TestPrintTree(t *testing.T) {
	reset(fsReqs)
	res := getFs(c)

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printTree(res, 0)

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

	assert.Contains(t, out, "> FOOBAR (TYPE: folder   ID: 0   PATH: /FOOBAR)")
	assert.NotNil(t, out)
}
