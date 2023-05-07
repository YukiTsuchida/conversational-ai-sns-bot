package service

import (
	"context"
	"fmt"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

type ReplyConversationService struct {
	sns              sns.SNS
	cmd              cmd.Cmd
	ai               ai.AI
	conversationRepo conversation.ConversationRepository
}

func (svc *ReplyConversationService) ReplyConversation(ctx context.Context, conversationID string, message string) error {
	// conversationが存在するか確認
	conversation, err := svc.conversationRepo.FetchByID(ctx, conversationID)
	if err != nil {
		return err
	}

	// conversationが中断されてたら終了
	if conversation.IsAborted() {
		return fmt.Errorf("this conversation is aborted")
	}

	// accountを取得
	account, err := svc.sns.FetchAccountByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}

	// messageからcmdを抽出する
	cmds := svc.cmd.ParseCmdsByMessage(message)

	// purposeコマンドを探す
	var purpose string
	for _, cmd := range cmds {
		if cmd.IsPostActionPurposeLog() {
			purpose, err = cmd.Option("log")
			if err != nil {
				return err
			}
			break
		}
	}

	// AIからのメッセージをログに追加
	err = svc.ai.AppendAIMessage(ctx, conversationID, message, purpose)
	if err != nil {
		return err
	}

	// cmdをSNSに送信して、どんどんログに追加していく
	for _, cmd := range cmds {
		if cmd.IsPostActionPurposeLog() {
			// purposeログは無視
			continue
		}
		snsRes, err := svc.sns.ExecuteCmd(ctx, account.ID(), &cmd)
		if err != nil {
			return err
		}
		nextMessage := svc.cmd.BuildNextMessage(&cmd, snsRes)
		err = svc.ai.AppendUserMessage(ctx, conversationID, nextMessage)
		if err != nil {
			return err
		}
	}

	// 会話履歴を結合して対話型AI用のリクエストを生成して送信する
	err = svc.ai.SendRequest(ctx, conversationID)
	if err != nil {
		return err
	}

	return nil
}

func NewReplyConversationService(sns sns.SNS, cmd cmd.Cmd, ai ai.AI, conversationRepo conversation.ConversationRepository) *ReplyConversationService {
	return &ReplyConversationService{sns, cmd, ai, conversationRepo}
}
