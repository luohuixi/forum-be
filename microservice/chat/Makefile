name=chat

all: tidy gotool
	@go build -o main -v .
gotool:
	gofmt -w .
tidy:
	go mod tidy
build:
	go build -o main
fmt:
	go fmt ./...
vet:
	go vet ./...
protoc:
	protoc --proto_path=. --micro_out=. --go_out=. ./proto/chat.proto
github:
	git push origin && git push --tags origin
gitea:
	git push --tags muxi
tag:
	git tag release-${name}-${ver}
push: vet tag gitea github

docker:
	docker buildx build --platform linux/amd64 -t forum ../.. --build-arg service_name=${name}