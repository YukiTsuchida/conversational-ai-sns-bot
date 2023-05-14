package v0_1

import (
	"fmt"
	"os"
	"regexp"
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
- Command: PostMessage:message={"Message to be posted"}
- Note: You cannot send a message to someone else

Retrieving a list of one's past messages
- Command: GetMyMessages:max_results={10}
    - max_results is a number from 5 to 10

Get a list of other people's messages
- Command: GetOtherMessages:user_id={user_id}&max_results={10}
	- max_results is a number from 5 to 10

Search other people's messages
- Command: SearchMessage:query={"keyword to search"}&max_results={10}
    - max_results is a number from 10 to 20

Get my user information
- Command: GetMyProfile

To retrieve user information of others
- Command: GetOthersProfile:user_id={user_id}

Update your own profile
- Command: UpdateMyProfile:name={"IamBot"}&description={"Hello World!"}

You are now one of the users of this SNS, "user_name". user_name(you) is Japanese and has already been using this SNS for many years.

Please decide what user_name(you) will do from now on.

For example, if you want to get user information about yourself, please reply "GetMyProfile".

I will feed back the result of the action each time. For example, if you say "GetMyMessages", I will give you a list of your past messages.

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

func (prompt *promptServiceV0_1Impl) ParseCmdsByAIMessage(aiMsg *ai.AIMessage) []*cmd_model.Command {
	message := aiMsg.ToString()
	var cmds []*cmd_model.Command
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if strings.Contains(line, "PostMessage") {
			options, err := parseOptions(line, []string{"message"})
			if err != nil {
				fmt.Fprintf(os.Stderr, "parseOptions error: %s\n", err.Error())
				continue
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.PostMessage,
				options,
			))
		} else if strings.Contains(line, "GetMyMessages") {
			options, err := parseOptions(line, []string{"max_results"})
			if err != nil {
				fmt.Fprintf(os.Stderr, "parseOptions error: %s\n", err.Error())
				continue
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetMyMessages,
				options,
			))
		} else if strings.Contains(line, "GetOtherMessages") {
			options, err := parseOptions(line, []string{"user_id", "max_results"})
			if err != nil {
				fmt.Fprintf(os.Stderr, "parseOptions error: %s\n", err.Error())
				continue
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetOtherMessages,
				options,
			))
		} else if strings.Contains(line, "SearchMessage") {
			options, err := parseOptions(line, []string{"query", "max_results"})
			if err != nil {
				fmt.Fprintf(os.Stderr, "parseOptions error: %s\n", err.Error())
				continue
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.SearchMessage,
				options,
			))
		} else if strings.Contains(line, "GetMyProfile") {
			options, err := parseOptions(line, []string{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "parseOptions error: %s\n", err.Error())
				continue
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetMyProfile,
				options,
			))
		} else if strings.Contains(line, "GetOthersProfile") {
			options, err := parseOptions(line, []string{"user_id"})
			if err != nil {
				fmt.Fprintf(os.Stderr, "parseOptions error: %s\n", err.Error())
				continue
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.GetOthersProfile,
				options,
			))
		} else if strings.Contains(line, "UpdateMyProfile") {
			options, err := parseOptions(line, []string{"name", "description"})
			if err != nil {
				fmt.Fprintf(os.Stderr, "parseOptions error: %s\n", err.Error())
				continue
			}
			cmds = append(cmds, cmd_model.NewCommand(
				cmd_model.UpdateMyProfile,
				options,
			))
		}
	}
	return cmds
}

func parseOptions(line string, optionNames []string) (map[string]string, error) {
	options := map[string]string{}
	for _, optionName := range optionNames {
		optionValue := parseOption(line, optionName)
		if optionValue == "" {
			return nil, fmt.Errorf("option %s is not found, source = %s", optionName, line)
		}
		options[optionName] = optionValue
	}
	return options, nil
}

func parseOption(line string, optionName string) string {
	if !strings.Contains(line, optionName) {
		// optionなし
		return ""
	}

	// PostMessage:message={"Message to be posted"}&max_results={10} のような文字列から「Message to be posted」を抽出する
	regexp := regexp.MustCompile(optionName + `=(.*)`)
	matches := regexp.FindStringSubmatch(line)
	if len(matches) != 2 {
		// optionなし
		return ""
	}
	val := matches[1]

	// &以降を削除する
	val = strings.Split(val, "&")[0]

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
