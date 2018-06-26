package sdf

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/jroehl/go-suitesync/suitetalk"
)

// Sync sync files with remote dir
func Sync(bash BashExec, client suitetalk.HTTPClient, src string, dest string, bidirectional bool, rh []lib.Hash) (al []string, ar []string, dl []string, dr []string, err error) {
	src = lib.AbsolutePath(src)
	fmt.Println(lib.CheckDir(src))
	if err := lib.CheckDir(src); err != nil {
		return nil, nil, nil, nil, err
	}
	if !strings.HasPrefix(dest, "/") {
		return nil, nil, nil, nil, errors.New("destination has to be absolute")
	}

	if rh == nil {
		rh = getRemoteHash(bash, client, filepath.Clean(filepath.Join(dest, lib.Credentials[lib.HashFile])))
	}

	if len(rh) < 1 {
		lib.PrResultf("Initializing new sync dir in \"%s\":\n", dest)
		lib.PrNoticef("(uploading files and creating hashfile)\n")

		Upload(bash, []string{src}, dest)
		UpdateHashFile(bash, src, dest, true, []string{})
		return nil, nil, nil, nil, err
	}

	lf, _ := lib.DirContent(src, dest, true, true)

	if lib.IsVerbose {
		lib.PrNoticef("Sync (bidirectional=%t)\n", bidirectional)
		lib.PrettyHash("Local files", lf)
		lib.PrettyHash("Remote hash", rh)
	}

	sfp, sfh := mapkeys(lf)
	hhp, hhh := mapkeys(rh)

	fmt.Println(sfh)
	fmt.Println(hhh)

	// added local
	_, ali := lib.Difference(sfh, hhh)
	for _, i := range ali {
		al = append(al, lib.NormalizeRootPath(lf[i].Path, src, dest))
	}
	lib.PrettyList("Added/altered local", al)

	// deleted locally
	dls, _ := lib.Difference(hhp, sfp)
	lib.PrettyList("Deleted local", dls)

	// if bidirectional {
	// 	res, err := suitetalk.ListFiles(http.DefaultClient, src)
	// 	if err != nil {
	// 		lib.PrFatalf("%s\n", err.Error())
	// 	}
	// 	rf := []string{}
	// 	for _, r := range res {
	// 		rf = append(rf, r.Path)
	// 	}
	// 	lib.PrettyList("Remote files", rf)

	// 	// added remote
	// 	ars, _ := lib.Difference(rf, hhp)
	// 	lib.PrettyList("Added remote", ars)

	// 	// deleted remote
	// 	drs, _ := lib.Difference(hhp, rf)
	// 	lib.PrettyList("Deleted remote", drs)

	// 	lib.PrNoticef("Execution skipped\n")
	// 	lib.PrResultf("(bidirectional mode not implemented yet)\n")

	// 	// non reachable code
	// 	// added remote
	// 	if len(ars) > 0 {
	// 		s := []string{}
	// 		for _, a := range ars {
	// 			s = append(s, lib.NormalizeRootPath(a, src, dest))
	// 		}
	// 		DownloadFiles(ars, src)
	// 	}

	// 	// deleted remote
	// 	if len(drs) > 0 {
	// 		for _, d := range drs {
	// 			lib.Remove(lib.NormalizeRootPath(d, src, dest))
	// 		}
	// 		lib.PrettyList("Deleted remote", dr)
	// 	}
	// }

	// added local
	if len(al) > 0 {
		Upload(bash, al, dest)
	}

	// deleted locally
	du := []string{}
	if len(dls) > 0 {
		r, _ := suitetalk.DeleteRequest(client, dls)
		if lib.IsVerbose {
			lib.PrintResponse("Delete results", r)
		}
		nf := []string{}
		for _, s := range r {
			if !s.Successful && !s.NotFound {
				du = append(du, s.ID)
			} else if s.NotFound {
				nf = append(nf, s.ID)
			}
		}
		if len(nf) > 0 {
			// these files are treated as successful when updating the hash
			lib.PrettyList("NOT FOUND (treated as successful for hashfile update)", nf)
		}
	}

	UpdateHashFile(bash, src, dest, true, du)

	return al, ar, dl, dr, err
}
