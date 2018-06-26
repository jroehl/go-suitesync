package sdf

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/jroehl/go-suitesync/suitetalk"
)

// Download paths from filecabinet to local filesystem
func Download(bash BashExec, c suitetalk.HTTPClient, paths []string, dest string) (downloads []FileTransfer, err error) {
	dest = lib.AbsolutePath(dest)
	var dirs []FileTransfer
	for _, p := range paths {
		if !strings.HasPrefix(p, "/") {
			return nil, errors.New("source paths have to be absolute")
		}
		item, err := suitetalk.GetPath(c, p)
		if err != nil {
			return nil, err
		}
		if item.IsDir {
			var dls []FileTransfer
			dls, dirs = getDirDownloads(c, p, dest)
			downloads = append(downloads, dls...)
		} else {
			downloads = append(downloads, FileTransfer{Path: filepath.Base(p), Src: p, Dest: dest})
		}
	}
	if len(downloads) > 0 {
		processDownloads(bash, downloads, dirs)
	}
	return downloads, nil
}

// processDownloads download the files specified in the paths array to the specified directory
func processDownloads(bash BashExec, downloads, dirs []FileTransfer) (downloaded, created []FileTransfer) {
	tmp := lib.MkTempDir()
	tmpDir := filepath.Dir(tmp)
	tmpBase := filepath.Base(tmp)
	pr := CreateAccountCustomizationProject(tmpBase, tmpDir)

	var paths []string
	for _, d := range downloads {
		paths = append(paths, d.Src)
	}

	for _, d := range dirs {
		// create empty dirs
		os.MkdirAll(filepath.Join(d.Dest, d.Path), os.ModePerm)
	}

	importFiles(bash, paths, pr.Dir)
	for _, d := range downloads {
		src := filepath.Join(pr.FileCabinet, d.Src)
		dest := filepath.Join(d.Dest, d.Path)
		os.MkdirAll(filepath.Dir(dest), os.ModePerm)
		lib.Copy(src, dest)
		downloaded = append(downloaded, FileTransfer{Dest: dest, Src: src})
	}

	if lib.IsVerbose {
		c := []string{}
		for _, d := range downloaded {
			c = append(c, strings.Join([]string{d.Src, d.Dest}, " -> "))
		}
		lib.PrettyList("Copied", c)
	}
	lib.Remove(tmp)

	for _, d := range dirs {
		paths = append(paths, d.Src)
	}
	lib.PrettyList("Downloaded", paths)
	return downloaded, dirs
}

// getDirDownloads from directory to use in download
func getDirDownloads(c suitetalk.HTTPClient, src string, dest string) (downloads, dirs []FileTransfer) {
	res, err := suitetalk.ListFiles(c, src)
	if err != nil {
		lib.PrFatalf("%s\n", err.Error())
	}
	for _, r := range res {
		p := filepath.Clean(strings.Replace(r.Path, src, string(filepath.Separator), 1))
		if !r.IsDir {
			downloads = append(downloads, FileTransfer{Path: p, Dest: dest, Src: r.Path})
		} else {
			dirs = append(dirs, FileTransfer{Path: p, Dest: dest, Src: r.Path})
		}
	}
	return downloads, dirs
}
