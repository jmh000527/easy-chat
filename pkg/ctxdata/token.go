package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identify = "github/jmh000527"

// GetJwtToken 根据给定的参数生成一个JWT令牌。
// secretKey: 用于签名JWT的密钥。
// iat: JWT的发行时间（Unix时间戳）。
// seconds: JWT的过期时间，单位为秒。
// uid: 用户唯一标识，将被包含在JWT的声明中。
// 返回生成的JWT令牌字符串，以及可能的错误。
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	// 初始化JWT声明
	claims := make(jwt.MapClaims)
	// 设置JWT的过期时间，基于发行时间和指定的秒数
	claims["exp"] = iat + seconds
	// 设置JWT的发行时间
	claims["iat"] = iat
	// 设置用户唯一标识
	claims[Identify] = uid

	// 创建一个新的JWT令牌，使用HS256算法签名
	token := jwt.New(jwt.SigningMethodHS256)
	// 设置令牌的声明
	token.Claims = claims

	// 使用密钥签名令牌并返回
	return token.SignedString([]byte(secretKey))
}
