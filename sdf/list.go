package sdf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/jroehl/go-suitesync/suitetalk"
)

//  paths from filecabinet to local filesystem
func List(bash BashExec, c suitetalk.HTTPClient, dir string) (err error) {
	dir = lib.AbsolutePath(dir)
	if !strings.HasPrefix(dir, "/") {
		return errors.New("dir path has to be absolute")
	}

	item, err := suitetalk.GetPath(c, dir)
	if err != nil {
		return err
	}
	if item.IsDir {
		res, err := suitetalk.ListFiles(c, dir)
		fmt.Println(res)
		if err != nil {
			lib.PrFatalf("%s\n", err.Error())
		}
	} else {
		return errors.New("path has to be a directory")
	}
	return nil
}
