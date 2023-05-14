package prompt

import (
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
)

type Service interface {
	BuildSystemMessage() *ai.SystemMessage
	BuildUserMessageCommandNotFoundResult() *ai.UserMessage
	BuildUserMessagePostMessageResult(snsResponse *sns.PostMessageResponse) *ai.UserMessage
	BuildUserMessageGetMyMessagesResult(snsResponse *sns.GetMyMessagesResponse) *ai.UserMessage
	BuildUserMessageGetOtherMessagesResult(snsResponse *sns.GetOtherMessagesResponse) *ai.UserMessage
	BuildUserMessageSearchMessageResult(snsResponse *sns.SearchMessageResponse) *ai.UserMessage
	BuildUserMessageGetMyProfileResult(snsResponse *sns.GetMyProfileResponse) *ai.UserMessage
	BuildUserMessageGetOthersProfileResult(snsResponse *sns.GetOthersProfileResponse) *ai.UserMessage
	BuildUserMessageUpdateMyProfileResult(snsResponse *sns.UpdateMyProfileResponse) *ai.UserMessage
	ParseCmdsByAIMessage(aiMsg *ai.AIMessage) []*cmd.Command
}
