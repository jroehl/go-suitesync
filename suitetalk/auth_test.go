package suitetalk

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/jroehl/go-suitesync/lib"
	"github.com/stretchr/testify/assert"
)

const (
	consumerKey    = "consumer_key"
	consumerSecret = "consumer_secret"
	account        = "account"
	tokenID        = "tokenid"
	tokenSecret    = "tokensecret"
)

func TestGetAuthSuiteTalk(t *testing.T) {
	lib.Credentials = make(map[string]string)
	lib.Credentials[lib.Account] = account
	lib.Credentials[lib.ConsumerKey] = consumerKey
	lib.Credentials[lib.ConsumerSecret] = consumerSecret
	lib.Credentials[lib.TokenID] = tokenID
	lib.Credentials[lib.TokenSecret] = tokenSecret

	res, err := GetAuthSuiteTalk(HmacSha256)

	assert.Equal(t, account, res.Account)
	assert.Equal(t, consumerKey, res.ConsumerKey)
	assert.Equal(t, tokenID, res.Token)
	assert.Equal(t, HmacSha256, res.Algorithm)
	assert.NotEmpty(t, res.Nonce)
	assert.NotEmpty(t, res.Signature)
	assert.NotEmpty(t, res.Timestamp)
	assert.Nil(t, err)
}

func TestGetAuthSuiteTalkFail(t *testing.T) {
	_, err := GetAuthSuiteTalk("foobar")
	assert.Equal(t, "Algorithm not known", err.Error())
	assert.NotNil(t, err)
}

func TestPercentEncode(t *testing.T) {
	res := PercentEncode("123.456789!-")
	assert.Equal(t, "123.456789%21-", res)
}

func TestAddTokenHeader(t *testing.T) {

	doc := etree.NewDocument()
	envelope := doc.CreateElement("soap:Envelope")
	addTokenHeader(envelope)

	s, _ := doc.WriteToString()

	assert.Contains(t, s, "<soap:Envelope>")
	assert.Contains(t, s, "<tokenPassport")
	assert.Contains(t, s, "<ns1:account>")
	assert.Contains(t, s, "<ns1:consumerKey>")
	assert.Contains(t, s, "<ns1:token>")
	assert.Contains(t, s, "<ns1:nonce>")
	assert.Contains(t, s, "<ns1:timestamp>")
	assert.Contains(t, s, "<ns1:signature")
	assert.NotNil(t, doc)
}
