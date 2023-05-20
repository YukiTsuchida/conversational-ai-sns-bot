package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/log"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/log/chatgpt_3_5_turbo"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/usecases"
	"github.com/ggicci/httpin"
)

type ViewConversationLogRequest struct {
	Id       string `in:"query=id;required"`
	Page     int    `in:"query=page;default=0"`
	Size     int    `in:"query=size;default=10"`
	Sort     string `in:"query=sort;default=asc"`
	Timezone string `in:"query=timezone;default=Asia/Tokyo"`
}

func ViewConversationLog(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(httpin.Input).(*ViewConversationLogRequest)

		err := req.validateParameter()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := template.ParseFiles("/app/http/template/simple_conversation_log_viewer.html")
		if err != nil {
			http.Error(w, "failed to view_conversation_log: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// DI
		var conversationRepo repositories.Conversation = repositories.NewConversation(db)
		var logSvc log.Service = chatgpt_3_5_turbo.NewLogServiceImpl(db)
		var viewConversationLog = usecases.NewViewConversationLog(logSvc, conversationRepo)

		data, err := viewConversationLog.Execute(r.Context(), conversation.NewID(req.Id), req.Page, req.Size, req.Sort, req.Timezone)
		if err != nil {
			http.Error(w, "failed to view_conversation_log: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, data); err != nil {
			http.Error(w, "failed to view_conversation_log: "+err.Error(), http.StatusInternalServerError)
		}

	}
}

func (req *ViewConversationLogRequest) validateParameter() error {
	if req.Page < 0 {
		return fmt.Errorf("a negative value was specified for 'page'. Please specify a non-negative integer.: %d", req.Page)
	}

	if req.Size < 1 || 500 < req.Size {
		return fmt.Errorf("an inappropriate value was specified for 'size'. Please specify an integer between 1 and 500.: %d", req.Page)
	}

	if req.Sort != "asc" && req.Sort != "desc" {
		return fmt.Errorf("an inappropriate value was specified for 'sort'. Please specify either 'asc' or 'desc'.: %s", req.Sort)
	}

	return nil
}
