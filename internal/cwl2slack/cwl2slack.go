package cwl2slack

import (
	"fmt"
	"slices"
	"strconv"

	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tomozo6/cwl2slack/pkg/slack"
)

type Cwl2slack struct {
	Mode      string
	Theashold float64
	Cwld      *events.CloudwatchLogsData
}

// NewCwl2slackはCwl2slackのコンストラクタ
func NewCwl2slack(m string, t float64, c *events.CloudwatchLogsData) (*Cwl2slack, error) {
	// Modeが想定外の値の場合はエラーを返す
	modes := []string{"plain", "slowquery"}

	if !slices.Contains(modes, m) {
		return nil, fmt.Errorf("invalid mode: %s", m)
	}

	return &Cwl2slack{
		Mode:      m,
		Theashold: t,
		Cwld:      c,
	}, nil
}

// Slack通知に必要なペイロードの配列を返します
func (c *Cwl2slack) GetSlackPayloads() (*[]slack.Payload, error) {

	switch c.Mode {
	case "plain":
		return c.getPlainPayloads()
	case "slowquery":
		return c.getSlowQueryPayloads()
	default:
		return nil, fmt.Errorf("invalid mode: %s", c.Mode)
	}
}

// plainモードのSlack通知に必要なペイロードの配列を返します
// (plainモードはメッセージを結合するので、配列の長さは必ず1になります
// 配列にする意味は無いのですが、他のモードとの互換性を保つために配列にしています)
func (c *Cwl2slack) getPlainPayloads() (*[]slack.Payload, error) {

	// ログイベントのメッセージを取得します
	// ログイベントが複数ある場合は全て取得して結合します
	messages := make([]string, len(c.Cwld.LogEvents))
	for i, e := range c.Cwld.LogEvents {
		messages[i] = e.Message
	}
	joinedMessages := strings.Join(messages, "\n")

	return &[]slack.Payload{
		{
			Username:  "CloudWatch Logs",
			IconEmoji: ":robot_face:",
			Attachments: []slack.Attachment{
				{
					Title:  ":rotating_light:CloudWatchLogsにてアラートを検知しました",
					Color:  "danger",
					Footer: "post by cwl2slack",
					Fields: []slack.Field{
						{
							Title: "Log Group",
							Value: c.Cwld.LogGroup,
							Short: false,
						},
						{
							Title: "Log Stream",
							Value: c.Cwld.LogStream,
							Short: false,
						},
						{
							Title: "Log Messages",
							Value: "```\n" + joinedMessages + "\n```",
							Short: false,
						},
					},
				},
			},
		},
	}, nil
}

func (c *Cwl2slack) getSlowQueryPayloads() (*[]slack.Payload, error) {

	// ログイベントの数だけペイロードを作成します
	payloads := make([]slack.Payload, len(c.Cwld.LogEvents))

	//	ログイベントのメッセージを取得します
	for i, e := range c.Cwld.LogEvents {

		// スロークエリーの情報を取得します
		sq, err := NewSlowQuery(e.Message)
		if err != nil {
			return nil, err
		}

		// スロークエリーの実行時間が閾値を超えていない場合はスキップします
		if sq.QueryTime < c.Theashold {
			continue
		}

		payloads[i] = slack.Payload{
			Username:  "CloudWatch Logs",
			IconEmoji: ":turtle:",
			Attachments: []slack.Attachment{
				{
					Title:  fmt.Sprintf(":rotating_light:ロググループ %s にて閾値を超えたスロークエリーが検知されました", c.Cwld.LogGroup),
					Color:  "danger",
					Footer: "post by cwl2slack",
					Fields: []slack.Field{
						{
							Title: "タイムスタンプ",
							Value: sq.Time,
							Short: true,
						},
						{
							Title: "クエリ実行ユーザ",
							Value: sq.User,
							Short: true,
						},
						{
							Title: "クエリ実行時間",
							Value: strconv.FormatFloat(sq.QueryTime, 'f', -1, 64),
							Short: true,
						},
						{
							Title: "通知閾値",
							Value: "",
							Short: true,
						},
						{
							Title: "ロック取得までの時間",
							Value: sq.LockTime,
							Short: true,
						},
						{
							Title: "クライアントへ送信した行数",
							Value: sq.RowsSent,
							Short: true,
						},
						{
							Title: "クエリ実行時にスキャンした行数",
							Value: sq.RowsExamined,
							Short: true,
						},
						{
							Title: "実行したクエリ",
							Value: "```\n" + sq.Query + "\n```",
							Short: false,
						},
					},
				},
			},
		}
	}
	return &payloads, nil
}
