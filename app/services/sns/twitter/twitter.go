package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/conversations"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/twitteraccounts"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
)

var _ sns.Service = (*snsServiceTwitterImpl)(nil)

type snsServiceTwitterImpl struct {
	db *ent.Client
}

func (sns *snsServiceTwitterImpl) FetchAccountByID(ctx context.Context, accountID *sns_model.AccountID) (*sns_model.Account, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
	if err != nil {
		return nil, err
	}
	conversationIDInt, err := account.QueryConversation().FirstID(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// レコードが見つからないケースは問題ない
			return sns_model.NewAccount(account.TwitterAccountID, nil), nil
		}
		return nil, err
	}
	conversationID := conversation.NewID(strconv.Itoa(conversationIDInt))
	return sns_model.NewAccount(account.TwitterAccountID, conversationID), nil
}

func (sns *snsServiceTwitterImpl) FetchAccountByConversationID(ctx context.Context, conversationID *conversation.ID) (*sns_model.Account, error) {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return nil, err
	}
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.HasConversationWith(conversations.IDEQ(conversationIDInt))).First(ctx)
	if err != nil {
		return nil, err
	}
	return sns_model.NewAccount(account.TwitterAccountID, conversationID), nil
}

func (sns *snsServiceTwitterImpl) CreateAccount(ctx context.Context, accountID *sns_model.AccountID, credential sns_model.Credential) error {
	c, ok := credential.(*sns_model.OAuth2Credential)
	if !ok {
		return fmt.Errorf("oauth2 credential parse failed")
	}

	accessToken, refreshToken := c.GetTokens()
	_, err := sns.db.TwitterAccounts.Create().
		SetTwitterAccountID(accountID.ToString()).
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sns *snsServiceTwitterImpl) GiveAccountConversationID(ctx context.Context, accountID *sns_model.AccountID, conversationID *conversation.ID) error {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return err
	}
	err = sns.db.TwitterAccounts.Update().SetConversationID(conversationIDInt).Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (sns *snsServiceTwitterImpl) RemoveAccountConversationID(ctx context.Context, accountID *sns_model.AccountID) error {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
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
	Text string `json:"text"`
}

func (sns *snsServiceTwitterImpl) ExecutePostMessageCmd(ctx context.Context, accountID *sns_model.AccountID, cmd *cmd.PostMessageCommand) (*sns_model.PostMessageResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("https://api.twitter.com/2/tweets")
	if err != nil {
		return nil, err
	}

	r := postMessageRequest{}
	r.Text = cmd.Message()

	b, _ := json.Marshal(r)

	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken)) // ToDo: アクセストークンのリフレッシュを実装する #30

	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: errorハンドリングもっと丁寧にする
	if resp.StatusCode != http.StatusCreated {
		// respをdumpする
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(b))
		}
		return sns_model.NewPostMessageResponse("post tweet failed"), fmt.Errorf("twitter API error")
	}

	return sns_model.NewPostMessageResponse(""), nil
}

type getMyMessagesResponse struct {
	Data []struct {
		Text string `json:"text"`
	} `json:"data"`
}

func (sns *snsServiceTwitterImpl) ExecuteGetMyMessagesCmd(ctx context.Context, accountID *sns_model.AccountID, cmd *cmd.GetMyMessagesCommand) (*sns_model.GetMyMessagesResponse, error) {
	maxResults := cmd.MaxResults()
	if maxResults < 5 || 10 < maxResults {
		maxResults = 5
	}
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := sns.getUserIDByAccessToken(ctx, account.AccessToken)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets", userID))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("max_results", fmt.Sprintf("%d", maxResults)) //min:5, max: 100
	q.Add("exclude", "retweets,replies")
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

	// TODO: errorハンドリングもっと丁寧にする
	if resp.StatusCode != http.StatusOK {
		// respをdumpする
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(b))
		}
		return sns_model.NewGetMyMessagesResponse(nil, "get my messages failed"), fmt.Errorf("twitter API error")
	}

	var j getMyMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	var messages []string
	for _, m := range j.Data {
		messages = append(messages, m.Text)
	}

	return sns_model.NewGetMyMessagesResponse(messages, ""), nil
}

type getOtherMessagesResponse struct {
	Data []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

func (sns *snsServiceTwitterImpl) ExecuteGetOtherMessagesCmd(ctx context.Context, accountID *sns_model.AccountID, cmd *cmd.GetOtherMessagesCommand) (*sns_model.GetOtherMessagesResponse, error) {
	if !validateUserID(cmd.UserID()) {
		return sns_model.NewGetOtherMessagesResponse(nil, "Request error: user_id must consists of digits only"), nil
	}
	maxResults := cmd.MaxResults()
	if maxResults < 5 || 10 < maxResults {
		maxResults = 5
	}
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(cmd.UserID())
	fmt.Println(cmd.MaxResults())

	u, err := url.Parse(fmt.Sprintf("https://api.twitter.com/2/users/%s/tweets", cmd.UserID()))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("max_results", fmt.Sprintf("%d", maxResults)) //min:5, max: 100
	q.Add("exclude", "retweets,replies")
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

	// TODO: errorハンドリングもっと丁寧にする
	if resp.StatusCode != http.StatusOK {
		// respをdumpする
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(b))
		}
		if resp.StatusCode == http.StatusNotFound {
			return sns_model.NewGetOtherMessagesResponse(nil, "user_id not found."), nil
		}
		if resp.StatusCode == http.StatusBadRequest {
			return sns_model.NewGetOtherMessagesResponse(nil, "One or more parameters to your request was invalid."), nil
		}
		return nil, fmt.Errorf("twitter API error")
	}

	var j getOtherMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	var messages []sns_model.GetOtherMessagesMessage
	for _, m := range j.Data {
		messages = append(messages, sns_model.NewGetOtherMessagesMessage(m.ID, m.Text))
	}

	return sns_model.NewGetOtherMessagesResponse(messages, ""), nil
}

