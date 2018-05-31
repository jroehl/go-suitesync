package lib

// global variables
var (
	// Credentials for export
	Creds Credentials
	// paths to deps
	Dependencies        string
	SdfCli              string
	SdfCliCreateProject string
	CliCache            string
	Restlet             string
	// IsVerbose variable for export
	IsVerbose = false
	// YesRequired - methods that require the YES prompt
	YesRequired = []string{
		"importfiles",
	}
)

type Credentials struct {
	Account        string
	Email          string
	Password       string
	Realm          string
	Rootpath       string
	Script         string
	Deployment     string
	Role           string
	Hashfile       string
	TokenID        string
	TokenSecret    string
	ConsumerKey    string
	ConsumerSecret string
	Url            string
	CliToken       string
}

type Responses struct {
	Successful   []Response `json:"successful"`
	Unsuccessful []Response `json:"unsuccessful"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Path    string `json:"path"`
	Status  string `json:"status"`
	Id      string `json:id,omitempty`
	Type    string `json:type,omitempty`
}
