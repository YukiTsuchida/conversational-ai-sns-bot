package chatgpt_3_5_turbo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/conversations"
	ai_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/grpc"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent/chatgpt35turboconversationlog"

	"github.com/pkoukk/tiktoken-go"
)

const modelName = "gpt-3.5-turbo"
const tokenLimit = 4000 // 実際は4096だが、回答も含めて4096なので4000にしておく

var _ ai.Service = (*aiServiceChatGPT3_5TurboImpl)(nil)

type aiServiceChatGPT3_5TurboImpl struct {
	db *ent.Client
}

// CloudTasksに積むhttpリクエストのbody
type Request struct {
	CallBackUrl string    `json:"call_back_url"`
	AIModel     string    `json:"ai_model"`
	Temperature string    `json:"temperature"`
	Messages    []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

func (ai *aiServiceChatGPT3_5TurboImpl) SendRequest(ctx context.Context, conversationID *conversation.ID) error {
	conversationIDInt, err := conversationID.ToInt()
	if err != nil {
		return err
	}

	logs, err := ai.db.Chatgpt35TurboConversationLog.Query().Where(chatgpt35turboconversationlog.HasConversationWith(conversations.IDEQ(conversationIDInt))).Order(ent.Asc(chatgpt35turboconversationlog.FieldCreatedAt)).All(ctx)
	if err != nil {
		return err
	}

	cloudTasksHost := config.CLOUDTASKS_HOST()     // local環境ではエミュレータを使うので環境変数からhostを指定する
	cloudtasksParent := config.CLOUDTASKS_PARENT() // 「projects/%s/locations/%s」の部分
	cloudtasksChild := "/queues/" + conversationID.ToString()

	var client *cloudtasks.Client
	var clientOpt option.ClientOption = nil
	if cloudTasksHost != "" {
		conn, _ := grpc.Dial(cloudTasksHost+":8123", grpc.WithInsecure())
		clientOpt = option.WithGRPCConn(conn)
		client, err = cloudtasks.NewClient(ctx, clientOpt)
		if err != nil {
			return err
		}
	} else {
		client, err = cloudtasks.NewClient(ctx)
		if err != nil {
			return err
		}
	}

	queuePath := cloudtasksParent + cloudtasksChild
	_, err = client.GetQueue(ctx, &taskspb.GetQueueRequest{Name: queuePath})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			// queueが存在しなければ作る
			createQueueRequest := taskspb.CreateQueueRequest{
				Parent: cloudtasksParent,
				Queue:  &taskspb.Queue{Name: queuePath, RateLimits: &taskspb.RateLimits{MaxDispatchesPerSecond: config.CONVERSATION_RATE_PER_SECOND()}},
			}
			_, err = client.CreateQueue(ctx, &createQueueRequest)
			if err != nil {
				return err
			}
			fmt.Println("queue created" + queuePath)
		} else {
			return err
		}
	}

	// 送信するメッセージを作成
	numTokens := 0
	messages := []Message{}

	// systemメッセージを最初に載せる
	messages = append(messages, Message{
		Role:    "system",
		Message: logs[0].Message,
	})
	if t, err := calcToken(logs[0].Message); err != nil {
		return err
	} else {
		numTokens += t
	}

	insert := func(a []Message, index int, value Message) []Message {
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
			return err
		}
		if numTokens+t > tokenLimit {
			// トークン数超過したらbreak
			break
		}
		numTokens += t
		messages = insert(messages, 1, Message{
			Role:    logs[i].Role.String(),
			Message: logs[i].Message,
		})
	}

	fmt.Println(numTokens)

	for _, message := range messages {
		fmt.Println(message.Role + " " + message.Message)
	}

	request := Request{
		CallBackUrl: config.SELF_HOST() + "/conversations/" + conversationID.ToString() + "/reply",
		AIModel:     modelName,
		Temperature: "0.7",
		Messages:    messages,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	createTaskRequest := taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        config.REQUESTOR_HOST() + "/openai_chat_gpt",
					Body:       body,
				},
			},
		}}
	_, err = client.CreateTask(context.Background(), &createTaskRequest)
	if err != nil {
		return err
	}

	return nil
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
