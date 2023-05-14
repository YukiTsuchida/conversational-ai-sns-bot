package conversation

import (
	"context"
	"strconv"
	"strings"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/conversations"
	conversation_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/conversation"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
)

// 本来インタフェース切るまでもないが、service側にdbを直接さわるロジックを書きたくないので別packageに切り出した
type ConversationRepository interface {
	Create(ctx context.Context, aiModel string, snsType string, cmdVersion string) (string, error)
	FetchByID(ctx context.Context, conversationID string) (*conversation_model.Conversation, error)
	Abort(ctx context.Context, conversationID string, reason string) error
}

type conversationRepository struct {
	db *ent.Client
}

var _ ConversationRepository = (*conversationRepository)(nil)

func (conversationRepo *conversationRepository) Create(ctx context.Context, aiModel string, snsType string, cmdVersion string) (string, error) {
	// ent用にデータを整形する、entでは「.」を使えないため全て「_」に置き換える
	aiModelEnt := conversations.AiModel(strings.ReplaceAll(aiModel, ".", "_"))
	snsTypeEnt := conversations.SnsType(snsType)
	cmdVersionEnt := conversations.CmdVersion(strings.ReplaceAll(cmdVersion, ".", "_"))

	c, err := conversationRepo.db.Conversations.Create().SetAiModel(aiModelEnt).SetSnsType(snsTypeEnt).SetCmdVersion(cmdVersionEnt).Save(ctx)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(c.ID), nil
}

func (conversationRepo *conversationRepository) FetchByID(ctx context.Context, conversationID string) (*conversation_model.Conversation, error) {
	conversationIDInt, err := strconv.Atoi(conversationID)
	c, err := conversationRepo.db.Conversations.Get(ctx, conversationIDInt)
	if err != nil {
		return nil, err
	}

	// entでは「.」を使えないため全て「_」に置き換えて入れている、ここで「.」に戻す
	aiModel := strings.ReplaceAll(c.AiModel.String(), "_", ".")
	snsType := c.SnsType.String()
	cmdVersion := strings.ReplaceAll(c.CmdVersion.String(), "_", ".")
	isAborted := c.IsAborted

	conversation := conversation_model.NewConversation(
		conversationID,
		aiModel,
		snsType,
		cmdVersion,
		isAborted,
	)
	return conversation, nil
}

func (conversationRepo *conversationRepository) Abort(ctx context.Context, conversationID string, reason string) error {
	conversationIDInt, err := strconv.Atoi(conversationID)
	if err != nil {
		return err
	}

	_, err = conversationRepo.db.Conversations.UpdateOneID(conversationIDInt).SetIsAborted(true).SetAbortReason(reason).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewConversationRepository(db *ent.Client) ConversationRepository {
	return &conversationRepository{db}
}
