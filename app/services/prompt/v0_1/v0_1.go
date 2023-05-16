package v0_1

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
	cmd_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/prompt"
)

var _ prompt.Service = (*promptServiceV0_1Impl)(nil)

type promptServiceV0_1Impl struct {
}

func (prompt *promptServiceV0_1Impl) BuildSystemMessage() *ai.SystemMessage {
	return ai.NewSystemMessage(`
A certain social networking site allows you to perform the following actions

Post a message
- Command: {"action":"PostMessage","options":{"message":"Message to be posted"}}
- Note: You cannot send a message to someone else

Retrieving a list of one's past messages
- Command: {"action":"GetMyMessages","options":{"max_results":5}}
    - max_results is a number from 5 to 10

Get a list of other people's messages
- Command: {"action":"GetOtherMessages","options":{"user_id":"foo","max_results":5}}
	- max_results is a number from 5 to 10

Search other people's messages
- Command: {"action":"SearchMessage","options":{"query":"keyword to search","max_results":10}}
    - max_results is a number from 10 to 20

Get my user information
- Command: {"action":"GetMyProfile"}

To retrieve user information of others
- Command: {"action":"GetOthersProfile","options":{"user_id":"foo"}}

You are now one of the users of this SNS, "user_name". user_name(you) is Japanese and has already been using this SNS for many years.

Please decide what user_name(you) will do from now on.

For example, if you want to get user information about yourself, please reply "{"action":"GetMyProfile"}".

I will feed back the result of the action each time. For example, if you say "{"action":"GetMyMessages","options":{"max_results":5}}", I will give you a list of your past 5 messages.

Your goal is to get closer to more people through social networking. What would you like to do first? From this point on, you can only talk to me on command.
	`)
}

