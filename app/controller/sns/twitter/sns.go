package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"bytes"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/conversations"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent/twitteraccounts"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"

	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

var _ sns.SNS = (*snsTwitterImpl)(nil)

type snsTwitterImpl struct {
	db *ent.Client
}

func (sns *snsTwitterImpl) FetchAccountByID(ctx context.Context, accountID string) (*sns_model.Account, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil, err
	}
	conversationID, err := account.QueryConversation().FirstID(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// レコードが見つからないケースは問題ない
			return sns_model.NewAccount(account.TwitterAccountID, ""), nil
		}
		return nil, err
	}
	conversationIDStr := strconv.Itoa(conversationID)
	return sns_model.NewAccount(account.TwitterAccountID, conversationIDStr), nil
}

func (sns *snsTwitterImpl) FetchAccountByConversationID(ctx context.Context, conversationID string) (*sns_model.Account, error) {
	conversationIDInt, err := strconv.Atoi(conversationID)
	if err != nil {
		return nil, err
	}
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.HasConversationWith(conversations.IDEQ(conversationIDInt))).First(ctx)
	if err != nil {
		return nil, err
	}
	return sns_model.NewAccount(account.TwitterAccountID, conversationID), nil
}

func (sns *snsTwitterImpl) CreateAccount(ctx context.Context, accountID string, credential sns_model.Credential) error {
	c, ok := credential.(*sns_model.OAuth2Credential)
	if !ok {
		return fmt.Errorf("oauth2 credential parse failed")
	}

	accessToken, refreshToken := c.GetTokens()
	_, err := sns.db.TwitterAccounts.Create().
		SetTwitterAccountID(accountID).
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sns *snsTwitterImpl) GiveAccountConversationID(ctx context.Context, accountID string, conversationID string) error {
	conversationIDInt, err := strconv.Atoi(conversationID)
	if err != nil {
		return err
	}
	err = sns.db.TwitterAccounts.Update().SetConversationID(conversationIDInt).Where(twitteraccounts.TwitterAccountIDEQ(accountID)).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sns *snsTwitterImpl) RemoveAccountConversationID(ctx context.Context, accountID string) error {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return err
	}
	_, err = account.Update().ClearConversation().Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

type postMessageRequest struct {
	Data struct{
		Text string `json:"text"`
	}`json:"data"`
}

func (sns *snsTwitterImpl) ExecutePostMessageCmd(ctx context.Context, accountID string, cmd *cmd.PostMessageCommand) (*sns_model.PostMessageResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil,err
	}

	u,err := url.Parse("https://api.twitter.com/2/tweets")
	if err != nil {
		return nil,err
	}

	r := postMessageRequest{}
	r.Data.Text = cmd.Message()

	b,_ := json.Marshal(r)

	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken))

	
	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO: twitter error handling
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get my twitter account failed")
	}

	return sns_model.NewPostMessageResponse(""), nil
}

type getMyMessagesResponse struct {
	Data []struct{
		Text string `json:"text"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) ExecuteGetMyMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetMyMessagesCommand) (*sns_model.GetMyMessagesResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil,err
	}

	userID,err := sns.getUserIDByAccessToken(ctx,account.AccessToken)
	if err != nil {
		return nil, err
	}

	u,err := url.Parse(fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets",userID))
	if err != nil {
		return nil,err
	}

	q := u.Query()
	q.Add("max_results",fmt.Sprintf("%d",cmd.MaxResults())) //min:5, max: 100
	q.Add("exclude","retweets,replies")
	u.RawQuery = q.Encode()

	newReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken))

	
	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO: twitter error handling
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get my twitter account failed")
	}

	var j getMyMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	var messages []string
	for _,m := range j.Data {
		messages = append(messages, m.Text)
	}

	return sns_model.NewGetMyMessagesResponse(messages,""), nil
}

func (sns *snsTwitterImpl) ExecuteGetOtherMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetOtherMessagesCommand) (*sns_model.GetOtherMessagesResponse, error) {
	return nil, nil
}

type searchMessagesResponse struct {
	Data []struct{
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) ExecuteSearchMessageCmd(ctx context.Context, accountID string, cmd *cmd.SearchMessageCommand) (*sns_model.SearchMessageResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil,err
	}

	u,err := url.Parse("https://api.twitter.com/2/tweets/search/recent")
	if err != nil {
		return nil,err
	}

	q := u.Query()
	q.Add("query",fmt.Sprintf("%s -is:retweet",cmd.Query()))
	q.Add("max_results",fmt.Sprintf("%d",cmd.MaxResults()))
	u.RawQuery = q.Encode()

	newReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken))

	
	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO: twitter error handling
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get my twitter account failed")
	}

	var j searchMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	var messages []sns_model.SearchMessageMessage
	for _,m := range j.Data {
		messages = append(messages, sns_model.NewSearchMessageMessage(m.ID,m.Text))
	}

	return sns_model.NewSearchMessageResponse(messages,""), nil
}

type getMyProfileResponse struct {
	Data struct {
		UserName string `json:"username"`
		Description string `json:"description"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) ExecuteGetMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetMyProfileCommand) (*sns_model.GetMyProfileResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil,err
	}

	u,err := url.Parse("https://api.twitter.com/2/users/me")
	if err != nil {
		return nil,err
	}

	q := u.Query()
	q.Add("user.fields","description")
	u.RawQuery = q.Encode()

	newReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken))
	
	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO: twitter error handling
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get my twitter account failed")
	}

	var j getMyProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	return sns_model.NewGetMyProfileResponse(j.Data.UserName,j.Data.Description,""), nil
}

type getOthersProfileResponse struct {
	Data struct {
		ID string `json:"id"`
		UserName string `json:"username"`
		Description string `json:"description"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) ExecuteGetOthersProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetOthersProfileCommand) (*sns_model.GetOthersProfileResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil,err
	}

	u,err := url.Parse(fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s",cmd.UserID()))
	if err != nil {
		return nil,err
	}

	q := u.Query()
	q.Add("user.fields","description")
	u.RawQuery = q.Encode()

	newReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken))
	
	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//TODO: twitter error handling
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get other twitter account failed")
	}

	var j getOthersProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	//MEMO: userIDは1123444555みたいなやつですか？
	return sns_model.NewGetOthersProfileResponse(j.Data.ID,j.Data.UserName,j.Data.Description,""), nil
}
func (sns *snsTwitterImpl) ExecuteUpdateMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.UpdateMyProfileCommand) (*sns_model.UpdateMyProfileResponse, error) {
	return nil, nil
}

type getUserIDByAccessTokenResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) getUserIDByAccessToken(ctx context.Context, accessToken string) (string, error) {
	u,err := url.Parse("https://api.twitter.com/2/users/me")
	if err != nil {
		return "",err
	}

	newReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	
	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//TODO: twitter error handling
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get my twitter account failed")
	}

	var j getUserIDByAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return "", err
	}

	return j.Data.ID, nil
}

func NewSNSTwitterImpl(db *ent.Client) sns.SNS {
	return &snsTwitterImpl{db}
}
