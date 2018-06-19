package lib

import (
	"testing"
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
	if creds[CliToken] != encryptionResult {
		t.Errorf("Encryption failed, got: %s, want: %s.", creds[CliToken], encryptionResult)
	}
}

func TestDecryptCliToken(t *testing.T) {
	creds := make(map[string]string)
	creds[CliToken] = encryptionResult
	DecryptCliToken(currentKey, creds)
	if creds[TokenID] != tokenID {
		t.Errorf("Decryption failed, TokenID => got: %s, want: %s.", creds[TokenID], tokenID)
	}
	if creds[TokenSecret] == "" {
		t.Errorf("Decryption failed, TokenSecret => got: %s, want: %s.", creds[TokenSecret], tokenSecret)
	}
	if creds[ConsumerKey] == "" {
		t.Errorf("Decryption failed, ConsumerKey => got: %s.", creds[ConsumerKey])
	}
	if creds[ConsumerSecret] == "" {
		t.Errorf("Decryption failed, ConsumerSecret => got: %s.", creds[ConsumerSecret])
	}
}

func TestDecryptCliTokenFail(t *testing.T) {
	creds := make(map[string]string)
	creds[CliToken] = encryptionResult
	err := DecryptCliToken("non-existing-key", creds)
	if err == nil {
		t.Errorf("Fail-decryption failed, Error => got: %s, want: %s.", "nil", "\"NSCONF_CLITOKEN\" does not contain secrets for \"non, existing, key\"")
	}
}

// func TestDecryptCliToken
