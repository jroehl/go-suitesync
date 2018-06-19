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
	if lib.Credentials[lib.Password] == "" {
		lib.PrFatalf("Environment variable \"%s\" is undefined\n", lib.Password)
	}
	f, err := buildFlags(nil)
	if err != nil {
		log.Fatal(err)
	}
	res := execute(lib.SdfCli, "issuetoken", f, strings.Join([]string{lib.Credentials[lib.Password], "\n"}, ""), "", false)
	lib.PrNoticeF(res)
}

// call sdfcli command
func Sdf(command string, flags []Flag, ignore bool) string {
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
		lib.PrHeaderF("Executed command:\n")
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
			lib.PrNoticeF("%s\n", stderrStr)
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
	fileBaseSuffix := ""
	projectName := ""
	projectType := ""
	deployXML := ""
	manifestXML := ""

	os.MkdirAll(path.Join(params.Path, projectName, "Objects"), os.ModePerm)
	switch kind {
	case "1":
		fileBaseSuffix = filepath.Join("FileCabinet", "SuiteScripts")
		projectType = "ACCOUNTCUSTOMIZATION"
		projectName = params.Name
		os.MkdirAll(path.Join(params.Path, projectName, "AccountConfiguration"), os.ModePerm)
		os.MkdirAll(path.Join(params.Path, projectName, fileBaseSuffix, ".attributes"), os.ModePerm)
		os.MkdirAll(path.Join(params.Path, projectName, fileBaseSuffix, ".attributes"), os.ModePerm)
		os.MkdirAll(path.Join(params.Path, projectName, "FileCabinet", "Templates", "Marketing Templates"), os.ModePerm)
		os.MkdirAll(path.Join(params.Path, projectName, "FileCabinet", "Templates", "E-mail Templates"), os.ModePerm)

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

		os.MkdirAll(path.Join(params.Path, projectName, "InstallationPreferences"), os.ModePerm)
		os.MkdirAll(path.Join(params.Path, projectName, fileBaseSuffix, projectName, ".attributes"), os.ModePerm)

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

	ioutil.WriteFile(filepath.Join(params.Path, projectName, "deploy.xml"), []byte(deployXML), os.ModePerm)
	ioutil.WriteFile(filepath.Join(params.Path, projectName, "manifest.xml"), []byte(manifestXML), os.ModePerm)

	if lib.IsVerbose {
		lib.PrNoticeF("Creating project\t\"%s %s\"\n", projectType, projectName)
	}

	dir := filepath.Join(params.Path, projectName)

	return Project{
		Type:     projectType,
		Dir:      dir,
		Filebase: filepath.Clean(path.Join(dir, fileBaseSuffix)),
		Name:     projectName,
		Params:   params,
	}
}
