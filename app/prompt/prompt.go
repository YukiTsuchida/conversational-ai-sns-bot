package prompt

import (
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
)

type Prompt interface {
	BuildFirstMessage() string
	BuildNextMessageCommandNotFound() string
	BuildNextMessagePostMessage(snsResponse *sns.PostMessageResponse) string
	BuildNextMessageGetMyMessages(snsResponse *sns.GetMyMessagesResponse) string
	BuildNextMessageGetOtherMessages(snsResponse *sns.GetOtherMessagesResponse) string
	BuildNextMessageSearchMessage(snsResponse *sns.SearchMessageResponse) string
	BuildNextMessageGetMyProfile(snsResponse *sns.GetMyProfileResponse) string
	BuildNextMessageGetOthersProfile(snsResponse *sns.GetOthersProfileResponse) string
	BuildNextMessageUpdateMyProfile(snsResponse *sns.UpdateMyProfileResponse) string
	ParseCmdsByMessage(message string) []*cmd.Command
}
