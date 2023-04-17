package common

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed assets
var assetsFS embed.FS

// Assets contains project assets.
var Assets = http.FS(assetsFS)

func init() {
	fmt.Println(1111, Assets)
}
