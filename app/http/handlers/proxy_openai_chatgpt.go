package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai/chatgpt_3_5_turbo"

	"github.com/go-resty/resty/v2"
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

type ControllerReplyRequest struct {
	Message    string `json:"message"`
	ErrMessage string `json:"err_message"`
}

type ChatGPTAPIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func ProxyOpenAIChatGPTHandler() func(w http.ResponseWriter, r *http.Request) {

	// ToDo: ここの処理をservices/aiに移す #26

	return func(w http.ResponseWriter, r *http.Request) {
		var req chatgpt_3_5_turbo.Request // ToDo: 実装に依存してしまってるので治す
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

		var controllerReplyReq ControllerReplyRequest
		chatGPTAPIResp, err := reqChatGPTAPI(jsonStr)
		if err == nil {
			controllerReplyReq.Message = chatGPTAPIResp.Choices[0].Message.Content
			controllerReplyReq.ErrMessage = ""
		} else {
			controllerReplyReq.Message = ""
			controllerReplyReq.ErrMessage = "ChatGPT API request error: " + err.Error()
			internalOpenAIChatGPTRequestError(err)
		}

		jsonStr, err = json.Marshal(controllerReplyReq)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newReq, err := http.NewRequest(
			"POST",
			req.CallBackUrl,
			bytes.NewBuffer([]byte(jsonStr)),
		)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		_, err = client.Do(newReq)
		if err != nil {
			internalOpenAIChatGPTRequestError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func reqChatGPTAPI(jsonStr []byte) (*ChatGPTAPIResponse, error) {

	var chatGPTAPIResp ChatGPTAPIResponse

	client := resty.New()

	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	client.
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+config.CHATGPT_API_KEY()).
		SetBody(jsonStr).
		SetResult(&chatGPTAPIResp). // or SetResult(AuthSuccess{}).
		Post("https://api.openai.com/v1/chat/completions")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("ChatGPT API request error: status code %d, body %s", resp.StatusCode(), string(resp.Body()))
	}
	return &chatGPTAPIResp, nil
}

func internalOpenAIChatGPTRequestError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] OpenAIChatGPTRequestHandler() error: %s\n", err.Error())
}
