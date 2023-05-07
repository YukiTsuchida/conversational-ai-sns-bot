package twitter

import (
	"context"
	"fmt"
	"strconv"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/conversations"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/twitteraccounts"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"

	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

var _ sns.SNS = (*snsTwitterImpl)(nil)

type snsTwitterImpl struct {
	db *ent.Client
}

func (sns *snsTwitterImpl) FetchAccountByID(ctx context.Context, accountID string) (*sns_model.Account, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil, err
	}
	conversationID, err := account.QueryConversation().FirstID(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// レコードが見つからないケースは問題ない
			return sns_model.NewAccount(account.TwitterAccountID, ""), nil
		}
		return nil, err
	}
	conversationIDStr := strconv.Itoa(conversationID)
	return sns_model.NewAccount(account.TwitterAccountID, conversationIDStr), nil
}

func (sns *snsTwitterImpl) FetchAccountByConversationID(ctx context.Context, conversationID string) (*sns_model.Account, error) {
	conversationIDInt, err := strconv.Atoi(conversationID)
	if err != nil {
		return nil, err
	}
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.HasConversationWith(conversations.IDEQ(conversationIDInt))).First(ctx)
	if err != nil {
		return nil, err
	}
	return sns_model.NewAccount(account.TwitterAccountID, conversationID), nil
}

func (sns *snsTwitterImpl) CreateAccount(ctx context.Context, accountID string, credential sns_model.Credential) error {
	c, ok := credential.(*sns_model.OAuth2Credential)
	if !ok {
		return fmt.Errorf("oauth2 credential parse failed")
	}

	accessToken, refreshToken := c.GetTokens()
	_, err := sns.db.TwitterAccounts.Create().
		SetTwitterAccountID(accountID).
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sns *snsTwitterImpl) GiveAccountConversationID(ctx context.Context, accountID string, conversationID string) error {
	conversationIDInt, err := strconv.Atoi(conversationID)
	if err != nil {
		return err
	}
	err = sns.db.TwitterAccounts.Update().SetConversationID(conversationIDInt).Where(twitteraccounts.TwitterAccountIDEQ(accountID)).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sns *snsTwitterImpl) ExecutePostMessageCmd(ctx context.Context, accountID string, cmd *cmd.PostMessageCommand) (*sns_model.PostMessageResponse, error) {
	return nil, nil
}
func (sns *snsTwitterImpl) ExecuteGetMyMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetMyMessagesCommand) (*sns_model.GetMyMessagesResponse, error) {
	return nil, nil
}
func (sns *snsTwitterImpl) ExecuteGetOtherMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetOtherMessagesCommand) (*sns_model.GetOtherMessagesResponse, error) {
	return nil, nil
}
func (sns *snsTwitterImpl) ExecuteSearchMessageCmd(ctx context.Context, accountID string, cmd *cmd.SearchMessageCommand) (*sns_model.SearchMessageResponse, error) {
	return nil, nil
}
func (sns *snsTwitterImpl) ExecuteGetMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetMyProfileCommand) (*sns_model.GetMyProfileResponse, error) {
	return nil, nil
}
func (sns *snsTwitterImpl) ExecuteGetOthersProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetOthersProfileCommand) (*sns_model.GetOthersProfileResponse, error) {
	return nil, nil
}
func (sns *snsTwitterImpl) ExecuteUpdateMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.UpdateMyProfileCommand) (*sns_model.UpdateMyProfileResponse, error) {
	return nil, nil
}

func NewSNSTwitterImpl(db *ent.Client) sns.SNS {
	return &snsTwitterImpl{db}
}
