package encypt

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math"

	"github.com/erfanshekari/go-talk/events"
)

const (
	rsaKeySize = 4096
)

func GenerateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, rsaKeySize)
}

func encrypt(i []byte, key *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, key, i)
}

func decrypt(i []byte, key *rsa.PrivateKey) ([]byte, error) {
	h, err := base64.StdEncoding.DecodeString(string(i))
	if err != nil {
		return nil, err
	}
	b, err := rsa.DecryptPKCS1v15(rand.Reader, key, h)
	if err != nil {
		return nil, err
	}
	return b, err
}

func encryptAddChunk(input []byte, key *rsa.PublicKey, chunks *[][]byte) error {
	val, err := encrypt(input, key)
	if err != nil {
		return err
	}
	*chunks = append(*chunks, val)
	return nil
}

func Encrypt(input []byte, key *rsa.PublicKey) (*events.EncryptedJson, error) {
	var chunks [][]byte
	totalLen := len(input)
	keySize := key.Size() - 11
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
	var ej events.EncryptedJson
	if len(chunks) > 1 {
		ej.Type = events.ByteArray
		ej.Content = chunks
	} else {
		ej.Type = events.Byte
		ej.Content = chunks[0]
	}
	return &ej, nil
}

func Decrypt(input []byte, k *rsa.PrivateKey) (*[]byte, error) {
	var encryptedJson events.EncryptedJson
	err := json.Unmarshal(input, &encryptedJson)
	if err != nil {
		return nil, err
	}
	var decryptedByte *[]byte
	switch encryptedJson.Type {
	case events.Byte:
		data, ok := encryptedJson.Content.(string)
		if !ok {
			return nil, errors.New("Can't parse Content to string")
		}
		db, err := decrypt([]byte(data), k)
		if err != nil {
			return nil, err
		}
		decryptedByte = &db
		break
	case events.ByteArray:
		data, ok := encryptedJson.Content.([]string)
		if !ok {
			return nil, errors.New("Can't parse Content to []string")
		}
		var decrypted string
		for _, val := range data {
			db, err := decrypt([]byte(val), k)
			if err != nil {
				return nil, err
			}
			decrypted += string(db)
		}
		db_ := []byte(decrypted)
		decryptedByte = &db_
		break
	}
	return decryptedByte, nil
}
