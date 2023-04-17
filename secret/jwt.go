package secret

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

func CreateJWT(privateKeyPath string, claims jwt.MapClaims) (string, error) {

	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}

func ParseJWT(publicKeyPath string, tokenString string) (jwt.MapClaims, error) {

	// 从文件中读取公钥
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, err
}
