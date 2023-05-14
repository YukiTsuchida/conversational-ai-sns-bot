package sns

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
)

type Service interface {
	FetchAccountByID(ctx context.Context, accountID string) (*sns.Account, error)
	FetchAccountByConversationID(ctx context.Context, conversationID *conversation.ID) (*sns.Account, error)
	CreateAccount(ctx context.Context, accountID string, credential sns.Credential) error // このinterfaceもうちょっとなんとかしたい
	GiveAccountConversationID(ctx context.Context, accountID string, conversationID *conversation.ID) error
	RemoveAccountConversationID(ctx context.Context, accountID string) error
	ExecutePostMessageCmd(ctx context.Context, accountID string, cmd *cmd.PostMessageCommand) (*sns.PostMessageResponse, error)
	ExecuteGetMyMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetMyMessagesCommand) (*sns.GetMyMessagesResponse, error)
	ExecuteGetOtherMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetOtherMessagesCommand) (*sns.GetOtherMessagesResponse, error)
	ExecuteSearchMessageCmd(ctx context.Context, accountID string, cmd *cmd.SearchMessageCommand) (*sns.SearchMessageResponse, error)
	ExecuteGetMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetMyProfileCommand) (*sns.GetMyProfileResponse, error)
	ExecuteGetOthersProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetOthersProfileCommand) (*sns.GetOthersProfileResponse, error)
	ExecuteUpdateMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.UpdateMyProfileCommand) (*sns.UpdateMyProfileResponse, error)
}
