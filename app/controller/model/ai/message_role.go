package ai

type MessageRole string

const (
	System = MessageRole("System") // ChatGPTでいうところのsystem roleに該当する、プロンプトに常に渡したいメッセージ
	User   = MessageRole("User")   // ChatGPTでいうところのuser roleに該当する、質問など
	Ai     = MessageRole("Ai")     // ChatGPTでいうところのassistant roleに該当する、対話型AIからの回答など
)
