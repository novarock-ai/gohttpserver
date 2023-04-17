package secret

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecret(t *testing.T) {
	test := assert.New(t)
	// 获取当前文件所在的绝对路径
	pwd, _ := os.Getwd()
	pwd = filepath.Join(pwd, "..", "var", "secret")
	p1, p2, err := CreatePEM(pwd)
	if err != nil {
		fmt.Println(err)
	}

	test.Equal(p1, filepath.Join(pwd, "public.pem"))
	test.Equal(p2, filepath.Join(pwd, "private.pem"))
}
