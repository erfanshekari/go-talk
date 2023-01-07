package encypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"math"
)

func GenerateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

func encrypt(i []byte, key *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key, i, nil)
}

func encryptAddChunk(input []byte, key *rsa.PublicKey, chunks *[][]byte) error {
	val, err := encrypt(input, key)
	if err != nil {
		return err
	}
	*chunks = append(*chunks, val)
	return nil
}

func Encrypt(input []byte, key *rsa.PublicKey) ([][]byte, error) {
	var chunks [][]byte
	totalLen := len(input)
	keySize := key.Size() - 66
	var encryptedLen int
	if totalLen > keySize {
		for i := 0; i < (int(math.Ceil(float64(totalLen)))); i++ {
			if (encryptedLen + keySize) > totalLen {
				// last part smaller then max size
				err := encryptAddChunk(input[encryptedLen:], key, &chunks)
				if err != nil {
					return nil, err
				}
				break
			} else if i == 0 {
				// first part
				err := encryptAddChunk(input[0:keySize], key, &chunks)
				if err != nil {
					return nil, err
				}
				encryptedLen += keySize
				continue
			}
			err := encryptAddChunk(input[encryptedLen:(encryptedLen+keySize)], key, &chunks)
			if err != nil {
				return nil, err
			}
			encryptedLen += keySize
		}
	} else {
		val, err := encrypt(input, key)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, val)
	}
	return chunks, nil
}
