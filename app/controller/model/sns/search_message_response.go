package sns

type SearchMessageMessage struct {
	userID  string
	message string
}

func NewSearchMessageMessage(userID string, message string) SearchMessageMessage {
	return SearchMessageMessage{userID, message}
}

func (message SearchMessageMessage) UserID() string {
	return message.userID
}

func (message SearchMessageMessage) Message() string {
	return message.message
}

type SearchMessageResponse struct {
	messages  []SearchMessageMessage
	errReason string
}

func NewSearchMessageResponse(messages []SearchMessageMessage, errReason string) *SearchMessageResponse {
	return &SearchMessageResponse{messages, errReason}
}

func (response SearchMessageResponse) Messages() []SearchMessageMessage {
	return response.messages
}

func (response SearchMessageResponse) ErrorOccured() bool {
	return response.errReason != ""
}

func (response SearchMessageResponse) ErrReason() string {
	return response.errReason
}
