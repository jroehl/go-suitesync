package lib

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestIsCommandAvailableTrue(t *testing.T) {
	res := isCommandAvailable("bash")
	assert.True(t, res)
}

func TestIsCommandAvailableFalse(t *testing.T) {
	res := isCommandAvailable("foobar_")
	assert.False(t, res)
}

func TestFindDir(t *testing.T) {
	res := FindDir("../tests/fs", "sub.*")
	assert.NotZero(t, res)
	assert.Equal(t, "../tests/fs/subdir", res[0])
}

func TestSed(t *testing.T) {
	f := "../tests/t"
	ioutil.WriteFile(f, []byte("foobar"), os.ModePerm)
	sed(f, "foo", "<FOO")
	b, _ := ioutil.ReadFile(f)
	os.Remove(f)
	assert.Equal(t, string(b), "<FOObar")
}

func TestOptDest(t *testing.T) {
	res := OptDest("foobar", "default")
	assert.Equal(t, "foobar", res)
}

func TestOptDestDef(t *testing.T) {
	res := OptDest("", "default")
	assert.Equal(t, "default", res)
}

func TestOptDestWd(t *testing.T) {
	res := OptDest("", "")
	assert.Contains(t, res, "/go-suitesync/lib")
}

func TestOptArgDest(t *testing.T) {
	c := cli.Args([]string{"foo", "bar"})
	res := OptArgDest(c, 1, "default")
	assert.Equal(t, "bar", res)
}

func TestOptArgDestDef(t *testing.T) {
	c := cli.Args([]string{"foo", "bar"})
	res := OptArgDest(c, 2, "default")
	assert.Equal(t, "default", res)
}

func TestCheckRequired(t *testing.T) {
	c := []string{"foo", "bar"}
	res := CheckRequired(c, 1, "arg")
	assert.Equal(t, "bar", res)
}

func TestCheckExistsFolder(t *testing.T) {
	res, err := CheckExists("../tests/fs/subdir")
	assert.True(t, res.IsDir())
	assert.Equal(t, "subdir", res.Name())
	assert.Nil(t, err)
}

func TestCheckExistsFile(t *testing.T) {
	res, err := CheckExists("../tests/fs/Bar Foo")
	assert.False(t, res.IsDir())
	assert.Equal(t, "Bar Foo", res.Name())
	assert.Nil(t, err)
}

func TestCheckExistsNot(t *testing.T) {
	res, err := CheckExists("../tests/fs/Snoo")
	assert.EqualError(t, err, "\"../tests/fs/Snoo\" does not exist")
	assert.Nil(t, res)
}

func TestCheckDir(t *testing.T) {
	err := CheckDir("../tests/fs/subdir")
	assert.Nil(t, err)
}

func TestCheckDirErr(t *testing.T) {
	err := CheckDir("../tests/fs/Bar Foo")
	assert.EqualError(t, err, "\"../tests/fs/Bar Foo\" is not a directory")
}

func TestCheckDirNotExist(t *testing.T) {
	err := CheckDir("../tests/fs/Snoo")
	assert.EqualError(t, err, "\"../tests/fs/Snoo\" does not exist")
}

func TestToJSON(t *testing.T) {
	h := []Hash{Hash{Name: "name", Path: "path", Hash: "hash"}}
	b := ToJSON(h)
	assert.Equal(t, "[{\"Path\":\"path\",\"Hash\":\"hash\",\"Name\":\"name\"}]", string(b))
}

func TestToJSONEmpty(t *testing.T) {
	h := []Hash{Hash{}}
	b := ToJSON(h)
	assert.Equal(t, "[{\"Path\":\"\",\"Hash\":\"\",\"Name\":\"\"}]", string(b))
}

func TestMkTempDir(t *testing.T) {
	dir := MkTempDir()
	f, _ := CheckExists(dir)
	assert.True(t, f.IsDir())
	assert.Nil(t, CheckDir(dir))
	os.Remove(dir)
	assert.NotNil(t, CheckDir(dir))
}

func TestRemove(t *testing.T) {
	dir := MkTempDir()
	ioutil.WriteFile(path.Join(dir, "testfile"), []byte("data"), os.ModePerm)
	assert.Nil(t, CheckDir(dir))
	Remove(dir)
	assert.NotNil(t, CheckDir(dir))
}

func TestArrayIncludes(t *testing.T) {
	a := []string{"foo", "bar"}
	s := "foo"
	assert.True(t, ArrayIncludes(s, a))
}

func TestArrayIncludesNot(t *testing.T) {
	a := []string{"foo", "bar"}
	s := "foobar"
	assert.False(t, ArrayIncludes(s, a))
}

