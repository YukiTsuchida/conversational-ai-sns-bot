package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
)

type SendConversation struct {
	aiSvc ai.Service
}

// ToDo: http/handlers側でまったく同じ型を定義してしまっているので共通化する
type ReplyConversationRequest struct {
	Message    string `json:"message"`
	ErrMessage string `json:"err_message"`
}

func (uc *SendConversation) Execute(ctx context.Context, conversationID *conversation.ID, callbackUrl string) error {
	// 対話型AIのAPIにリクエストを送信する

	aiRes, err := uc.aiSvc.SendRequest(ctx, conversationID)
	if err != nil {
		return err
	}

	// 対話型AIのレスポンスをcallbackする

	replyConversationReq := ReplyConversationRequest{
		Message:    aiRes.Message().ToString(),
		ErrMessage: aiRes.ErrMessage(),
	}

	jsonStr, err := json.Marshal(replyConversationReq)
	if err != nil {
		return err
	}

	newReq, err := http.NewRequest(
		"POST",
		callbackUrl,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return err
	}
	newReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(newReq)
	if err != nil {
		return err
	}
	return nil
}

func NewSendConversation(aiSvc ai.Service) *SendConversation {
	return &SendConversation{aiSvc}
}
