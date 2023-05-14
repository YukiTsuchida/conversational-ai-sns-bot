package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/sns/twitter"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/usecases"
)

type AbortTwitterConversationRequest struct {
	TwitterID string `json:"twitter_id"`
}

func AbortTwitterConversationHandler(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {

	sns := twitter.NewSNSTwitterImpl(db)
	conversationRepo := conversation.NewConversationRepository(db)
	abortConversationUsecase := usecases.NewAbortConversation(sns, conversationRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		var req AbortTwitterConversationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internalAbortTwitterConversationError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.TwitterID == "" {
			http.Error(w, "invalid request body param: twitter_id", http.StatusBadRequest)
			return
		}

		account, err := sns.FetchAccountByID(r.Context(), req.TwitterID)
		if err != nil {
			internalAbortTwitterConversationError(err)
			http.Error(w, "failed to fetch account: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if !account.IsInConversations() {
			http.Error(w, "account is not in conversation", http.StatusBadRequest)
			return
		}

		err = abortConversationUsecase.Execute(r.Context(), account.ConversationID(), "aborted by user")
		if err != nil {
			// ToDo: エラーの内容に応じてresponseを変える
			internalAbortTwitterConversationError(err)
			http.Error(w, "failed to abort_conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func internalAbortTwitterConversationError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] AbortTwitterConversationHandler() error: %s\n", err.Error())
}
