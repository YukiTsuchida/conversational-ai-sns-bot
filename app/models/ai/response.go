package ai

type Response struct {
	message    *AIMessage
	errMessage string
}

func NewResponse(message string, errMessage string) *Response {
	return &Response{NewAIMessage(message), errMessage}
}

func (r *Response) Message() *AIMessage {
	return r.message
}

func (r *Response) ErrMessage() string {
	return r.errMessage
}
