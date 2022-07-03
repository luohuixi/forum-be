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
buildx: check
	@echo -e "make tag=1.0.0 buildx\n"
	docker buildx build --platform linux/amd64 -t registry.cn-shenzhen.aliyuncs.com/muxi/forum-backend:$(tag) .
check:
	grep -qr "2020213675" ./microservice && echo -e "exist!!!!!!!!!!!\nexist!!!!!!!!!!!\nexist!!!!!!!!!!!\nexist!!!!!!!!!!!\nexist!!!!!!!!!!!\nexist!!!!!!!!!!!\nexist!!!!!!!!!!!\n" || echo "ok"

#push: tag gitea github
