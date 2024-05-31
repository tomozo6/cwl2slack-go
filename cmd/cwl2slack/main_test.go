package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
)

// テストだけで使用する関数
// CloudwatchLogsDataをCloudwatchLogsRawDataに変換します
func ConvertToRawData(d events.CloudwatchLogsData) (c events.CloudwatchLogsRawData, err error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	data, err := json.Marshal(d)
	if err != nil {
		return
	}

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}

	c.Data = base64.StdEncoding.EncodeToString(buf.Bytes())

	return
}

func TestHandler(t *testing.T) {
	testCases := []struct {
		name     string
		mode     string
		data     events.CloudwatchLogsData
		isNormal bool
		want     string
	}{
		{
			name: "plainモードの正常系",
			mode: "plain",
			data: events.CloudwatchLogsData{
				Owner:     "123456789012",
				LogGroup:  "testLogGroup",
				LogStream: "testLogStream",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{
						ID:        "12345678901234567890123456789012345678901234567890123456",
						Timestamp: 1716792813043,
						Message:   "[ERROR] First test message",
					},
					{
						ID:        "12345678901234567890123456789012345678901234567890123456",
						Timestamp: 1716792813043,
						Message:   "[ERROR] Second test message",
					},
				},
			},
			isNormal: true,
			want:     "cwl2slack executed successfully.",
		},
		{
			name: "slowqueryモードで、渡されたログがスロークエリログの形式の場合",
			mode: "slowquery",
			data: events.CloudwatchLogsData{
				Owner:     "123456789012",
				LogGroup:  "testLogGroup",
				LogStream: "testLogStream",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{
						ID:        "12345678901234567890123456789012345678901234567890123456",
						Timestamp: 1716792813043,
						Message:   "# Time: 2024-05-27T06:53:33.043104Z\n# User@Host: wsprodadminuser[wsprodadminuser] @ [172.17.0.178] Id: 1436601\n# Query_time: 4.275485 Lock_time: 0.000002 Rows_sent: 58 Rows_examined: 12158\nuse work_prod;\nSET timestamp=1716792808;\nSELECT `mp`.`prime_company_id` AS `primeCompanyId`, `mp`.`prime_company_name` AS `primeCompanyName`, `mp`.`project_id` AS `projectId`, `mp`.`project_name` AS `projectName`, `mt`.`meeting_name` AS `meetingName` FROM (SELECT project_id, ROW_NUMBER() OVER (PARTITION BY project_id) AS `rm` FROM `task` `task` WHERE `task`.`company_id` = '0000012183' AND `task`.`level` > 0 AND EXISTS (SELECT `task`.`task_id` FROM `task_result` `result` WHERE `result`.`task_id` = `task`.`task_id`) AND `task`.`status` = 1 AND `task`.`is_draft` = 0) `target` INNER JOIN `meeting_project` `mp` ON `mp`.`project_id` = target.project_id INNER JOIN `meeting` `mt` ON `mp`.`meeting_id` = `mt`.`meeting_id` WHERE target.rm = 1 AND `mp`.`prime_company_id` IN ('0000011131');",
					},
				},
			},
			isNormal: true,
			want:     "cwl2slack executed successfully.",
		},
		{
			name: "slowqueryモードで、渡されたログがスロークエリログの形式でない場合",
			mode: "slowquery",
			data: events.CloudwatchLogsData{
				Owner:     "123456789012",
				LogGroup:  "testLogGroup",
				LogStream: "testLogStream",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{
						ID:        "eventId1",
						Timestamp: 1440442987000,
						Message:   "aiueo",
					},
				},
			},
			isNormal: false,
			want:     "",
		},
	}

	// 環境変数設定
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	// テストケースの実行
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			// モードの設定(環境変数として指定する)
			t.Setenv("MODE", tt.mode)

			// テストケースのデータをCloudwatchLogsRawDataに変換
			d, _ := ConvertToRawData(tt.data)
			event := events.CloudwatchLogsEvent{
				AWSLogs: d,
			}

			got, err := handler(context.Background(), event)

			// 正常系のテストケース
			if tt.isNormal {
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("unexpected result: %+v", got)
				}
				// 異常系のテストケース
			} else {
				if err == nil {
					t.Fatalf("expected error, but got nil")
				}
				t.Logf("aiuel: %v", err)
			}
		})
	}

}

func TestSlowQueryMode(t *testing.T) {

	testCases := []struct {
		name       string
		data       events.CloudwatchLogsData
		wantNormal bool
		want       string
	}{
		{
			name: "渡されたログがスロークエリログの形式でない場合",
			data: events.CloudwatchLogsData{
				Owner:     "123456789012",
				LogGroup:  "testLogGroup",
				LogStream: "testLogStream",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{
						ID:        "eventId1",
						Timestamp: 1440442987000,
						Message:   "aiueo",
					},
				},
			},
			wantNormal: false,
			want:       "",
		},
		{
			name: "渡されたログがスロークエリログの形式の場合",
			data: events.CloudwatchLogsData{
				Owner:     "123456789012",
				LogGroup:  "testLogGroup",
				LogStream: "testLogStream",
				LogEvents: []events.CloudwatchLogsLogEvent{
					{
						ID:        "12345678901234567890123456789012345678901234567890123456",
						Timestamp: 1716792813043,
						Message:   "# Time: 2024-05-27T06:53:33.043104Z\n# User@Host: wsprodadminuser[wsprodadminuser] @ [172.17.0.178] Id: 1436601\n# Query_time: 4.275485 Lock_time: 0.000002 Rows_sent: 58 Rows_examined: 12158\nuse work_prod;\nSET timestamp=1716792808;\nSELECT `mp`.`prime_company_id` AS `primeCompanyId`, `mp`.`prime_company_name` AS `primeCompanyName`, `mp`.`project_id` AS `projectId`, `mp`.`project_name` AS `projectName`, `mt`.`meeting_name` AS `meetingName` FROM (SELECT project_id, ROW_NUMBER() OVER (PARTITION BY project_id) AS `rm` FROM `task` `task` WHERE `task`.`company_id` = '0000012183' AND `task`.`level` > 0 AND EXISTS (SELECT `task`.`task_id` FROM `task_result` `result` WHERE `result`.`task_id` = `task`.`task_id`) AND `task`.`status` = 1 AND `task`.`is_draft` = 0) `target` INNER JOIN `meeting_project` `mp` ON `mp`.`project_id` = target.project_id INNER JOIN `meeting` `mt` ON `mp`.`meeting_id` = `mt`.`meeting_id` WHERE target.rm = 1 AND `mp`.`prime_company_id` IN ('0000011131');",
					},
				},
			},
			wantNormal: true,
			want:       "aiueo",
		},
	}

	// 環境変数設定
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MODE", "slowquery")

	ctx := context.Background()

	// テストケースの実行
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			// テストケースのデータをCloudwatchLogsRawDataに変換
			d, _ := ConvertToRawData(tt.data)
			event := events.CloudwatchLogsEvent{
				AWSLogs: d,
			}

			got, err := handler(ctx, event)

			// 正常系のテストケース
			if tt.wantNormal {
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("unexpected result: %+v", got)
				}
				// 異常系のテストケース
			} else {
				if err == nil {
					t.Fatalf("expected error, but got nil")
				}
				t.Logf("aiuel: %v", err)
			}
		})
	}

}

// func TestMain(t *testing.T) {
// main()
// }
