package sns

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
)

type SNS interface {
	FetchAccountByID(ctx context.Context, accountID string) (*sns.Account, error)
	FetchAccountByConversationID(ctx context.Context, conversationID string) (*sns.Account, error)
	CreateAccount(ctx context.Context, accountID string, credential sns.Credential) error // このinterfaceもうちょっとなんとかしたい
	GiveAccountConversationID(ctx context.Context, accountID string, conversationID string) error
	ExecuteCmd(ctx context.Context, accountID string, cmd *cmd.Command) (*sns.Response, error)
}
