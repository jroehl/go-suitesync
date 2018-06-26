//+build !test

package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/mholt/archiver"
)

// InitEnv init the environment used by the cli
func InitEnv(skipAuthReq bool) (err error) {
	creds, err := sanitizeCredentials(skipAuthReq)
	if err != nil {
		return err
	}
	Credentials = creds

	Dependencies = path.Join(CurrentDir, ".dependencies")
	SdfCli = path.Join(Dependencies, "sdfcli")
	CliCache = path.Join(Dependencies, ".cliCache")

	if _, err := os.Stat(Dependencies); os.IsNotExist(err) {
		err = initDependencies()
	}

	if _, err := os.Stat(SdfCli); os.IsNotExist(err) {
		Remove(Dependencies)
		err = initDependencies()
	}
	return err
}

func initDependencies() error {
	PrNoticef("\nSetting up wrapped sdfcli\n\n")
	os.MkdirAll(Dependencies, os.ModePerm)
	var p, s string

	files := []string{}
	javaExists := os.Getenv("JAVA_HOME") != "" && isCommandAvailable("java")
	if !javaExists {
		switch runtime.GOOS {
		case "darwin":
			p = JavaPlatformMac
			s = JavaSubDirMac
		case "linux":
			p = JavaPlatformLinux
			s = JavaSubDirLinux
		default:
			return fmt.Errorf("only \"MacOS\" and \"Linux\" are supported - not \"%s\"", runtime.GOOS)
		}
		files = append(files, downloadFile(Dependencies, strings.Join([]string{JavaBaseURL, p}, ""), true, false))
	}
	mavenExists := isCommandAvailable("mvn")

	files = append(files, downloadDependencies(!mavenExists)...)

	PrNoticef("\nExtracting %v files\n", len(files))

	for _, f := range files {
		if strings.HasSuffix(f, ".jar") {
			archiver.Zip.Open(f, Dependencies)
		} else if strings.HasSuffix(f, ".tar.gz") {
			archiver.TarGz.Open(f, Dependencies)
		}
	}

	if !javaExists {
		javaDir := FindDir(Dependencies, "jre*.*")[0]
		javaHome := path.Join(javaDir, s)
		sed(SdfCli, "mvn", strings.Join([]string{"JAVA_HOME=", javaHome, " mvn"}, ""))
	}

	if !mavenExists {
		mavenDir := FindDir(Dependencies, "-maven*")[0]
		mavenBin := path.Join(mavenDir, "bin", "mvn")
		sed(SdfCli, "mvn", mavenBin)
	}

	sed(SdfCli, "/webdev/sdf/sdk/", path.Join(Dependencies))

	cmd := exec.Command(path.Join(SdfCli))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sdfcli setup failed\n%s", err.Error())
	}

	CheckCliCache()
	PrNoticef("\nSetup completed\n\n")
	return nil
}

func downloadDependencies(downloadMaven bool) (files []string) {
	files = []string{
		downloadFile(Dependencies, URLSdfCore, false, false),
		downloadFile(Dependencies, URLSdfIde, false, false),
		downloadFile(path.Join(Dependencies, "sdfcli-supplemental_18_1.tar.gz"), URLSdfSupplemental, false, true),
	}
	if downloadMaven {
		files = append(files, downloadFile(Dependencies, URLMaven, false, false))
	}
	return
}
