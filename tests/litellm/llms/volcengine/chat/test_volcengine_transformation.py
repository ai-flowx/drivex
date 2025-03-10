import json
import os
import sys
from unittest.mock import AsyncMock, MagicMock, patch

import pytest

sys.path.insert(
    0, os.path.abspath("../../../../..")
)  # Adds the parent directory to the system path
import litellm


@pytest.mark.parametrize("is_async", [True, False])
@pytest.mark.asyncio
async def test_volcengine_request_format(is_async):
    """
    Test that Volcengine requests are formatted correctly with the proper endpoint and parameters
    for both synchronous and asynchronous calls
    """
    litellm._turn_on_debug()

    # Set up the test parameters
    api_key = "key"
    api_base = "https://ark.cn-beijing.volces.com/api/v3"
    model = "volcengine/deepseek-v3-241226"
    messages = [
        {"role": "system", "content": "you are ai assistant"},
        {"role": "user", "content": "Hello!"}
    ]

    if is_async:
        # Mock AsyncHTTPHandler.post method for async test
        with patch(
            "litellm.llms.custom_httpx.llm_http_handler.AsyncHTTPHandler.post"
        ) as mock_post:
            # Set up mock response
            mock_post.return_value = AsyncMock()

            # Call the acompletion function
            try:
                await litellm.acompletion(
                    custom_llm_provider="volcengine",
                    api_key=api_key,
                    api_base=api_base,
                    model=model,
                    messages=messages,
                )
            except Exception as e:
                # We expect an exception since we're mocking the response
                pass

    else:
        # Mock HTTPHandler.post method for sync test
        with patch(
            "litellm.llms.custom_httpx.llm_http_handler.HTTPHandler.post"
        ) as mock_post:
            # Set up mock response
            mock_post.return_value = MagicMock()

            # Call the completion function
            try:
                litellm.completion(
                    custom_llm_provider="volcengine",
                    api_key=api_key,
                    api_base=api_base,
                    model=model,
                    messages=messages,
                )
            except Exception as e:
                # We expect an exception since we're mocking the response
                pass

    # Verify the request was made with the correct parameters
    mock_post.assert_called_once()
    call_args = mock_post.call_args
    print("sync request call=", json.dumps(call_args.kwargs, indent=4, default=str))

    # Check URL
    assert api_base in call_args.kwargs["url"]

    # Check headers
    assert api_key in call_args.kwargs["headers"]["Authorization"]

    # Check request body
    request_body = json.loads(call_args.kwargs["data"])
    assert (
        request_body["model"] == "deepseek-v3-241226"
    )  # Model name should be stripped of provider prefix
    assert request_body["messages"] == messages
