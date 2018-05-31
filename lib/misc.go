package lib

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// Difference of slice1 to slice2
func Difference(slice1 []string, slice2 []string) (dStr []string, dInt []int) {
	for i, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			dStr = append(dStr, s1)
			dInt = append(dInt, i)
		}
	}

	return dStr, dInt
}

// check if all needed params exist and gather env variables
func GetCredentials() Credentials {
	godotenv.Load()

	a := os.Getenv("NSCONF_ACCOUNT")
	if a == "" {
		PrFatalf("Environment variable \"NSCONF_ACCOUNT\" is undefined\n")
	}
	e := os.Getenv("NSCONF_EMAIL")
	if e == "" {
		PrFatalf("Environment variable \"NSCONF_EMAIL\" is undefined\n")
	}
	ti := os.Getenv("NSCONF_TOKEN_ID")
	if ti == "" {
		PrFatalf("Environment variable \"NSCONF_TOKEN_ID\" is undefined\n")
	}
	ts := os.Getenv("NSCONF_TOKEN_SECRET")
	if ts == "" {
		PrFatalf("Environment variable \"NSCONF_TOKEN_SECRET\" is undefined\n")
	}
	ck := os.Getenv("NSCONF_CONSUMER_KEY")
	if ck == "" {
		PrFatalf("Environment variable \"NSCONF_CONSUMER_KEY\" is undefined\n")
	}
	cs := os.Getenv("NSCONF_CONSUMER_SECRET")
	if cs == "" {
		PrFatalf("Environment variable \"NSCONF_CONSUMER_SECRET\" is undefined\n")
	}
	pw := os.Getenv("NSCONF_PASSWORD")
	ct := os.Getenv("NSCONF_CLITOKEN")
	if ct == "" {
		PrFatalf("Environment variable \"NSCONF_CLITOKEN\" is undefined\n")
	}
	r := os.Getenv("NSCONF_REALM")
	if r == "" {
		PrFatalf("Environment variable \"NSCONF_REALM\" is undefined\n")
	}
	rt := os.Getenv("NSCONF_ROOTPATH")
	if rt == "" {
		rt = "/SuiteScripts"
	}
	hf := os.Getenv("NSCONF_HASHFILE")
	if hf == "" {
		hf = "hashes.json"
	}
	s := os.Getenv("NSCONF_SCRIPT")
	if s == "" {
		s = "customscript_node_suitesync_restlet"
	}
	d := os.Getenv("NSCONF_DEPLOYMENT")
	if d == "" {
		d = "1"
	}
	ro := os.Getenv("NSCONF_ROLE")
	if ro == "" {
		ro = "3"
	}

	u := r
	if !strings.HasPrefix(r, "system.") {
		u = strings.Join([]string{"system.", r}, "")
	}

	return Credentials{
		Account:        a,
		Email:          e,
		Realm:          r,
		Rootpath:       rt,
		Script:         s,
		Deployment:     d,
		Role:           ro,
		Hashfile:       hf,
		TokenID:        ti,
		TokenSecret:    ts,
		ConsumerKey:    ck,
		ConsumerSecret: cs,
		Password:       pw,
		Url:            u,
		CliToken:       ct,
	}
}

// mk temporary dir
func MkTempDir() string {
	dir, err := ioutil.TempDir("/var/tmp", "suitesync")
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// check if array includes
func ArrayIncludes(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// remove directory
func Remove(p string) {
	err := os.RemoveAll(p)
	if err != nil {
		PrFatalf("Error removing %s\n%s", p, err)
	}
}

// convert interface to json
func ToJson(p interface{}) []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func CheckDir(p string) {
	s, err := CheckExists(p)
	if err != nil {
		PrFatalf("CheckDir error %s", err)
	}
	if !s.IsDir() {
		PrFatalf("\"%s\" is not a directory", p)
	}
}

// exists returns whether the given file or directory exists or not
func CheckExists(p string) (os.FileInfo, error) {
	s, err := os.Stat(p)
	if os.IsNotExist(err) {
		PrFatalf("\"%s\" does not exist", p)
		return s, nil
	}
	return s, err
}
