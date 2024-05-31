package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tomozo6/cwl2slack/internal/cwl2slack"
	"github.com/tomozo6/cwl2slack/pkg/myutil"
	"github.com/tomozo6/cwl2slack/pkg/slack"
)

func handler(ctx context.Context, event events.CloudwatchLogsEvent) (string, error) {
	// 環境変数の設定
	slackURL := os.Getenv("SLACK_WEBHOOK_URL")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	mode := os.Getenv("MODE")
	threshold := os.Getenv("THRESHOLD")

	t, err := myutil.StrconvParseFloat(threshold, 64)
	if err != nil {
		return "", err
	}

	// 与えられたイベントをパースする
	cwld, err := event.AWSLogs.Parse()
	if err != nil {
		return "", err
	}

	// cwl2slackインスタンスの作成
	c, err := cwl2slack.NewCwl2slack(mode, t, &cwld)
	if err != nil {
		return "", err
	}

	// slackインスタンスの作成
	s := slack.Slack{
		URL:     slackURL,
		Channel: slackChannel,
	}

	// Slack通知に必要なペイロードを取得
	payloads, err := c.GetSlackPayloads()
	if err != nil {
		return "", err
	}

	// Slack通知
	for _, p := range *payloads {
		err = s.SendNotification(p)
		if err != nil {
			return "", fmt.Errorf("slack notification failed: %s", err)

		}
	}

	return "cwl2slack executed successfully.", nil
}

func main() {
	fmt.Println("Hello, World!")
	lambda.Start(handler)
}
