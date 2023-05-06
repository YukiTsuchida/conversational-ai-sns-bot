package handler

import (
	"net/http"
	"net/url"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/config"
)

func LoginTwitterAccountHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := url.Values{}

		params.Add("response_type", "code")
		params.Add("client_id", config.TWITTER_CLIENT_ID())
		params.Add("redirect_uri", config.TWITTER_CALLBACK_URL())
		params.Add("scope", "tweet.read users.read offline.access")
		params.Add("state", "abc")
		params.Add("code_challenge", "aaa")
		params.Add("code_challenge_method", "plain")

		apiURL := "http://twitter.com/i/oauth2/authorize?" + params.Encode()

		http.Redirect(w, r, apiURL, http.StatusFound)
	}
}
