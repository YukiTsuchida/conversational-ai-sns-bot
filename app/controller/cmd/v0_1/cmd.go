package v0_1

import (
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/cmd"
	cmd_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
)

var _ cmd.Cmd = (*cmdV0_1Impl)(nil)

type cmdV0_1Impl struct {
}

func (cmd *cmdV0_1Impl) BuildFirstMessage() string {
	return `
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
    - max_results is a number from 10 to 10

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
	`
}

func (cmd *cmdV0_1Impl) BuildNextMessage(req *cmd_model.Command, res *sns.Response) string {
	return ""
}

func (cmd *cmdV0_1Impl) ParseCmdsByMessage(message string) []cmd_model.Command {
	return nil
}

func NewCmdV0_1Impl() cmd.Cmd {
	return &cmdV0_1Impl{}
}
