# AGENTS

## Repo shape (important)
- This repo is a Go multi-module monorepo: root module `github.com/Muxi-X/forum-be` plus one module per service under `microservice/*` (`gateway`, `user`, `post`, `chat`, `feed`, `ocr`).
- Shared packages live at root (`model/`, `pkg/`, `config/`, `client/`, `log/`); service modules use `replace forum => ../../` to import them.
- Real service entrypoints are `microservice/*/main.go`.

## Build and run commands
- Build a single service from its directory: `make` (each service Makefile runs `go mod tidy`, `gofmt -w .`, then `go build -o main`).
- Run service binary locally from service dir after build: `./main`.
- Feed has two modes: `./main` (feed service) and `./main -sub` (subscribe worker).
- Gateway Swagger generation is from repo root: `make swag` (runs `cd microservice/gateway && swag init`).

## Config loading and env gotchas
- Config source order is code-defined in `config/config.go`: try Nacos first (`NACOSDSN`), then fallback to local file.
- Without `NACOSDSN`, services load `./conf/config.yaml` by default (relative to current working dir), except feed sub-mode loads `./conf/config_sub.yaml`.
- Example configs are in `configs/forum_be_*/example.yaml`; local `conf/config*.yaml` is not committed by default.
- `.env` is preloaded in most services via `godotenv.Load()` in `init()`.
- Service discovery depends on etcd; `MICRO_REGISTRY_ADDRESS` is used in deployment/local compose.

## Codegen and generated files
- Proto code is generated per service via `make protoc` in each `microservice/<svc>` directory.
- Proto generation requires `protoc-gen-micro` plugin (`go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest`).
- Gateway Swagger docs are generated artifacts in `microservice/gateway/docs/` and are imported by router code.

## Testing reality
- Test files exist, but not all are pure unit tests (some user auth tests call external auth endpoints).
- Prefer targeted tests first (example: `go test ./pkg/...` at root, or `go test ./...` inside one service module you changed).

## Deployment/release conventions
- CI is Drone (`.drone.yml`), triggered by tags `release-*`.
- Service image builds are tag-pattern specific: `release-user-*`, `release-gateway-*`, `release-post-*`, `release-chat-*`.
- Root `Makefile` builds/pushes per-service images with env vars: `SERVICE`, `REGISTRY`, `IMAGE_TAG`.

## Domain-specific pitfalls worth remembering
- Casbin is initialized in `dao.Init()` of gateway/user/post; permission changes should go through code paths, not direct DB edits.
- OCR service includes a persistent Python worker; see `microservice/ocr/README.md` for required env vars and runtime dirs.

## Rule
- Always reply by Chinese, unless the user explicitly requests English.
- Update AGENTS.md if necessary.