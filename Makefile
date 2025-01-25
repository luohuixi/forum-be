build:
	go build -o main
proto:
	protoc --go_out=plugins=micro:. ./proto/status.proto
github:
	git push origin && git push --tags origin
gitea:
	git push --tags muxi
tag:
	git tag release-${name}-${ver}
buildx:
	@echo -e "make tag=1.0.0 buildx\n"
	docker buildx build --platform linux/amd64 -t registry.cn-shenzhen.aliyuncs.com/muxi/forum-backend:$(tag) .

#push: tag gitea github
