package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/queue"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/queue/cloud_tasks"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	ai_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai/chatgpt_3_5_turbo"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/prompt"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/prompt/v0_1"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns/twitter"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/usecases"
	"github.com/go-chi/chi/v5"
)

type ReplyConversationRequest struct {
	Message    string `json:"message"`
	ErrMessage string `json:"err_message"`
}

func ReplyConversationHandler(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {

	conversationRepo := repositories.NewConversation(db)

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
		conversationID := conversation.NewID(chi.URLParam(r, "id"))

		fmt.Println(req.Message)

		// DIするためにconversationを取得する
		conversation, err := conversationRepo.FetchByID(r.Context(), conversationID)
		if err != nil {
			// ToDo: エラーをハンドリングしてレスポンスを変える
			internalReplyConversationError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// DIは一旦ここでやる
		var queueSvc queue.Service = cloud_tasks.NewQueueServiceCloudTasksImpl()
		var snsSvc sns.Service
		var aiSvc ai.Service
		var promptSvc prompt.Service
		if conversation.SNSType() == "twitter" {
			snsSvc = twitter.NewSNSServiceTwitterImpl(db)
		} else {
			http.Error(w, "invalid sns_type: "+conversation.SNSType(), http.StatusBadRequest)
		}
		if conversation.AIModel() == "gpt-3.5-turbo" {
			aiSvc = chatgpt_3_5_turbo.NewAIServiceChatGPT3_5TurboImpl(db)
		} else {
			http.Error(w, "invalid ai_model: "+conversation.AIModel(), http.StatusBadRequest)
			return
		}
		if conversation.CmdVersion() == "v0.1" {
			promptSvc = v0_1.NewPromptServiceV0_1Impl()
		} else {
			http.Error(w, "invalid ai_model: "+conversation.CmdVersion(), http.StatusBadRequest)
			return
		}
		replyConversationUsecase := usecases.NewReplyConversation(snsSvc, promptSvc, aiSvc, queueSvc, conversationRepo)
		abortConversationUsecase := usecases.NewAbortConversation(snsSvc, conversationRepo)

		// エラーがあればconversationをabortする
		if req.ErrMessage != "" {
			abortReason := req.ErrMessage
			err = abortConversationUsecase.Execute(r.Context(), conversationID, abortReason)
			if err != nil {
				internalReplyConversationError(err)
				http.Error(w, "failed to abort_conversation: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ok")
			return
		}

		err = replyConversationUsecase.Execute(r.Context(), conversationID, ai_model.NewAIMessage(req.Message))
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
	fmt.Fprintf(os.Stderr, "[ERROR] ReplyConversationHandler() error: %s\n", err.Error())
}
