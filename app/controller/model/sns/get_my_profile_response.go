package sns

type GetMyProfileResponse struct {
	name        string
	description string
	errReason   string
}

func NewGetMyProfileResponse(name string, description string, errReason string) *GetMyProfileResponse {
	return &GetMyProfileResponse{name, description, errReason}
}

func (response GetMyProfileResponse) Name() string {
	return response.name
}

func (response GetMyProfileResponse) Description() string {
	return response.description
}

func (response GetMyProfileResponse) ErrorOccured() bool {
	return response.errReason != ""
}

func (response GetMyProfileResponse) ErrReason() string {
	return response.errReason
}
