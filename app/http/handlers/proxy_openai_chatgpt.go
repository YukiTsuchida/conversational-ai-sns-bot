package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai/chatgpt_3_5_turbo"
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
	newReq, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return nil, err
	}
	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Authorization", "Bearer "+config.CHATGPT_API_KEY())

	client := &http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBuf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("ChatGPT API request error: status code %d, body %s", resp.StatusCode, string(bodyBuf))
	}

	// ChatGPT APIのレスポンスをパースする
	var chatGPTAPIResp ChatGPTAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&chatGPTAPIResp)
	if err != nil {
		return nil, err
	}
	return &chatGPTAPIResp, nil
}

func internalOpenAIChatGPTRequestError(err error) {
	fmt.Fprintf(os.Stderr, "[ERROR] OpenAIChatGPTRequestHandler() error: %s\n", err.Error())
}
