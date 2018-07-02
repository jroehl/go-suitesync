package lib

import (
	"strings"
)

const (
	TokenID        = "NSCONF_TOKEN_ID"
	TokenSecret    = "NSCONF_TOKEN_SECRET"
	ConsumerKey    = "NSCONF_CONSUMER_KEY"
	ConsumerSecret = "NSCONF_CONSUMER_SECRET"
	Password       = "NSCONF_PASSWORD"
	Account        = "NSCONF_ACCOUNT"
	Email          = "NSCONF_EMAIL"
	Realm          = "NSCONF_REALM"
	CliToken       = "NSCONF_CLITOKEN"
	RootPath       = "NSCONF_ROOTPATH"
	HashFile       = "NSCONF_HASHFILE"
	Role           = "NSCONF_ROLE"
	URL            = "NSCONF_URL"

	RestletScriptID         = "customscript_suitesync_restlet"
	RestletScriptDeployment = "1"

	SdfCliConsumerKey    = "517e56cc85d90498a49eafea76a3f59e8933957bdf3fc69f7c1280d023d7b4e9175f76aaa3c4b44d4d8a16dcd04fb98087753aa1d4098338ab8b09b99716d9647b10d12e05eca4c487cc65cc96c7b22c"
	SdfCLiConsumerSecret = "b086f032d8b869915d38d192157a023442d7b3cfac14670b2731800120da3c70475a411f8ce0061a22a35b14cf257c26767374a6a8acb8ad2a24bf19fa80873c7b10d12e05eca4c487cc65cc96c7b22c"
	SdfCliPw             = "&Get File List"

	JavaVersion       = "8u171"
	JavaBuildNumber   = "11"
	JavaPlatformLinux = "linux-x64.tar.gz"
	JavaSubDirLinux   = ""
	JavaPlatformMac   = "macosx-x64.tar.gz"
	JavaSubDirMac     = "/Contents/Home"

	// dependency urls
	URLSdfCore         = "https://system.netsuite.com/download/ide/update_18_1/plugins/com.netsuite.ide.core_2018.1.2.jar"
	URLSdfIde          = "https://system.netsuite.com/download/ide/update_18_1/plugins/com.netsuite.ide.eclipse.ws_2018.1.2.jar"
	URLSdfSupplemental = "https://system.netsuite.com/core/media/media.nl?id=95083164&c=NLCORP&h=37e6a602c5c4fc0fb3e3&_xt=.gz"
	URLMaven           = "http://artfiles.org/apache.org/maven/maven-3/3.5.3/binaries/apache-maven-3.5.3-bin.tar.gz"
)

// global variables
var (
	Credentials map[string]string
	// paths to deps
	CurrentDir   string
	Dependencies string
	SdfCli       string
	CliCache     string
	RestletTar   string
	// IsVerbose variable for export
	IsVerbose = false

	// Whitelisted filenames hat are not included while uploading the files
	Whitelisted = []string{"error.log"}
	// constructed java base url
	JavaBaseURL = strings.Join([]string{"http://download.oracle.com/otn-pub/java/jdk/", JavaVersion, "-b", JavaBuildNumber, "/512cd62ec5174c3487ac17c61aaa89e8/jre-", JavaVersion, "-"}, "")
)

type Responses struct {
	Successful   []Response `json:"successful"`
	Unsuccessful []Response `json:"unsuccessful"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Path    string `json:"path"`
	Status  string `json:"status"`
	ID      string `json:"id"`
	Type    string `json:"type"`
}
