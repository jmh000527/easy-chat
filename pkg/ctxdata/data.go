package ctxdata

import (
	"context"
)

// GetUId 从上下文中尝试获取用户ID。
// 如果上下文中存在以Identify为键的值，并且该值的类型为string，则返回该字符串值。
// 如果上下文中没有以Identify为键的值，或该值的类型不是string，则返回空字符串。
// 参数:
//
//	ctx - 上下文对象，用于传递请求范围内的值。
//
// 返回值:
//
//	用户ID的字符串表示，如果无法获取则为空字符串。
func GetUId(ctx context.Context) string {
	// 尝试从上下文中获取以Identify为键的值，并断言其类型为string。
	if u, ok := ctx.Value(Identify).(string); ok {
		// 如果获取成功，返回用户ID。
		return u
	}
	// 如果获取失败，返回空字符串。
	return ""
}
