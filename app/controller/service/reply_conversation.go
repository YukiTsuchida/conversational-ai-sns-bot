package service

import (
	"context"
	"fmt"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/conversation"
	cmd_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

type ReplyConversationService struct {
	sns              sns.SNS
	cmd              cmd.Cmd
	ai               ai.AI
	conversationRepo conversation.ConversationRepository
}

func (svc *ReplyConversationService) ReplyConversation(ctx context.Context, conversationID string, message string) error {
	// conversationが存在するか確認
	conversation, err := svc.conversationRepo.FetchByID(ctx, conversationID)
	if err != nil {
		return err
	}

	// conversationが中断されてたら終了
	if conversation.IsAborted() {
		return fmt.Errorf("this conversation is aborted")
	}

	// accountを取得
	account, err := svc.sns.FetchAccountByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}

	// messageからcmdを抽出する
	cmds := svc.cmd.ParseCmdsByMessage(message)

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
	err = svc.ai.AppendAIMessage(ctx, conversationID, message, purpose)
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
			snsRes, err := svc.sns.ExecutePostMessageCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = svc.cmd.BuildNextMessagePostMessage(snsRes)
		} else if cmd.IsGetMyMessages() {
			maxResults, err := cmd.OptionInInt("max_results")
			if err != nil {
				// ToDo: max_resultsにint以外が入っていた場合もエラーになるのでここはハンドリングしたほうがいい
				return err
			}
			cmd := cmd_model.NewGetMyMessagesCommand(maxResults)
			snsRes, err := svc.sns.ExecuteGetMyMessagesCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = svc.cmd.BuildNextMessageGetMyMessages(snsRes)
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
			snsRes, err := svc.sns.ExecuteGetOtherMessagesCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = svc.cmd.BuildNextMessageGetOtherMessages(snsRes)
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
			snsRes, err := svc.sns.ExecuteSearchMessageCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = svc.cmd.BuildNextMessageSearchMessage(snsRes)
		} else if cmd.IsGetMyProfile() {
			cmd := cmd_model.NewGetMyProfileCommand()
			snsRes, err := svc.sns.ExecuteGetMyProfileCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = svc.cmd.BuildNextMessageGetMyProfile(snsRes)
		} else if cmd.IsGetOthersProfile() {
			userID, err := cmd.Option("user_id")
			if err != nil {
				return err
			}
			cmd := cmd_model.NewGetOthersProfileCommand(userID)
			snsRes, err := svc.sns.ExecuteGetOthersProfileCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = svc.cmd.BuildNextMessageGetOthersProfile(snsRes)
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
			snsRes, err := svc.sns.ExecuteUpdateMyProfileCmd(ctx, account.ID(), cmd)
			if err != nil {
				return err
			}
			nextMessage = svc.cmd.BuildNextMessageUpdateMyProfile(snsRes)
		} else {
			// 存在しない
			continue
		}
		err = svc.ai.AppendUserMessage(ctx, conversationID, nextMessage)
		if err != nil {
			return err
		}
	}

	if len(cmds) == 0 {
		// cmdがない場合は、メッセージが間違ってるよって教える
		nextMessage := svc.cmd.BuildNextMessageCommandNotFound()
		err = svc.ai.AppendUserMessage(ctx, conversationID, nextMessage)
		if err != nil {
			return err
		}
	}

	// 会話履歴を結合して対話型AI用のリクエストを生成して送信する
	err = svc.ai.SendRequest(ctx, conversationID)
	if err != nil {
		return err
	}

	return nil
}

func NewReplyConversationService(sns sns.SNS, cmd cmd.Cmd, ai ai.AI, conversationRepo conversation.ConversationRepository) *ReplyConversationService {
	return &ReplyConversationService{sns, cmd, ai, conversationRepo}
}
