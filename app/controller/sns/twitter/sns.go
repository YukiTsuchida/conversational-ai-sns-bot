package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

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
	Text string `json:"text"`
}

func (sns *snsTwitterImpl) ExecutePostMessageCmd(ctx context.Context, accountID string, cmd *cmd.PostMessageCommand) (*sns_model.PostMessageResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
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
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", account.AccessToken))

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

func (sns *snsTwitterImpl) ExecuteGetMyMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetMyMessagesCommand) (*sns_model.GetMyMessagesResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
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
	q.Add("max_results", fmt.Sprintf("%d", cmd.MaxResults())) //min:5, max: 100
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

func (sns *snsTwitterImpl) ExecuteGetOtherMessagesCmd(ctx context.Context, accountID string, cmd *cmd.GetOtherMessagesCommand) (*sns_model.GetOtherMessagesResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
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
	q.Add("max_results", fmt.Sprintf("%d", cmd.MaxResults())) //min:5, max: 100
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
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) ExecuteSearchMessageCmd(ctx context.Context, accountID string, cmd *cmd.SearchMessageCommand) (*sns_model.SearchMessageResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("https://api.twitter.com/2/tweets/search/recent")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("query", fmt.Sprintf("%s -is:retweet", cmd.Query()))
	q.Add("max_results", fmt.Sprintf("%d", cmd.MaxResults()))
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
		messages = append(messages, sns_model.NewSearchMessageMessage(m.ID, m.Text))
	}

	return sns_model.NewSearchMessageResponse(messages, ""), nil
}

type getMyProfileResponse struct {
	Data struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) ExecuteGetMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetMyProfileCommand) (*sns_model.GetMyProfileResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
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

	return sns_model.NewGetMyProfileResponse(j.Data.Name, j.Data.Description, ""), nil
}

type getOthersProfileResponse struct {
	Data struct {
		ID          string `json:"id"`
		UserName    string `json:"username"`
		Description string `json:"description"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) ExecuteGetOthersProfileCmd(ctx context.Context, accountID string, cmd *cmd.GetOthersProfileCommand) (*sns_model.GetOthersProfileResponse, error) {
	account, err := sns.db.TwitterAccounts.Query().Where(twitteraccounts.TwitterAccountIDEQ(accountID)).First(ctx)
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
		return sns_model.NewGetOthersProfileResponse("", "", "", fmt.Sprintf("userID (%s) not found", cmd.UserID())), fmt.Errorf("get other twitter account failed")
	}

	var j getOthersProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	//MEMO: userIDは1123444555みたいなやつですか？
	return sns_model.NewGetOthersProfileResponse(j.Data.ID, j.Data.UserName, j.Data.Description, ""), nil
}
func (sns *snsTwitterImpl) ExecuteUpdateMyProfileCmd(ctx context.Context, accountID string, cmd *cmd.UpdateMyProfileCommand) (*sns_model.UpdateMyProfileResponse, error) {
	return sns_model.NewUpdateMyProfileResponse("this command is not implemented."), nil
}

type getUserIDByAccessTokenResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (sns *snsTwitterImpl) getUserIDByAccessToken(ctx context.Context, accessToken string) (string, error) {
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

func NewSNSTwitterImpl(db *ent.Client) sns.SNS {
	return &snsTwitterImpl{db}
}
