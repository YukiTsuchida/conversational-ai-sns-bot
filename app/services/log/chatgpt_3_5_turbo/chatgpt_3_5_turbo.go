package chatgpt_3_5_turbo

import (
	"context"
	"fmt"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/chatgpt35turboconversationlog"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/log"
)

var _ log.Service = (*logServiceImpl)(nil)

type logServiceImpl struct {
	db *ent.Client
}

func NewLogServiceImpl(db *ent.Client) log.Service {
	return &logServiceImpl{db}
}

func (log *logServiceImpl) CountMessageLog(ctx context.Context, conversationID *conversation.ID) (int, error) {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return -1, err
	}
	return log.db.Chatgpt35TurboConversationLog.Query().Where(chatgpt35turboconversationlog.IDEQ(conversationIDInt)).Count(ctx)
}

func (log *logServiceImpl) FetchMessageLogs(ctx context.Context, conversationID *conversation.ID, page int, size int, sort string) ([]*conversation.ConversationLog, error) {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return nil, err
	}

	orderOpt := ent.Asc(chatgpt35turboconversationlog.FieldCreatedAt)
	if sort == "desc" {
		orderOpt = ent.Desc(chatgpt35turboconversationlog.FieldCreatedAt)
	}

	queryResult, err := log.db.Chatgpt35TurboConversationLog.Query().Where(chatgpt35turboconversationlog.IDEQ(conversationIDInt)).Limit(size).Offset(page * size).Order(orderOpt).All(ctx)
	if err != nil {
		return nil, err
	}

	var logs []*conversation.ConversationLog
	for _, v := range queryResult {
		logs = append(logs, conversation.NewConversationLog(
			fmt.Sprint(v.ID),
			v.Message,
			v.Purpose,
			v.Role.String(),
			v.CreatedAt,
		))
	}

	return logs, nil
}
