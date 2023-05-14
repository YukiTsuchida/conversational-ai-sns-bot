package ai

type Message struct {
	message string
}

func (msg *Message) ToString() string {
	return msg.message
}

func (msg *Message) Append(v string) {
	msg.message += v
}

type SystemMessage struct {
	Message
}

func NewSystemMessage(msg string) *SystemMessage {
	return &SystemMessage{Message{msg}}
}

type UserMessage struct {
	Message
}

func NewUserMessage(msg string) *UserMessage {
	return &UserMessage{Message{msg}}
}

type AIMessage struct {
	Message
}

func NewAIMessage(msg string) *AIMessage {
	return &AIMessage{Message{msg}}
}
