package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/model/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/sns/twitter"
)

func CallbackTwitterAccountHandler(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {
	// DI
	sns := twitter.NewSNSTwitterImpl(db)

	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		stateCookie, err := r.Cookie("state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if s := r.URL.Query().Get("state"); s != stateCookie.Value {
			http.Error(w, "state is not match", http.StatusBadRequest)
			return
		}
		// cookieの削除
		stateCookie.MaxAge = -1
		http.SetCookie(w, stateCookie)

		cvCookie, err := r.Cookie("code_verifier")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// cookieの削除
		cvCookie.MaxAge = -1
		http.SetCookie(w, cvCookie)

		// codeとtokenの交換
		tokenResp, err := getTwitterOAuth2Token(r.Context(), code, cvCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 自分のaccount情報の取得
		accountResp, err := getMyTwitterAccount(r.Context(), tokenResp.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		credential := sns_model.NewOAuth2Credential(tokenResp.AccessToken, tokenResp.RefreshToken)

		// DBに保存
		if err := sns.CreateAccount(r.Context(), accountResp.Data.UserName, credential); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	}
}

type twitterOAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func getTwitterOAuth2Token(ctx context.Context, code, codeVerifier string) (*twitterOAuth2TokenResponse, error) {
	u, err := url.Parse("https://api.twitter.com/2/oauth2/token")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("client_id", config.TWITTER_CLIENT_ID())
	q.Add("redirect_uri", config.TWITTER_CALLBACK_URL())
	q.Add("code_verifier", codeVerifier)
	u.RawQuery = q.Encode()

	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.SetBasicAuth(config.TWITTER_CLIENT_ID(), config.TWITTER_CLIENT_SECRET())

	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get token failed")
	}

	var j twitterOAuth2TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	return &j, nil
}

type myTwitterAccountResponse struct {
	Data struct {
		UserName string `json:"username"`
	} `json:"data"`
}

func getMyTwitterAccount(ctx context.Context, token string) (*myTwitterAccountResponse, error) {
	apiURL := "https://api.twitter.com/2/users/me"

	newReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	c := http.Client{}
	resp, err := c.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get my twitter account failed")
	}

	var j myTwitterAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	return &j, nil
}
