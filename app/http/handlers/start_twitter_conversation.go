package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/queue"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/queue/cloud_tasks"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/prompt"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/usecases"

	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai/chatgpt_3_5_turbo"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/prompt/v0_1"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns/twitter"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
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
		var queueSvc queue.Service = cloud_tasks.NewQueueServiceCloudTasksImpl()
		var conversationRepo repositories.Conversation = repositories.NewConversation(db)
		var snsSvc sns.Service = twitter.NewSNSServiceTwitterImpl(db)
		var aiSvc ai.Service
		var promptSvc prompt.Service
		if req.AIModel == "gpt-3.5-turbo" {
			aiSvc = chatgpt_3_5_turbo.NewAIServiceChatGPT3_5TurboImpl(db)
		} else {
			http.Error(w, "invalid request body param: ai_model", http.StatusBadRequest)
			return
		}
		if req.CmdVersion == "v0.1" {
			promptSvc = v0_1.NewPromptServiceV0_1Impl()
		} else {
			http.Error(w, "invalid request body param: cmd_version", http.StatusBadRequest)
			return
		}

		startConvarsationUsecase := usecases.NewStartConversation(snsSvc, promptSvc, aiSvc, queueSvc, conversationRepo)

		err = startConvarsationUsecase.Execute(r.Context(), sns_model.NewAccountID(req.TwitterID), req.AIModel, "twitter", req.CmdVersion)
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
