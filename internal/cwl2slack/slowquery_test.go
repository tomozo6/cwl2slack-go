package cwl2slack

import (
	"fmt"
	"testing"
)

func TestNewSlowQuery(t *testing.T) {

	// テストケース
	data := "# Time: 2023-10-22T02:57:55.655927Z\n# User@Host: xxxxxxxxx[xxxxxxxxx] @ [10.13.103.170] Id: 2638113\n# Query_time: 35.549734 Lock_time: 0.000164 Rows_sent: 1 Rows_examined: 15535\nuse xxxx;\nSET timestamp=1697943475;\nSELECT SLEEP(20);"
	// fmt.Println(data)

	s, err := NewSlowQuery(data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", s)

}
