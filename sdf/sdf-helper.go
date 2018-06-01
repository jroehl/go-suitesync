package sdf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/jroehl/go-suitesync/lib"
	"github.com/jroehl/go-suitesync/restlet"
)

// Flag flag structure
type Flag struct {
	F string
	A string
}

// Sync sync files with remote dir
func Sync(src string, dest string, bidirectional bool) (al []string, ar []string, dl []string, dr []string) {
	rh := getRemoteHash(filepath.Clean(filepath.Join(dest, lib.Creds.Hashfile)))
	lib.CheckDir(src)

	if len(rh) < 1 {
		lib.PrNoticeF("Initializing new sync dir in \"%s\":", dest)
		lib.PrNoticeF("(uploading files and creating hashfile)")

		UploadDir(src, dest)
		UpdateHashFile(src, dest, true, []string{})
		return nil, nil, nil, nil
	}

	lf := lib.DirContent(src, dest, true, true)

	if lib.IsVerbose {
		lib.PrNoticeF("Sync (bidirectional=%t)", bidirectional)
		lib.PrettyHash("Local files", lf)
		lib.PrettyHash("Remote hash", rh)
	}

	sfp, sfh := mapkeys(lf)
	hhp, hhh := mapkeys(rh)

	// added local
	_, ali := lib.Difference(sfh, hhh)
	for _, i := range ali {
		al = append(al, lib.NormalizeRootPath(lf[i].Path, src, dest))
	}
	lib.PrettyList("Added local", al)

	// deleted locally
	dls, _ := lib.Difference(hhp, sfp)
	lib.PrettyList("Deleted local", dls)

	if bidirectional {
		rf := ListFiles(dest)
		lib.PrettyList("Remote files", rf)

		// added remote
		ars, _ := lib.Difference(rf, hhp)
		lib.PrettyList("Added remote", ars)

		// deleted remote
		drs, _ := lib.Difference(hhp, rf)
		lib.PrettyList("Deleted remote", drs)

		lib.PrNoticeF("Execution skipped")
		lib.PrNoticeF("(bidirectional mode not implemented yet)")

		// non reachable code
		// added remote
		// if len(ars) > 0 {
		// 	s := []string{}
		// 	for _, a := range ars {
		// 		s = append(s, lib.NormalizeRootPath(a, src, dest))
		// 	}
		// 	DownloadFiles(ars, src)
		// }

		// // deleted remote
		// if len(drs) > 0 {
		// 	for _, d := range drs {
		// 		lib.Remove(lib.NormalizeRootPath(d, src, dest))
		// 	}
		// 	lib.PrettyList("Deleted remote", dr)
		// }
	}

	// added local
	if len(al) > 0 {
		UploadFiles(src, al, dest)
	}

	// deleted locally
	du := []string{}
	if len(dls) > 0 {
		r := restlet.Delete(dls).Unsuccessful
		if lib.IsVerbose {
			fmt.Println("Delete result", r)
		}
		nf := []string{}
		for _, s := range r {
			if s.Code != 404 {
				du = append(du, s.Path)
			} else {
				nf = append(nf, s.Path)
			}
		}
		if len(nf) > 0 {
			// these files are treated as successful when updating the hash
			lib.PrettyList("NOT FOUND (treated as successful for hash creation)", nf)
		}
	}

	UpdateHashFile(src, dest, true, du)

	return al, ar, dl, dr
}

// UploadRestlet upload the restlet and do deployment
func UploadRestlet() {
	deployProject(lib.Restlet)
	lib.PrNoticeF("Restlet (\"%s\") deployed", lib.Restlet)
}

// ListFiles list the files of a specific directory
func ListFiles(src string) []string {
	res := Sdf("listfiles", []Flag{Flag{F: "folder", A: src}}, false)
	re := regexp.MustCompile(strings.Join([]string{src, "\\/.*"}, ""))
	mt := re.FindAllString(res, -1)
	lib.PrettyList("Content", mt)
	return mt
}

// UploadFiles upload files to SuiteScripts directory, retaining the directory structure of the files
func UploadFiles(root string, files []string, dest string) []string {
	tmp := lib.MkTempDir()
	pr := SdfCreateAccountCustomizationProject("temp", tmp)
	nr := filepath.Clean(root)

	copied := []string{}
	for _, f := range files {
		nf := strings.Replace(filepath.Clean(f), nr, "", -1)
		sf := filepath.Clean(filepath.Join(nr, nf))

		lib.CheckExists(sf)
		nfb := filepath.Clean(pr.Filebase)
		nrd := strings.Replace(dest, lib.Creds.Rootpath, "", -1)
		rf := filepath.Clean(filepath.Join("", nrd, nf))
		df := filepath.Clean(filepath.Join(nfb, rf))
		os.MkdirAll(filepath.Dir(df), os.ModePerm)
		lib.Copy(sf, df)
		copied = append(copied, rf)
	}
	deployProject(pr.Dir)
	lib.Remove(tmp)
	lib.PrettyList("Uploaded", copied)
	return copied
}

