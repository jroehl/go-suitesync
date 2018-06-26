package sdf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	Type        string
	Dir         string
	FileCabinet string
	Name        string
	Params      ProjectParams
}

type FileTransfer struct {
	Root string
	Path string
	Dest string
	Src  string
}

type BashExec interface {
	Command(name string, arg ...string) *exec.Cmd
}

// GenerateToken sdf cli token
func GenerateToken(bash BashExec, password string) (res string, err error) {
	token := lib.CheckCliCache()
	if string(token) != "" {
		return "", fmt.Errorf("Clitoken seems to be set up, aborting\n\"%s\"", token)
	}
	f, _ := buildFlags(nil)
	execute(bash, lib.SdfCli, "issuetoken", f, strings.Join([]string{password, "\n"}, ""), "", false)
	con, _ := ioutil.ReadFile(lib.CliCache)
	res = string(con)
	lib.PrResultf("\nToken\n%s\n", res)
	return res, nil
}

// Command call sdfcli command
func Command(bash BashExec, command string, flags []Flag, ignore bool) string {
	f, err := buildFlags(flags)
	if err != nil {
		panic(err)
	}
	res := execute(bash, lib.SdfCli, command, f, "YES\n", "", ignore)
	return res
}

func execute(bash BashExec, bin string, cmd string, flags string, prompt string, dir string, ignore bool) string {
	cmdStr := strings.Join([]string{bin, cmd, flags}, " ")

	if lib.IsVerbose {
		lib.PrHeaderf("\nExecuted command:\n")
		fmt.Printf("%s\n\n", cmdStr)
	}

	proc := bash.Command(bin, cmd, flags)

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
		lib.PrFatalf("\n\"%s\" failed with %s\n%s\n%s\n", cmdStr, cmderr, stdoutStr, stderrStr)
	}

	if stderrStr != "" {
		if ignore {
			lib.PrNoticef("\n%s\n", stderrStr)
		} else {
			lib.PrFatalf("\n%s\n", stderrStr)
		}
	}

	return stdoutStr
}

// CreateAccountCustomizationProject create an sdf account customization project
func CreateAccountCustomizationProject(name string, path string) Project {
	return sdfCreateProject("1", ProjectParams{
		Name: name,
		Path: path,
	})
}

// CreateSuiteAppProject create an sdf suite app project
func CreateSuiteAppProject(name string, path string, id string, version string, publisherId string) Project {
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
	fileBaseSuffix := ""
	projectName := ""
	projectType := ""
	deployXML := ""
	manifestXML := ""

	var dir string
	switch kind {
	case "1":
		fileBaseSuffix = filepath.Join("FileCabinet", "SuiteScripts")
		projectType = "ACCOUNTCUSTOMIZATION"
		projectName = params.Name
		dir = filepath.Join(params.Path, projectName)
		os.MkdirAll(filepath.Join(dir, "AccountConfiguration"), os.ModePerm)
		os.MkdirAll(filepath.Join(dir, fileBaseSuffix, ".attributes"), os.ModePerm)
		os.MkdirAll(filepath.Join(dir, "FileCabinet", "Templates", "Marketing Templates"), os.ModePerm)
		os.MkdirAll(filepath.Join(dir, "FileCabinet", "Templates", "E-mail Templates"), os.ModePerm)

		deployXML = `
			<deploy>
				<configuration>
						<path>~/AccountConfiguration/*</path>
				</configuration>
				<files>
						<path>~/FileCabinet/SuiteScripts/*</path>
				</files>
				<objects>
						<path>~/Objects/*</path>
				</objects>
			</deploy>
		`
		manifestXML = fmt.Sprintf(`
			<manifest projecttype="ACCOUNTCUSTOMIZATION">
					<projectname>%s</projectname>
					<frameworkversion>1.0</frameworkversion>
					<dependencies>
							<features>
									<feature required="true">CUSTOMRECORDS</feature>
									<feature required="true">SERVERSIDESCRIPTING</feature>
									<feature required="false">CREATESUITEBUNDLES</feature>
							</features>
					</dependencies>
			</manifest>
		`, params.Name)

	case "2":
		fileBaseSuffix = filepath.Join("FileCabinet", "SuiteApps")
		projectType = "SUITEAPP"
		projectName = strings.Join([]string{params.PublisherID, params.ID}, ".")
		dir = filepath.Join(params.Path, projectName)
		os.MkdirAll(filepath.Join(dir, "InstallationPreferences"), os.ModePerm)
		os.MkdirAll(filepath.Join(dir, fileBaseSuffix, projectName, ".attributes"), os.ModePerm)

		deployXML = fmt.Sprintf(`
			<deploy>
				<files>
						<path>~/FileCabinet/SuiteApps/%s/*</path>
				</files>
				<objects>
						<path>~/Objects/*</path>
				</objects>
			</deploy>
		`, projectName)

		manifestXML = fmt.Sprintf(`
			<manifest projecttype="SUITEAPP">
					<publisherid>%s</publisherid>
					<projectid>%s</projectid>
					<projectname>%s</projectname>
					<projectversion>%s<projectversion>
					<frameworkversion>1.0</frameworkversion>
					<dependencies>
							<features>
									<feature required="true">CUSTOMRECORDS</feature>
									<feature required="true">SERVERSIDESCRIPTING</feature>
									<feature required="false">CREATESUITEBUNDLES</feature>
							</features>
					</dependencies>
			</manifest>
		`, params.PublisherID, params.ID, params.Name, params.Version)

	default:
		lib.PrFatalf("Project type has to be either \"1\" or \"2\"!\n")
	}

	os.Remove(filepath.Clean(filepath.Join(params.Path, projectName)))

	os.MkdirAll(filepath.Join(dir, "Objects"), os.ModePerm)
	ioutil.WriteFile(filepath.Join(dir, "deploy.xml"), []byte(deployXML), os.ModePerm)
	ioutil.WriteFile(filepath.Join(dir, "manifest.xml"), []byte(manifestXML), os.ModePerm)

	if lib.IsVerbose {
		lib.PrNoticef("Creating project \"%s %s\"\n", projectType, projectName)
	}

	return Project{
		Type:        projectType,
		Dir:         dir,
		FileCabinet: filepath.Clean(filepath.Join(dir, "FileCabinet")),
		Name:        projectName,
		Params:      params,
	}
}
