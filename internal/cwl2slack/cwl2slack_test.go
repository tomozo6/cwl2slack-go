package cwl2slack

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tomozo6/cwl2slack/pkg/slack"
)

func TestNewCwl2slack(t *testing.T) {

	testCloudwatchLogsData := events.CloudwatchLogsData{}

	testCases := []struct {
		name      string
		mode      string
		threshold float64
		cwld      *events.CloudwatchLogsData
		isNormal  bool
		want      *Cwl2slack
	}{
		{
			name:      "[正常系]plainMode",
			mode:      "plain",
			threshold: 0,
			cwld:      &testCloudwatchLogsData,
			isNormal:  true,
			want:      &Cwl2slack{Mode: "plain", Theashold: 0, Cwld: &testCloudwatchLogsData},
		},
		{
			name:      "[正常系]slowqueryMode",
			mode:      "slowquery",
			threshold: 0,
			cwld:      &testCloudwatchLogsData,
			isNormal:  true,
			want:      &Cwl2slack{Mode: "slowquery", Theashold: 0, Cwld: &testCloudwatchLogsData},
		},
		{
			name:      "[異常系]正しくないmodeが指定された場合",
			mode:      "slowqueryy",
			threshold: 0,
			cwld:      &events.CloudwatchLogsData{},
			isNormal:  false,
			want:      nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewCwl2slack(tt.mode, tt.threshold, tt.cwld)

			// 正常系のテストケース
			if tt.isNormal {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("unexpected result: %v", got)
				}
				// 異常系のテストケース
			} else {
				if err == nil {
					t.Fatalf("expected error, but got nil")
				}
			}
		})
	}
}

func TestGetPlainPayload(t *testing.T) {
	testCloudwatchLogsData := events.CloudwatchLogsData{
		LogGroup:  "testLogGroup",
		LogStream: "testLogStream",
		LogEvents: []events.CloudwatchLogsLogEvent{
			{
				Message: "message1",
			},
			{
				Message: "message2",
			},
		},
	}

	// AttachmentsのFields以外は固定値なのでFieldsのみをテストします
	testCases := []struct {
		name     string
		isNormal bool
		want     []slack.Field
	}{
		{
			name:     "[正常系]テスト",
			isNormal: true,
			want: []slack.Field{
				{
					Title: "Log Group",
					Value: "testLogGroup",
					Short: false,
				},
				{
					Title: "Log Stream",
					Value: "testLogStream",
					Short: false,
				},
				{
					Title: "Log Messages",
					Value: "```\nmessage1\nmessage2\n```",
					Short: false,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// cwl2slackインスタンスの作成
			c, _ := NewCwl2slack("plain", 0, &testCloudwatchLogsData)

			// テスト対象のメソッドを実行してFields部分を取得
			p, err := c.getPlainPayloads()
			got := (*p)[0].Attachments[0].Fields

			// 正常系のテストケース
			if tt.isNormal {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("\n got: %+v;\nwant: %+v", got, tt.want)
				}
				// 異常系のテストケース
			} else {
				if err == nil {
					t.Fatalf("expected error, but got nil")
				}
			}
		})
	}
}
func TestGetSlowQueryPayload(t *testing.T) {
	testCloudwatchLogsData := events.CloudwatchLogsData{
		LogGroup:  "testLogGroup",
		LogStream: "testLogStream",
		LogEvents: []events.CloudwatchLogsLogEvent{
			{
				Message: "# Time: 2024-05-27T06:53:33.043104Z\n# User@Host: wsprodadminuser[wsprodadminuser] @ [172.17.0.178] Id: 1436601\n# Query_time: 4.275485 Lock_time: 0.000002 Rows_sent: 58 Rows_examined: 12158\nuse work_prod;\nSET timestamp=1716792808;\nSELECT `mp`.`prime_company_id` AS `primeCompanyId`, `mp`.`prime_company_name` AS `primeCompanyName`, `mp`.`project_id` AS `projectId`, `mp`.`project_name` AS `projectName`, `mt`.`meeting_name` AS `meetingName` FROM (SELECT project_id, ROW_NUMBER() OVER (PARTITION BY project_id) AS `rm` FROM `task` `task` WHERE `task`.`company_id` = '0000012183' AND `task`.`level` > 0 AND EXISTS (SELECT `task`.`task_id` FROM `task_result` `result` WHERE `result`.`task_id` = `task`.`task_id`) AND `task`.`status` = 1 AND `task`.`is_draft` = 0) `target` INNER JOIN `meeting_project` `mp` ON `mp`.`project_id` = target.project_id INNER JOIN `meeting` `mt` ON `mp`.`meeting_id` = `mt`.`meeting_id` WHERE target.rm = 1 AND `mp`.`prime_company_id` IN ('0000011131');",
			},
		},
	}

	// AttachmentsのFields以外は固定値なのでFieldsのみをテストします
	testCases := []struct {
		name      string
		theashold float64
		isNormal  bool
	}{
		{
			name:      "[正常系]クエリー実行時間がしきい値を超えている場合",
			theashold: 4,
			isNormal:  true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// cwl2slackインスタンスの作成
			c, _ := NewCwl2slack("slowquery", tt.theashold, &testCloudwatchLogsData)

			// テスト対象のメソッドを実行してFields部分を取得
			p, err := c.getPlainPayloads()
			got := len(*p)
			t.Logf("got: %v", len(got))

			// 正常系のテストケース
			if tt.isNormal {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				// if !reflect.DeepEqual(got, tt.want) {
				// t.Fatalf("\n got: %+v;\nwant: %+v", got, tt.want)
				// }
				// 異常系のテストケース
			} else {
				if err == nil {
					t.Fatalf("expected error, but got nil")
				}
			}
		})
	}
}
