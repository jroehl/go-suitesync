package sdf

import (
	"testing"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/stretchr/testify/assert"
)

func TestDeployProject(t *testing.T) {
	res := deployProject(e, "test-project")
	assert.Equal(t, " deploy  -project test-project -np  -url  -email  -account  -role  \n", res)
}

func TestImportFiles(t *testing.T) {
	res := importFiles(e, []string{"path1", "path2", "path3"}, "test-project")
	assert.Equal(t, " importfiles  -paths \"path1\" \"path2\" \"path3\"  -p test-project -url  -email  -account  -role  \n", res)
}

func TestBuildFlags(t *testing.T) {
	res, err := buildFlags([]Flag{Flag{F: "flag1", A: "flag1val"}, Flag{F: "flag2"}})
	assert.Nil(t, err)
	assert.Equal(t, " -flag1 flag1val -flag2  -url  -email  -account  -role  ", res)
}

func TestMapKeys(t *testing.T) {
	sp, sh := mapkeys([]lib.Hash{lib.Hash{Path: "path1", Hash: "hash1"}, lib.Hash{Path: "path2", Hash: "hash2"}})
	assert.ElementsMatch(t, []string{"path1", "path2"}, sp)
	assert.ElementsMatch(t, []string{"hash1", "hash2"}, sh)
}
