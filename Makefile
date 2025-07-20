# Makefile（根目录）
.PHONY: docker-build docker-push swag

# 构建指定微服务镜像，SERVICE、REGISTRY、IMAGE_TAG 通过环境变量传入
docker-build:
	@echo "Building image for service: $(SERVICE)"
	docker build \
	  --build-arg service_name=$(SERVICE) \
	  -t $(REGISTRY)/forum_be_$(SERVICE):$(IMAGE_TAG) \
	  .

# 推送镜像
docker-push:
	@echo "Pushing image for service: $(SERVICE)"
	docker push $(REGISTRY)/forum_be_$(SERVICE):$(IMAGE_TAG)

# 生成 Swagger 文档（示例）
swag:
	cd microservice/gateway && swag init
