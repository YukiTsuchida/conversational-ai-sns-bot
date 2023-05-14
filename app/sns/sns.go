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
	RemoveAccountConversationID(ctx context.Context, accountID string) error
	ExecutePostMessageCmd(ctx context.Context, accountID string, cmd *cmd.PostMessageCommand) (*sns.PostMessageResponse, error)
	ExecuteGetMyMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetMyMessagesCommand) (*sns.GetMyMessagesResponse, error)
	ExecuteGetOtherMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetOtherMessagesCommand) (*sns.GetOtherMessagesResponse, error)
	ExecuteSearchMessageCmd(ctx context.Context, accountID string, cmd *cmd.SearchMessageCommand) (*sns.SearchMessageResponse, error)
	ExecuteGetMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetMyProfileCommand) (*sns.GetMyProfileResponse, error)
	ExecuteGetOthersProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetOthersProfileCommand) (*sns.GetOthersProfileResponse, error)
	ExecuteUpdateMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.UpdateMyProfileCommand) (*sns.UpdateMyProfileResponse, error)
}
