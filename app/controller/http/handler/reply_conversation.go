package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai/chatgpt_3_5_turbo"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd/v0_1"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/service"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns/twitter"
	"github.com/go-chi/chi/v5"
)

type ReplyConversationRequest struct {
	Message    string `json:"message"`
	ErrMessage string `json:"err_message"`
}

func ReplyConversationHandler(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {

	conversationRepo := conversation.NewConversationRepository(db)

	return func(w http.ResponseWriter, r *http.Request) {
		var req ReplyConversationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internalReplyConversationError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Message == "" && req.ErrMessage == "" {
			http.Error(w, "request body params were empty: message, err", http.StatusBadRequest)
			return
		}
		conversationID := chi.URLParam(r, "id")

		// DIするためにconversationを取得する
		conversation, err := conversationRepo.FetchByID(r.Context(), conversationID)
		if err != nil {
			// ToDo: エラーをハンドリングしてレスポンスを変える
			internalReplyConversationError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// DIは一旦ここでやる
		var sns sns.SNS
		var ai ai.AI
		var cmd cmd.Cmd
		if conversation.SNSType() == "twitter" {
			sns = twitter.NewSNSTwitterImpl(db)
		} else {
			http.Error(w, "invalid sns_type: "+conversation.SNSType(), http.StatusBadRequest)
		}
		if conversation.AIModel() == "gpt-3.5-turbo" {
			ai = chatgpt_3_5_turbo.NewAIChatGPT3_5TurboImpl(db)
		} else {
			http.Error(w, "invalid ai_model: "+conversation.AIModel(), http.StatusBadRequest)
			return
		}
		if conversation.CmdVersion() == "v0.1" {
			cmd = v0_1.NewCmdV0_1Impl()
		} else {
			http.Error(w, "invalid ai_model: "+conversation.CmdVersion(), http.StatusBadRequest)
			return
		}
		replyConversationService := service.NewReplyConversationService(sns, cmd, ai, conversationRepo)
		abortConversationService := service.NewAbortConversationService(sns, conversationRepo)

		// エラーがあればconversationをabortする
		if req.ErrMessage != "" {
			abortReason := req.ErrMessage
			err = abortConversationService.AbortConversation(r.Context(), conversationID, abortReason)
			if err != nil {
				internalReplyConversationError(err)
				http.Error(w, "failed to abort_conversation: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ok")
			return
		}

		err = replyConversationService.ReplyConversation(r.Context(), conversationID, req.Message)
		if err != nil {
			// ToDo: エラーの内容に応じてresponseを変える
			internalReplyConversationError(err)
			http.Error(w, "failed to reply_conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func internalReplyConversationError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] StartTwitterConversationHandler() error: %s\n", err.Error())
}