func (prompt *promptServiceV0_1Impl) BuildUserMessagePostMessageResult(res *sns.PostMessageResponse) (nextMessage *ai.UserMessage) {
	defer func() {
		nextMessage.Append("\nWhat do you want to do next?\nplease talk to me on command.")
	}()
	if res == nil {
		// ないとは思うが、SNSのレスポンスがない場合
		return ai.NewUserMessage(`
		No response found from SNS.
		`)
	}

	cmdSuccess := !res.ErrorOccured()

	if cmdSuccess {
		return ai.NewUserMessage(`
		Your message has been posted.
		`)
	} else {
		return ai.NewUserMessage(`
		Failed to post your message.

		reason: ` + res.ErrReason())
	}
}
func (prompt *promptServiceV0_1Impl) BuildUserMessageGetMyMessagesResult(res *sns.GetMyMessagesResponse) (nextMessage *ai.UserMessage) {
	defer func() {
		nextMessage.Append("\nWhat do you want to do next?\nplease talk to me on command.")
	}()
	if res == nil {
		// ないとは思うが、SNSのレスポンスがない場合
		return ai.NewUserMessage(`
		No response found from SNS.
		`)
	}

	cmdSuccess := !res.ErrorOccured()

	if cmdSuccess {
		nextMessage = ai.NewUserMessage(`
		Below is a list of messages you posted.

		`)
		for _, message := range res.Messages() {
			nextMessage.Append(fmt.Sprintf("- message=\"%s\"\n", message))
		}
		return nextMessage
	} else {
		return ai.NewUserMessage(`
		Failed to get your messages.

		reason: ` + res.ErrReason())
	}
}
func (prompt *promptServiceV0_1Impl) BuildUserMessageGetOtherMessagesResult(res *sns.GetOtherMessagesResponse) (nextMessage *ai.UserMessage) {
	defer func() {
		nextMessage.Append("\nWhat do you want to do next?\nplease talk to me on command.")
	}()
	if res == nil {
		// ないとは思うが、SNSのレスポンスがない場合
		return ai.NewUserMessage(`
		No response found from SNS.
		`)
	}

	cmdSuccess := !res.ErrorOccured()

	if cmdSuccess {
		nextMessage = ai.NewUserMessage(`
		The following messages were found.

		`)
		for _, message := range res.Messages() {
			nextMessage.Append(fmt.Sprintf("- user_id=%s, message=\"%s\"\n", message.UserID(), message.Message()))
		}
		return nextMessage
	} else {
		return ai.NewUserMessage(`
		Failed to get messages.

		reason: ` + res.ErrReason())
	}

}
func (prompt *promptServiceV0_1Impl) BuildUserMessageSearchMessageResult(res *sns.SearchMessageResponse) (nextMessage *ai.UserMessage) {
	defer func() {
		nextMessage.Append("\nWhat do you want to do next?\nplease talk to me on command.")
	}()
	if res == nil {
		// ないとは思うが、SNSのレスポンスがない場合
		return ai.NewUserMessage(`
		No response found from SNS.
		`)
	}

	cmdSuccess := !res.ErrorOccured()

	if cmdSuccess {
		nextMessage = ai.NewUserMessage(`
		The following messages were found.

		`)
		for _, message := range res.Messages() {
			nextMessage.Append(fmt.Sprintf("- user_id=%s, message=\"%s\"\n", message.UserID(), message.Message()))
		}
		return nextMessage
	} else {
		return ai.NewUserMessage(`
		Failed to search messages.

		reason: ` + res.ErrReason())
	}
}
func (prompt *promptServiceV0_1Impl) BuildUserMessageGetMyProfileResult(res *sns.GetMyProfileResponse) (nextMessage *ai.UserMessage) {
	defer func() {
		nextMessage.Append("\nWhat do you want to do next?\nplease talk to me on command.")
	}()
	if res == nil {
		// ないとは思うが、SNSのレスポンスがない場合
		return ai.NewUserMessage(`
		No response found from SNS.
		`)
	}

	cmdSuccess := !res.ErrorOccured()

	if cmdSuccess {
		nextMessage = ai.NewUserMessage(`
		This is your profile.

		`)
		nextMessage.Append(fmt.Sprintf("- name=\"%s\"\n", res.Name()))
		nextMessage.Append(fmt.Sprintf("- description=\"%s\"\n", res.Description()))
		return nextMessage
	} else {
		return ai.NewUserMessage(`
		Failed to get your profile.

		reason: ` + res.ErrReason())
	}
}
func (prompt *promptServiceV0_1Impl) BuildUserMessageGetOthersProfileResult(res *sns.GetOthersProfileResponse) (nextMessage *ai.UserMessage) {
	defer func() {
		nextMessage.Append("\nWhat do you want to do next?\nplease talk to me on command.")
	}()
	if res == nil {
		// ないとは思うが、SNSのレスポンスがない場合
		return ai.NewUserMessage(`
		No response found from SNS.
		`)
	}

	cmdSuccess := !res.ErrorOccured()

	if cmdSuccess {
		nextMessage = ai.NewUserMessage(`
		The following profile was found.

		`)
		nextMessage.Append(fmt.Sprintf("- user_id=\"%s\"\n", res.UserID()))
		nextMessage.Append(fmt.Sprintf("- name=\"%s\"\n", res.Name()))
		nextMessage.Append(fmt.Sprintf("- description=\"%s\"\n", res.Description()))
		return nextMessage
	} else {
		return ai.NewUserMessage(`
		Failed to get profile.

		reason: ` + res.ErrReason())
	}
}
func (prompt *promptServiceV0_1Impl) BuildUserMessageUpdateMyProfileResult(res *sns.UpdateMyProfileResponse) (nextMessage *ai.UserMessage) {
	defer func() {
		nextMessage.Append("\nWhat do you want to do next?\nplease talk to me on command.")
	}()
	if res == nil {
		// ないとは思うが、SNSのレスポンスがない場合
		return ai.NewUserMessage(`
		No response found from SNS.
		`)
	}

	cmdSuccess := !res.ErrorOccured()

	if cmdSuccess {
		return ai.NewUserMessage(`
		Your profile has been updated.
		`)
	} else {
		return ai.NewUserMessage(`
		Failed to update profile.

		reason: ` + res.ErrReason())
	}
}

func (prompt *promptServiceV0_1Impl) BuildUserMessageCommandNotFoundResult() (nextMessage *ai.UserMessage) {
	return ai.NewUserMessage(`
		Command not found.
		What do you want to do next?
		please talk to me on only command.
		`)
}

type aiAction struct {
	Action  string `json:"action"`
	Options struct {
		UserID     string `json:"user_id"`
		Query      string `json:""`
		Message    string `json:"message"`
		MaxResults int    `json:"max_results"`
	} `json:"options"`
}

