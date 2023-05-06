package sns

// credentialは入れない
type Account struct {
	id             string
	conversationID string
}

func NewAccount(id string, conversationID string) *Account {
	return &Account{id, conversationID}
}

func (account Account) IsInConversations() bool {
	if account.conversationID == "" {
		return false
	}
	return true
}
