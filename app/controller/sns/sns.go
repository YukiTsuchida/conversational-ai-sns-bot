package sns

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
)

type SNS interface {
	GetAccountById(ctx context.Context, accountId string) (*sns.Account, error)
	CreateAccount(ctx context.Context, accountId string, credential string) error // このinterfaceもうちょっとなんとかしたい
	GiveAccountConversationId(ctx context.Context, conversationId string) error
	ExecuteCmd(ctx context.Context, cmd *cmd.Command) (*sns.Response, error)
}
