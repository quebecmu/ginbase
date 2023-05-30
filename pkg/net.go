package pkg

import (
	"bytes"
	"github.com/xxandjg/ginbase/global"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// ReadRequest 将网络请求后得到的流转换为字符串
func ReadRequest(closer io.ReadCloser) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(closer)
	if err != nil {
		return ""
	}
	return buf.String()
}

// PasswordHash 密码加密
func PasswordHash(pwd string) (string, global.Error) {
	by, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", global.FAIL
	}

	return string(by), global.SUCCESS
}

// PasswordVerify 密码验证
func PasswordVerify(pwd, hash string) (bool, global.Error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if err != nil {
		return false, global.FAIL
	}
	return true, global.SUCCESS
}
