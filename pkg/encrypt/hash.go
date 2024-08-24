package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// Md5 返回给定字符串的 MD5 哈希值。
//
// 参数:
//   - str: 需要计算 MD5 哈希值的字节切片。
//
// 返回值:
//   - string: 计算得到的 MD5 哈希值，以十六进制字符串形式返回。
func Md5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// GenPasswordHash 生成给定密码的 bcrypt 哈希值。
//
// 参数:
//   - password: 需要加密的密码，类型为字节切片。
//
// 返回值:
//   - []byte: 生成的 bcrypt 哈希值，类型为字节切片。
//   - error: 如果生成哈希值过程中出现错误，则返回该错误；否则返回 nil。
func GenPasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// ValidatePasswordHash 验证明文密码与 bcrypt 哈希密码是否匹配。
//
// 参数:
//   - password: 明文密码，类型为字符串。
//   - hashed: bcrypt 哈希后的密码，类型为字符串。
//
// 返回值:
//   - bool: 如果密码匹配返回 true；否则返回 false。
func ValidatePasswordHash(password string, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}
