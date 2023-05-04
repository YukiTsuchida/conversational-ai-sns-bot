package service

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	ai_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

type StartConversationService struct {
	sns sns.SNS
	cmd cmd.Cmd
	ai  ai.AI
}

func (svc *StartConversationService) StartConversation(ctx context.Context, accountId string) error {
	// accountが存在するか確認
	account, err := svc.sns.GetAccountById(ctx, accountId)
	if err != nil {
		return err
	}

	// accountが会話中でないか確認
	if account.IsInConversations() {
		return nil // ToDo: 既に会話中ですよエラーを返す
	}

	// Todo: conversationを新規作成する
	conversationId := "test"

	// accountにconversation_idを付与する(accountとconversationを1:1対応させたいため)
	err = svc.sns.GiveAccountConversationId(ctx, conversationId)
	if err != nil {
		return err
	}

	// 最初に送る文章を生成する
	msg := svc.cmd.BuildFirstMessage()

	// 会話履歴をDBに積む
	err = svc.ai.SaveMessageLog(ctx, *msg, ai_model.System)
	if err != nil {
		return err
	}

	// 会話履歴を結合してAI用のリクエストを生成して送信する
	err = svc.ai.SendRequest(ctx, conversationId)
	if err != nil {
		return err
	}

	return nil
}

func NewStartConversationService(sns sns.SNS, cmd cmd.Cmd, ai ai.AI) *StartConversationService {
	return &StartConversationService{sns, cmd, ai}
}
