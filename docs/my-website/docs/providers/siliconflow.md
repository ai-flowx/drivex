import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# SiliconFlow
https://siliconflow.cn/

**Set `siliconflow/` as a prefix when sending completion requests**

## API Key
```python
# env variable
os.environ['SILICONFLOW_API_KEY']
```

## Sample Usage
```python
from litellm import completion
import os

os.environ['SILICONFLOW_API_KEY'] = ""
response = completion(
    model="siliconflow/deepseek-ai/DeepSeek-R1-Distill-Qwen-32B",
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

os.environ['SILICONFLOW_API_KEY'] = ""
response = completion(
    model="siliconflow/deepseek-ai/DeepSeek-R1-Distill-Qwen-32B",
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

os.environ['SILICONFLOW_API_KEY'] = ""
resp = completion(
    model="siliconflow/deepseek-ai/DeepSeek-R1-Distill-Qwen-32B",
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
  - model_name: deepseek-ai/DeepSeek-R1-Distill-Qwen-32B
    litellm_params:
        model: siliconflow/deepseek-ai/DeepSeek-R1-Distill-Qwen-32B
        api_key: os.environ/SILICONFLOW_API_KEY
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
    "model": "siliconflow/deepseek-ai/DeepSeek-R1-Distill-Qwen-32B",
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
