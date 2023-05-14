package sns

import "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

// 例えばtwitterであれば@hogeの「hoge」が入っているようなイメージ、DBのincrementalなIDが入っているわけではないことに注意する
type AccountID struct {
	id string
}

func NewAccountID(id string) *AccountID {
	return &AccountID{id}
}

func (accountID *AccountID) ToString() string {
	return accountID.id
}

// credentialは入れない
type Account struct {
	AccountID
	conversationID *conversation.ID
}

func NewAccount(id string, conversationID *conversation.ID) *Account {
	return &Account{AccountID{id}, conversationID}
}

func (account Account) IsInConversations() bool {
	return account.conversationID != nil
}

func (account Account) ConversationID() *conversation.ID {
	return account.conversationID
}
