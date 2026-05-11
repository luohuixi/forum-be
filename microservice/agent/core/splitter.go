package core

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	markdownsplitter "github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	semanticsplitter "github.com/cloudwego/eino-ext/components/document/transformer/splitter/semantic"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func MarkdownSplitter(ctx context.Context, postId uint32, content string, splitterSize uint32) ([]*schema.Document, error) {
	transformer, err := markdownsplitter.NewHeaderSplitter(ctx, &markdownsplitter.HeaderConfig{
		Headers:     MarkdownSize(int(splitterSize)),
		TrimHeaders: true,
		IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
			return fmt.Sprintf("%s-%04d", originalID, splitIndex+1)
		},
	})
	if err != nil {
		return nil, err
	}

	return splitter(ctx, transformer, postId, content)
}

func SemanticSpliter(ctx context.Context, postId uint32, content string, splitterSize uint32) ([]*schema.Document, error) {
	transformer, err := semanticsplitter.NewSplitter(ctx, &semanticsplitter.Config{
		Embedding:    NewLimitedEmbedder(Embedder, 10),
		BufferSize:   1,
		MinChunkSize: int(splitterSize),
		Separators:   []string{".", "?", "!", "。", "！", "？"},
		Percentile:   0.8,
		IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
			return fmt.Sprintf("%s-%04d", originalID, splitIndex+1)
		},
	})
	if err != nil {
		return nil, err
	}

	return splitter(ctx, transformer, postId, markdownToText(content))
}

func splitter(ctx context.Context, transformer document.Transformer, postId uint32, content string) ([]*schema.Document, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	originalID := contentHash(content)
	splitDocs, err := transformer.Transform(ctx, []*schema.Document{{
		ID:      originalID,
		Content: content,
	}})
	if err != nil {
		return nil, err
	}

	for _, doc := range splitDocs {
		if doc.MetaData == nil {
			doc.MetaData = make(map[string]any)
		}
		doc.MetaData["chunk_id"] = doc.ID
		doc.MetaData["post_id"] = postId
	}

	return splitDocs, nil
}

func contentHash(content string) string {
	sum := sha1.Sum([]byte(content))
	return hex.EncodeToString(sum[:])
}

var markdownHeaderLevels = []struct {
	header string
	level  string
}{
	{"#", "h1"},
	{"##", "h2"},
	{"###", "h3"},
	{"####", "h4"},
	{"#####", "h5"},
	{"######", "h6"},
}

func MarkdownSize(maxLevel int) map[string]string {
	if maxLevel <= 0 {
		maxLevel = 6
	}

	headers := make(map[string]string)
	for i, item := range markdownHeaderLevels {
		if i+1 > maxLevel {
			break
		}
		headers[item.header] = item.level
	}
	return headers
}

// markdownToText 将Markdown内容转换为纯文本
func markdownToText(markdownText string) string {
	source := []byte(markdownText)
	doc := goldmark.New().Parser().Parse(text.NewReader(source))

	var b strings.Builder
	appendMarkdownText(&b, doc, source)
	return normalizePlainText(b.String())
}

func appendMarkdownText(b *strings.Builder, node ast.Node, source []byte) {
	switch n := node.(type) {
	case *ast.Image:
		return
	case *ast.Text:
		b.Write(n.Text(source))
		if n.HardLineBreak() {
			b.WriteByte('\n')
		} else if n.SoftLineBreak() {
			b.WriteByte(' ')
		}
		return
	case *ast.String:
		b.Write(n.Value)
		return
	case *ast.CodeBlock:
		writeTextBlock(b, blockText(n, source))
		return
	case *ast.FencedCodeBlock:
		writeTextBlock(b, blockText(n, source))
		return
	}

	isBlock := node.Type() == ast.TypeBlock
	if isBlock {
		writeBlockSeparator(b)
	}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		appendMarkdownText(b, child, source)
	}
	if isBlock {
		writeBlockSeparator(b)
	}
}

func writeBlockSeparator(b *strings.Builder) {
	if b.Len() == 0 {
		return
	}
	current := b.String()
	if strings.HasSuffix(current, "\n") {
		return
	}
	b.WriteByte('\n')
}

func normalizePlainText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	previousBlank := true
	for _, line := range lines {
		line = strings.TrimRight(line, " \t\r")
		if strings.TrimSpace(line) == "" {
			if len(out) == 0 || previousBlank {
				continue
			}
			out = append(out, "")
			previousBlank = true
			continue
		}
		out = append(out, line)
		previousBlank = false
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return strings.Join(out, "\n")
}

func blockText(node ast.Node, source []byte) string {
	lines := node.Lines()
	if lines == nil || lines.Len() == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		line := strings.TrimRight(string(seg.Value(source)), "\r\n")
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
	}
	return b.String()
}

func writeTextBlock(b *strings.Builder, content string) {
	content = strings.TrimSpace(content)
	if content == "" {
		return
	}
	if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") {
		b.WriteByte('\n')
	}
	b.WriteString(content)
	b.WriteByte('\n')
}
