package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/config"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns/twitter"
)

func CallbackTwitterAccountHandler(db *ent.Client) func(w http.ResponseWriter,r *http.Request) {
	// DI
	sns := twitter.NewSNSTwitterImpl(db)
	
	return func(w http.ResponseWriter, r *http.Request) {
		// stateの検証
		if s:=r.URL.Query().Get("state");s != "abc" {
			http.Error(w,"state is not match",http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")

		// codeとtokenの交換
		tokenResp,err := getTwitterOAuth2Token(r.Context(),code)
		if err != nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}

		// 自分のaccount情報の取得
		accountResp,err := getMyTwitterAccount(r.Context(),tokenResp.AccessToken)
		if err != nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}

		credential := sns_model.NewOAuth2Credential(tokenResp.AccessToken,tokenResp.RefreshToken)
		if err := sns.CreateAccount(r.Context(),accountResp.Data.UserName,credential);err!=nil{
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w,"ok")
	}
}

type twitterOAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func getTwitterOAuth2Token(ctx context.Context, code string) (*twitterOAuth2TokenResponse,error) {
	newParam := url.Values{}

	newParam.Add("code",code)
	newParam.Add("grant_type","authorization_code")
	newParam.Add("client_id",config.TWITTER_CLIENT_ID())
	newParam.Add("redirect_uri",config.TWITTER_CALLBACK_URL())
	newParam.Add("code_verifier","aaa")

	apiURL := "https://api.twitter.com/2/oauth2/token?" + newParam.Encode()

	newReq,err := http.NewRequestWithContext(ctx, http.MethodPost,apiURL,nil)
	if err != nil {
		return nil,err
	}

	newReq.Header.Set("Content-Type","application/x-www-form-urlencoded")
	newReq.SetBasicAuth(config.TWITTER_CLIENT_ID(),config.TWITTER_CLIENT_SECRET())

	c := http.Client{}
	resp,err := c.Do(newReq)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,fmt.Errorf("get token failed")
	}

	var j twitterOAuth2TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&j);err!=nil{
		return nil,err
	}

	return &j,nil
}

type myTwitterAccountResponse struct {
	Data struct{
		UserName string `json:"username"`
	} `json:"data"`
}

func getMyTwitterAccount(ctx context.Context,token string) (*myTwitterAccountResponse, error){
	apiURL := "https://api.twitter.com/2/users/me"

	newReq,err := http.NewRequestWithContext(ctx,http.MethodGet,apiURL,nil)
	if err != nil {
		return nil,err
	}

	newReq.Header.Set("Content-Type","application/x-www-form-urlencoded")
	newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s",token))

	c := http.Client{}
	resp,err := c.Do(newReq)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,fmt.Errorf("get my twitter account failed")
	}


	var j myTwitterAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&j);err!=nil{
		return nil,err
	}

	return &j,nil
}
