package sns

type GetOthersProfileResponse struct {
	userID      string
	name        string
	description string
	errReason   string
}

func NewGetOthersProfileResponse(userID string, name string, description string, errReason string) *GetOthersProfileResponse {
	return &GetOthersProfileResponse{userID, name, description, errReason}
}

func (response GetOthersProfileResponse) UserID() string {
	return response.userID
}

func (response GetOthersProfileResponse) Name() string {
	return response.name
}

func (response GetOthersProfileResponse) Description() string {
	return response.description
}

func (response GetOthersProfileResponse) ErrorOccured() bool {
	return response.errReason != ""
}

func (response GetOthersProfileResponse) ErrReason() string {
	return response.errReason
}
