package agentctx

import (
	"context"
	"sync"
)

type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	mu               sync.Mutex
}

func (t *TokenUsage) Add(prompt, completion, total int) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.PromptTokens += prompt
	t.CompletionTokens += completion
	t.TotalTokens += total
}

func (t *TokenUsage) Snapshot() TokenUsage {
	if t == nil {
		return TokenUsage{}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return TokenUsage{
		PromptTokens:     t.PromptTokens,
		CompletionTokens: t.CompletionTokens,
		TotalTokens:      t.TotalTokens,
	}
}

type tokenUsageKey struct{}

func WithTokenUsage(ctx context.Context) context.Context {
	return context.WithValue(ctx, tokenUsageKey{}, &TokenUsage{
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
	})
}

func TokenUsageFromContext(ctx context.Context) *TokenUsage {
	if v, ok := ctx.Value(tokenUsageKey{}).(*TokenUsage); ok {
		return v
	}
	return nil
}

func FinalTokenLogFields(ctx context.Context) map[string]any {
	usage := TokenUsageFromContext(ctx).Snapshot()
	return map[string]any{
		"tokens": map[string]any{
			"prompt_tokens":     usage.PromptTokens,
			"completion_tokens": usage.CompletionTokens,
			"total_tokens":      usage.TotalTokens,
		},
	}
}
