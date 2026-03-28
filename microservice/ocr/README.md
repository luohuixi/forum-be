OCR 服务用于识别统一身份认证验证码。

本地启动默认读取同目录 `.env`，当前已提供一组开发机可用的覆盖值：

- `FORUM_OCR_OCR_MODELSCOPE_PYTHON_BIN=python3`
- `FORUM_OCR_OCR_MODELSCOPE_WORKSPACE=/tmp/forum-ocr-workspace`
- `FORUM_OCR_OCR_MODELSCOPE_RUNTIME_DIR=/tmp/forum-ocr-runtime`
- `FORUM_OCR_OCR_MODELSCOPE_SKIP_SELF_CHECK=true`

直接启动：

```bash
cd microservice/ocr
go run .
```

说明：

- `workspace` 是 Python 进程工作目录。
- `runtime_dir` 是临时运行目录，用来放 helper 脚本和临时图片。
- `skip_self_check=true` 只是不在启动时做环境自检，不代表实际 OCR 依赖已经装齐。
- 如果要验证真实识别能力，请先安装 `torch`、`torchvision`、`modelscope` 等依赖，再把 `skip_self_check` 改成 `false`。
