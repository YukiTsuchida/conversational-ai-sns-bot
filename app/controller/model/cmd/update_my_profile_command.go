package cmd

type UpdateMyProfileCommand struct {
	name        string
	description string
}

func NewUpdateMyProfileCommand(name string, description string) *UpdateMyProfileCommand {
	return &UpdateMyProfileCommand{name, description}
}

func (command UpdateMyProfileCommand) Name() string {
	return command.name
}

func (command UpdateMyProfileCommand) Description() string {
	return command.description
}
