package ai

type Message struct {
	message string
}

func NewMessage(message string) Message {
	return Message{message: message}
}

func (msg *Message) ToString() string {
	return msg.message
}

func (msg *Message) Append(v string) {
	msg.message += v
}

type Role string

const (
	RoleSystem Role = "system"
	RoleUser   Role = "user"
	RoleAI     Role = "ai"
)

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
