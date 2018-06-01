package lib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Hash struct {
	Path string
	Hash string
	Name string
}

// get the hash of a file
func getHash(p string) string {
	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// normalize the remote path and return localpath
func NormalizeRootPath(p string, root string, remote string) string {
	return filepath.Clean(strings.Replace(path.Join(root, p), remote, "", 1))
}

// get all files from a directory and their hashes
func DirContent(dir string, prefix string, hash bool, exclude bool) []Hash {
	fileList := []Hash{}
	if prefix == "" {
		prefix = "/SuiteScripts"
	}
	err := filepath.Walk(dir, func(pt string, f os.FileInfo, err error) error {
		if (exclude && (strings.HasPrefix(f.Name(), ".") || strings.Contains(pt, "/."))) || f.IsDir() {
			return nil
		}
		h := ""
		if hash {
			h = getHash(pt)
		}
		p := filepath.Clean(path.Join(prefix, strings.Replace(pt, dir, "", 1)))
		fileList = append(fileList, Hash{Name: f.Name(), Path: p, Hash: h})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return fileList
}