type searchMessagesResponse struct {
	Data []struct {
		AuthorID string `json:"author_id"`
		Text     string `json:"text"`
	} `json:"data"`
}

func (sns *snsServiceTwitterImpl) ExecuteSearchMessageCmd(ctx context.Context, accountID *sns_model.AccountID, cmd *cmd.SearchMessageCommand) (*sns_model.SearchMessageResponse, error) {
	maxResults := cmd.MaxResults()
	if maxResults < 10 || 20 < maxResults {
		maxResults = 10
	}
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("https://api.twitter.com/2/tweets/search/recent")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("query", fmt.Sprintf("%s -is:retweet", cmd.Query()))
	q.Add("max_results", fmt.Sprintf("%d", maxResults))
	q.Add("tweet.fields", "author_id")
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

	// TODO: errorハンドリングもっと丁寧にする
	if resp.StatusCode != http.StatusOK {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(b))
		}
		if resp.StatusCode == http.StatusBadRequest {
			return sns_model.NewSearchMessageResponse(nil, "invalid request."), nil
		}
		return nil, fmt.Errorf("twitter API error")
	}

	var j searchMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	var messages []sns_model.SearchMessageMessage
	for _, m := range j.Data {
		messages = append(messages, sns_model.NewSearchMessageMessage(m.AuthorID, m.Text))
	}

	return sns_model.NewSearchMessageResponse(messages, ""), nil
}

type getMyProfileResponse struct {
	Data struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"data"`
}

func (sns *snsServiceTwitterImpl) ExecuteGetMyProfileCmd(ctx context.Context, accountID *sns_model.AccountID, cmd *cmd.GetMyProfileCommand) (*sns_model.GetMyProfileResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("https://api.twitter.com/2/users/me")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("user.fields", "description")
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

	// TODO: errorハンドリングもっと丁寧にする
	if resp.StatusCode != http.StatusOK {
		// respをdumpする
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(b))
		}
		return sns_model.NewGetMyProfileResponse("", "", "get my profile failed"), fmt.Errorf("twitter API error")
	}

	var j getMyProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	// return sns_model.NewGetMyProfileResponse(j.Data.Name, j.Data.Description, ""), nil
	return sns_model.NewGetMyProfileResponse(j.Data.Name, "my description", ""), nil // 一時的にdescriptionは封鎖する
}

type getOthersProfileResponse struct {
	Data struct {
		ID          string `json:"id"`
		UserName    string `json:"username"`
		Description string `json:"description"`
	} `json:"data"`
}

func (sns *snsServiceTwitterImpl) ExecuteGetOthersProfileCmd(ctx context.Context, accountID *sns_model.AccountID, cmd *cmd.GetOthersProfileCommand) (*sns_model.GetOthersProfileResponse, error) {
	if !validateUserID(cmd.UserID()) {
		return sns_model.NewGetOthersProfileResponse("", "", "", "Request error: user_id must consists of digits only"), nil
	}
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID.ToString())).First(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf("https://api.twitter.com/2/users/%s", cmd.UserID()))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("user.fields", "description")
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

	// TODO: errorハンドリングもっと丁寧にする
	if resp.StatusCode != http.StatusOK {
		// respをdumpする
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(b))
		}
		if resp.StatusCode == http.StatusBadRequest {
			return sns_model.NewGetOthersProfileResponse("", "", "", "One or more parameters to your request was invalid."), nil
		}
		return sns_model.NewGetOthersProfileResponse("", "", "", fmt.Sprintf("userID (%s) not found", cmd.UserID())), fmt.Errorf("get other twitter account failed")
	}

	var j getOthersProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	//MEMO: userIDは1123444555みたいなやつですか？
	return sns_model.NewGetOthersProfileResponse(j.Data.ID, j.Data.UserName, j.Data.Description, ""), nil
}
func (sns *snsServiceTwitterImpl) ExecuteUpdateMyProfileCmd(ctx context.Context, accountID *sns_model.AccountID, cmd *cmd.UpdateMyProfileCommand) (*sns_model.UpdateMyProfileResponse, error) {
	return sns_model.NewUpdateMyProfileResponse("this command is not implemented."), nil
}

type getUserIDByAccessTokenResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (sns *snsServiceTwitterImpl) getUserIDByAccessToken(ctx context.Context, accessToken string) (string, error) {
	u, err := url.Parse("https://api.twitter.com/2/users/me")
	if err != nil {
		return "", err
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

	// TODO: errorハンドリングもっと丁寧にする
	if resp.StatusCode != http.StatusOK {
		// respをdumpする
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(b))
		}
		return "", fmt.Errorf("twitter API error")
	}

	var j getUserIDByAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return "", err
	}

	return j.Data.ID, nil
}

func validateUserID(userID string) bool {
	// userIDが数字だけで構成されていることを確認する
	r := regexp.MustCompile(`^[0-9]+$`)
	return r.MatchString(userID)
}

func NewSNSServiceTwitterImpl(db *ent.Client) sns.Service {
	return &snsServiceTwitterImpl{db}
}
