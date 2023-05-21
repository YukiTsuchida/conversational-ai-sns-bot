package control

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/http/handlers"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/go-chi/chi/v5"
)

func (c *ConversationController) ReplyConversationHandler(w http.ResponseWriter, r *http.Request) {
	var req handlers.ReplyConversationRequest
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

	conversation, err := c.ConversationRepo.FetchByID(r.Context(), conversationID)
	if err != nil {
		// ToDo: エラーをハンドリングしてレスポンスを変える
		internalReplyConversationError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	service := c.getService(conversation.SNSType())

	fmt.Println(service)
}
