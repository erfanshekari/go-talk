package encypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func ExportPubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return pubKeyPem
}

func ExportPrvKeyAsPEMStr(pubkey *rsa.PrivateKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(pubkey),
		},
	))
	return pubKeyPem
}
