package sns

type Credential interface{}

type OAuth2Credential struct {
	accessToken  string
	refreshToken string
}

func (c *OAuth2Credential) GetTokens() (string, string) {
	return c.accessToken, c.refreshToken
}

func NewOAuth2Credential(accessToken, refreshToken string) *OAuth2Credential {
	return &OAuth2Credential{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}
