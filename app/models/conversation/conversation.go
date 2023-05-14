package conversation

import "strconv"

type ID struct {
	id string
}

func NewID(id string) *ID {
	return &ID{id: id}
}

func (id *ID) ToString() string {
	return id.id
}

func (id *ID) ToInt() (int, error) {
	v, err := strconv.Atoi(id.id)
	if err != nil {
		return 0, err
	}
	return v, nil
}

type Conversation struct {
	ID
	aiModel    string
	snsType    string
	cmdVersion string
	isAborted  bool
}

func NewConversation(conversationID string, aiModel string, snsType string, cmdVersion string, isAborted bool) *Conversation {
	return &Conversation{
		ID:         ID{id: conversationID},
		aiModel:    aiModel,
		snsType:    snsType,
		cmdVersion: cmdVersion,
		isAborted:  isAborted,
	}
}

func (c *Conversation) AIModel() string {
	return c.aiModel
}

func (c *Conversation) SNSType() string {
	return c.snsType
}

func (c *Conversation) CmdVersion() string {
	return c.cmdVersion
}

func (c *Conversation) IsAborted() bool {
	return c.isAborted
}
