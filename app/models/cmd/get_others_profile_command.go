package cmd

type GetOthersProfileCommand struct {
	userID string
}

func NewGetOthersProfileCommand(userID string) *GetOthersProfileCommand {
	return &GetOthersProfileCommand{userID}
}

func (command GetOthersProfileCommand) UserID() string {
	return command.userID
}
