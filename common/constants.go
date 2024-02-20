package common

import "os"

const YAMLCONF = ".ghs.yml"
const SecretPath = "./var/secret"

var PublicKeyPath = SecretPath + "/public.pem"
var PrivateKeyPath = SecretPath + "/private.pem"

func init() {
	customPublicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
	if customPublicKeyPath != "" {
		PublicKeyPath = customPublicKeyPath
	}
	customPrivateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if customPrivateKeyPath != "" {
		PrivateKeyPath = customPrivateKeyPath
	}
}
