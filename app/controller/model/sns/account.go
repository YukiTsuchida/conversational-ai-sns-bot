package sns

// credentialは入れない
type Account struct {
	id             string
	conversationID string
}

func (account Account) IsInConversations() bool {
	if account.conversationID == "" {
		return false
	}
	return true
}
