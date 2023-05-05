package v0_1

import (
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	cmd_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
)

var _ cmd.Cmd = (*cmdV0_1Impl)(nil)

type cmdV0_1Impl struct {
}

func (cmd *cmdV0_1Impl) BuildFirstMessage() string {
	return ""
}
func (cmd *cmdV0_1Impl) BuildNextMessage(snsResponse *sns.Response) string {
	return ""
}
func (cmd *cmdV0_1Impl) ParseCmdByMessage(message string) *cmd_model.Command {
	return nil
}

func NewCmdV0_1Impl() cmd.Cmd {
	return &cmdV0_1Impl{}
}