func TestDifference(t *testing.T) {
	a := []string{"foo", "bar", "baz"}
	b := []string{"foo", "bar"}
	difS, difI := Difference(a, b)
	assert.ElementsMatch(t, difS, []string{"baz"})
	assert.ElementsMatch(t, difI, []int{2})
}

func TestCopy(t *testing.T) {
	dir := MkTempDir()
	np := path.Join(dir, "Bar Foo")
	err := Copy("../tests/fs/Bar Foo", np)
	assert.Nil(t, err)
	f, err := CheckExists(np)
	assert.Nil(t, err)
	assert.False(t, f.IsDir())
	Remove(dir)
}

func TestCopyFailSrc(t *testing.T) {
	err := Copy("../tests/fs/B", "")
	assert.EqualError(t, err, "open ../tests/fs/B: no such file or directory")
}

func TestCopyFailDest(t *testing.T) {
	err := Copy("../tests/fs/Bar Foo", "../tests")
	assert.EqualError(t, err, "open ../tests: is a directory")
}

func TestCheckCliCache(t *testing.T) {
	CliCache = ".clicache"
	Credentials = map[string]string{}
	Credentials[CliToken] = "123456"
	res := CheckCliCache()
	assert.Equal(t, "123456", string(res))
	b, _ := ioutil.ReadFile(CliCache)
	assert.Equal(t, "123456", string(b))
	assert.NotNil(t, res)
	os.Remove(CliCache)
}

func TestCheckCliCacheNil(t *testing.T) {
	CliCache = ".clicache"
	Credentials = map[string]string{}
	Credentials[CliToken] = ""
	res := CheckCliCache()
	assert.Equal(t, "", string(res))
	assert.Zero(t, string(res))
	os.Remove(CliCache)
}

func TestSanitizeCredentials(t *testing.T) {
	os.Setenv(Account, "account")
	os.Setenv(Email, "email")
	os.Setenv(Realm, "realm")
	os.Setenv(CliToken, "b8d3596db885fd1c4be5cccccd89a78ebb97147409d0c8e1832014c7cd242e753158638f039f11bd27784269942e24e67b10d12e05eca4c487cc65cc96c7b22c")
	creds, err := sanitizeCredentials(false)
	assert.Equal(t, "3", creds[Role])
	assert.Equal(t, "1811dfc40afbf177e53a21eb8cd2d3b95ff6281a41fe803964ea7d60e6cefb3e", creds[ConsumerSecret])
	assert.Equal(t, "b8d3596db885fd1c4be5cccccd89a78ebb97147409d0c8e1832014c7cd242e753158638f039f11bd27784269942e24e67b10d12e05eca4c487cc65cc96c7b22c", creds[CliToken])
	assert.Equal(t, "email", creds[Email])
	assert.Equal(t, "account", creds[Account])
	assert.Equal(t, "realm", creds[Realm])
	assert.Equal(t, "hashes.json", creds[HashFile])
	assert.Equal(t, "system.realm", creds[URL])
	assert.Equal(t, "tokenid", creds[TokenID])
	assert.Equal(t, "tokensecret", creds[TokenSecret])
	assert.Equal(t, "6da57bf05a6247fc876c6d228184ff487760a382a43ac7e93eaff743803d22ac", creds[ConsumerKey])
	assert.Equal(t, "", creds[Password])
	assert.Nil(t, err)
}

func TestSanitizeCredentialsCK(t *testing.T) {
	IsVerbose = true
	os.Setenv(Account, "account")
	os.Setenv(Email, "email")
	os.Setenv(Realm, "realm")
	os.Setenv(ConsumerKey, "consumerKey")
	os.Setenv(ConsumerSecret, "consumerSecret")
	os.Setenv(TokenID, "tokenId")
	os.Setenv(TokenSecret, "tokenSecret")
	creds, err := sanitizeCredentials(false)
	assert.Equal(t, "3", creds[Role])
	assert.Equal(t, "consumerSecret", creds[ConsumerSecret])
	assert.Equal(t, "b8d3596db885fd1c4be5cccccd89a78ebb97147409d0c8e1832014c7cd242e75b3a05f774a062f02aa5a7fd090fe35d47b10d12e05eca4c487cc65cc96c7b22c", creds[CliToken])
	assert.Equal(t, "email", creds[Email])
	assert.Equal(t, "account", creds[Account])
	assert.Equal(t, "realm", creds[Realm])
	assert.Equal(t, "hashes.json", creds[HashFile])
	assert.Equal(t, "system.realm", creds[URL])
	assert.Equal(t, "tokenId", creds[TokenID])
	assert.Equal(t, "tokenSecret", creds[TokenSecret])
	assert.Equal(t, "consumerKey", creds[ConsumerKey])
	assert.Equal(t, "", creds[Password])
	assert.Nil(t, err)
}

