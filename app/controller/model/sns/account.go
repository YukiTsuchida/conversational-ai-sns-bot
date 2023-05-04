package sns

// credentialは入れない
type Account struct {
	id             string
	conversationId string
}

func (account Account) IsInConversations() bool {
	if account.conversationId == "" {
		return false
	}
	return true
}
