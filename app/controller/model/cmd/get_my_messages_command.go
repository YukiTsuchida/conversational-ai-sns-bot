package cmd

type GetMyMessagesCommand struct {
	maxResults int
}

func NewGetMyMessagesCommand(maxResults int) *GetMyMessagesCommand {
	return &GetMyMessagesCommand{maxResults}
}

func (command GetMyMessagesCommand) MaxResults() int {
	return command.maxResults
}
