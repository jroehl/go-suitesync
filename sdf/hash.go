package sdf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jroehl/go-suitesync/suitetalk"

	"github.com/jroehl/go-suitesync/lib"
)

// UpdateHashFile create and upload the hashes to ns filecabinet
func UpdateHashFile(bash BashExec, src string, dest string, exclude bool, failed []string) (hashes []lib.Hash, hashfile string) {
	c, _ := lib.DirContent(src, dest, true, exclude)

	for _, ct := range c {
		if !lib.ArrayIncludes(ct.Path, failed) {
			// exclude the failed hashes
			hashes = append(hashes, ct)
		}
	}

	if lib.IsVerbose {
		lib.PrettyHash("Updated hash", hashes)
	}

	tmp := lib.MkTempDir()
	tmpDest := filepath.Join(tmp, dest)
	tmpFile := filepath.Join(tmpDest, lib.Credentials[lib.HashFile])
	os.MkdirAll(tmpDest, os.ModePerm)
	err := ioutil.WriteFile(tmpFile, lib.ToJSON(hashes), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// TODO check if necessary
	// suitetalk.DeleteRequest(http.DefaultClient, []string{filepath.Clean(filepath.Join("", dest, lib.Credentials[lib.HashFile]))})

	uploaded, err := Upload(bash, []string{tmpFile}, dest)
	lib.Remove(tmp)
	if err != nil {
		lib.PrFatalf("\n%s\n", err.Error())
	}
	hashfile = filepath.Join(uploaded[0].Dest, uploaded[0].Path)
	return hashes, hashfile
}

// get the hash files content
func getRemoteHash(bash BashExec, client suitetalk.HTTPClient, location string) (c []lib.Hash) {
	dest := lib.MkTempDir()
	h, err := Download(bash, client, []string{location}, dest)
	if len(h) == 0 || err != nil {
		// no hashfile found
		lib.Remove(dest)
		return
	}

	raw, _ := ioutil.ReadFile(filepath.Join(h[0].Dest, h[0].Path))

	lib.Remove(dest)

	json.Unmarshal(raw, &c)
	return
}
