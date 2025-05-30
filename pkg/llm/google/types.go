package google

import (
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/sonirico/mcphost/pkg/llm"
)

type ToolCall struct {
	genai.FunctionCall

	toolCallID int
}

func (t *ToolCall) GetName() string {
	if t == nil {
		return ""
	}
	return t.Name
}

func (t *ToolCall) GetArguments() map[string]any {
	if t == nil {
		return nil
	}
	return t.Args
}

func (t *ToolCall) GetID() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("Tool<%d>", t.toolCallID)
}

type Message struct {
	*genai.Candidate

	toolCallID int
}

func (m *Message) GetRole() string {
	if m == nil || m.Candidate == nil || m.Candidate.Content == nil {
		return ""
	}
	return m.Candidate.Content.Role
}

func (m *Message) GetContent() string {
	if m == nil || m.Candidate == nil || m.Candidate.Content == nil || m.Candidate.Content.Parts == nil {
		return ""
	}
	var sb strings.Builder
	for _, part := range m.Candidate.Content.Parts {
		if text, ok := part.(genai.Text); ok {
			sb.WriteString(string(text))
		}
	}
	return sb.String()
}

func (m *Message) GetToolCalls() []llm.ToolCall {
	if m == nil || m.Candidate == nil {
		return nil
	}
	var calls []llm.ToolCall
	for i, call := range m.Candidate.FunctionCalls() {
		calls = append(calls, &ToolCall{call, m.toolCallID + i})
	}
	return calls
}

func (m *Message) IsToolResponse() bool {
	if m == nil || m.Candidate == nil || m.Candidate.Content == nil || m.Candidate.Content.Parts == nil {
		return false
	}
	for _, part := range m.Candidate.Content.Parts {
		if _, ok := part.(*genai.FunctionResponse); ok {
			return true
		}
	}
	return false
}

func (m *Message) GetToolResponseID() string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("Tool<%d>", m.toolCallID)
}

func (m *Message) GetUsage() (input int, output int) {
	if m == nil {
		return 0, 0
	}
	return 0, 0
}
