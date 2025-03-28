#!/bin/bash

curl --location 'http://localhost:4000/chat/completions' \
    --header 'Authorization: Bearer sk-y2gMsalbslZHYUp0Sl8DUQ' \
    --header 'Content-Type: application/json' \
    --data '{
        "model": "aliyun-deepseek-r1-distill-llama-70b",
        "messages": [
            {
                "role": "user",
                "content": "what llm are you"
            }
        ]
    }'

curl --location 'http://localhost:4000/chat/completions' \
    --header 'Authorization: Bearer sk-y2gMsalbslZHYUp0Sl8DUQ' \
    --header 'Content-Type: application/json' \
    --data '{
        "model": "gemini/gemini-2.0-flash",
        "messages": [
            {
                "role": "user",
                "content": "what llm are you"
            }
        ]
    }'

curl --location 'http://localhost:4000/chat/completions' \
    --header 'Authorization: Bearer sk-y2gMsalbslZHYUp0Sl8DUQ' \
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

curl --location 'http://localhost:4000/chat/completions' \
    --header 'Authorization: Bearer sk-y2gMsalbslZHYUp0Sl8DUQ' \
    --header 'Content-Type: application/json' \
    --data '{
        "model": "volcengine-deepseek-v3-241226",
        "messages": [
            {
                "role": "user",
                "content": "what llm are you"
            }
        ]
    }'
