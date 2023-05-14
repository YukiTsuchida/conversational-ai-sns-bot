package cmd

type GetOtherMessagesCommand struct {
	userID     string
	maxResults int
}

func NewGetOtherMessagesCommand(userID string, maxResults int) *GetOtherMessagesCommand {
	return &GetOtherMessagesCommand{userID, maxResults}
}

func (command GetOtherMessagesCommand) UserID() string {
	return command.userID
}

func (command GetOtherMessagesCommand) MaxResults() int {
	return command.maxResults
}