func TestSanitizeCredentialsOpts(t *testing.T) {
	IsVerbose = true
	os.Setenv(Account, "account")
	os.Setenv(Email, "email")
	os.Setenv(Realm, "realm")
	os.Setenv(HashFile, "foobar.json")
	os.Setenv(ConsumerKey, "")
	os.Setenv(ConsumerSecret, "")
	os.Setenv(TokenID, "")
	os.Setenv(TokenSecret, "")
	os.Setenv(Role, "3")
	os.Setenv(Password, "password")
	os.Setenv(CliToken, "b8d3596db885fd1c4be5cccccd89a78ebb97147409d0c8e1832014c7cd242e753158638f039f11bd27784269942e24e67b10d12e05eca4c487cc65cc96c7b22c")
	creds, err := sanitizeCredentials(false)
	assert.Equal(t, "3", creds[Role])
	assert.Equal(t, "1811dfc40afbf177e53a21eb8cd2d3b95ff6281a41fe803964ea7d60e6cefb3e", creds[ConsumerSecret])
	assert.Equal(t, "b8d3596db885fd1c4be5cccccd89a78ebb97147409d0c8e1832014c7cd242e753158638f039f11bd27784269942e24e67b10d12e05eca4c487cc65cc96c7b22c", creds[CliToken])
	assert.Equal(t, "email", creds[Email])
	assert.Equal(t, "account", creds[Account])
	assert.Equal(t, "realm", creds[Realm])
	assert.Equal(t, "foobar.json", creds[HashFile])
	assert.Equal(t, "system.realm", creds[URL])
	assert.Equal(t, "tokenid", creds[TokenID])
	assert.Equal(t, "tokensecret", creds[TokenSecret])
	assert.Equal(t, "6da57bf05a6247fc876c6d228184ff487760a382a43ac7e93eaff743803d22ac", creds[ConsumerKey])
	assert.Equal(t, "password", creds[Password])
	assert.Nil(t, err)
}

func TestSanitizeCredentialsErrCreds(t *testing.T) {
	os.Setenv(Account, "account")
	os.Setenv(Email, "email")
	os.Setenv(CliToken, "")
	os.Setenv(Realm, "realm")
	os.Setenv(ConsumerKey, "")
	os.Setenv(ConsumerSecret, "")
	os.Setenv(TokenID, "")
	os.Setenv(TokenSecret, "")
	creds, err := sanitizeCredentials(false)
	assert.EqualError(t, err, "either \"NSCONF_CLITOKEN\" or \"NSCONF_CONSUMER_KEY\", \"NSCONF_CONSUMER_SECRET\", \"NSCONF_TOKEN_ID\" and \"NSCONF_TOKEN_SECRET\" env variables have to be defined")
	assert.Nil(t, creds)
}

func TestSanitizeCredentialsErrSkip(t *testing.T) {
	os.Setenv(Account, "account")
	os.Setenv(Email, "email")
	os.Setenv(CliToken, "")
	os.Setenv(Realm, "realm")
	os.Setenv(ConsumerKey, "")
	os.Setenv(ConsumerSecret, "")
	os.Setenv(TokenID, "")
	os.Setenv(TokenSecret, "")
	creds, err := sanitizeCredentials(true)
	assert.Nil(t, err)
	assert.NotNil(t, creds)
}

func TestSanitizeCredentialsErrEnv(t *testing.T) {
	os.Setenv(Account, "account")
	os.Setenv(Email, "email")
	os.Setenv(Realm, "")

	creds, err := sanitizeCredentials(false)
	assert.EqualError(t, err, "environment variable \"NSCONF_REALM\" is undefined")
	assert.Nil(t, creds)
}

func TestAbsolutePath1(t *testing.T) {
	pwd, _ := os.Getwd()
	p := AbsolutePath("")
	assert.Equal(t, pwd, p)
}

func TestAbsolutePath2(t *testing.T) {
	pwd, _ := os.Getwd()
	p := AbsolutePath("./")
	assert.Equal(t, pwd, p)
}

func TestAbsolutePath3(t *testing.T) {
	pwd, _ := os.Getwd()
	p := AbsolutePath(".dependencies")
	assert.Equal(t, filepath.Join(pwd, ".dependencies"), p)
}
