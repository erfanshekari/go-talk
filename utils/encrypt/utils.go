package encypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func ExportPubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	publicKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return publicKeyPem
}

func ExportPrvKeyAsPEMStr(pubkey *rsa.PrivateKey) string {
	privateKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(pubkey),
		},
	))
	return privateKeyPem
}

func ParsePublicKey(pk string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pk))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		pubkey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return pubkey, nil
	}
	pubkey := publicKey.(*rsa.PublicKey)
	return pubkey, nil
}
