package cmd

import (
	"fmt"
	"strconv"
)

type Command struct {
	cmdType Type
	options map[string]string
}

func NewCommand(cmdType Type, options map[string]string) *Command {
	return &Command{
		cmdType: cmdType,
		options: options,
	}
}

func (c *Command) IsPostActionPurposeLog() bool {
	return c.cmdType == PostActionPurposeLog
}

func (c *Command) IsPostMessage() bool {
	return c.cmdType == PostMessage
}

func (c *Command) IsGetMyMessages() bool {
	return c.cmdType == GetMyMessages
}

func (c *Command) IsGetOtherMessages() bool {
	return c.cmdType == GetOtherMessages
}

func (c *Command) IsSearchMessage() bool {
	return c.cmdType == SearchMessage
}

func (c *Command) IsGetMyProfile() bool {
	return c.cmdType == GetMyProfile
}

func (c *Command) IsGetOthersProfile() bool {
	return c.cmdType == GetOthersProfile
}

func (c *Command) IsUpdateMyProfile() bool {
	return c.cmdType == UpdateMyProfile
}

func (c *Command) Option(optionName string) (string, error) {
	v, found := c.options[optionName]
	if !found {
		return "", fmt.Errorf("option %s is not found", optionName)
	}
	return v, nil
}

func (c *Command) OptionInInt(optionName string) (int, error) {
	v, found := c.options[optionName]
	if !found {
		return 0, fmt.Errorf("option %s is not found", optionName)
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return vInt, nil
}
