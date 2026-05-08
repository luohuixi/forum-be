package core

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"forum/log"
	"forum/pkg/agentctx"
	"forum/pkg/tracer"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	toolcallback "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type callbackStartKey struct{}

const logMaxRuneCount = 800

type LoggingCallback struct {
	dev bool
}

func NewLoggingCallback() *LoggingCallback {
	return &LoggingCallback{dev: viper.GetString("env") == "dev"}
}

//func (h *LoggingCallback) Needed(_ context.Context, _ *callbacks.RunInfo, timing callbacks.CallbackTiming) bool {
//	if timing == callbacks.TimingOnStart {
//		return h.dev
//	}
//
//	return timing == callbacks.TimingOnError || timing == callbacks.TimingOnEnd
//}

func (h *LoggingCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	if !h.dev || !isLoggableInput(input, info) {
		return context.WithValue(ctx, callbackStartKey{}, time.Now())
	}

	log.Info("agent callback start",
		zap.String("component", fmt.Sprint(info.Component)),
		zap.String("name", info.Name),
		zap.Any("message", callbackInputString(input)),
		zap.String("trace_id", tracer.GetTraceId(ctx)),
	)

	return context.WithValue(ctx, callbackStartKey{}, time.Now())
}

func (h *LoggingCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	if usage := callbackTokenSummary(output); usage != nil {
		agentctx.TokenUsageFromContext(ctx).Add(
			usage.PromptTokens,
			usage.CompletionTokens,
			usage.TotalTokens,
		)
	}

	if !h.dev || !isLoggableOutput(output, info) {
		return context.WithValue(ctx, callbackStartKey{}, time.Now())
	}

	log.Info("agent callback end",
		zap.String("component", fmt.Sprint(info.Component)),
		zap.String("name", info.Name),
		zap.Duration("duration", callbackDuration(ctx)),
		zap.Any("message", callbackOutputString(output)),
		zap.String("trace_id", tracer.GetTraceId(ctx)),
	)

	return ctx
}

func (h *LoggingCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	log.Error("agent callback error",
		zap.String("component", fmt.Sprint(info.Component)),
		zap.String("name", info.Name),
		zap.Duration("duration", callbackDuration(ctx)),
		zap.Error(err),
		zap.String("trace_id", tracer.GetTraceId(ctx)),
	)

	return ctx
}

func (h *LoggingCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	return ctx
}

func (h *LoggingCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	return ctx
}

func isLoggableInput(input callbacks.CallbackInput, info *callbacks.RunInfo) bool {
	return (model.ConvCallbackInput(input) != nil || toolcallback.ConvCallbackInput(input) != nil) && !isHiddenComponent(info)
}

func isLoggableOutput(output callbacks.CallbackOutput, info *callbacks.RunInfo) bool {
	return (model.ConvCallbackOutput(output) != nil || toolcallback.ConvCallbackOutput(output) != nil) && !isHiddenComponent(info)
}

func isHiddenComponent(info *callbacks.RunInfo) bool {
	switch fmt.Sprint(info.Component) {
	case "Graph", "ToolsNode", "ChatTemplate", "Retriever", "Embedding":
		return true
	default:
		return false
	}
}

func callbackDuration(ctx context.Context) time.Duration {
	started, ok := ctx.Value(callbackStartKey{}).(time.Time)
	if !ok {
		return 0
	}
	return time.Since(started)
}

type logMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func callbackInputString(input callbacks.CallbackInput) any {
	if toolInput := toolcallback.ConvCallbackInput(input); toolInput != nil {
		return struct {
			Arguments string `json:"arguments"`
		}{Arguments: truncateString(toolInput.ArgumentsInJSON)}
	}

	modelInput := model.ConvCallbackInput(input)
	if modelInput == nil {
		return truncateAny(fmt.Sprintf("%#v", input))
	}

	messages := make([]logMessage, 0, len(modelInput.Messages))
	for _, msg := range modelInput.Messages {
		messages = append(messages, logMessage{
			Role:    string(msg.Role),
			Content: truncateString(msg.Content),
		})
	}

	return struct {
		Messages []logMessage `json:"messages"`
	}{Messages: messages}
}

func callbackOutputString(output callbacks.CallbackOutput) any {
	if toolOutput := toolcallback.ConvCallbackOutput(output); toolOutput != nil {
		return struct {
			Response string `json:"response"`
		}{Response: truncateString(toolOutput.Response)}
	}

	modelOutput := model.ConvCallbackOutput(output)
	if modelOutput == nil || modelOutput.Message == nil {
		return truncateAny(fmt.Sprintf("%#v", output))
	}

	return struct {
		Message logMessage `json:"message"`
	}{Message: logMessage{
		Role:    string(modelOutput.Message.Role),
		Content: truncateString(modelOutput.Message.Content),
	}}
}

func callbackTokenSummary(output callbacks.CallbackOutput) *agentctx.TokenUsage {
	modelOutput := model.ConvCallbackOutput(output)
	if modelOutput == nil || modelOutput.Message == nil || modelOutput.Message.ResponseMeta == nil || modelOutput.Message.ResponseMeta.Usage == nil {
		return nil
	}

	return &agentctx.TokenUsage{
		PromptTokens:     modelOutput.Message.ResponseMeta.Usage.PromptTokens,
		CompletionTokens: modelOutput.Message.ResponseMeta.Usage.CompletionTokens,
		TotalTokens:      modelOutput.Message.ResponseMeta.Usage.TotalTokens,
	}
}

func truncateAny(v string) string {
	return truncateString(v)
}

func truncateString(s string) string {
	if s == "" {
		return s
	}
	if utf8.RuneCountInString(s) <= logMaxRuneCount {
		return s
	}

	r := []rune(s)
	return string(r[:logMaxRuneCount]) + "...<truncated>"
}

func InitCallbacks() {
	callbacks.AppendGlobalHandlers(NewLoggingCallback())
}
