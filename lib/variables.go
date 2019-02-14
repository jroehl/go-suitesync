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
	HashFile       = "NSCONF_HASHFILE"
	Role           = "NSCONF_ROLE"
	URL            = "NSCONF_URL"

	SdfCliConsumerKey    = "517e56cc85d90498a49eafea76a3f59e8933957bdf3fc69f7c1280d023d7b4e9175f76aaa3c4b44d4d8a16dcd04fb98087753aa1d4098338ab8b09b99716d9647b10d12e05eca4c487cc65cc96c7b22c"
	SdfCLiConsumerSecret = "b086f032d8b869915d38d192157a023442d7b3cfac14670b2731800120da3c70475a411f8ce0061a22a35b14cf257c26767374a6a8acb8ad2a24bf19fa80873c7b10d12e05eca4c487cc65cc96c7b22c"
	SdfCliPw             = "&Get File List"

	JavaVersion       = "8u181"
	JavaBuildNumber   = "13"
	JavaPlatformLinux = "linux-x64.tar.gz"
	JavaSubDirLinux   = ""
	JavaPlatformMac   = "macosx-x64.tar.gz"
	JavaSubDirMac     = "/Contents/Home"

	// URLSdfCore dependency urls
	URLSdfCore         = "https://system.netsuite.com/download/ide/update_18_1/plugins/com.netsuite.ide.core_2018.1.2.jar"
	URLSdfIde          = "https://system.netsuite.com/download/ide/update_18_1/plugins/com.netsuite.ide.eclipse.ws_2018.1.2.jar"
	URLSdfSupplemental = "https://system.netsuite.com/core/media/media.nl?id=95083164&c=NLCORP&h=37e6a602c5c4fc0fb3e3&_xt=.gz"
	URLMaven           = "http://artfiles.org/apache.org/maven/maven-3/3.5.4/binaries/apache-maven-3.5.4-bin.tar.gz"
)

// global variables
var (
	Credentials map[string]string
	// paths to deps
	CurrentDir   string
	Dependencies string
	SdfCli       string
	CliCache     string
	// IsVerbose variable for export
	IsVerbose = false
	IsDebug   = false

	// Whitelisted filenames hat are not included while uploading the files
	Whitelisted = []string{"error.log"}
	// constructed java base url
	JavaBaseURL = strings.Join([]string{"http://download.oracle.com/otn-pub/java/jdk/", JavaVersion, "-b", JavaBuildNumber, "/96a7b8442fe848ef90c96a2fad6ed6d1/jre-", JavaVersion, "-"}, "")
)

// DeleteResult struct
type DeleteResult struct {
	Successful bool
	NotFound   bool
	Code       string
	Message    string
	ID         string
	Type       string
}

// Meta struct
type Meta struct {
	Successful   bool
	TotalPages   int
	TotalRecords int
	SearchID     string
}

// SearchResult struct
type SearchResult struct {
	InternalID string
	Parent     string
	Name       string
	Children   []*SearchResult
	IsDir      bool
	Path       string
}

type Attr struct {
	Key   string
	Value string
}

type SearchValue struct {
	Inner string
	Attrs []Attr
}
type SearchFilter struct {
	Tag          string
	Operator     string
	SearchValues []SearchValue
}
