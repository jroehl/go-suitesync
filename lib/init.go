package lib

import (
	"log"
	"os"
	"os/exec"
	"path"
)

func InitEnv() {
	Creds = GetCredentials()

	Dependencies = path.Join(".dependencies")
	SdfCli = path.Join(Dependencies, "sdfcli")
	SdfCliCreateProject = path.Join(Dependencies, "sdfcli-createproject")
	CliCache = path.Join(Dependencies, ".clicache")
	Restlet = path.Join(Dependencies, "restlet", "project")

	if _, err := os.Stat(Dependencies); os.IsNotExist(err) {
		PrNoticeF("No \".dependencies\" directory found")
		initDependencies()
	}

	if _, err := os.Stat(SdfCli); os.IsNotExist(err) {
		log.Fatal(err)
	}

	if _, err := os.Stat(SdfCliCreateProject); os.IsNotExist(err) {
		log.Fatal(err)
	}
}

func initDependencies() {
	cmd := exec.Command("./init.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
