package lib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
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
func getHash(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", nil
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// NormalizeRootPath normalize the remote path and return localpath
func NormalizeRootPath(p string, root string, remote string) string {
	return filepath.Clean(strings.Replace(path.Join(root, p), remote, "", 1))
}

// DirContent get all files from a directory and their hashes
func DirContent(dir string, prefix string, hash bool, exclude bool) ([]Hash, error) {
	if err := CheckDir(dir); err != nil {
		return nil, err
	}
	fileList := []Hash{}
	err := filepath.Walk(dir, func(pt string, f os.FileInfo, err error) error {
		if f.IsDir() ||
			(exclude &&
				(strings.HasPrefix(f.Name(), ".") ||
					strings.HasPrefix(pt, ".") ||
					strings.Contains(pt, "/.") ||
					ArrayIncludes(f.Name(), Whitelisted))) {
			return nil
		}
		h := ""
		if hash {
			h, err = getHash(pt)
			if err != nil {
				return err
			}
		}
		p := filepath.Clean(path.Join(prefix, strings.Replace(pt, dir, "", 1)))
		fileList = append(fileList, Hash{Name: f.Name(), Path: p, Hash: h})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}
