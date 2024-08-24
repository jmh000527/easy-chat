package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identify = "github/jmh000527"

// GetJwtToken 根据给定的参数生成一个 JWT 令牌。
//
// 参数:
//   - secretKey: 用于签名 JWT 令牌的密钥，类型为字符串。
//   - iat: 令牌的发行时间（以 Unix 时间戳表示），类型为 int64。
//   - seconds: 令牌的有效时间长度（以秒为单位），类型为 int64。
//   - uid: 用户唯一标识符，类型为字符串。
//
// 返回值:
//   - string: 生成的 JWT 令牌，类型为字符串。
//   - error: 如果生成令牌时出现错误，返回该错误；否则返回 nil。
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	// 初始化 JWT 声明
	claims := make(jwt.MapClaims)

	// 设置 JWT 的过期时间，基于发行时间和指定的秒数
	claims["exp"] = iat + seconds

	// 设置 JWT 的发行时间
	claims["iat"] = iat

	// 设置用户唯一标识符，后续通过 GetUid 方法从 token 中获取 uid
	claims[Identify] = uid

	// 创建一个新的 JWT 令牌，使用 HS256 算法签名
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置令牌的声明
	token.Claims = claims

	// 使用密钥签名令牌并返回
	return token.SignedString([]byte(secretKey))
}
