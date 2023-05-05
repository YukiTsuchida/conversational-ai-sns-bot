package cmd

import (
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
)

type Cmd interface {
	BuildFirstMessage() string
	BuildNextMessage(cmd *cmd.Command, snsResponse *sns.Response) string
	ParseCmdsByMessage(message string) []cmd.Command
}
