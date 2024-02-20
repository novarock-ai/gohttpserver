package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

func CreatePEM(path string) (public, private string, err error) {
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return "", "", err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	privateKeyFilePath := filepath.Join(path, "private.pem")
	err = os.WriteFile(privateKeyFilePath, privateKeyPEM, 0600)
	if err != nil {
		return "", "", err
	}

	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)

	publicKeyFilePath := filepath.Join(path, "public.pem")
	err = os.WriteFile(publicKeyFilePath, publicKeyPEM, 0644)
	if err != nil {
		return "", "", err
	}

	privateKeyFilePath, err = filepath.Abs(privateKeyFilePath)
	if err != nil {
		return "", "", err
	}
	publicKeyFilePath, err = filepath.Abs(publicKeyFilePath)
	if err != nil {
		return "", "", err
	}
	return publicKeyFilePath, privateKeyFilePath, nil
}
