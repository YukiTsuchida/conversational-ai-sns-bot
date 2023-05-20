package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/usecases"
	"github.com/go-chi/chi/v5"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai/chatgpt_3_5_turbo"
)

type SendConversationRequest struct {
	CallBackUrl string `json:"call_back_url"`
}

func SendConversationHandler(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {

	conversationRepo := repositories.NewConversation(db)

	return func(w http.ResponseWriter, r *http.Request) {
		var req SendConversationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internalSendConversationError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		conversationID := conversation.NewID(chi.URLParam(r, "id"))

		if conversationID.ToString() == "" {
			http.Error(w, "invalid request body param: conversation_id", http.StatusBadRequest)
			return
		}
		if req.CallBackUrl == "" {
			http.Error(w, "invalid request body param: call_back_url", http.StatusBadRequest)
			return
		}

		conversation, err := conversationRepo.FetchByID(r.Context(), conversationID)
		if err != nil {
			internalSendConversationError(err)
			http.Error(w, "failed to fetch conversation: id: "+conversationID.ToString()+" "+err.Error(), http.StatusInternalServerError)
			return
		}

		var aiSvc ai.Service
		if conversation.AIModel() == "gpt-3.5-turbo" {
			aiSvc = chatgpt_3_5_turbo.NewAIServiceChatGPT3_5TurboImpl(db)
		} else {
			http.Error(w, "invalid ai_model: "+conversation.AIModel(), http.StatusBadRequest)
			return
		}

		uc := usecases.NewSendConversation(aiSvc)

		err = uc.Execute(r.Context(), &conversation.ID, req.CallBackUrl)
		if err != nil {
			internalSendConversationError(err)
			http.Error(w, "failed to send conversation: id: "+conversationID.ToString()+" "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func internalSendConversationError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] SendConversationHandler() error: %s\n", err.Error())
}
