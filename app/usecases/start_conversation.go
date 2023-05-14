package usecases

import (
	"context"
	"fmt"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/prompt"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/sns"
)

type StartConversation struct {
	sns              sns.SNS
	prompt           prompt.Prompt
	ai               ai.AI
	conversationRepo repositories.Conversation
}

func (uc *StartConversation) Execute(ctx context.Context, accountID string, aiModel string, snsType string, cmdVersion string) error {
	// accountが存在するか確認
	account, err := uc.sns.FetchAccountByID(ctx, accountID)
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
	err = uc.sns.GiveAccountConversationID(ctx, accountID, conversationID)
	if err != nil {
		return err
	}

	// 最初に送る文章を生成する
	msg := uc.prompt.BuildFirstMessage()

	// 会話履歴を追加する
	err = uc.ai.AppendSystemMessage(ctx, conversationID, msg)
	if err != nil {
		return err
	}

	// 会話履歴を結合して対話型AI用のリクエストを生成して送信する
	err = uc.ai.SendRequest(ctx, conversationID)
	if err != nil {
		return err
	}

	return nil
}

func NewStartConversation(sns sns.SNS, prompt prompt.Prompt, ai ai.AI, conversationRepo repositories.Conversation) *StartConversation {
	return &StartConversation{sns, prompt, ai, conversationRepo}
}
