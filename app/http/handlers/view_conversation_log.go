package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai/chatgpt_3_5_turbo"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/usecases"
	"github.com/ggicci/httpin"
)

type ViewConversationLogRequest struct {
	ConversationID string `in:"query=conversation_id;required"`
	Page           int    `in:"query=page;default=0"`
	Size           int    `in:"query=size;default=10"`
	Sort           string `in:"query=sort;default=asc"`
	Timezone       string `in:"query=timezone;default=Asia/Tokyo"`
}

func ViewConversationLog(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(httpin.Input).(*ViewConversationLogRequest)

		timezone, err := req.convertParameter()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// DI
		var conversationRepo repositories.Conversation = repositories.NewConversation(db)
		var logSvc ai.Service = chatgpt_3_5_turbo.NewAIServiceChatGPT3_5TurboImpl(db)
		var generateSimpleConversationLogHtml = usecases.NewGenerateSimpleConversationLogHtml(logSvc, conversationRepo)

		htmlStr, err := generateSimpleConversationLogHtml.Execute(r.Context(), conversation.NewID(req.ConversationID), req.Page, req.Size, conversation.Sort(req.Sort), timezone)
		if err != nil {
			internalViewConversationLogError(err)
			http.Error(w, "failed to view_conversation_log: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(htmlStr))
	}
}

func (req *ViewConversationLogRequest) convertParameter() (*time.Location, error) {
	if req.Page < 0 {
		return nil, fmt.Errorf("a negative value was specified for 'page'. Please specify a non-negative integer.: %d", req.Page)
	}

	if req.Size < 1 || 500 < req.Size {
		return nil, fmt.Errorf("an inappropriate value was specified for 'size'. Please specify an integer between 1 and 500.: %d", req.Page)
	}

	if req.Sort != string(conversation.SortAsc) && req.Sort != string(conversation.SortDesc) {
		return nil, fmt.Errorf("an inappropriate value was specified for 'sort'. Please specify either 'asc' or 'desc'.: %s", req.Sort)
	}

	timezone, err := time.LoadLocation(req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("an inappropriate value was specified for the timeZone parameter.: %s", req.Timezone)
	}

	return timezone, nil
}

func internalViewConversationLogError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] ViewConversationLogHandler() error: %s\n", err.Error())
}
