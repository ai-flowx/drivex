import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Aliyun
https://cn.aliyun.com/product/bailian

**Set `aliyun/` as a prefix when sending completion requests**

## API Key
```python
# env variable
os.environ['ALIYUN_API_KEY']
```

## Sample Usage
```python
from litellm import completion
import os

os.environ['ALIYUN_API_KEY'] = ""
response = completion(
    model="aliyun/deepseek-r1-distill-llama-70b",
    messages=[
       {"role": "user", "content": "hello from litellm"}
   ],
)
print(response)
```

## Sample Usage - Streaming
```python
from litellm import completion
import os

os.environ['ALIYUN_API_KEY'] = ""
response = completion(
    model="aliyun/deepseek-r1-distill-llama-70b",
    messages=[
       {"role": "user", "content": "hello from litellm"}
   ],
    stream=True
)

for chunk in response:
    print(chunk)
```

<Tabs>
<TabItem value="sdk" label="SDK">

```python
from litellm import completion
import os

os.environ['ALIYUN_API_KEY'] = ""
resp = completion(
    model="aliyun/deepseek-r1-distill-llama-70b",
    messages=[{"role": "user", "content": "Tell me a joke."}],
)

print(
    resp.choices[0].message.reasoning_content
)
```

</TabItem>
<TabItem value="proxy" label="PROXY">

1. Setup config.yaml

```yaml
model_list:
  - model_name: deepseek-r1-distill-llama-70b
    litellm_params:
        model: aliyun/deepseek-r1-distill-llama-70b
        api_key: os.environ/ALIYUN_API_KEY
```

2. Run proxy

```bash
python litellm/proxy/main.py
```

3. Test it!

```bash
curl -L -X POST 'http://0.0.0.0:4000/v1/chat/completions' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer sk-1234' \
-d '{
    "model": "deepseek-r1-distill-llama-70b",
    "messages": [
      {
        "role": "user",
        "content": [
          {
            "type": "text",
            "text": "Hi, how are you ?"
          }
        ]
      }
    ]
}'
```

</TabItem>

</Tabs>
