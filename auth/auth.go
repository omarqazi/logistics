package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"io"
)

func GeneratePrivateKey(bits int) (p *rsa.PrivateKey, e error) {
	p, e = rsa.GenerateKey(rand.Reader, bits)
	if e != nil {
		return
	}
	p.Precompute()
	e = p.Validate()
	return
}

func StringForPrivateKey(p *rsa.PrivateKey) string {
	bytes := x509.MarshalPKCS1PrivateKey(p)
	bsf := base64.URLEncoding.EncodeToString(bytes)
	return bsf
}

func PrivateKeyFromString(pem string) (*rsa.PrivateKey, error) {
	der, err := base64.URLEncoding.DecodeString(pem)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(der)
	return key, err
}

func StringForPublicKey(p *rsa.PublicKey) (string, error) {
	bytes, err := x509.MarshalPKIXPublicKey(p)
	if err != nil {
		return "", err
	}

	der := base64.URLEncoding.EncodeToString(bytes)
	return der, nil
}

func PublicKeyFromString(pk string) (*rsa.PublicKey, error) {
	der, err := base64.URLEncoding.DecodeString(pk)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(der)
	rsaKey, ok := key.(*rsa.PublicKey)
	if ok {
		return rsaKey, err
	}
	return nil, err
}

func SignMessageWithKey(key *rsa.PrivateKey, msg string) (sig []byte, err error) {
	hashFunction := sha256.New()
	io.WriteString(hashFunction, msg)
	hashSum := hashFunction.Sum(nil)

	sig, err = rsa.SignPSS(rand.Reader, key, crypto.SHA256, hashSum, nil)
	return
}

func ValidateSignatureForMessage(msg string, sig []byte, pub *rsa.PublicKey) (err error) {
	hashFunction := sha256.New()
	io.WriteString(hashFunction, msg)
	hashSum := hashFunction.Sum(nil)

	err = rsa.VerifyPSS(pub, crypto.SHA256, hashSum, sig, nil)
	return
}
