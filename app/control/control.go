package control

import (
	"fmt"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
)

type ConversationController struct {
	TwitterService   sns.Service
	ConversationRepo repositories.Conversation
}

func (c *ConversationController) getService(snsType string) sns.Service {
	switch snsType {
	case "twitter":
		return c.TwitterService
	}
	return nil
}

func internalReplyConversationError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] ReplyConversationHandler() error: %s\n", err.Error())
}
