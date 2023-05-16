package sns

type GetMyMessagesResponse struct {
	messages  []string
	errReason string
}

func NewGetMyMessagesResponse(messages []string, errReason string) *GetMyMessagesResponse {
	return &GetMyMessagesResponse{messages, errReason}
}

func (response *GetMyMessagesResponse) AppendMessage(message string) {
	response.messages = append(response.messages, message)
}

func (response *GetMyMessagesResponse) Messages() []string {
	return response.messages
}

func (response *GetMyMessagesResponse) ErrorOccured() bool {
	return response.errReason != ""
}

func (response *GetMyMessagesResponse) ErrReason() string {
	return response.errReason
}
