package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"
)

func LoginTwitterAccountHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse("http://twitter.com/i/oauth2/authorize")
		if err != nil {
			http.Error(w, "url.Parse failed", http.StatusBadRequest)
			return
		}

		//csrf対策, client側で検証する
		state := randomString(20)

		sCookie := &http.Cookie{
			Name:     "state",
			Value:    state,
			Path:     "/", //cookieを送信するpathを自動で絞り込む
			Expires:  time.Now().Add(60 * time.Minute),
			Secure:   config.IsProd(), //trueの場合httpsのみ送信
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, sCookie)

		// csrf,認可コード横取り攻撃対策, 認可サーバで検証
		codeVerifier := randomString(20)
		b := sha256.Sum256([]byte(codeVerifier))
		codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b[:])

		cCookie := &http.Cookie{
			Name:     "code_verifier",
			Value:    codeVerifier,
			Path:     "/", //cookieを送信するpathを自動で絞り込む
			Expires:  time.Now().Add(60 * time.Minute),
			Secure:   config.IsProd(), //trueの場合httpsのみ送信
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cCookie)

		q := u.Query()
		q.Add("response_type", "code")
		q.Add("client_id", config.TWITTER_CLIENT_ID())
		q.Add("redirect_uri", config.TWITTER_CALLBACK_URL())
		q.Add("scope", "tweet.read tweet.write users.read offline.access")
		q.Add("state", state)
		q.Add("code_challenge", codeChallenge) //次のリクエスト時にcode_verifierを認可サーバ側でhash化し比較, 同一のユーザーかの確認
		q.Add("code_challenge_method", "S256")
		u.RawQuery = q.Encode()

		http.Redirect(w, r, u.String(), http.StatusFound)
	}
}

func randomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
