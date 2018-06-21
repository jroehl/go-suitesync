package lib

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mholt/archiver"
)

func InitEnv() {
	Credentials = sanitizeCredentials()

	Dependencies = path.Join(CurrentDir, ".dependencies")
	SdfCli = path.Join(Dependencies, "sdfcli")
	CliCache = path.Join(Dependencies, ".clicache")
	RestletTar = path.Join(Dependencies, "restlet.tar.gz")

	if _, err := os.Stat(Dependencies); os.IsNotExist(err) {
		initDependencies()
	}

	if _, err := os.Stat(SdfCli); os.IsNotExist(err) {
		Remove(Dependencies)
		initDependencies()
	}

}

func initDependencies() {
	PrNoticeF("\nSetting up wrapped sdfcli\n\n")
	os.MkdirAll(Dependencies, os.ModePerm)
	var p, s string

	rp := path.Join(CurrentDir, "restlet.tar.gz")
	if err := Copy(rp, path.Join(Dependencies, "restlet.tar.gz")); err != nil {
		log.Fatal(err)
	}

	files := []string{}
	javaExists := os.Getenv("JAVA_HOME") != ""
	if !javaExists {
		switch runtime.GOOS {
		case "darwin":
			p = JavaPlatformMac
			s = JavaSubDirMac
		case "linux":
			p = JavaPlatformLinux
			s = JavaSubDirLinux
		default:
			PrFatalf("Only \"MacOS\" and \"Linux\" are supported - not \"%s\"", runtime.GOOS)
		}
		files = append(files, downloadFile(Dependencies, strings.Join([]string{JavaBaseURL, p}, ""), true, false))
	}
	mavenExists := isCommandAvailable("mvn")

	files = append(files, downloadDependencies(!mavenExists)...)

	PrNoticeF("\nExtracting %v files\n", len(files))

	for _, f := range files {
		if strings.HasSuffix(f, ".jar") {
			archiver.Zip.Open(f, Dependencies)
		} else if strings.HasSuffix(f, ".tar.gz") {
			archiver.TarGz.Open(f, Dependencies)
		}
	}

	if !javaExists {
		javaDir := find(Dependencies, "jre*.*")[0]
		javaHome := path.Join(javaDir, s)
		sed(SdfCli, "mvn", strings.Join([]string{"JAVA_HOME=", javaHome, " mvn"}, ""))
	}

	if !mavenExists {
		mavenDir := find(Dependencies, "-maven*")[0]
		mavenBin := path.Join(mavenDir, "bin", "mvn")
		sed(SdfCli, "mvn", mavenBin)
	}

	sed(SdfCli, "/webdev/sdf/sdk/", path.Join(Dependencies))

	cmd := exec.Command(path.Join(SdfCli))
	if err := cmd.Run(); err != nil {
		PrFatalf("sdfcli setup failed\n%s\n", err.Error())
	}

	checkCliCache()
	PrNoticeF("\nSetup completed\n\n")
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

type EnvVar struct {
	env string
	def string
}

// check if all needed params exist and gather env variables
func sanitizeCredentials() (creds map[string]string) {
	godotenv.Load()
	creds = make(map[string]string)

	optionals := []string{TokenID, TokenSecret, ConsumerKey, ConsumerSecret, Password, CliToken}
	mandatories := []string{Account, Email, Realm}
	defaults := []EnvVar{
		EnvVar{env: RootPath, def: "/SuiteScripts"},
		EnvVar{env: HashFile, def: "hashes.json"},
		EnvVar{env: Role, def: "3"},
	}

	for _, o := range optionals {
		creds[o] = os.Getenv(o)
	}

	for _, m := range mandatories {
		creds[m] = os.Getenv(m)
		if creds[m] == "" {
			PrFatalf("Environment variable \"%s\" is undefined\n", m)
		}
	}

	for _, d := range defaults {
		creds[d.env] = os.Getenv(d.env)
		if creds[d.env] == "" {
			creds[d.env] = d.def
		}
	}

	// prefix realm with "system."
	creds[URL] = creds[Realm]
	if !strings.HasPrefix(creds[URL], "system.") {
		creds[URL] = strings.Join([]string{"system.", creds[URL]}, "")
	}

	currentKey := strings.Join([]string{creds[URL], creds[Account], creds[Email], creds[Role]}, "-")
	if creds[TokenID] != "" && creds[TokenSecret] != "" && creds[ConsumerKey] != "" && creds[ConsumerSecret] != "" {
		PrNoticeF("Using \"%s\", \"%s\", \"%s\" and \"%s\" for authentication\n", ConsumerKey, ConsumerSecret, TokenID, TokenSecret)
		EncryptCliToken(currentKey, creds)
	} else if creds[CliToken] != "" {
		PrNoticeF("Using \"%s\" for authentication\n", CliToken)
		if err := DecryptCliToken(currentKey, creds); err != nil {
			PrFatalf(err.Error())
		}
	} else {
		PrFatalf("Either \"%s\" or \"%s\", \"%s\", \"%s\" and \"%s\" env variables have to be defined\n", CliToken, ConsumerKey, ConsumerSecret, TokenID, TokenSecret)
	}

	return
}

func checkCliCache() []byte {
	b, err := ioutil.ReadFile(CliCache)
	if err != nil || string(b) == "" {
		// file does not exist or is empty
		b = []byte(Credentials[CliToken])
		ioutil.WriteFile(CliCache, b, os.ModePerm)
	}
	return b
}
