package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/service"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns/twitter"
)

type RegisterTwitterAccountRequest struct {
	TwitterID   string `json:"twitter_id"`
	BearerToken string `json:"bearer_token"`
}

func RegisterTwitterAccountHandler(db *ent.Client) func(w http.ResponseWriter, r *http.Request) {
	// DI
	sns := twitter.NewSNSTwitterImpl(db)
	registerAccountService := service.NewRegisterAccountService(sns)

	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterTwitterAccountRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internalRegisterTwitterAccountError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.TwitterID == "" {
			http.Error(w, "invalid request body param: twitter_id", http.StatusBadRequest)
			return
		}
		if req.BearerToken == "" {
			http.Error(w, "invalid request body param: bearer_token", http.StatusBadRequest)
			return
		}

		err = registerAccountService.RegisterAccount(r.Context(), req.TwitterID, req.BearerToken)
		if err != nil {
			// ToDo: エラーの内容に応じてresponseを変える
			internalRegisterTwitterAccountError(err)
			http.Error(w, "failed to register_account: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func internalRegisterTwitterAccountError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] RegisterTwitterAccountHandler() error: %s\n", err.Error())
}
