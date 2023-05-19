package cloud_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/queue"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/grpc"
)

var _ queue.Service = (*queueCloudTasksImpl)(nil)

type queueCloudTasksImpl struct {
}

// ToDo: http/handlers側でまったく同じ型を定義してしまっているので共通化する
type SendConversationRequest struct {
	CallBackUrl string `json:"call_back_url"`
}

func (q *queueCloudTasksImpl) Push(ctx context.Context, conversationID *conversation.ID) error {
	cloudTasksHost := config.CLOUDTASKS_HOST()     // local環境ではエミュレータを使うので環境変数からhostを指定する
	cloudtasksParent := config.CLOUDTASKS_PARENT() // 「projects/%s/locations/%s」の部分
	cloudtasksChild := "/queues/" + conversationID.ToString()

	var err error
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

	request := SendConversationRequest{
		CallBackUrl: config.SELF_HOST() + "/conversations/" + conversationID.ToString() + "/reply",
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
					Url:        config.REQUESTOR_HOST() + "/conversations/" + conversationID.ToString() + "/send",
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

func NewQueueServiceCloudTasksImpl() queue.Service {
	return &queueCloudTasksImpl{}
}
