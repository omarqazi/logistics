package auth

import (
	"testing"
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
