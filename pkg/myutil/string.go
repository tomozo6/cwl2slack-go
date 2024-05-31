package myutil

import (
	"strconv"
)

// strconvParseFloatは文字列をfloat64に変換します。
// strconv.ParseFloatと違い、空文字列の場合は0を返します。
func StrconvParseFloat(str string, bitSize int) (float64, error) {
	// 空文字列の場合は0を返す
	if str == "" {
		return 0, nil
	}

	// 文字列をfloat64に変換
	f, err := strconv.ParseFloat(str, bitSize)
	if err != nil {
		return f, err
	}

	return f, nil
}
