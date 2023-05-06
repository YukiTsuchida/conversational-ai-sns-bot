package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/conversation"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai/chatgpt_3_5_turbo"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd/v0_1"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/service"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns/twitter"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
)

type StartTwitterConversationRequest struct {
	TwitterID  string `json:"twitter_id"`
	AIModel    string `json:"ai_model"`
	CmdVersion string `json:"cmd_version"`
}

func StartTwitterConversationHandler(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var req StartTwitterConversationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internalStartTwitterConversationError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.TwitterID == "" {
			http.Error(w, "invalid request body param: twitter_id", http.StatusBadRequest)
			return
		}
		if req.AIModel == "" {
			http.Error(w, "invalid request body param: ai_model", http.StatusBadRequest)
			return
		}
		if req.CmdVersion == "" {
			http.Error(w, "invalid request body param: cmd_version", http.StatusBadRequest)
			return
		}

		// DIは一旦ここでやる
		var conversationRepo conversation.ConversationRepository = conversation.NewConversationRepository(db)
		var sns sns.SNS = twitter.NewSNSTwitterImpl(db)
		var ai ai.AI
		var cmd cmd.Cmd
		if req.AIModel == "gpt-3.5-turbo" {
			ai = chatgpt_3_5_turbo.NewAIChatGPT3_5TurboImpl(db)
		} else {
			http.Error(w, "invalid request body param: ai_model", http.StatusBadRequest)
			return
		}
		if req.CmdVersion == "v0.1" {
			cmd = v0_1.NewCmdV0_1Impl()
		} else {
			http.Error(w, "invalid request body param: cmd_version", http.StatusBadRequest)
			return
		}

		startConvarsationService := service.NewStartConversationService(sns, cmd, ai, conversationRepo)

		err = startConvarsationService.StartConversation(r.Context(), req.TwitterID, req.AIModel, "twitter", req.CmdVersion)
		if err != nil {
			// ToDo: エラーの内容に応じてresponseを変える
			internalStartTwitterConversationError(err)
			http.Error(w, "failed to start_conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func internalStartTwitterConversationError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] StartTwitterConversationHandler() error: %s\n", err.Error())
}