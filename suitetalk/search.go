package suitetalk

import (
	"log"

	"github.com/beevik/etree"
)

func search(body *etree.Element, qType string, searchBy string) {

	search := body.CreateElement("search")
	search.CreateAttr("xmlns", messages)

	searchRecord := search.CreateElement("searchRecord")
	searchRecord.CreateAttr("xmlns:q1", filecabinet)
	searchRecord.CreateAttr("xsi:type", qType)

	// if searchBy != nil {
	// 	criteria := searchRecord.CreateElement("q1:criteria")
	// 	critBasic := criteria.CreateElement("q1:basic")
	// 	name := critBasic.CreateElement("name")
	// 	name.CreateAttr("operator", "is")
	// 	name.CreateAttr("xmlns", common)
	// 	searchValue := name.CreateElement("searchValue")
	// 	searchValue.CreateAttr("xmlns", "urn:core_2018_1.platform.webservices.netsuite.com")
	// 	searchValue.CreateCharData(searchString)
	// }

	columns := searchRecord.CreateElement("q1:columns")
	colBasic := columns.CreateElement("q1:basic")
	colBasic.CreateElement("parent").CreateAttr("xmlns", common)
	colBasic.CreateElement("internalId").CreateAttr("xmlns", common)
	colBasic.CreateElement("name").CreateAttr("xmlns", common)
}

func parseSearch(xml []byte) (arr []SearchResult) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xml); err != nil {
		log.Fatal(err)
	}
	searchRowList := doc.FindElement("soapenv:Envelope/soapenv:Body/searchResponse/platformCore:searchResult/platformCore:searchRowList")

	if searchRowList == nil {
		log.Fatal("NO_RESULT")
	}
	for _, el := range searchRowList.FindElements("platformCore:searchRow/docFileCab:basic") {
		sr := SearchResult{}
		if sv := el.FindElement("platformCommon:internalId/platformCore:searchValue"); sv != nil {
			sr.InternalID = sv.SelectAttrValue("internalId", "NOT_FOUND")
		}
		if p := el.FindElement("platformCommon:parent/platformCore:searchValue"); p != nil {
			sr.Parent = p.SelectAttrValue("internalId", "NOT_FOUND")
		}
		if n := el.FindElement("platformCommon:name/platformCore:searchValue"); n != nil {
			sr.Name = n.Text()
		}
		arr = append(arr, sr)
	}

	return
}
