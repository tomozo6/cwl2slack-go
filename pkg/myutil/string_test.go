package myutil

import (
	"testing"
)

func TestStrconvParseFloat(t *testing.T) {
	testCases := []struct {
		name     string
		str      string
		isNormal bool
		want     float64
	}{
		{
			name:     "文字列\"\"の場合",
			str:      "",
			isNormal: true,
			want:     0,
		},
		{
			name:     "文字列\"0\"の場合",
			str:      "0",
			isNormal: true,
			want:     0,
		},
		{
			name:     "文字列\"1\"の場合",
			str:      "1",
			isNormal: true,
			want:     1,
		},
		{
			name:     "文字列が数字でない場合",
			str:      "aiueo",
			isNormal: false,
			want:     0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrconvParseFloat(tt.str, 64)

			// 正常系のテストケース
			if tt.isNormal {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Fatalf("unexpected result: %v", got)
				}
				// 異常系のテストケース
			} else {
				if err == nil {
					t.Fatalf("unexpected result: %v", got)
				}
			}
		})
	}
}
