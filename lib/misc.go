package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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

// MkTempDir create temporary dir
func MkTempDir() string {
	dir, err := ioutil.TempDir("/var/tmp", "suitesync")
	if err != nil {
		log.Fatal(err)
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
		PrFatalf("Error removing %s\n%s", p, err)
	}
}

// ToJson convert interface to json
func ToJson(p interface{}) []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

// CheckDir check if path is directory
func CheckDir(p string) error {
	s, err := CheckExists(p)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return errors.New(fmt.Sprintf("\"%s\" is not a directory", p))
	}
	return nil
}

// CheckExists returns whether the given file or directory exists or not
func CheckExists(p string) (os.FileInfo, error) {
	s, err := os.Stat(p)
	if os.IsNotExist(err) {
		return s, errors.New(fmt.Sprintf("\"%s\" does not exist", p))
	}
	return s, err
}

func sed(file, old, new string) {
	data, _ := ioutil.ReadFile(file)
	output := strings.Replace(string(data), old, new, 1)
	ioutil.WriteFile(file, []byte(output), os.ModePerm)
}

func find(pathS, rx string) (found []string) {
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
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
