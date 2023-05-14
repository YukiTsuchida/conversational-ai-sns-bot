package repositories

import (
	"context"
	"strconv"
	"strings"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/conversations"
	conversation_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
)

type Conversation interface {
	Create(ctx context.Context, aiModel string, snsType string, cmdVersion string) (*conversation_model.ID, error)
	FetchByID(ctx context.Context, conversationID *conversation_model.ID) (*conversation_model.Conversation, error)
	Abort(ctx context.Context, conversationID *conversation_model.ID, reason string) error
}

type conversation struct {
	db *ent.Client
}

var _ Conversation = (*conversation)(nil)

func (conversationRepo *conversation) Create(ctx context.Context, aiModel string, snsType string, cmdVersion string) (*conversation_model.ID, error) {
	// ent用にデータを整形する、entでは「.」を使えないため全て「_」に置き換える
	aiModelEnt := conversations.AiModel(strings.ReplaceAll(aiModel, ".", "_"))
	snsTypeEnt := conversations.SnsType(snsType)
	cmdVersionEnt := conversations.CmdVersion(strings.ReplaceAll(cmdVersion, ".", "_"))

	c, err := conversationRepo.db.Conversations.Create().SetAiModel(aiModelEnt).SetSnsType(snsTypeEnt).SetCmdVersion(cmdVersionEnt).Save(ctx)
	if err != nil {
		return nil, err
	}
	return conversation_model.NewID(strconv.Itoa(c.ID)), nil
}

func (conversationRepo *conversation) FetchByID(ctx context.Context, conversationID *conversation_model.ID) (*conversation_model.Conversation, error) {
	conversationIDInt, err := conversationID.ToInt()
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
		conversationID.ToString(),
		aiModel,
		snsType,
		cmdVersion,
		isAborted,
	)
	return conversation, nil
}

func (conversationRepo *conversation) Abort(ctx context.Context, conversationID *conversation_model.ID, reason string) error {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return err
	}

	_, err = conversationRepo.db.Conversations.UpdateOneID(conversationIDInt).SetIsAborted(true).SetAbortReason(reason).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewConversation(db *ent.Client) Conversation {
	return &conversation{db}
}
