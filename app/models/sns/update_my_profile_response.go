package sns

type UpdateMyProfileResponse struct {
	errReason string
}

func NewUpdateMyProfileResponse(errReason string) *UpdateMyProfileResponse {
	return &UpdateMyProfileResponse{errReason}
}

func (response UpdateMyProfileResponse) ErrorOccured() bool {
	return response.errReason != ""
}

func (response UpdateMyProfileResponse) ErrReason() string {
	return response.errReason
}
