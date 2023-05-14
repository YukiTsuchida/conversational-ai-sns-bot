package cmd

type PostMessageCommand struct {
	message string
}

func NewPostMessageCommand(message string) *PostMessageCommand {
	return &PostMessageCommand{message}
}

func (command PostMessageCommand) Message() string {
	return command.message
}
