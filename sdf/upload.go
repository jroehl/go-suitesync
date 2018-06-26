package sdf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
)

// Upload src paths to dest
func Upload(bash BashExec, sources []string, dest string) (uploads []FileTransfer, err error) {
	var dirs []string
	for _, s := range sources {
		s := lib.AbsolutePath(s)
		if !strings.HasPrefix(dest, "/") {
			return nil, errors.New("destination has to be absolute")
		}
		f, err := lib.CheckExists(s)
		if err != nil {
			return nil, err
		}
		switch f.IsDir() {
		case true:
			dirs = append(dirs, s)
		case false:
			b := strings.Join([]string{string(filepath.Separator), filepath.Base(s)}, "")
			r := filepath.Clean(filepath.Dir(s))
			uploads = append(uploads, FileTransfer{
				Root: r,
				Path: b,
				Dest: dest,
				Src:  s,
			})
		}
	}
	for _, d := range dirs {
		uploads = append(uploads, getDirUploads(d, dest)...)
	}

	if len(uploads) > 0 {
		_, err := processUploads(bash, uploads, dest)
		if err != nil {
			return nil, err
		}
	}
	return uploads, nil
}

// processUploads upload files to SuiteScripts directory, retaining the directory structure of the files
func processUploads(bash BashExec, uploads []FileTransfer, dest string) (copied []string, err error) {
	tmp := lib.MkTempDir()
	tmpDir := filepath.Dir(tmp)
	tmpBase := filepath.Base(tmp)
	pr := CreateAccountCustomizationProject(tmpBase, tmpDir)

	for _, u := range uploads {

		if _, err := lib.CheckExists(u.Src); err != nil {
			return nil, err
		}

		pathDir := filepath.Dir(u.Path)
		tempDir := filepath.Join(pr.FileCabinet, dest, pathDir)

		// create temporary directory structure
		os.MkdirAll(tempDir, os.ModePerm)

		pathBase := filepath.Base(u.Path)
		err := lib.Copy(u.Src, filepath.Join(tempDir, pathBase))
		if err != nil {
			return nil, err
		}
		remote := filepath.Join(u.Dest, u.Path)
		copied = append(copied, remote)
	}
	deployProject(bash, pr.Dir)
	lib.Remove(tmp)
	lib.PrettyList("Uploaded", copied)
	return copied, nil
}

// getDirUploads from directory to use in upload
func getDirUploads(src string, dest string) (uploads []FileTransfer) {
	r := filepath.Clean(src)
	c, _ := lib.DirContent(src, dest, true, true)
	for _, h := range c {
		nr := lib.NormalizeRootPath(h.Path, src, dest)
		fmt.Println(nr)
		s := strings.Replace(nr, r, "", 1)
		uploads = append(uploads, FileTransfer{Src: nr, Root: r, Path: s, Dest: dest})
	}
	return uploads
}
