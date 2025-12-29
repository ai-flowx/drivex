"""
Transformation logic from Cohere's /v1/rerank format to NebulaCoder's `/v1/rerank` format.
"""

import uuid
from typing import Any, Dict, List, Optional, Tuple, Union

from httpx import URL, Response

from litellm import LlmProviders
from litellm.llms.base_llm.chat.transformation import LiteLLMLoggingObj
from litellm.llms.base_llm.rerank.transformation import BaseRerankConfig
from litellm.secret_managers.main import get_secret_str
from litellm.types.rerank import (
    OptionalRerankParams,
    RerankBilledUnits,
    RerankResponse,
    RerankResponseMeta,
    RerankTokens,
)


class NebulaCoderRerankConfig(BaseRerankConfig):
    def get_supported_cohere_rerank_params(self, model: str) -> list:
        return [
            "query",
            "top_n",
            "documents",
            "return_documents",
        ]

    def map_cohere_rerank_params(
        self,
        non_default_params: dict,
        model: str,
        drop_params: bool,
        query: str,
        documents: List[Union[str, Dict[str, Any]]],
        custom_llm_provider: Optional[str] = None,
        top_n: Optional[int] = None,
        rank_fields: Optional[List[str]] = None,
        return_documents: Optional[bool] = True,
        max_chunks_per_doc: Optional[int] = None,
        max_tokens_per_doc: Optional[int] = None,
    ) -> OptionalRerankParams:
        optional_params = {}
        supported_params = self.get_supported_cohere_rerank_params(model)
        for k, v in non_default_params.items():
            if k in supported_params:
                optional_params[k] = v
        return OptionalRerankParams(
            **optional_params,
        )

    def get_complete_url(
        self,
        api_base: Optional[str],
        model: str,
        optional_params: Optional[dict] = None,
    ) -> str:
        """
        Override get_complete_url to handle NebulaCoder-specific URL construction.
        If api_base already ends with '/rerank', don't append it again.
        """
        if api_base is None:
            api_base = get_secret_str("NEBULACODER_API_BASE") or "https://api.nebulacoder.ai/v1"

        if api_base.endswith('/rerank'):
            return api_base
        else:
            return f"{api_base}/rerank"

    def validate_environment(
        self,
        headers: dict,
        model: str,
        api_key: Optional[str] = None,
        optional_params: Optional[dict] = None,
    ) -> dict:
        if api_key is None:
            api_key = get_secret_str("NEBULACODER_API_KEY") or get_secret_str("NEBULACODER_TOKEN")

        if api_key is None:
            raise ValueError(
                "NebulaCoder API key is required. Set via `api_key` parameter or `NEBULACODER_API_KEY` environment variable."
            )

        return {
            "accept": "application/json",
            "content-type": "application/json",
            "authorization": f"Bearer {api_key}",
        }

    def transform_rerank_request(
        self, model: str, optional_rerank_params: OptionalRerankParams, headers: Dict
    ) -> Dict:
        return {"model": model, **optional_rerank_params}

    def transform_rerank_response(
        self,
        model: str,
        raw_response: Response,
        model_response: RerankResponse,
        logging_obj: LiteLLMLoggingObj,
        api_key: Optional[str] = None,
        request_data: Dict = {},
        optional_params: Dict = {},
        litellm_params: Dict = {},
    ) -> RerankResponse:
        if raw_response.status_code != 200:
            raise Exception(raw_response.text)

        logging_obj.post_call(original_response=raw_response.text)

        _json_response = raw_response.json()

        # Handle usage information - NebulaCoder only returns total_tokens
        usage = _json_response.get("usage", {})
        _billed_units = RerankBilledUnits(**usage) if usage else RerankBilledUnits()
        _tokens = RerankTokens(**usage) if usage else RerankTokens()
        rerank_meta = RerankResponseMeta(billed_units=_billed_units, tokens=_tokens)

        _results: Optional[List[dict]] = _json_response.get("results")

        if _results is None:
            raise ValueError(f"No results found in the response={_json_response}")

        # Transform results to match LiteLLM's expected format
        # NebulaCoder returns: {"document": {"text": "...", "multi_modal": null}}
        # LiteLLM expects: {"document": {"text": "..."}}
        transformed_results = []
        for result in _results:
            transformed_result = {
                "index": result["index"],
                "relevance_score": result["relevance_score"]
            }

            # Extract only the text field from document
            if "document" in result and isinstance(result["document"], dict):
                transformed_result["document"] = {"text": result["document"].get("text", "")}

            transformed_results.append(transformed_result)

        return RerankResponse(
            id=_json_response.get("id") or str(uuid.uuid4()),
            results=transformed_results,  # type: ignore
            meta=rerank_meta,
        )

    def _get_openai_compatible_provider_info(
        self,
        api_base: Optional[str],
        api_key: Optional[str],
    ) -> Tuple[str, Optional[str], Optional[str]]:
        """
        Returns:
            Tuple[str, Optional[str], Optional[str]]:
                - custom_llm_provider: str
                - api_base: str
                - dynamic_api_key: str
        """
        api_base = (
            api_base or get_secret_str("NEBULACODER_API_BASE") or "https://api.nebulacoder.ai/v1"
        )  # type: ignore
        dynamic_api_key = api_key or (
            get_secret_str("NEBULACODER_API_KEY")
            or get_secret_str("NEBULACODER_TOKEN")
        )
        return LlmProviders.NEBULACODER.value, api_base, dynamic_api_key
