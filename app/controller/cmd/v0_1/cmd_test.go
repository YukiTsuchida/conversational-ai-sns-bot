package v0_1

import "testing"

func TestV0_1_parseOption(t *testing.T) {
	type TestCase struct {
		line       string
		optionName string
		want       string
	}
	tests := []TestCase{
		{
			`PostMessage:message={"Message to be posted"}?max_results={10}`,
			"message",
			`Message to be posted`,
		},
		{
			`PostMessage:message={"Message to be posted"}?max_results={10}`,
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
