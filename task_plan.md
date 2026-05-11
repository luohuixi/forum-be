# Task Plan

## Goal
Refactor agent knowledge ingestion so `AddKnowledge` stores markdown content directly into the vector database in chunks. Remove vector write behavior from `tool/rag.go` and keep RAG retrieval-only.

## Phases

### Phase 1 — Confirm inputs and existing paths
- [x] Inspect current proto, gateway request, `addKnowledge`, and RAG write/retrieval paths.
- [x] Identify what must stay vs. be removed.

### Phase 2 — Implement markdown ingestion
- [x] Update gateway AddKnowledge request to accept markdown content.
- [x] Implement markdown chunking and vector upsert in agent `AddKnowledge`.
- [x] Keep chunk metadata minimal and useful for retrieval.

### Phase 3 — Remove write path from RAG tool
- [x] Delete vector-store write tool and related helpers from `tool/rag.go`.
- [x] Keep retrieval tool and index config helpers needed for search.

### Phase 4 — Verify
- [ ] Run diagnostics on changed files.
- [ ] Run agent module tests.
- [ ] Check for compile issues in gateway agent handler.

## Decisions
- Markdown chunking should preserve heading/paragraph boundaries when possible.
- AddKnowledge should be a direct ingestion path, not agent-assisted.
- RAG tools should be read-only for retrieval.

## Risks
- Gateway request shape may be out of sync with existing clients.
- Chunking logic must avoid breaking code blocks and headings too aggressively.
