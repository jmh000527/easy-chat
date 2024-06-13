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
func GenUid(dsn string) string {
	// 如果 w 为空，初始化 w。
	if w == nil {
		Init(dsn)
	}

	// 生成一个新的唯一标识符并返回。
	return fmt.Sprintf("%#016x", w.Next())
}

// CombineId 将两个字符串标识符组合成一个新的字符串标识符。
func CombineId(aid, bid string) string {
	// 将两个标识符排序后，以 "_"" 分隔拼接成一个新的字符串标识符。
	ids := []string{aid, bid}
	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b
	})
	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}
