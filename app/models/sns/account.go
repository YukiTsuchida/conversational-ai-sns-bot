package sns

import "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

// credentialは入れない
type Account struct {
	id             string
	conversationID *conversation.ID
}

func NewAccount(id string, conversationID *conversation.ID) *Account {
	return &Account{id, conversationID}
}

func (account Account) IsInConversations() bool {
	return account.conversationID != nil
}

func (account Account) ID() string {
	return account.id
}

func (account Account) ConversationID() *conversation.ID {
	return account.conversationID
}