func extractAIAction(message string) *aiAction {
	r := regexp.MustCompile(`\{"action":"[A-Za-z]*"(,"options":\{.*\})?\}`)
	matched := r.FindAllStringSubmatch(message, -1)
	if len(matched) == 0 || len(matched[0]) == 0 {
		return nil
	}

	var aiAction aiAction
	err := json.Unmarshal([]byte(matched[0][0]), &aiAction)
	if err != nil {
		// 基本的に起きてはならない
		log.Fatalf("json.Unmarshal error: %s\n", err.Error())
		return nil
	}
	return &aiAction
}

func (prompt *promptServiceV0_1Impl) ParseCmdsByAIMessage(aiMsg *ai.AIMessage) []*cmd_model.Command {
	message := aiMsg.ToString()
	var cmds []*cmd_model.Command
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		aiRes := extractAIAction(line)
		if aiRes == nil {
			continue
		}
		switch aiRes.Action {
		case "PostMessage":
			if aiRes.Options.Message == "" {
				fmt.Fprintf(os.Stderr, "PostMessage parse error: message is empty\n")
				break
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.PostMessage,
				map[string]string{
					"message": aiRes.Options.Message,
				},
			))
		case "GetMyMessages":
			if aiRes.Options.MaxResults == 0 {
				aiRes.Options.MaxResults = 5
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetMyMessages,
				map[string]string{
					"max_results": strconv.Itoa(aiRes.Options.MaxResults),
				},
			))
		case "GetOtherMessages":
			if aiRes.Options.UserID == "" {
				fmt.Fprintf(os.Stderr, "GetOtherMessages parse error: user_id is empty\n")
				break
			}
			if aiRes.Options.MaxResults == 0 {
				aiRes.Options.MaxResults = 5
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetOtherMessages,
				map[string]string{
					"user_id":     aiRes.Options.UserID,
					"max_results": strconv.Itoa(aiRes.Options.MaxResults),
				},
			))
		case "SearchMessage":
			if aiRes.Options.Query == "" {
				fmt.Fprintf(os.Stderr, "SearchMessage parse error: query is empty\n")
				break
			}
			if aiRes.Options.MaxResults == 0 {
				aiRes.Options.MaxResults = 5
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.SearchMessage,
				map[string]string{
					"query":       aiRes.Options.Query,
					"max_results": strconv.Itoa(aiRes.Options.MaxResults),
				},
			))
		case "GetMyProfile":
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetMyProfile,
				map[string]string{},
			))
		case "GetOthersProfile":
			if aiRes.Options.UserID == "" {
				fmt.Fprintf(os.Stderr, "GetOthersProfile parse error: user_id is empty\n")
				break
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetOthersProfile,
				map[string]string{
					"user_id": aiRes.Options.UserID,
				},
			))
		default:
			fmt.Fprintf(os.Stderr, "unknown action: %s\n", aiRes.Action)
		}
	}
	return cmds
}

type option struct {
	name         string
	defaultValue string
}

func parseOptions(line string, options []option) (map[string]string, error) {
	retOptions := map[string]string{}
	for _, option := range options {
		optionValue := parseOption(line, option.name, option.defaultValue)
		if optionValue == "" {
			return nil, fmt.Errorf("option %s is not found, source = %s", option.name, line)
		}
		retOptions[option.name] = optionValue
	}
	return retOptions, nil
}

func parseOption(line string, optionName string, or string) string {
	if !strings.Contains(line, optionName) {
		// optionなし
		return or
	}

	// PostMessage:message={"Message to be posted"}&max_results={10} のような文字列から「Message to be posted」を抽出する
	regexp := regexp.MustCompile(optionName + `=(.*)`)
	matches := regexp.FindStringSubmatch(line)
	if len(matches) != 2 {
		// optionなし
		return or
	}
	val := matches[1]

	// &以降を削除する
	val = strings.Split(val, "&")[0]
	// )以降を削除する
	val = strings.Split(val, ")")[0]

	val = strings.TrimSuffix(val, `.`)

	// {}を削除する
	val = strings.TrimPrefix(val, "{")
	val = strings.TrimSuffix(val, "}")

	// "を削除する
	val = strings.TrimPrefix(val, `"`)
	val = strings.TrimSuffix(val, `"`)

	return val
}

func NewPromptServiceV0_1Impl() prompt.Service {
	return &promptServiceV0_1Impl{}
}