// UploadDir upload directory to SuiteScripts directory, retaining the directory structure of the files
func UploadDir(src string, dest string) []string {
	r := filepath.Clean(src)
	c := lib.DirContent(src, dest, true, true)
	f := []string{}
	for _, h := range c {
		f = append(f, lib.NormalizeRootPath(h.Path, src, dest))
	}
	return UploadFiles(r, f, dest)
}

// DownloadFiles download the files specified in the paths array to the specified directory
func DownloadFiles(paths []string, dest string) (srcs []string, dests []string) {
	if lib.IsVerbose {
		lib.PrettyList("DownloadFiles", paths)
		lib.PrNoticeF("Destination\t\t%s", dest)
	}
	tmp := lib.MkTempDir()
	pr := SdfCreateAccountCustomizationProject("temp", tmp)
	importFiles(paths, pr.Dir)
	for _, p := range paths {
		rel := strings.Replace(p, lib.Creds.Rootpath, "", -1)
		src := filepath.Join(strings.Replace(pr.Filebase, "/SuiteScripts", "", -1), p)
		dest := filepath.Join(dest, rel)
		os.MkdirAll(filepath.Dir(dest), os.ModePerm)
		lib.Copy(src, dest)
		srcs = append(srcs, src)
		dests = append(dests, dest)
	}
	if lib.IsVerbose {
		c := []string{}
		for i, s := range srcs {
			c = append(c, strings.Join([]string{s, dests[i]}, " -> "))
		}
		lib.PrettyList("Copied", c)
	}
	lib.Remove(tmp)
	lib.PrettyList("Downloaded", dests)
	return srcs, dests
}

// DownloadDir download all files from a folder to the specified directory
func DownloadDir(src string, dest string) ([]string, []string) {
	res := ListFiles(src)
	return DownloadFiles(res, dest)
}

// UpdateHashFile create and upload the hashes to ns filecabinet
func UpdateHashFile(src string, dest string, exclude bool, failed []string) (f []lib.Hash) {
	c := lib.DirContent(src, dest, true, exclude)

	for _, ct := range c {
		if !lib.ArrayIncludes(ct.Path, failed) {
			f = append(f, ct)
		}
	}

	nd := strings.Replace(dest, lib.Creds.Rootpath, "", -1)
	tmp := lib.MkTempDir()
	th := filepath.Clean(filepath.Join(tmp, "tmp-hash"))
	d := filepath.Clean(filepath.Join(th, nd))
	hf := filepath.Clean(filepath.Join(d, lib.Creds.Hashfile))

	if lib.IsVerbose {
		lib.PrettyHash("Updated hashes", f)
	}

	os.MkdirAll(d, os.ModePerm)
	err := ioutil.WriteFile(hf, lib.ToJson(f), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	restlet.Delete([]string{filepath.Clean(filepath.Join("", dest, lib.Creds.Hashfile))})

	uf := UploadFiles(th, []string{hf}, "")[0]
	lib.PrNoticeF("Uploaded\t\t%s", uf)
	lib.Remove(tmp)
	lib.Remove(d)
	return f
}

// get the hash files content
func getRemoteHash(location string) (c []lib.Hash) {
	tmp := lib.MkTempDir()
	_, h := DownloadFiles([]string{location}, tmp)

	raw, err := ioutil.ReadFile(h[0])
	if err != nil {
		// file does not exist
		return
	}

	json.Unmarshal(raw, &c)
	return
}

// deploy the project to netsuite
func deployProject(project string) string {
	return Sdf("deploy", []Flag{
		Flag{F: "project", A: project},
		Flag{F: "np", A: ""},
	}, false)
}

// import the files specified in the paths array to the specified project
func importFiles(paths []string, project string) string {
	t := template.Must(template.New("paths").Funcs(template.FuncMap{"clean": filepath.Clean}).Parse(`{{block "list" .}}{{range .}}"{{clean .}}" {{end}}{{end}}`))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, paths); err != nil {
		log.Fatal(err)
	}
	return Sdf("importfiles", []Flag{
		Flag{F: "paths", A: tpl.String()},
		Flag{F: "p", A: project},
	}, true)
}

// import all files from a folder to the specified project
func importAllFiles(src string, project string) string {
	f := ListFiles(src)
	return importFiles(f, project)
}

// build flags for sdf command
func buildFlags(flags []Flag) (string, error) {
	f := flags
	if f == nil {
		f = []Flag{}
	}
	f = append(flags,
		Flag{F: "url", A: lib.Creds.Url},
		Flag{F: "email", A: lib.Creds.Email},
		Flag{F: "account", A: lib.Creds.Account},
		Flag{F: "role", A: lib.Creds.Role},
	)
	t := template.Must(template.New("args").Parse(` {{block "list" .}}{{range .}}-{{.F}} {{.A}} {{end}}{{end}}`))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, f); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func mapkeys(arr []lib.Hash) (sp []string, sh []string) {
	for _, x := range arr {
		sp = append(sp, x.Path)
		sh = append(sh, x.Hash)
	}
	return
}
