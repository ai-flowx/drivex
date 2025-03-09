#!/bin/bash

curl --location 'http://localhost:4000/chat/completions' \
    --header 'Content-Type: application/json' \
    --data '{
    "model": "siliconflow-deepseek-ai-DeepSeek-R1-Distill-Qwen-32B",
    "messages": [
        {
            "role": "user",
            "content": "what llm are you"
        }
    ]
}'
