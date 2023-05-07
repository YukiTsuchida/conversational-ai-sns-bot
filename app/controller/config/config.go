package config

import (
	"os"
	"strconv"
)

func POSTGRES_HOST() string {
	return os.Getenv("POSTGRES_HOST")
}

func POSTGRES_PORT() string {
	return os.Getenv("POSTGRES_PORT")
}

func POSTGRES_USER() string {
	return os.Getenv("POSTGRES_USER")
}

func POSTGRES_PASSWORD() string {
	return os.Getenv("POSTGRES_PASSWORD")
}

func POSTGRES_DB() string {
	return os.Getenv("POSTGRES_DB")
}

func CLOUDTASKS_HOST() string {
	return os.Getenv("CLOUDTASKS_HOST")
}

func CLOUDTASKS_PARENT() string {
	return os.Getenv("CLOUDTASKS_PARENT")
}

func SELF_HOST() string {
	return os.Getenv("SELF_HOST")
}

func REQUESTOR_HOST() string {
	return os.Getenv("REQUESTOR_HOST")
}

// 1秒あたりの会話数、0.05に設定するとAIと1分間に3回の会話が可能
func CONVERSATION_RATE_PER_SECOND() float64 {
	tmp := os.Getenv("CONVERSATION_RATE_PER_SECOND")
	v, err := strconv.ParseFloat(tmp, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func SLEEP_TIME_FOR_REPLY_SECONDS() int {
	tmp := os.Getenv("SLEEP_TIME_FOR_REPLY_SECONDS")
	v, err := strconv.Atoi(tmp)
	if err != nil {
		panic(err)
	}
	return v
}

func CHATGPT_API_KEY() string {
	return os.Getenv("CHATGPT_API_KEY")
}

func TWITTER_CLIENT_ID() string {
	return os.Getenv("TWITTER_CLIENT_ID")
}

func TWITTER_CALLBACK_URL() string {
	return os.Getenv("TWITTER_CALLBACK_URL")
}

func TWITTER_CLIENT_SECRET() string {
	return os.Getenv("TWITTER_CLIENT_SECRET")
}

func ENVIRONMENT() string {
	return os.Getenv("ENVIRONMENT")
}

func IsDevelop() bool {
	return ENVIRONMENT() == "develop"
}

func IsProd() bool {
	return ENVIRONMENT() == "prod"
}
