package cmd

type PostActionPurposeLogCommand struct {
	log string
}

func NewPostActionPurposeLogCommand(log string) *PostActionPurposeLogCommand {
	return &PostActionPurposeLogCommand{log}
}

func (command PostActionPurposeLogCommand) Log() string {
	return command.log
}
