package v0_1

import (
	"testing"
)

func TestV0_1_parseOption(t *testing.T) {
	type TestCase struct {
		line       string
		optionName string
		want       string
	}
	tests := []TestCase{
		{
			`PostMessage(message={"Message to be posted"}&max_results={10})`,
			"message",
			`Message to be posted`,
		},
		{
			`PostMessage(message={"Message to be posted"}&max_results={10})`,
			"max_results",
			`10`,
		},
		{
			`GetOtherMessages(user_id=all&max_results=5)`,
			"user_id",
			`all`,
		},
		{
			`"GetOtherMessages(user_id={user_id}&max_results=10)" to retrieve messages from a user with a specific ID?`,
			"max_results",
			`10`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := parseOption(tt.line, tt.optionName); got != tt.want {
				t.Errorf("parseOption(%s,%s) = %v, want %v", tt.line, tt.optionName, got, tt.want)
			}
		})
	}
}
