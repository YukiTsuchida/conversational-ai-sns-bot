package sns

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
)

type SNS interface {
	GetAccountById(context context.Context, accountId string) (*sns.Account, error)
	CreateAccount(context context.Context, accountId string, credential string) error // このinterfaceもうちょっとなんとかしたい
	GiveAccountConversationId(context context.Context, conversationId string) error
	DoCmd(context context.Context, cmd *cmd.Command) (*sns.Response, error)
}
