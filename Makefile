# Makefile（根目录）
.PHONY: docker-build docker-push

docker-build:
	@# SERVICE、REGISTRY、IMAGE_TAG 通过环境变量传入
	docker build -t $(REGISTRY)/forum_be_$(SERVICE):$(IMAGE_TAG) .

docker-push:
	docker push $(REGISTRY)/forum_be_$(SERVICE):$(IMAGE_TAG)

swag:
	cd microservice/gateway && swag init
