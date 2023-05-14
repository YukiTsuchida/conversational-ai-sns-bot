package sns

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
)

type Service interface {
	FetchAccountByID(ctx context.Context, accountID *sns.AccountID) (*sns.Account, error)
	FetchAccountByConversationID(ctx context.Context, conversationID *conversation.ID) (*sns.Account, error)
	CreateAccount(ctx context.Context, accountID *sns.AccountID, credential sns.Credential) error // このinterfaceもうちょっとなんとかしたい #27
	GiveAccountConversationID(ctx context.Context, accountID *sns.AccountID, conversationID *conversation.ID) error
	RemoveAccountConversationID(ctx context.Context, accountID *sns.AccountID) error
	ExecutePostMessageCmd(ctx context.Context, accountID *sns.AccountID, cmd *cmd.PostMessageCommand) (*sns.PostMessageResponse, error)
	ExecuteGetMyMessagesCmd(ctx context.Context, accountID *sns.AccountID, cmd *cmd.GetMyMessagesCommand) (*sns.GetMyMessagesResponse, error)
	ExecuteGetOtherMessagesCmd(ctx context.Context, accountID *sns.AccountID, cmd *cmd.GetOtherMessagesCommand) (*sns.GetOtherMessagesResponse, error)
	ExecuteSearchMessageCmd(ctx context.Context, accountID *sns.AccountID, cmd *cmd.SearchMessageCommand) (*sns.SearchMessageResponse, error)
	ExecuteGetMyProfileCmd(ctx context.Context, accountID *sns.AccountID, cmd *cmd.GetMyProfileCommand) (*sns.GetMyProfileResponse, error)
	ExecuteGetOthersProfileCmd(ctx context.Context, accountID *sns.AccountID, cmd *cmd.GetOthersProfileCommand) (*sns.GetOthersProfileResponse, error)
	ExecuteUpdateMyProfileCmd(ctx context.Context, accountID *sns.AccountID, cmd *cmd.UpdateMyProfileCommand) (*sns.UpdateMyProfileResponse, error)
}
