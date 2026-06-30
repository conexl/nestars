"""FastAPI dependencies for building a per-request FragmentClient."""

from __future__ import annotations

import json
from typing import Annotated

from fastapi import Header
from FragmentAPI import FragmentClient


async def get_fragment_client(
    seed: Annotated[str, Header(alias="x-fragment-seed")],
    wallet_version: Annotated[str, Header(alias="x-fragment-wallet-version")] = "V5R1",
    cookies: Annotated[str | None, Header(alias="x-fragment-cookies")] = None,
    api_key: Annotated[str | None, Header(alias="x-fragment-api-key")] = None,
    timeout: Annotated[float, Header(alias="x-fragment-timeout")] = 60.0,
    stats_enabled: Annotated[bool, Header(alias="x-fragment-stats-enabled")] = True,
):
    """Create an async FragmentClient for the duration of the request.

    Credentials are read from HTTP headers so the API remains stateless and
    no secrets are stored on the server.

    Headers:
        x-fragment-seed:            required wallet mnemonic
        x-fragment-wallet-version:  V4R2 or V5R1 (default V5R1)
        x-fragment-cookies:         optional Fragment cookies as JSON string
                                    or semicolon-separated key=value pairs
        x-fragment-api-key:         optional TONAPI key
        x-fragment-timeout:         request timeout (default 60.0)
        x-fragment-stats-enabled:   telemetry flag (default true)
    """
    parsed_cookies = cookies
    if cookies is not None and cookies.strip().startswith("{"):
        parsed_cookies = json.loads(cookies)

    async with FragmentClient(
        seed=seed,
        wallet_version=wallet_version,
        cookies=parsed_cookies,
        api_key=api_key,
        timeout=timeout,
        stats_enabled=stats_enabled,
    ) as client:
        yield client
