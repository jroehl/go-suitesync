package sdf

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/jroehl/go-suitesync/lib"
)

// Flag flag structure
type Flag struct {
	F string
	A string
}

// deploy the project to netsuite
func deployProject(bash BashExec, project string) string {
	return Command(bash, "deploy", []Flag{
		Flag{F: "project", A: project},
		Flag{F: "np", A: ""},
	}, false)
}

// import the files specified in the paths array to the specified project
func importFiles(bash BashExec, paths []string, project string) string {
	t := template.Must(template.New("paths").Funcs(template.FuncMap{"clean": filepath.Clean}).Parse(`{{block "list" .}}{{range .}}"{{clean .}}" {{end}}{{end}}`))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, paths); err != nil {
		panic(err)
	}
	return Command(bash, "importfiles", []Flag{
		Flag{F: "paths", A: tpl.String()},
		Flag{F: "p", A: project},
	}, true)
}

// build flags for sdf command
func buildFlags(flags []Flag) (string, error) {
	flags = append(flags,
		Flag{F: "url", A: lib.Credentials[lib.URL]},
		Flag{F: "email", A: lib.Credentials[lib.Email]},
		Flag{F: "account", A: lib.Credentials[lib.Account]},
		Flag{F: "role", A: lib.Credentials[lib.Role]},
	)
	t := template.Must(template.New("args").Parse(` {{block "list" .}}{{range .}}-{{.F}} {{.A}} {{end}}{{end}}`))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, flags); err != nil {
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
