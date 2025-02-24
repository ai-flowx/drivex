# drivex

[![Build Status](https://github.com/ai-flowx/drivex/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/ai-flowx/drivex/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/ai-flowx/drivex)](https://goreportcard.com/report/github.com/ai-flowx/drivex)
[![License](https://img.shields.io/github/license/ai-flowx/drivex.svg)](https://github.com/ai-flowx/drivex/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/ai-flowx/drivex.svg)](https://github.com/ai-flowx/drivex/tags)



## Introduction

*drivex* is the ai api of [ai-flowx](https://github.com/ai-flowx) written in Go.



## Prerequisites

- Go >= 1.22.0



## Build

```bash
make build
```



## Run

```
./bin/drivex /path/to/config.json
```



## Config

```json
{
  "server_port": ":9090",
  "load_balancing": "random",
  "debug": false,
  "services": {
    "aliyun": [
      {
        "models": ["deepseek-r1"],
        "enabled": true,
        "credentials": {
          "api_key": "key"
        },
        "server_url":"https://dashscope.aliyuncs.com/compatible-mode/v1"
      }
    ],
    "siliconflow": [
      {
        "models": ["deepseek-ai/DeepSeek-R1-Distill-Llama-8B"],
        "enabled": true,
        "credentials": {
          "api_key": "key"
        },
        "server_url":"https://api.siliconflow.cn/v1"
      }
    ],
    "volcengine": [
      {
        "models": ["ep-20240612090709-hzjz5"],
        "enabled": true,
        "credentials": {
          "access_key": "key",
          "secret_key": "key"
        },
        "server_url":"https://ark.cn-beijing.volces.com/api/v3"
      }
    ]
  }
}
```



## License

Project License can be found [here](LICENSE).



## Reference

- [simple-one-api](https://github.com/fruitbars/simple-one-api)
