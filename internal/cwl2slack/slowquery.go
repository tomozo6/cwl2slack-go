package cwl2slack

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SlowQuery struct {
	Time         string
	User         string
	ID           string
	QueryTime    float64
	LockTime     string
	RowsSent     string
	RowsExamined string
	Query        string
}

// NewSlowQueryはスロークエリログテキストを解析し、SlowQueryインスタンスを返します。
// ログテキストが予想される形式と一致しない場合、エラーを返します
func NewSlowQuery(logText string) (*SlowQuery, error) {
	pattern := regexp.MustCompile(`# Time: ([\d-]+T[\d:.]+Z)\n# User@Host: ([\w-]+)\[\w+\] @ \[[\d.]+\] Id: (\d+)\n# Query_time: ([\d\.]+) Lock_time: ([\d\.]+) Rows_sent: (\d+) Rows_examined: (\d+)\n`)
	matches := pattern.FindStringSubmatch(logText)

	// 一致する項目の数が8でない場合ログテキストは予想される形式と一致しません。
	if len(matches) != 8 {
		return nil, fmt.Errorf("failed to parse log text")
	}

	// クエリ時間をfloat64に変換します。
	queryTime, err := strconv.ParseFloat(matches[4], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query time: %v", err)
	}

	// クエリ部分を抽出します。これは、ログテキストから一致した部分を削除することで行います。
	// この行は、元のログテキスト(logText)からマッチした部分(matches[0])を削除しています。
	// 具体的には、strings.TrimPrefix関数は、最初の引数の文字列が二番目の引数の文字列で始まる場合、その部分を削除した新しい文字列を返します。
	// このコードの場合、matches[0]は正規表現にマッチした全体の文字列（つまり、ログのメタデータ部分）を表しています。
	//  したがって、この行は元のログテキストからメタデータ部分を削除し、クエリ部分だけを抽出しています。
	query := strings.TrimPrefix(logText, matches[0])

	return &SlowQuery{
		Time:         matches[1],
		User:         matches[2],
		ID:           matches[3],
		QueryTime:    queryTime,
		LockTime:     matches[5],
		RowsSent:     matches[6],
		RowsExamined: matches[7],
		Query:        query,
	}, nil
}
