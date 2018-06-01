package lib

import (
	"log"
	"os"
	"path"
	"runtime"

	config "github.com/spf13/viper"
)

func InitConfig() {
	// set path and type
	config.AddConfigPath("./config")
	config.SetConfigType("yaml")

	// load the default config
	defaultConfig := "config.default"
	config.SetConfigName(defaultConfig)

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error while loading %s: %+v\n", defaultConfig, err)
	}

}

func InitEnv() {
	Creds = GetCredentials()

	Dependencies = path.Join(".dependencies")
	SdfCli = path.Join(Dependencies, "sdfcli")
	SdfCliCreateProject = path.Join(Dependencies, "sdfcli-createproject")
	CliCache = path.Join(Dependencies, ".clicache")
	Restlet = path.Join("restlet", "project")

	if _, err := os.Stat(SdfCli); os.IsNotExist(err) {
		log.Fatal(err)
	}

	if _, err := os.Stat(SdfCliCreateProject); os.IsNotExist(err) {
		log.Fatal(err)
	}
}

func InitDependencies() {
	platform := ""
	javaSubdir := ""
	switch runtime.GOOS {
	case "darwin":
		platform = "linux-x64.tar.gz"
	case "linux":
		platform = "macosx-x64.tar.gz"
		javaSubdir = "/Contents/Home"
	default:
		log.Fatalf("Only \"MacOS\" and \"Linux\" are supported - not \"%s\"", runtime.GOOS)
	}

	urls := config.GetStringMapString("urls")
	PrNoticeF("Downloading dependencies")

	for k, v := range urls {
		if err := DownloadFile(k, v); err != nil {
			log.Fatal(err)
		}
	}

	// PrNoticeF("Download finished")

	log.Fatal("foobar")
}
