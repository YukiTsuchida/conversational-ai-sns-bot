package conversation

import (
	"context"
	"strconv"
	"strings"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/conversations"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
)

// 本来インタフェース切るまでもないが、service側にdbを直接さわるロジックを書きたくないので別packageに切り出した
type ConversationRepository interface {
	Create(ctx context.Context, aiModel string, snsType string, cmdVersion string) (string, error)
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

func NewConversationRepository(db *ent.Client) ConversationRepository {
	return &conversationRepository{db}
}
