package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateKey(t *testing.T) {
	shittyKey, err := GeneratePrivateKey(1024)
	if err != nil {
		t.Fatal(err)
	}

	betterKey, ex := GeneratePrivateKey(2048)
	if ex != nil {
		t.Fatal(ex)
	}

	shittyKeyString := StringForPrivateKey(shittyKey)
	betterKeyString := StringForPrivateKey(betterKey)

	if betterKeyString == shittyKeyString {
		t.Error("Error: Two separately generated keys are equal")
	}

	if len(shittyKeyString) >= len(betterKeyString) {
		t.Error("Error: bits property of GeneratePrivateKey seems to be ignored")
	}
}

func TestParsePrivateKey(t *testing.T) {
	key, err := GeneratePrivateKey(1024)
	if err != nil {
		t.Fatal(err)
	}

	keyString := StringForPrivateKey(key)
	parsedKey, err := PrivateKeyFromString(keyString)
	if err != nil {
		t.Fatal("Error parsing key", err)
	}

	if err := parsedKey.Validate(); err != nil {
		t.Fatal("Error validating parsed key:", err)
	}

	if StringForPrivateKey(parsedKey) != keyString {
		t.Error("Error: Parsed private key does not match original")
	}
}

func TestParsePublicKey(t *testing.T) {
	privateKey, err := GeneratePrivateKey(1024)
	if err != nil {
		t.Fatal("Error generating private key", err)
	}

	publicKey := privateKey.PublicKey
	publicKeyString, err := StringForPublicKey(&publicKey)
	if err != nil {
		t.Fatal("Error marshaling public key:", err)
	}

	parsedKey, err := PublicKeyFromString(publicKeyString)
	if err != nil {
		t.Fatal("Error parsing marshaled public key:", err)
	}

	pks, err := StringForPublicKey(parsedKey)
	if err != nil {
		t.Fatal("Error marshaling public key", err)
	}

	if pks != publicKeyString {
		t.Error("Error: parsed public key does not match")
	}

	fmt.Println(publicKeyString)
}

func TestSign(t *testing.T) {
	privateKey, err := GeneratePrivateKey(1024)
	if err != nil {
		t.Fatal("Error generating private key:", err)
	}

	sig, err := SignMessageWithKey(privateKey, "hello world")
	if err != nil {
		t.Fatal("Error signing message", err)
	}

	fmt.Println(sig)

	publicKey := privateKey.PublicKey
	if err := ValidateSignatureForMessage("hello world", sig, &publicKey); err != nil {
		t.Fatal("Error verifying signature:", err)
	}
}

func TestToken(t *testing.T) {
	privateKey, err := GeneratePrivateKey(1024)
	if err != nil {
		t.Fatal("Error generating private key:", err)
	}

	token, err := NewToken(privateKey)
	if err != nil {
		t.Fatal("Error generating token:", err)
	}

	expirationDuration := 100 * time.Millisecond
	fmt.Println("Generated token", token)

	if !TokenValid(token, 1*time.Hour, &privateKey.PublicKey) {
		t.Fatal("Error: freshly generated token already invalid:", token)
	}

	time.Sleep(expirationDuration)

	if TokenValid(token, expirationDuration, &privateKey.PublicKey) {
		t.Fatal("Error: token still valid after expiration")
	}

	anotherKey, err := GeneratePrivateKey(1024)
	if err != nil {
		t.Fatal("Error generating private key:", err)
	}

	if TokenValid(token, 1*time.Hour, &anotherKey.PublicKey) {
		t.Fatal("Error: token valid for another public key")
	}
}
