package suitetalk

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/beevik/etree"
	"github.com/jroehl/go-suitesync/lib"
)

// DeleteRequest make delete request for filedirectory paths
func DeleteRequest(client HTTPClient, paths []string) ([]lib.DeleteResult, []*lib.SearchResult) {

	if lib.IsVerbose {
		lib.PrettyList("DeleteRequest issued", paths)
	}

	items := []*lib.SearchResult{}
	for _, p := range paths {
		np := strings.Split(p, "/")

		if it, err := GetPath(client, strings.Join(np[:len(np)-1], "/")); it != nil {
			items = append(items, it)
			if it.IsDir {
				items = append(items, FlattenChildren(it.Children)...)
			}
		} else if err != nil {
			lib.PrResultf(err.Error())
		}
	}

	// delete longest paths (deepest nested) items first
	sort.Slice(items, func(i, j int) bool {
		return len(items[i].Path) > len(items[j].Path)
	})
	// files have to be deleted first
	sort.Slice(items, func(i, j int) bool {
		return !items[i].IsDir && len(items[i].Path) < len(items[j].Path)
	})

	if lib.IsVerbose && len(items) > 0 {
		lib.PrNoticef("\nMarked for removal\n")
		for _, x := range items {
			lib.PrNoticef("   - %s\n", x.Path)
		}
		fmt.Println()
	}

	reqs, _ := soapDelete(items, chunkSize)
	var res []lib.DeleteResult
	for i, b := range reqs {
		if lib.IsVerbose {
			lib.PrNoticef("%d/%d: DeleteRequest issued - \"%s\"\n", i+1, len(reqs), deleteList)
		}
		r := doRequest(client, b, deleteList)
		parsed, err := parseSoapDelete(r)
		if err != nil {
			lib.PrFatalf("Error \"%s\" - %s", deleteList, err.Error())
		}
		res = append(res, parsed...)
		if lib.IsVerbose {
			lib.PrNoticef("%d/%d: DeleteRequest done - \"%s\"\n", i+1, len(reqs), deleteList)
		}
	}

	return res, items
}

// Delete delete items from filecabinet
func soapDelete(items []*lib.SearchResult, chunkSize int) (res [][]byte, docs []*etree.Document) {
	chunks := split(items, chunkSize)

	for _, c := range chunks {
		doc, body := soap()
		deleteList := body.CreateElement(deleteList)
		deleteList.CreateAttr("xmlns", messages)
		for _, i := range c {
			baseRef := deleteList.CreateElement("baseRef")
			var t string
			if i.IsDir {
				t = "folder"
			} else {
				t = "file"
			}
			baseRef.CreateAttr("type", t)
			baseRef.CreateAttr("internalId", i.InternalID)
			baseRef.CreateAttr("xsi:type", "q1:RecordRef")
			baseRef.CreateAttr("xmlns:q1", core)
		}
		bytes, err := doc.WriteToBytes()
		if err != nil {
			panic(err)
		}
		// doc.Indent(2)
		// doc.WriteTo(os.Stdout)

		docs = append(docs, doc)
		res = append(res, bytes)
	}

	return res, docs
}

func parseSoapDelete(xml []byte) (res []lib.DeleteResult, err error) {
	doc, err := parseByte(xml)

	if err != nil {
		return nil, err
	}

	writeResponseList := doc.FindElement("soapenv:Envelope/soapenv:Body/deleteListResponse/writeResponseList")
	if writeResponseList == nil {
		doc.Indent(2)
		doc.WriteTo(os.Stdout)
		err = errors.New("REQUEST_ERROR")
		return nil, err
	}
	for _, el := range writeResponseList.FindElements("writeResponse") {
		sr := lib.DeleteResult{}
		if sv := el.FindElement("platformCore:status"); sv != nil {
			sr.Successful = sv.SelectAttrValue("isSuccess", "false") == "true"
			if !sr.Successful {
				if c := sv.FindElement("platformCore:statusDetail/platformCore:code"); c != nil {
					sr.Code = c.Text()
				}
				if m := sv.FindElement("platformCore:statusDetail/platformCore:message"); m != nil {
					sr.Message = m.Text()
				}
				sr.NotFound = sr.Code == "RCRD_DSNT_EXIST" || sr.Code == "MEDIA_NOT_FOUND"
			} else {
				sr.NotFound = false
				sr.Code = "DELETED"
				sr.Message = "Record was successfully deleted"
			}
		}
		if br := el.FindElement("baseRef"); br != nil {
			sr.ID = br.SelectAttrValue("internalId", "")
			sr.Type = br.SelectAttrValue("type", "")
		}
		res = append(res, sr)
	}
	return res, nil
}

// split array in size lim parts
func split(buf []*lib.SearchResult, lim int) [][]*lib.SearchResult {
	var chunk []*lib.SearchResult
	chunks := make([][]*lib.SearchResult, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
