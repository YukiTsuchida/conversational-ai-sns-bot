package v0_1

import (
	"testing"
)

func TestV0_1_extractAIAction(t *testing.T) {
	type TestCase struct {
		line       string
		wantAction string
	}
	tests := []TestCase{
		{
			`{"action":"GetMyProfile"}`,
			`GetMyProfile`,
		},
		{
			`{"action":"PostMessage","options":{"message":"Hello everyone, I love cats! 🐱❤️","max_results":5}}`,
			`PostMessage`,
		},
		{
			`{"action":"PostMessage","options":{"message":"Hello everyone! I'm new here and I love cats. Nice to meet you all!"}}`,
			`PostMessage`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := extractAIAction(tt.line)
			if got == nil {
				t.Errorf("extractAIAction(%s) is nil", tt.line)
				return
			}
			if got.Action != tt.wantAction {
				t.Errorf("extractAIAction(%s) = %+v, want %v", tt.line, got, tt.wantAction)
				return
			}
		})
	}
}
