package slack

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestSendAttachements(t *testing.T) {
	// 環境変数設定
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}
	slackURL := os.Getenv("SLACK_WEBHOOK_URL")
	slackChannel := os.Getenv("SLACK_CHANNEL")

	// テストケース
	testCases := []struct {
		name     string
		payload  Payload
		isNormal bool
	}{
		{
			name: "正常系",
			payload: Payload{
				Username:  "CloudWatch Logs",
				IconEmoji: ":robot_face:",
				Channel:   "#dummy-channel",
				Attachments: []Attachment{
					{
						Title:  "AWS CloudWatch Logs",
						Color:  "danger",
						Footer: "post by cwl2slack",
						Fields: []Field{
							{
								Title: "Log Group",
								Value: "/aws/lambda/lambda-function",
								Short: true,
							},
							{
								Title: "Log Stream",
								Value: "2020/01/01/[$LATEST]xxxxxxxxxxxx",
								Short: true,
							},
						},
					},
				},
			},
			isNormal: true,
		},
	}

	// テストケースの実行
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			slack := Slack{URL: slackURL, Channel: slackChannel}
			err = slack.SendNotification(tt.payload)

			// 正常系のテストケース
			if tt.isNormal {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
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
