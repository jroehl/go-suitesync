package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

// EnvVar struct for env and default values
type EnvVar struct {
	env string
	def string
}

// Copy the src file to dst. Any existing file will be overwritten
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

// MkTempDir create temporary dir
func MkTempDir() string {
	dir, err := ioutil.TempDir("/var/tmp", "suitesync_")
	if err != nil {
		panic(err)
	}
	return dir
}

// ArrayIncludes check if array includes
func ArrayIncludes(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Remove remove directory
func Remove(p string) {
	err := os.RemoveAll(p)
	if err != nil {
		PrFatalf("Error removing %s\n%s", p, err.Error())
	}
}

// ToJSON convert interface to json
func ToJSON(p interface{}) []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bytes
}

// CheckDir check if path is directory and exists
func CheckDir(p string) error {
	s, err := CheckExists(p)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("\"%s\" is not a directory", p)
	}
	return nil
}

// CheckExists returns whether the given file or directory exists or not
func CheckExists(p string) (s os.FileInfo, err error) {
	s, err = os.Stat(p)
	if os.IsNotExist(err) {
		return s, fmt.Errorf("\"%s\" does not exist", p)
	}
	return s, err
}

// CheckRequired check if required arg exists
func CheckRequired(s []string, i int, arg string) string {
	if i >= len(s) {
		PrFatalf("\nRequired arg \"%s\" is missing\n", arg)
	}
	return s[i]
}

// OptArgDest set default for optional arg
func OptArgDest(args cli.Args, i int, def string) string {
	if i >= len(args) {
		return OptDest(def, "")
	}
	return args[i]
}

// OptDest set default for optional destination
func OptDest(d string, def string) (s string) {
	s = d
	if s == "" {
		s = def
	}
	if s == "" {
		s, _ = os.Getwd()
	}
	return
}

func sed(file, old, new string) {
	data, _ := ioutil.ReadFile(file)
	output := strings.Replace(string(data), old, new, 1)
	ioutil.WriteFile(file, []byte(output), os.ModePerm)
}

// FindDir find dirs recursively by regex
func FindDir(pathS, rx string) (found []string) {
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			r, err := regexp.MatchString(rx, f.Name())
			if err == nil && r {
				found = append(found, path)
			}
		}
		return nil
	})
	return found
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", strings.Join([]string{"command -v", name}, " "))
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// CheckCliCache test and set the cli cache token
func CheckCliCache() []byte {
	b, err := ioutil.ReadFile(CliCache)
	if err != nil || string(b) == "" {
		// file does not exist or is empty
		b = []byte(Credentials[CliToken])
		ioutil.WriteFile(CliCache, b, os.ModePerm)
	}
	return b
}

// check if all needed params exist and gather env variables
func sanitizeCredentials(skipAuthReq bool) (creds map[string]string, err error) {
	godotenv.Load()
	creds = make(map[string]string)

	optionals := []string{TokenID, TokenSecret, ConsumerKey, ConsumerSecret, Password, CliToken}
	required := []string{Account, Email, Realm}
	defaults := []EnvVar{
		EnvVar{env: HashFile, def: "hashes.json"},
		EnvVar{env: Role, def: "3"},
	}

	for _, o := range optionals {
		creds[o] = os.Getenv(o)
	}

	for _, r := range required {
		creds[r] = os.Getenv(r)
		if creds[r] == "" {
			return nil, fmt.Errorf("environment variable \"%s\" is undefined", r)
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
		if IsVerbose {
			PrNoticef("\nUsing \"%s\", \"%s\", \"%s\" and \"%s\" for authentication\n", ConsumerKey, ConsumerSecret, TokenID, TokenSecret)
		}
		if _, err := EncryptCliToken(currentKey, creds); err != nil {
			return nil, err
		}
	} else if creds[CliToken] != "" {
		if IsVerbose {
			PrNoticef("\nUsing \"%s\" for authentication\n", CliToken)
		}
		if err := DecryptCliToken(currentKey, creds); err != nil {
			return nil, err
		}
	} else if !skipAuthReq {
		return nil, fmt.Errorf("either \"%s\" or \"%s\", \"%s\", \"%s\" and \"%s\" env variables have to be defined", CliToken, ConsumerKey, ConsumerSecret, TokenID, TokenSecret)
	}

	return creds, nil
}

// AbsolutePath get absolute path if it is relative
func AbsolutePath(path string) string {
	pref := strings.Join([]string{".", string(filepath.Separator)}, "")
	if path == "" || strings.HasPrefix(path, pref) || !strings.HasPrefix(path, string(filepath.Separator)) {
		pwd, _ := os.Getwd()
		path = filepath.Join(pwd, path)
	}
	return path
}
