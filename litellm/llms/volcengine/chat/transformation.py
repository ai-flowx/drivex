"""
Translates from OpenAI's `/v1/chat/completions` to Volcengine's `/v3/chat/completions`
"""

from typing import List, Optional, Tuple

from litellm.litellm_core_utils.prompt_templates.common_utils import (
    handle_messages_with_content_list_to_str_conversion,
)
from litellm.secret_managers.main import get_secret_str
from litellm.types.llms.openai import AllMessageValues

from ...openai.chat.gpt_transformation import OpenAIGPTConfig


class VolcengineChatConfig(OpenAIGPTConfig):

    def _transform_messages(
        self, messages: List[AllMessageValues], model: str
    ) -> List[AllMessageValues]:
        """
        Volcengine does not support content in list format.
        """
        messages = handle_messages_with_content_list_to_str_conversion(messages)
        return super()._transform_messages(messages=messages, model=model)

    def _get_openai_compatible_provider_info(
        self, api_base: Optional[str], api_key: Optional[str]
    ) -> Tuple[Optional[str], Optional[str]]:
        api_base = (
            api_base
            or get_secret_str("VOLCENGINE_API_BASE")
            or "https://ark.cn-beijing.volces.com/api/v3"
        )  # type: ignore
        dynamic_api_key = api_key or get_secret_str("VOLCENGINE_API_KEY")
        return api_base, dynamic_api_key

    def get_complete_url(
        self,
        api_base: Optional[str],
        model: str,
        optional_params: dict,
        stream: Optional[bool] = None,
    ) -> str:
        """
        If api_base is not provided, use the default Volcengine /chat/completions endpoint.
        """
        if not api_base:
            api_base = "https://ark.cn-beijing.volces.com/api/v3"

        if not api_base.endswith("/chat/completions"):
            api_base = f"{api_base}/chat/completions"

        return api_base
