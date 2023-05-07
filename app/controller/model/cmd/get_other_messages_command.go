package cmd

type GetOtherMessagesCommand struct {
	userID string
}

func NewGetOtherMessagesCommand(userID string) *GetOtherMessagesCommand {
	return &GetOtherMessagesCommand{userID}
}

func (command GetOtherMessagesCommand) UserID() string {
	return command.userID
}
