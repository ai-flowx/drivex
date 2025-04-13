"""
Translates from OpenAI's `/v1/chat/completions` to Aliyun's `/v1/chat/completions`
"""

from typing import List, Optional, Tuple

from litellm.litellm_core_utils.prompt_templates.common_utils import (
    handle_messages_with_content_list_to_str_conversion,
)
from litellm.secret_managers.main import get_secret_str
from litellm.types.llms.openai import AllMessageValues

from ...openai.chat.gpt_transformation import OpenAIGPTConfig


class AliyunChatConfig(OpenAIGPTConfig):

    def _transform_messages(
        self, messages: List[AllMessageValues], model: str
    ) -> List[AllMessageValues]:
        """
        Aliyun does not support content in list format.
        """
        messages = handle_messages_with_content_list_to_str_conversion(messages)
        return super()._transform_messages(messages=messages, model=model)

    def _get_openai_compatible_provider_info(
        self, api_base: Optional[str], api_key: Optional[str]
    ) -> Tuple[Optional[str], Optional[str]]:
        api_base = (
            api_base
            or get_secret_str("ALIYUN_API_BASE")
            or "https://dashscope.aliyuncs.com/compatible-mode/v1"
        )  # type: ignore
        dynamic_api_key = api_key or get_secret_str("ALIYUN_API_KEY")
        return api_base, dynamic_api_key

    def get_complete_url(
        self,
        api_base: Optional[str],
        api_key: Optional[str],
        model: str,
        optional_params: dict,
        litellm_params: dict,
        stream: Optional[bool] = None,
    ) -> str:
        """
        Get the complete URL for the API call.

        Returns:
            str: The complete URL for the API call.
        """
        if api_base is None:
            api_base = "https://dashscope.aliyuncs.com/compatible-mode/v1"
        endpoint = "chat/completions"

        # Remove trailing slash from api_base if present
        api_base = api_base.rstrip("/")

        # Check if endpoint is already in the api_base
        if endpoint in api_base:
            return api_base

        return f"{api_base}/{endpoint}"
