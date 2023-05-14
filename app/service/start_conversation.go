package service

import (
	"context"
	"fmt"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/conversation"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

type StartConversationService struct {
	sns              sns.SNS
	cmd              cmd.Cmd
	ai               ai.AI
	conversationRepo conversation.ConversationRepository
}

func (svc *StartConversationService) StartConversation(ctx context.Context, accountID string, aiModel string, snsType string, cmdVersion string) error {
	// accountが存在するか確認
	account, err := svc.sns.FetchAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	// accountが会話中でないか確認
	if account.IsInConversations() {
		return fmt.Errorf("this account is now in the conversation. if you want to restart, please abort.")
	}

	// conversationを新規作成する
	conversationID, err := svc.conversationRepo.Create(ctx, aiModel, snsType, cmdVersion)
	if err != nil {
		return err
	}

	// accountにconversation_idを付与する(accountとconversationを1:1対応させたいため)
	err = svc.sns.GiveAccountConversationID(ctx, accountID, conversationID)
	if err != nil {
		return err
	}

	// 最初に送る文章を生成する
	msg := svc.cmd.BuildFirstMessage()

	// 会話履歴を追加する
	err = svc.ai.AppendSystemMessage(ctx, conversationID, msg)
	if err != nil {
		return err
	}

	// 会話履歴を結合して対話型AI用のリクエストを生成して送信する
	err = svc.ai.SendRequest(ctx, conversationID)
	if err != nil {
		return err
	}

	return nil
}

func NewStartConversationService(sns sns.SNS, cmd cmd.Cmd, ai ai.AI, conversationRepo conversation.ConversationRepository) *StartConversationService {
	return &StartConversationService{sns, cmd, ai, conversationRepo}
}
