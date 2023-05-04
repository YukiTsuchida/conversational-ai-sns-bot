package twitter

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"

	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

var _ sns.SNS = (*snsTwitterImpl)(nil)

type snsTwitterImpl struct {
}

func (sns *snsTwitterImpl) GetAccountById(context context.Context, accountId string) (*sns_model.Account, error) {
	return nil, nil
}

func (sns *snsTwitterImpl) CreateAccount(context context.Context, accountId string, credential string) error {
	return nil
}

func (sns *snsTwitterImpl) GiveAccountConversationId(context context.Context, conversationId string) error {
	return nil
}

func (sns *snsTwitterImpl) DoCmd(context context.Context, cmd *cmd.Command) (*sns_model.Response, error) {
	return nil, nil
}

func NewSNSTwitterImpl() *snsTwitterImpl {
	return &snsTwitterImpl{}
}
