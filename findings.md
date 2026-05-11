# Findings

## Current state
- `microservice/agent/proto/agent.proto` already defines `AddKnowledgeRequest` with only `content`.
- `microservice/agent/service/addKnowledge.go` is currently empty aside from logging.
- `microservice/agent/tool/rag.go` still contains both vector-store write and vector-search code.
- Gateway `AddKnowledgeRequest` still uses `post_id`, so it is out of sync with proto.

## Strategy
- Use heading/paragraph-aware markdown chunking with a size cap and small overlap.
- Include chunk metadata such as `chunk_order` and `section_path`.
- Keep RAG write logic out of tools; retrieval remains in the agent tool.

## Implementation notes
- Gateway `AddKnowledgeRequest` now accepts `content` instead of `post_id`.
- Agent `AddKnowledge` now writes markdown chunks directly via ES indexer.
- `tool/rag.go` now only exposes `vector_search`.
