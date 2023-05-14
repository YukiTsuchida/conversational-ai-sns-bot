package sns

type PostMessageResponse struct {
	errReason string
}

func NewPostMessageResponse(errReason string) *PostMessageResponse {
	return &PostMessageResponse{errReason}
}

func (response PostMessageResponse) ErrorOccured() bool {
	return response.errReason != ""
}

func (response PostMessageResponse) ErrReason() string {
	return response.errReason
}
