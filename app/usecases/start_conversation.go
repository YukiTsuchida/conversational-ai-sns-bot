package usecases

import (
	"context"
	"fmt"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/prompt"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/queue"

	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
)

type StartConversation struct {
	snsSvc           sns.Service
	promptSvc        prompt.Service
	aiSvc            ai.Service
	queueSvc         queue.Service
	conversationRepo repositories.Conversation
}

func (uc *StartConversation) Execute(ctx context.Context, accountID *sns_model.AccountID, aiModel string, snsType string, cmdVersion string) error {
	// accountが存在するか確認
	account, err := uc.snsSvc.FetchAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	// accountが会話中でないか確認
	if account.IsInConversations() {
		return fmt.Errorf("this account is now in the conversation. if you want to restart, please abort.")
	}

	// conversationを新規作成する
	conversationID, err := uc.conversationRepo.Create(ctx, aiModel, snsType, cmdVersion)
	if err != nil {
		return err
	}

	// accountにconversation_idを付与する(accountとconversationを1:1対応させたいため)
	err = uc.snsSvc.GiveAccountConversationID(ctx, accountID, conversationID)
	if err != nil {
		return err
	}

	// 最初に送る文章を生成する
	msg := uc.promptSvc.BuildSystemMessage()

	// 会話履歴を追加する
	err = uc.aiSvc.AppendSystemMessage(ctx, conversationID, msg)
	if err != nil {
		return err
	}

	// 対話型AIにリクエストを送信するためにqueueにリクエストを積む、queueはconversationID単位でレートリミットをしてくれる
	err = uc.queueSvc.Enqueue(ctx, conversationID)
	if err != nil {
		return err
	}

	return nil
}

func NewStartConversation(snsSvc sns.Service, promptSvc prompt.Service, aiSvc ai.Service, queueSvc queue.Service, conversationRepo repositories.Conversation) *StartConversation {
	return &StartConversation{snsSvc, promptSvc, aiSvc, queueSvc, conversationRepo}
}
