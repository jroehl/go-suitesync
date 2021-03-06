package suitetalk

import (
	"errors"
	"os"
	"strconv"

	"github.com/beevik/etree"
	"github.com/jroehl/go-suitesync/lib"
)

// SearchRequest make search request
func SearchRequest(client HTTPClient, command string) []lib.SearchResult {
	doc, body := soap()
	sf := false
	switch command {
	case searchFolder:
		sf = true
		soapSearch(body, folderSearchAdvanced)
	case searchFile:
		soapSearch(body, fileSearchAdvanced)
	default:
		lib.PrFatalf("Command \"%s\" not implemented", command)
	}

	bytes, err := doc.WriteToBytes()
	if err != nil {
		panic(err)
	}

	if lib.IsVerbose {
		lib.PrNoticef("SearchRequest issued - \"%s\"\n", command)
	}

	// doc.Indent(2)
	// doc.WriteTo(os.Stdout)

	res := doRequest(client, bytes, search)
	parsed, meta, err := parseSoapSearch(res, sf, false)
	if err != nil {
		lib.PrFatalf("Error \"%s\" - %s", search, err.Error())
	}
	// if results are paginated
	if meta.TotalPages > 1 {
		for i := 2; i < meta.TotalPages+1; i++ {
			b, _ := soapSearchMore(i, meta.SearchID)
			r := doRequest(client, b, searchMoreWithID)
			p, _, err := parseSoapSearch(r, sf, true)
			if err != nil {
				lib.PrFatalf("Error \"%s\" id \"%s\" - %s", searchMoreWithID, meta.SearchID, err.Error())
			}
			parsed = append(parsed, p...)
		}
	}

	return parsed
}

// SearchMore search more by search id
func soapSearchMore(idx int, id string) ([]byte, *etree.Document) {
	doc, body := soap()

	search := body.CreateElement(searchMoreWithID)

	searchID := search.CreateElement("searchId")
	searchID.CreateCharData(id)

	pageIndex := search.CreateElement("pageIndex")
	pageIndex.CreateCharData(strconv.Itoa(idx))

	bytes, err := doc.WriteToBytes()
	if err != nil {
		panic(err)
	}

	return bytes, doc
}

// Search soap search in filecabinet
func soapSearch(body *etree.Element, qType string) {

	search := body.CreateElement("search")
	search.CreateAttr("xmlns", messages)

	searchRecord := search.CreateElement("searchRecord")
	searchRecord.CreateAttr("xmlns:q1", filecabinet)
	searchRecord.CreateAttr("xsi:type", qType)

	columns := searchRecord.CreateElement("q1:columns")
	colBasic := columns.CreateElement("q1:basic")
	switch qType {
	case folderSearchAdvanced:
		colBasic.CreateElement("parent").CreateAttr("xmlns", common)
	case fileSearchAdvanced:
		colBasic.CreateElement("folder").CreateAttr("xmlns", common)
	}
	colBasic.CreateElement("internalId").CreateAttr("xmlns", common)
	colBasic.CreateElement("name").CreateAttr("xmlns", common)
}

func parseSoapSearch(xml []byte, searchFolder, searchMore bool) (res []lib.SearchResult, m lib.Meta, err error) {
	doc, err := parseByte(xml)

	if err != nil {
		return nil, m, err
	}

	// doc.Indent(2)
	// doc.WriteTo(os.Stdout)

	var searchResult *etree.Element
	if !searchMore {
		searchResult = doc.FindElement("soapenv:Envelope/soapenv:Body/searchResponse/platformCore:searchResult")
	} else {
		searchResult = doc.FindElement("soapenv:Envelope/soapenv:Body/searchMoreWithIdResponse/platformCore:searchResult")
	}
	if searchResult == nil {
		doc.Indent(2)
		doc.WriteTo(os.Stdout)
		err := errors.New("REQUEST_ERROR")
		return nil, m, err
	}
	if s := searchResult.FindElement("platformCore:status"); s != nil {
		m.Successful = s.SelectAttrValue("isSuccess", "false") == "true"
	}
	if tr := searchResult.FindElement("platformCore:totalRecords"); tr != nil {
		i, _ := strconv.Atoi(tr.Text())
		m.TotalRecords = i
	}
	if tp := searchResult.FindElement("platformCore:totalPages"); tp != nil {
		i, _ := strconv.Atoi(tp.Text())
		m.TotalPages = i
	}
	if si := searchResult.FindElement("platformCore:searchId"); si != nil {
		m.SearchID = si.Text()
	}

	searchRowList := searchResult.FindElement("platformCore:searchRowList")
	for _, el := range searchRowList.FindElements("platformCore:searchRow/docFileCab:basic") {
		var sr lib.SearchResult
		if sv := el.FindElement("platformCommon:internalId/platformCore:searchValue"); sv != nil {
			sr.InternalID = sv.SelectAttrValue("internalId", "NOT_FOUND")
		}
		var p *etree.Element
		sr.IsDir = searchFolder
		if searchFolder {
			p = el.FindElement("platformCommon:parent/platformCore:searchValue")
		} else {
			p = el.FindElement("platformCommon:folder/platformCore:searchValue")
		}
		if p != nil {
			sr.Parent = p.SelectAttrValue("internalId", "NOT_FOUND")
		}

		if n := el.FindElement("platformCommon:name/platformCore:searchValue"); n != nil {
			sr.Name = n.Text()
		}
		res = append(res, sr)
	}

	return res, m, nil
}
