package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	encryptionResult = "b8d3596db885fd1c4be5cccccd89a78ebb97147409d0c8e1832014c7cd242e753158638f039f11bd27784269942e24e67b10d12e05eca4c487cc65cc96c7b22c"
	currentKey       = "system.realm-account-email-3"
	tokenID          = "tokenid"
	tokenSecret      = "tokensecret"
)

func TestEncryptCliToken(t *testing.T) {
	creds := make(map[string]string)
	creds[TokenID] = tokenID
	creds[TokenSecret] = tokenSecret
	EncryptCliToken(currentKey, creds)
	assert.Equal(t, encryptionResult, creds[CliToken])
}

func TestDecryptCliToken(t *testing.T) {
	creds := make(map[string]string)
	creds[CliToken] = encryptionResult
	DecryptCliToken(currentKey, creds)
	assert.Equal(t, tokenID, creds[TokenID])
	assert.Equal(t, tokenSecret, creds[TokenSecret])
	assert.NotZero(t, creds[ConsumerKey])
	assert.NotZero(t, creds[ConsumerSecret])
}

func TestDecryptCliTokenFail(t *testing.T) {
	creds := make(map[string]string)
	creds[CliToken] = encryptionResult
	err := DecryptCliToken("non-existing-key", creds)
	assert.EqualError(t, err, "\"NSCONF_CLITOKEN\" does not contain secrets for \"non, existing, key\"")
}
