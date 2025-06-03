"""
Translates from OpenAI's `/v1/chat/completions` to AIS's `/v1/chat/completions`
"""
import logging

from typing import List, Optional, Tuple

from litellm.litellm_core_utils.prompt_templates.common_utils import (
    handle_messages_with_content_list_to_str_conversion,
    strip_none_values_from_message,
)
from litellm.secret_managers.main import get_secret_str
from litellm.types.llms.openai import AllMessageValues, ChatCompletionRequest

from ...openai.chat.gpt_transformation import OpenAIGPTConfig

logger = logging.getLogger(__name__)


class AISChatConfig(OpenAIGPTConfig):

    def _transform_messages(
        self, messages: List[AllMessageValues], model: str
    ) -> List[AllMessageValues]:
        messages = handle_messages_with_content_list_to_str_conversion(messages)

        new_messages: List[AllMessageValues] = []
        for m in messages:
            m = strip_none_values_from_message(m)  # prevents 'extra_forbidden' error
            new_messages.append(m)

        return new_messages

    def _get_openai_compatible_provider_info(
        self, api_base: Optional[str], api_key: Optional[str]
    ) -> Tuple[Optional[str], Optional[str]]:
        api_base = (
            api_base
            or get_secret_str("AIS_API_BASE")
            or "https://api.ais.ai/v1"
        )  # type: ignore
        dynamic_api_key = api_key or get_secret_str("AIS_API_KEY")
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
            api_base = "https://api.ais.ai/v1"
        endpoint = "chat/completions"

        # Remove trailing slash from api_base if present
        api_base = api_base.rstrip("/")

        # Check if endpoint is already in the api_base
        if endpoint in api_base:
            return api_base

        return f"{api_base}/{endpoint}"

    def transform_request(
            self,
            model: str,
            messages: List[AllMessageValues],
            optional_params: dict,
            litellm_params: dict,
            headers: dict,
        ) -> dict:
            if "max_retries" in optional_params:
                logger.warning("`max_retries` is not supported. It will be ignored.")
                optional_params.pop("max_retries", None)
            messages = self._transform_messages(messages=messages, model=model)
            return dict(
                ChatCompletionRequest(
                    model=model, messages=messages, **optional_params
                )
            )
