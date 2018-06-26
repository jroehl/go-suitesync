package suitetalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/jroehl/go-suitesync/lib"
)

type Auth struct {
	Nonce       string
	Timestamp   string
	Signature   string
	Account     string
	ConsumerKey string
	Token       string
	Algorithm   string
}

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyz123456789"
	HmacSha1    = "HMAC-SHA1"
	HmacSha256  = "HMAC-SHA256"
)
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// GetAuthSuiteTalk create oauth credentials
func GetAuthSuiteTalk(algorithm string) (Auth, error) {

	a := Auth{}
	a.Algorithm = algorithm

	a.Nonce = randStringBytesMaskImprSrc(16)
	a.Timestamp = getTimestamp()

	r := lib.Credentials

	a.Account = r[lib.Account]
	a.ConsumerKey = r[lib.ConsumerKey]
	a.Token = r[lib.TokenID]

	baseString := strings.Join([]string{
		a.Account,
		a.ConsumerKey,
		a.Token,
		a.Nonce,
		a.Timestamp,
	}, "&")

	key := strings.Join([]string{r[lib.ConsumerSecret], r[lib.TokenSecret]}, "&")

	sign, err := compute(algorithm, baseString, key)
	a.Signature = sign
	return a, err
}

// PercentEncode percent encodes a string according to RFC 3986 2.1.
func PercentEncode(input string) string {
	var buf bytes.Buffer
	for _, b := range []byte(input) {
		// if in unreserved set
		if shouldEscape(b) {
			buf.Write([]byte(fmt.Sprintf("%%%02X", b)))
		} else {
			// do not escape, write byte as-is
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

// shouldEscape returns false if the byte is an unreserved character that
// should not be escaped and true otherwise, according to RFC 3986 2.1.
func shouldEscape(c byte) bool {
	// RFC3986 2.3 unreserved characters
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '.', '_', '~':
		return false
	}
	// all other bytes must be escaped
	return true
}

func compute(algorithm string, message string, secret string) (string, error) {
	var a func() hash.Hash
	if algorithm == HmacSha1 {
		a = sha1.New
	} else if algorithm == HmacSha256 {
		a = sha256.New
	} else {
		return "", errors.New("Algorithm not known")
	}
	h := hmac.New(a, []byte(secret))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, Cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			Cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(Cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		Cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func getTimestamp() string {
	now := time.Now()
	nanos := now.UnixNano()
	return strconv.Itoa(int(nanos / 1000000000))
}

func addTokenHeader(header *etree.Element) {

	tokenPassport := header.CreateElement("tokenPassport")
	tokenPassport.CreateAttr("xmlns", messages)
	tokenPassport.CreateAttr("xmlns:ns1", core)

	a, _ := GetAuthSuiteTalk(HmacSha256)

	tokenPassport.CreateElement("ns1:account").CreateCharData(a.Account)
	tokenPassport.CreateElement("ns1:consumerKey").CreateCharData(a.ConsumerKey)
	tokenPassport.CreateElement("ns1:token").CreateCharData(a.Token)
	tokenPassport.CreateElement("ns1:nonce").CreateCharData(a.Nonce)
	tokenPassport.CreateElement("ns1:timestamp").CreateCharData(a.Timestamp)

	signature := tokenPassport.CreateElement("ns1:signature")
	signature.CreateAttr("algorithm", a.Algorithm)
	signature.CreateCharData(a.Signature)
}
