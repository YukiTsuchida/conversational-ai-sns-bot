package cmd

type SearchMessageCommand struct {
	query      string
	maxResults string
}

func NewSearchMessageCommand(query string, maxResults string) *SearchMessageCommand {
	return &SearchMessageCommand{query, maxResults}
}

func (command SearchMessageCommand) Query() string {
	return command.query
}

func (command SearchMessageCommand) MaxResults() string {
	return command.maxResults
}
