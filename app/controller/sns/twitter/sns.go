package twitter

import (
	"context"

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
	return nil, nil
}

func (sns *snsTwitterImpl) CreateAccount(ctx context.Context, accountID string, credential string) error {
	_, err := sns.db.TwitterAccounts.Create().SetTwitterAccountID(accountID).SetBearerToken(credential).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sns *snsTwitterImpl) GiveAccountConversationID(ctx context.Context, conversationID string) error {
	return nil
}

func (sns *snsTwitterImpl) ExecuteCmd(ctx context.Context, cmd *cmd.Command) (*sns_model.Response, error) {
	return nil, nil
}

func NewSNSTwitterImpl(db *ent.Client) sns.SNS {
	return &snsTwitterImpl{db}
}
