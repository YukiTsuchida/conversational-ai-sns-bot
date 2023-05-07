package sns

type GetOtherMessagesMessage struct {
	userID  string
	message string
}

func NewGetOtherMessagesMessage(userID string, message string) GetOtherMessagesMessage {
	return GetOtherMessagesMessage{userID, message}
}

func (message GetOtherMessagesMessage) UserID() string {
	return message.userID
}

func (message GetOtherMessagesMessage) Message() string {
	return message.message
}

type GetOtherMessagesResponse struct {
	messages  []GetOtherMessagesMessage
	errReason string
}

func NewGetOtherMessagesResponse(messages []GetOtherMessagesMessage, errReason string) *GetOtherMessagesResponse {
	return &GetOtherMessagesResponse{messages, errReason}
}

func (response GetOtherMessagesResponse) Messages() []GetOtherMessagesMessage {
	return response.messages
}

func (response GetOtherMessagesResponse) ErrorOccured() bool {
	return response.errReason != ""
}

func (response GetOtherMessagesResponse) ErrReason() string {
	return response.errReason
}
