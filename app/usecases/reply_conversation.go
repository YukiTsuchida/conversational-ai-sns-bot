package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/prompt"

	cmd_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
)

type ReplyConversation struct {
	snsSvc           sns.Service
	promptSvc        prompt.Service
	aiSvc            ai.Service
	conversationRepo repositories.Conversation
}

func (uc *ReplyConversation) Execute(ctx context.Context, conversationID string, message string) error {
	// conversationが存在するか確認
	conversation, err := uc.conversationRepo.FetchByID(ctx, conversationID)
	if err != nil {
		return err
	}

	// conversationが中断されてたら終了
	if conversation.IsAborted() {
		return fmt.Errorf("this conversation is aborted")
	}

	// accountを取得
	account, err := uc.snsSvc.FetchAccountByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}

	// messageからcmdを抽出する
	cmds := uc.promptSvc.ParseCmdsByMessage(message)

	// purposeコマンドを探す
	var purpose string
	for _, cmd := range cmds {
		if cmd.IsPostActionPurposeLog() {
			purpose, err = cmd.Option("log")
			if err != nil {
				return err
			}
			break
		}
	}

	// AIからのメッセージをログに追加
	err = uc.aiSvc.AppendAIMessage(ctx, conversationID, message, purpose)
	if err != nil {
		return err
	}

	// cmdをSNSに送信して、どんどんログに追加していく
	// 複数のコマンドを同時に実行できるようにしているが、基本的には1つのコマンドしか実行されない想定(プロンプトで縛ってるだけなのでAIがプロンプトを無視したらその限りではない)
	for _, cmd := range cmds {
		if cmd.IsPostActionPurposeLog() {
			// purposeログは無視
			continue
		}
		var nextMessage string
		if cmd.IsPostMessage() {
			message, err := cmd.Option("message")
			if err != nil {
				// 基本的にパース時にエラーが出るのでオプションは必ず存在する
				return err
			}
			cmd := cmd_model.NewPostMessageCommand(message)
			snsRes, err := uc.snsSvc.ExecutePostMessageCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = uc.promptSvc.BuildNextMessagePostMessage(snsRes)
		} else if cmd.IsGetMyMessages() {
			maxResults, err := cmd.OptionInInt("max_results")
			if err != nil {
				// ToDo: max_resultsにint以外が入っていた場合もエラーになるのでここはハンドリングしたほうがいい
				return err
			}
			cmd := cmd_model.NewGetMyMessagesCommand(maxResults)
			snsRes, err := uc.snsSvc.ExecuteGetMyMessagesCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = uc.promptSvc.BuildNextMessageGetMyMessages(snsRes)
		} else if cmd.IsGetOtherMessages() {
			userID, err := cmd.Option("user_id")
			if err != nil {
				return err
			}
			maxResults, err := cmd.OptionInInt("max_results")
			if err != nil {
				// ToDo: max_resultsにint以外が入っていた場合もエラーになるのでここはハンドリングしたほうがいい
				return err
			}
			cmd := cmd_model.NewGetOtherMessagesCommand(userID, maxResults)
			snsRes, err := uc.snsSvc.ExecuteGetOtherMessagesCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = uc.promptSvc.BuildNextMessageGetOtherMessages(snsRes)
		} else if cmd.IsSearchMessage() {
			query, err := cmd.Option("query")
			if err != nil {
				return err
			}
			maxResults, err := cmd.OptionInInt("max_results")
			if err != nil {
				// ToDo: max_resultsにint以外が入っていた場合もエラーになるのでここはハンドリングしたほうがいい
				return err
			}
			cmd := cmd_model.NewSearchMessageCommand(query, maxResults)
			snsRes, err := uc.snsSvc.ExecuteSearchMessageCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = uc.promptSvc.BuildNextMessageSearchMessage(snsRes)
		} else if cmd.IsGetMyProfile() {
			cmd := cmd_model.NewGetMyProfileCommand()
			snsRes, err := uc.snsSvc.ExecuteGetMyProfileCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = uc.promptSvc.BuildNextMessageGetMyProfile(snsRes)
		} else if cmd.IsGetOthersProfile() {
			userID, err := cmd.Option("user_id")
			if err != nil {
				return err
			}
			cmd := cmd_model.NewGetOthersProfileCommand(userID)
			snsRes, err := uc.snsSvc.ExecuteGetOthersProfileCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = uc.promptSvc.BuildNextMessageGetOthersProfile(snsRes)
		} else if cmd.IsUpdateMyProfile() {
			name, err := cmd.Option("name")
			if err != nil {
				return err
			}
			description, err := cmd.Option("description")
			if err != nil {
				return err
			}
			cmd := cmd_model.NewUpdateMyProfileCommand(name, description)
			snsRes, err := uc.snsSvc.ExecuteUpdateMyProfileCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = uc.promptSvc.BuildNextMessageUpdateMyProfile(snsRes)
		} else {
			// 存在しない
			continue
		}
		err = uc.aiSvc.AppendUserMessage(ctx, conversationID, nextMessage)
		if err != nil {
			return err
		}
	}

	if len(cmds) == 0 {
		// cmdがない場合は、メッセージが間違ってるよって教える
		nextMessage := uc.promptSvc.BuildNextMessageCommandNotFound()
		err = uc.aiSvc.AppendUserMessage(ctx, conversationID, nextMessage)
		if err != nil {
			return err
		}
	}

	time.Sleep(time.Duration(config.SLEEP_TIME_FOR_REPLY_SECONDS()) * time.Second)

	// 会話履歴を結合して対話型AI用のリクエストを生成して送信する
	err = uc.aiSvc.SendRequest(ctx, conversationID)
	if err != nil {
		return err
	}

	return nil
}

func NewReplyConversation(snsSvc sns.Service, promptSvc prompt.Service, aiSvc ai.Service, conversationRepo repositories.Conversation) *ReplyConversation {
	return &ReplyConversation{snsSvc, promptSvc, aiSvc, conversationRepo}
}
