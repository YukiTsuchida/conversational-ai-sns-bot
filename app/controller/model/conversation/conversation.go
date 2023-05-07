package conversation

type Conversation struct {
	conversationID string
	aiModel        string
	snsType        string
	cmdVersion     string
	isAborted      bool
}

func NewConversation(conversationID string, aiModel string, snsType string, cmdVersion string, isAborted bool) *Conversation {
	return &Conversation{
		conversationID: conversationID,
		aiModel:        aiModel,
		snsType:        snsType,
		cmdVersion:     cmdVersion,
		isAborted:      isAborted,
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
