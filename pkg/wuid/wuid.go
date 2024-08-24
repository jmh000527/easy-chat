package wuid

import (
	"database/sql"
	"fmt"
	"github.com/edwingeng/wuid/mysql/wuid"
	"sort"
	"strconv"
)

// w 是全局的 wuid.WUID 对象，用于生成唯一标识符（UID）。
var w *wuid.WUID

// Init 初始化 wuid.WUID 对象。
//
// 功能描述:
//   - 该函数根据提供的数据源名称（DSN）初始化一个 wuid.WUID 对象，并从 MySQL 数据库加载初始值。
//
// 参数:
//   - dsn: 数据源名称，用于连接 MySQL 数据库的 DSN 字符串。
func Init(dsn string) {
	// newDB 是一个函数，用于创建数据库连接。
	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		return db, true, nil
	}

	// 创建一个 wuid.WUID 对象，并从 MySQL 加载初始值。
	w = wuid.NewWUID("default", nil)
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

// GenUid 生成一个唯一标识符（UID）。
//
// 功能描述:
//   - 该函数生成一个新的唯一标识符。如果全局 wuid.WUID 对象未初始化，则会首先进行初始化。
//
// 参数:
//   - dsn: 数据源名称，用于连接 MySQL 数据库的 DSN 字符串（仅在初始化时使用）。
//
// 返回值:
//   - string: 生成的唯一标识符，以十六进制字符串形式返回。
func GenUid(dsn string) string {
	// 如果 w 为空，初始化 w。
	if w == nil {
		Init(dsn)
	}

	// 生成一个新的唯一标识符并返回。
	return fmt.Sprintf("%#016x", w.Next())
}

// CombineId 将两个字符串标识符组合成一个新的字符串标识符。
//
// 功能描述:
//   - 该函数接收两个字符串标识符，将它们按数值大小排序后拼接成一个新的标识符，返回格式为 "小ID_大ID"。
//
// 参数:
//   - aid: 第一个标识符，类型为字符串。
//   - bid: 第二个标识符，类型为字符串。
//
// 返回值:
//   - string: 组合后的新标识符，格式为 "小ID_大ID"。
func CombineId(aid, bid string) string {
	// 将两个标识符排序后，以 "_" 分隔拼接成一个新的字符串标识符。
	ids := []string{aid, bid}
	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b
	})
	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}
