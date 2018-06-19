package lib

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// EncryptCliToken encrypt the given key and secrets
func EncryptCliToken(currentKey string, creds map[string]string) (string, error) {
	data := []byte(strings.Join(
		[]string{
			currentKey,
			strings.Join([]string{creds[TokenID], creds[TokenSecret]}, "&"),
		}, "=",
	))
	encrypted, err := encrypt([]byte(pad(SdfCliPw, 16)), data)
	creds[CliToken] = encrypted
	return encrypted, err
}

// DecryptCliToken decrypt the given token by using the key
func DecryptCliToken(currentKey string, creds map[string]string) error {
	hx, err := hex.DecodeString(creds[CliToken])
	if err != nil {
		return err
	}
	decrypted, err := decrypt([]byte(pad(SdfCliPw, 16)), hx)
	if err != nil {
		return err
	}
	r := bufio.NewReader(strings.NewReader(string(decrypted)))
	s, e := readln(r)
	for e == nil {
		if s != "" && !strings.HasPrefix(s, "#") {
			sp := strings.Split(s, "=")
			if len(sp) == 2 {
				k, v := sp[0], sp[1]
				if currentKey == k {
					tSp := strings.Split(v, "&")
					creds[TokenID], creds[TokenSecret] = tSp[0], tSp[1]

					dCK, _ := hex.DecodeString(SdfCliConsumerKey)
					dCS, _ := hex.DecodeString(SdfCLiConsumerSecret)
					rCL, err := decrypt([]byte(pad(SdfCliPw, 16)), dCK)
					rCS, err := decrypt([]byte(pad(SdfCliPw, 16)), dCS)
					if err != nil {
						return err
					}
					creds[ConsumerKey] = string(rCL)
					creds[ConsumerSecret] = string(rCS)
				}
			}
		}
		s, e = readln(r)
	}
	if creds[TokenID] == "" {
		return errors.New(fmt.Sprintf("\"%s\" does not contain secrets for \"%s\"", CliToken, strings.Replace(currentKey, "-", ", ", -1)))
	}
	return nil
}

// pad string to fixed length
func pad(str string, length int) (padded string) {
	padded = str
	for len(padded) < length {
		padded = padded + str
	}
	return padded[:length]
}

// readln read line from reader
func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

// decrypt aes/ecb decrypt string
func decrypt(passphrase, data []byte) ([]byte, error) {
	cipher, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted, nil
}

// pKCS5Padding pad pkcs5 string
func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// newECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		PrFatalf("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		PrFatalf("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// encrypt aes/ecb encrypt string
func encrypt(passphrase, data []byte) (string, error) {
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return "", err
	}
	ecb := newECBEncrypter(block)
	content := data
	content = pKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	encoded := hex.EncodeToString(crypted)
	return encoded, nil
}
