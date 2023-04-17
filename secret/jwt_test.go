package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	test := assert.New(t)

	publicKeyPath, privateKeyPath, _ := CreatePEM("../var/secret")

	token, err := CreateJWT(privateKeyPath, map[string]interface{}{
		"sub": "1234567890",
	})
	if err != nil {
		test.Fail(err.Error())
	}

	claims, err := ParseJWT(publicKeyPath, token)
	if err != nil {
		test.Fail(err.Error())
	}
	test.Equal(claims["sub"], "1234567890")

}
