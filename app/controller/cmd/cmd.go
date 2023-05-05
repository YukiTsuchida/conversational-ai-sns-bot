package cmd

import (
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
)

type Cmd interface {
	BuildFirstMessage() string
	BuildNextMessage(snsResponse *sns.Response) string
	ParseCmdByMessage(message string) *cmd.Command
}
