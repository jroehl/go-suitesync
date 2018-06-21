package suitetalk

import (
	"github.com/beevik/etree"
	"github.com/jroehl/go-suitesync/lib"
)

func addTokenHeader(header *etree.Element) {

	tokenPassport := header.CreateElement("tokenPassport")
	tokenPassport.CreateAttr("xmlns", "urn:messages_2018_1.platform.webservices.netsuite.com")
	tokenPassport.CreateAttr("xmlns:ns1", "urn:core_2018_1.platform.webservices.netsuite.com")

	a, _ := lib.GetAuthSuiteTalk(lib.HmacSha256)

	tokenPassport.CreateElement("ns1:account").CreateCharData(a.Account)
	tokenPassport.CreateElement("ns1:consumerKey").CreateCharData(a.ConsumerKey)
	tokenPassport.CreateElement("ns1:token").CreateCharData(a.Token)
	tokenPassport.CreateElement("ns1:nonce").CreateCharData(a.Nonce)
	tokenPassport.CreateElement("ns1:timestamp").CreateCharData(a.Timestamp)

	signature := tokenPassport.CreateElement("ns1:signature")
	signature.CreateAttr("algorithm", a.Algorithm)
	signature.CreateCharData(a.Signature)
}
