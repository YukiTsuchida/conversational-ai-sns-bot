package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/config"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai/chatgpt_3_5_turbo"
)

// ChatGPT APIに投げるリクエストのbody型
// Ref: https://platform.openai.com/docs/api-reference/chat/create
type ChatGPTAPIRequest struct {
	Model       string              `json:"model"`
	Temperature float64             `json:"temperature"`
	Messages    []ChatGPTAPIMessage `json:"messages"`
}

type ChatGPTAPIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func OpenAIChatGPTRequestHandler() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var req chatgpt_3_5_turbo.Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.AIModel == "" {
			http.Error(w, "invalid request body param: ai_model", http.StatusBadRequest)
			return
		}
		if req.CallBackUrl == "" {
			http.Error(w, "invalid request body param: call_back_url", http.StatusBadRequest)
			return
		}
		if len(req.Messages) == 0 {
			http.Error(w, "invalid request body param: messages", http.StatusBadRequest)
			return
		}
		if req.Temperature == "" {
			http.Error(w, "invalid request body param: temperature", http.StatusBadRequest)
			return
		}

		// ChatGPT APIに投げるリクエストを組み立てる
		temperatureFloat64, err := strconv.ParseFloat(req.Temperature, 64)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		chatGPTAPIReq := ChatGPTAPIRequest{
			Model:       req.AIModel,
			Temperature: temperatureFloat64,
			Messages:    []ChatGPTAPIMessage{},
		}
		for _, message := range req.Messages {
			chatGPTAPIReq.Messages = append(chatGPTAPIReq.Messages, ChatGPTAPIMessage{
				Role:    message.Role,
				Content: message.Message,
			})
		}

		jsonStr, err := json.Marshal(chatGPTAPIReq)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newReq, err := http.NewRequest(
			"POST",
			"https://api.openai.com/v1/chat/completions",
			bytes.NewBuffer([]byte(jsonStr)),
		)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newReq.Header.Set("Content-Type", "application/json")
		newReq.Header.Set("Authorization", "Bearer "+config.CHATGPT_API_KEY())

		client := &http.Client{}
		resp, err := client.Do(newReq)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// 本来はここでChatGPTAPIのレスポンスをパースして、call_back_urlに返すが、一旦繋ぎ込みはせずに結果だけログに出力する
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("%q", dump)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func internalOpenAIChatGPTRequestError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] OpenAIChatGPTRequestHandler() error: %s\n", err.Error())
}
