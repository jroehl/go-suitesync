package sdf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
)

type ProjectParams struct {
	Name        string
	Path        string
	ID          string
	Version     string
	PublisherID string
}

type Project struct {
	Type     string
	Dir      string
	Filebase string
	Name     string
	Params   ProjectParams
}

func GenerateToken() {
	if lib.Creds.Password == "" {
		lib.PrFatalf("Environment variable \"NSCONF_PASSWORD\" is undefined\n")
	}
	f, err := buildFlags(nil)
	if err != nil {
		log.Fatal(err)
	}
	res := execute(lib.SdfCli, "issuetoken", f, strings.Join([]string{lib.Creds.Password, "\n"}, ""), "", false)
	lib.PrNoticeF(res)
}

func CheckCliCache() []byte {
	b, err := ioutil.ReadFile(lib.CliCache)
	if err != nil {
		// file does not exist
		b = []byte(lib.Creds.CliToken)
		ioutil.WriteFile(lib.CliCache, b, os.ModePerm)
	}
	return b
}

// call sdfcli command
func Sdf(command string, flags []Flag, ignore bool) string {
	CheckCliCache()
	f, err := buildFlags(flags)
	if err != nil {
		log.Fatal(err)
	}

	res := execute(lib.SdfCli, command, f, "YES\n", "", ignore)
	return res
}

func execute(bin string, cmd string, flags string, prompt string, dir string, ignore bool) string {
	cmdStr := strings.Join([]string{bin, cmd, flags}, " ")

	if lib.IsVerbose {
		fmt.Println()
		lib.PrHeaderF("Executed command:")
		fmt.Printf("  %s\n", cmdStr)
		fmt.Println()
	}

	proc := exec.Command(bin, cmd, flags)
	if dir != "" {
		proc.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	if prompt != "" {
		proc.Stdin = bytes.NewBuffer([]byte(prompt))
	}
	proc.Stdout = &stdout
	proc.Stderr = &stderr
	cmderr := proc.Run()
	stdoutStr, stderrStr := string(stdout.Bytes()), string(stderr.Bytes())

	if cmderr != nil {
		lib.PrFatalf("\"%s\" failed with %s\n%s\n%s\n", cmdStr, cmderr, stdoutStr, stderrStr)
	}

	if stderrStr != "" {
		if ignore {
			lib.PrNoticeF(stderrStr)
		} else {
			lib.PrFatalf("%s\n", stderrStr)
		}
	}

	return stdoutStr
}

// SdfCreateAccountCustomizationProject create an sdf account customization project
func SdfCreateAccountCustomizationProject(name string, path string) Project {
	return sdfCreateProject("1", ProjectParams{
		Name: name,
		Path: path,
	})
}

// SdfCreateSuiteAppProject create an sdf suite app project
func SdfCreateSuiteAppProject(name string, path string, id string, version string, publisherId string) Project {
	return sdfCreateProject("2", ProjectParams{
		Name:        name,
		Path:        path,
		ID:          id,
		Version:     version,
		PublisherID: publisherId,
	})
}

// create an sdf project
func sdfCreateProject(kind string, params ProjectParams) Project {
	sequence := []string{}
	fileBaseSuffix := ""
	projectName := ""
	projectType := ""
	switch kind {
	case "1":
		sequence = append(sequence, params.Name, "")
		fileBaseSuffix = filepath.Join("FileCabinet", "SuiteScripts")
		projectType = "ACCOUNTCUSTOMIZATION"
		projectName = params.Name
	case "2":
		sequence = append(sequence, params.PublisherID, params.ID, params.Name, params.Version, "")
		fileBaseSuffix = filepath.Join("FileCabinet", "SuiteApps")
		projectType = "SUITEAPP"
		projectName = strings.Join([]string{params.PublisherID, params.ID}, ".")
	default:
		lib.PrFatalf("Project type has to be either \"1\" or \"2\"!\n")
	}

	s := strings.Join(sequence, "\n")
	prompt := strings.Join([]string{kind, s}, " ")

	os.Remove(filepath.Clean(filepath.Join(params.Path, projectName)))

	if lib.IsVerbose {
		lib.PrNoticeF("Creating project\t\"%s %s\"\n", projectType, projectName)
	}
	execute("sh", lib.SdfCliCreateProject, "", prompt, params.Path, false)

	dir := filepath.Join(params.Path, projectName)

	return Project{
		Type:     projectType,
		Dir:      dir,
		Filebase: filepath.Clean(path.Join(dir, fileBaseSuffix)),
		Name:     projectName,
		Params:   params,
	}
}
