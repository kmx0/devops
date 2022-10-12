package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func ReadPublicKey(filepath string) (*rsa.PublicKey, error) {
	var publicKey *rsa.PublicKey
	pubData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return publicKey, err
	}
	// сертификат разбора

	block, _ := pem.Decode([]byte(pubData))
	cert, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return publicKey, err
	}

	publicKey = cert.(*rsa.PublicKey)
	return publicKey, nil
}

func ReadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	pfxData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	// сертификат разбора
	block, _ := pem.Decode(pfxData)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Ошибка разбора сертификата" + err.Error())
		return nil, err
	}
	return privateKey, nil
}
func EncryptData(publicKey rsa.PublicKey, data []byte) (encryptedBytes []byte, err error) {

	msgLen := len(data)
	hash := sha256.New()
	random := rand.Reader

	step := publicKey.Size() - 2*hash.Size() - 2
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, &publicKey, data[start:finish], nil)
		if err != nil {
			return nil, err
		}
		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}
	return encryptedBytes, nil
}

func DecryptData(privateKey *rsa.PrivateKey, encryptedBytes []byte) (decryptedBytes []byte, err error) {
	msgLen := len(encryptedBytes)
	step := privateKey.PublicKey.Size()
	hash := sha256.New()
	random := rand.Reader

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, privateKey, encryptedBytes[start:finish], nil)
		if err != nil {
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}
	return decryptedBytes, nil
}
