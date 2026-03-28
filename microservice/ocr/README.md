OCR microservice for forum internal captcha recognition.

Runtime config should stay minimal. `python_bin` and `workspace` are fixed by the
image layout and are not expected in Nacos or local YAML. Only keep business or
operational knobs such as `model`, `timeout_ms`, and `self_check_timeout_ms`.

Deployment compose uses the versioned image:
- `registry.cn-shenzhen.aliyuncs.com/muxi/forum_be_ocr:v1.0.0`

Build and push OCR with the same root workflow as other services:

```bash
cd forum-be
make docker-build SERVICE=ocr REGISTRY=registry.cn-shenzhen.aliyuncs.com/muxi IMAGE_TAG=v1.0.0
make docker-push SERVICE=ocr REGISTRY=registry.cn-shenzhen.aliyuncs.com/muxi IMAGE_TAG=v1.0.0
```
