package chatgpt_3_5_turbo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/conversations"
	ai_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"github.com/go-resty/resty/v2"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/chatgpt35turboconversationlog"

	"github.com/pkoukk/tiktoken-go"
)

const modelName = "gpt-3.5-turbo"
const tokenLimit = 3600 // 実際は4096だが、回答も含めて4096なので少なめにしておく
const temperature = 1.0

var _ ai.Service = (*aiServiceChatGPT3_5TurboImpl)(nil)

type aiServiceChatGPT3_5TurboImpl struct {
	db *ent.Client
}

// CloudTasksに積むhttpリクエストのbody
// type Request struct {
// 	CallBackUrl string    `json:"call_back_url"`
// 	AIModel     string    `json:"ai_model"`
// 	Temperature string    `json:"temperature"`
// 	Messages    []Message `json:"messages"`
// }

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

func (ai *aiServiceChatGPT3_5TurboImpl) SendRequest(ctx context.Context, conversationID *conversation.ID) (*ai_model.Response, error) {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return nil, err
	}

	logs, err := ai.db.Chatgpt35TurboConversationLog.Query().Where(chatgpt35turboconversationlog.HasConversationWith(conversations.IDEQ(conversationIDInt))).Order(ent.Asc(chatgpt35turboconversationlog.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	// 送信するメッセージを作成
	numTokens := 0
	messages := []ChatGPTAPIMessage{}

	// systemメッセージを最初に載せる
	messages = append(messages, ChatGPTAPIMessage{
		Role:    "system",
		Content: logs[0].Message,
	})
	if t, err := calcToken(logs[0].Message); err != nil {
		return nil, err
	} else {
		numTokens += t
	}

	insert := func(a []ChatGPTAPIMessage, index int, value ChatGPTAPIMessage) []ChatGPTAPIMessage {
		if len(a) == index { // nil or empty slice or after last element
			return append(a, value)
		}
		a = append(a[:index+1], a[index:]...) // index < len(a)
		a[index] = value
		return a
	}

	// logsを逆順にloop回す
	for i := len(logs) - 1; i != 0; i-- { // systemメッセージは飛ばす
		t, err := calcToken(logs[i].Message)
		if err != nil {
			return nil, err
		}
		if numTokens+t > tokenLimit {
			// トークン数超過したらbreak
			break
		}
		numTokens += t
		messages = insert(messages, 1, ChatGPTAPIMessage{
			Role:    logs[i].Role.String(),
			Content: logs[i].Message,
		})
	}

	fmt.Println(numTokens)

	for _, message := range messages {
		fmt.Println(message.Role + " " + message.Content)
	}

	// ChatGPT APIに投げるリクエストを組み立てる
	chatGPTAPIReq := ChatGPTAPIRequest{
		Model:       modelName,
		Temperature: temperature,
		Messages:    messages,
	}

	jsonStr, err := json.Marshal(chatGPTAPIReq)
	if err != nil {
		return nil, err
	}
	chatGPTAPIResp, err := reqChatGPTAPI(jsonStr)
	if err == nil {
		return ai_model.NewResponse(chatGPTAPIResp.Choices[0].Message.Content, ""), nil
	} else {
		fmt.Fprintf(os.Stderr, "ChatGPT API request error: %v\n", err)
		return ai_model.NewResponse("", "ChatGPT API request error: "+err.Error()), nil
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

func calcToken(text string) (int, error) {
	tkm, err := tiktoken.EncodingForModel(modelName)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return 0, err
	}

	// encode
	token := tkm.Encode(text, nil, nil)

	// num_tokens
	return len(token), nil
}

func (ai *aiServiceChatGPT3_5TurboImpl) AppendSystemMessage(ctx context.Context, conversationID *conversation.ID, message *ai_model.SystemMessage) error {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return err
	}
	_, err = ai.db.Chatgpt35TurboConversationLog.Create().SetConversationID(conversationIDInt).SetMessage(message.ToString()).SetRole(chatgpt35turboconversationlog.RoleSystem).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ai *aiServiceChatGPT3_5TurboImpl) AppendUserMessage(ctx context.Context, conversationID *conversation.ID, message *ai_model.UserMessage) error {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return err
	}
	_, err = ai.db.Chatgpt35TurboConversationLog.Create().SetConversationID(conversationIDInt).SetMessage(message.ToString()).SetRole(chatgpt35turboconversationlog.RoleUser).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ai *aiServiceChatGPT3_5TurboImpl) AppendAIMessage(ctx context.Context, conversationID *conversation.ID, message *ai_model.AIMessage, purpose string) error {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return err
	}
	_, err = ai.db.Chatgpt35TurboConversationLog.Create().SetConversationID(conversationIDInt).SetMessage(message.ToString()).SetRole(chatgpt35turboconversationlog.RoleAssistant).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewAIServiceChatGPT3_5TurboImpl(db *ent.Client) ai.Service {
	return &aiServiceChatGPT3_5TurboImpl{db}
}
