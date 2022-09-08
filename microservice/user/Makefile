name=user

build:
	go build -o main
fmt:
	go fmt ./...
vet:
	go vet ./...
protoc:
	protoc --proto_path=. --micro_out=. --go_out=. ./proto/user.proto
github:
	git push origin && git push --tags origin
gitea:
	git push --tags muxi
tag:
	git tag release-${name}-${ver}
push: vet tag gitea github