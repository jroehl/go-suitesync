package suitetalk

import (
	"fmt"
	"strings"

	"github.com/jroehl/go-suitesync/lib"
)

var (
	Cache      = make(map[string][]*lib.SearchResult)
	Pathlookup = make(map[string]*lib.SearchResult)
)

// ListFiles list the files of a specific directory in the filecabinet
func ListFiles(client HTTPClient, dir string) ([]*lib.SearchResult, error) {
	res, err := GetPath(client, dir)
	if err != nil {
		return nil, err
	}
	if !res.IsDir {
		return nil, fmt.Errorf("\"%s\" is not a directory", dir)
	}
	return FlattenChildren(res.Children), nil
}

// FlattenChildren recursively flatten the children of a directory in the filecabinet
func FlattenChildren(children []*lib.SearchResult) (cs []*lib.SearchResult) {
	cs = append(cs, children...)
	for _, c := range children {
		cs = append(cs, FlattenChildren(c.Children)...)
	}
	return cs
}

// GetPath get the item from the lookup map
func GetPath(client HTTPClient, path string) (res *lib.SearchResult, err error) {
	getFs(client, path)
	res = Pathlookup[path]
	if res == nil {
		return nil, fmt.Errorf("\nNo result for \"%s\"\n\n", path)
	}
	return res, nil
}

func getFs(client HTTPClient, path string) []*lib.SearchResult {
	if Cache[path] == nil || Pathlookup == nil {
		dirID := getDirectoryID(client, path)
		fmt.Println(dirID)
		var (
			fif []lib.SearchFilter
			fof []lib.SearchFilter
		)
		if dirID != "@NONE@" {
			fif = []lib.SearchFilter{
				lib.SearchFilter{
					Tag:      "folder",
					Operator: "anyOf",
					SearchValues: []lib.SearchValue{
						lib.SearchValue{
							Inner: "",
							Attrs: []lib.Attr{
								lib.Attr{"internalId", dirID},
								lib.Attr{"type", "folder"},
							},
						},
					},
				},
			}
			fof = []lib.SearchFilter{
				lib.SearchFilter{
					Tag:      "predecessor",
					Operator: "anyOf",
					SearchValues: []lib.SearchValue{
						lib.SearchValue{
							Inner: "",
							Attrs: []lib.Attr{
								lib.Attr{"internalId", dirID},
								lib.Attr{"type", "folder"},
							},
						},
					},
				},
			}
		}
		folders := SearchRequest(client, searchFolder, fof)
		files := SearchRequest(client, searchFile, fif)
		fmt.Println(folders)
		fmt.Println(files)
		merged := append(folders, files...)
		Cache[path] = sliceToTree(merged)
		extractPaths(Cache[path], "")
	}
	return Cache[path]
}

func getDirectoryID(client HTTPClient, path string) string {
	parentID := "@NONE@" // first search is in root dir
	fmt.Println(path)
	for _, p := range strings.Split(strings.Trim(path, "/"), "/") {
		fmt.Println(p)
		f := []lib.SearchFilter{
			lib.SearchFilter{
				Tag:      "name",
				Operator: "is",
				SearchValues: []lib.SearchValue{
					lib.SearchValue{
						Inner: p,
					},
				},
			},
			lib.SearchFilter{
				Tag:      "parent",
				Operator: "anyOf",
				SearchValues: []lib.SearchValue{
					lib.SearchValue{
						Inner: "",
						Attrs: []lib.Attr{
							lib.Attr{"internalId", parentID},
							lib.Attr{"type", "folder"},
						},
					},
				},
			},
		}
		res := SearchRequest(client, searchFolder, f)
		if len(res) < 1 || res[0].InternalID == "" {
			lib.PrFatalf("Error \"%s\" - \"%s\" not found", searchFolder, p)
		}
		fmt.Println(res)
		parentID = res[0].InternalID
	}
	return parentID
}

func extractPaths(srs []*lib.SearchResult, path string) map[string]*lib.SearchResult {
	for _, s := range srs {
		p := path
		p = strings.Join([]string{p, s.Name}, "/")
		s.Path = p
		Pathlookup[p] = s
		extractPaths(s.Children, p)
	}
	return Pathlookup
}

// sliceToTree convert slice of parent/child entries to tree map
func sliceToTree(slice []lib.SearchResult) (tree []*lib.SearchResult) {
	nodes := make(map[string]*lib.SearchResult)
	var xx *lib.SearchResult
	for _, x := range slice {
		tmp := x
		xx = &tmp
		nodes[xx.InternalID] = xx
	}

	for _, v := range nodes {
		if v.Parent == "" {
			tree = append(tree, v)
		} else if nodes[v.Parent] != nil {
			nodes[v.Parent].Children = append(nodes[v.Parent].Children, v)
		}
	}
	return tree
}

// printTree print tree to console
func printTree(arr []*lib.SearchResult, i int) {
	for _, c := range arr {
		var t string
		if c.IsDir {
			t = "folder"
		} else {
			t = "file"
		}
		fmt.Printf("%s> %s (TYPE: %s   ID: %s   PATH: %s)\n", strings.Repeat("|   ", i), c.Name, t, c.InternalID, c.Path)
		if len(c.Children) > 0 {
			printTree(c.Children, i+1)
		}
	}
}
