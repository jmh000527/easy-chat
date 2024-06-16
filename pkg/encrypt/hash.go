package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// Md5 返回给定字符串的MD5哈希值
// 参数:
//
//	str: 需要进行MD5哈希的字符串
//
// 返回值:
//
//	字符串形式的MD5哈希值
func Md5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// GenPasswordHash 生成给定密码的bcrypt哈希值
// 参数:
//
//	password: 需要加密的明文密码
//
// 返回值:
//
//	加密后的哈希密码字节切片
//	错误: 在加密过程中可能出现的错误
//
// 说明:
//
//	使用bcrypt的默认成本值对密码进行加密
//
// GenPasswordHash hash加密
func GenPasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// ValidatePasswordHash 验证明文密码和bcrypt哈希密码是否匹配
// 参数:
//
//	password: 用户输入的明文密码
//	hashed: 存储的bcrypt哈希密码
//
// 返回值:
//
//	匹配成功返回true，否则返回false
//
// 说明:
//
//	通过bcrypt的CompareHashAndPassword函数验证明文密码和哈希密码是否匹配
//
// ValidatePasswordHash hash校验
func ValidatePasswordHash(password string, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}
