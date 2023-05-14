package cmd

type SearchMessageCommand struct {
	query      string
	maxResults int
}

func NewSearchMessageCommand(query string, maxResults int) *SearchMessageCommand {
	return &SearchMessageCommand{query, maxResults}
}

func (command SearchMessageCommand) Query() string {
	return command.query
}

func (command SearchMessageCommand) MaxResults() int {
	return command.maxResults
}
