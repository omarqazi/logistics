package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
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
