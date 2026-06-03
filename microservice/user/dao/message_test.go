package dao

import (
	"encoding/json"
	"testing"
)

func TestDeduplicateInteractionMessages(t *testing.T) {
	messages := []string{
		mustMessage(t, map[string]interface{}{
			"id":           "collection-old",
			"send_user_id": "113",
			"post_id":      "444",
			"comment_id":   "0",
			"type":         "collection",
			"read":         true,
			"created_at":   "2026-06-03 19:46:21",
		}),
		mustMessage(t, map[string]interface{}{
			"id":           "collection-new",
			"send_user_id": "113",
			"post_id":      "444",
			"comment_id":   "0",
			"type":         "collection",
			"read":         false,
			"created_at":   "2026-06-03 19:47:21",
		}),
		mustMessage(t, map[string]interface{}{
			"id":           "like-one",
			"send_user_id": "113",
			"post_id":      "444",
			"comment_id":   "0",
			"type":         "like",
			"read":         false,
			"created_at":   "2026-06-03 19:48:21",
		}),
	}

	got := deduplicateInteractionMessages(messages)
	if len(got) != 2 {
		t.Fatalf("expected two interaction notifications, got %d", len(got))
	}

	collection := mustDecodeMessage(t, got[0])
	if collection["id"] != "collection-new" {
		t.Fatalf("expected newest collection notification to win, got %v", collection["id"])
	}
	if collection["read"] != false {
		t.Fatalf("expected merged collection notification to remain unread")
	}
}

func TestMergeInteractionIntoMessagesCompressesSameInteraction(t *testing.T) {
	messages := []string{
		mustMessage(t, map[string]interface{}{
			"id":           "old-newest",
			"send_user_id": "113",
			"post_id":      "444",
			"comment_id":   "0",
			"type":         "collection",
			"read":         true,
			"created_at":   "2026-06-03 19:47:21",
		}),
		mustMessage(t, map[string]interface{}{
			"id":           "old-older",
			"send_user_id": "113",
			"post_id":      "444",
			"comment_id":   "0",
			"type":         "collection",
			"read":         false,
			"created_at":   "2026-06-03 19:46:21",
		}),
		mustMessage(t, map[string]interface{}{
			"id":           "other",
			"send_user_id": "113",
			"post_id":      "21",
			"comment_id":   "0",
			"type":         "collection",
			"read":         false,
			"created_at":   "2026-06-03 19:45:21",
		}),
	}
	next := map[string]interface{}{
		"id":           "new-request",
		"send_user_id": "113",
		"post_id":      "444",
		"comment_id":   "0",
		"type":         "collection",
		"read":         false,
		"created_at":   "2026-06-03 19:49:21",
	}

	got, err := mergeInteractionIntoMessages(messages, next)
	if err != nil {
		t.Fatalf("mergeInteractionIntoMessages returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected compressed list length 2, got %d", len(got))
	}

	merged := mustDecodeMessage(t, got[0])
	if merged["id"] != "old-newest" {
		t.Fatalf("expected stable notification id old-newest, got %v", merged["id"])
	}
	if merged["read"] != false {
		t.Fatalf("expected repeated interaction to become unread")
	}
	if merged["created_at"] != "2026-06-03 19:49:21" {
		t.Fatalf("expected newest created_at, got %v", merged["created_at"])
	}
}

func mustMessage(t *testing.T, data map[string]interface{}) string {
	t.Helper()
	msg, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}
	return string(msg)
}

func mustDecodeMessage(t *testing.T, msg string) map[string]interface{} {
	t.Helper()
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &data); err != nil {
		t.Fatalf("unmarshal message: %v", err)
	}
	return data
}
