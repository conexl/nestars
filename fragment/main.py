"""FastAPI wrapper exposing all FragmentClient capabilities.

Run with:
    source .venv/bin/activate
    uvicorn main:app --reload

OpenAPI/Swagger UI is available at /docs by default.
"""

from __future__ import annotations

import json
from dataclasses import asdict
from typing import Annotated, Any

from fastapi import Depends, FastAPI, HTTPException
from FragmentAPI import FragmentClient

import config
import schemas
from dependencies import get_fragment_client


app = FastAPI(
    title=config.settings.APP_NAME,
    description=config.settings.APP_DESCRIPTION,
    version=config.settings.APP_VERSION,
    docs_url=config.settings.DOCS_URL,
    redoc_url=config.settings.REDOC_URL,
    openapi_url=config.settings.OPENAPI_URL,
)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------


def _serialize(obj: Any) -> Any:
    """Serialize dataclass/list/datetime objects to JSON-compatible structures."""
    from dataclasses import is_dataclass
    from datetime import datetime

    if is_dataclass(obj) and not isinstance(obj, type):
        return asdict(obj)
    if isinstance(obj, list):
        return [_serialize(item) for item in obj]
    if isinstance(obj, dict):
        return {k: _serialize(v) for k, v in obj.items()}
    if isinstance(obj, datetime):
        return obj.isoformat()
    return obj


async def _call(client: FragmentClient, method: str, *args, **kwargs) -> Any:
    """Call a client method and translate library exceptions to HTTP errors."""
    try:
        return await getattr(client, method)(*args, **kwargs)
    except Exception as exc:  # noqa: BLE001
        raise HTTPException(status_code=400, detail=str(exc)) from exc


# ---------------------------------------------------------------------------
# Authentication
# ---------------------------------------------------------------------------


@app.post("/auth/authenticate", response_model=schemas.AuthResponse, tags=["Auth"])
async def authenticate(payload: schemas.AuthRequest):
    """Authenticate with Fragment via TON wallet and Telegram.

    Returns session cookies that can later be passed in the
    x-fragment-cookies header for full-mode requests.
    """
    try:
        cookies = await FragmentClient.authenticate(
            seed=payload.seed,
            wallet_version=payload.wallet_version,
            phone=payload.phone,
            print_qr=payload.print_qr,
            timeout=payload.timeout,
        )
    except Exception as exc:  # noqa: BLE001
        raise HTTPException(status_code=400, detail=str(exc)) from exc

    # Return cookies in a structured response instead of profile because auth
    # itself does not need a client; the caller can use the cookies elsewhere.
    return {"cookies": cookies}


# ---------------------------------------------------------------------------
# Wallet & Profile
# ---------------------------------------------------------------------------


@app.get("/wallet", response_model=schemas.WalletInfo, tags=["Wallet"])
async def get_wallet(client: Annotated[FragmentClient, Depends(get_fragment_client)]):
    """Return wallet address, state, TON and USDT balance."""
    return _serialize(await _call(client, "get_wallet"))


@app.get("/profile", response_model=schemas.ProfileInfo, tags=["Profile"])
async def get_profile(client: Annotated[FragmentClient, Depends(get_fragment_client)]):
    """Return Fragment account profile information (full mode only)."""
    return _serialize(await _call(client, "get_profile"))


# ---------------------------------------------------------------------------
# Recipients
# ---------------------------------------------------------------------------


@app.post("/recipients/stars", response_model=schemas.RecipientInfo | None, tags=["Recipients"])
async def get_stars_recipient(
    payload: schemas.RecipientUsernameRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Resolve a Telegram user for a Stars gift (full mode only)."""
    return _serialize(await _call(client, "get_stars_recipient", payload.username))


@app.post("/recipients/premium", response_model=schemas.RecipientInfo | None, tags=["Recipients"])
async def get_premium_recipient(
    payload: schemas.PremiumRecipientRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Resolve a Telegram user for a Premium gift (full mode only)."""
    return _serialize(await _call(client, "get_premium_recipient", payload.username, payload.months))


@app.post("/recipients/ads-topup", response_model=schemas.RecipientInfo | None, tags=["Recipients"])
async def get_ads_topup_recipient(
    payload: schemas.RecipientUsernameRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Resolve a Telegram user for an Ads TON top-up (full mode only)."""
    return _serialize(await _call(client, "get_ads_topup_recipient", payload.username))


@app.post("/recipients/giveaway/stars", response_model=schemas.RecipientInfo | None, tags=["Recipients"])
async def get_giveaway_stars_recipient(
    payload: schemas.GiveawayStarsRecipientRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Resolve a channel for a Stars giveaway (full mode only)."""
    return _serialize(await _call(client, "get_giveaway_stars_recipient", payload.channel, payload.winners, payload.amount))


@app.post("/recipients/giveaway/premium", response_model=schemas.RecipientInfo | None, tags=["Recipients"])
async def get_giveaway_premium_recipient(
    payload: schemas.GiveawayPremiumRecipientRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Resolve a channel for a Premium giveaway (full mode only)."""
    return _serialize(await _call(client, "get_giveaway_premium_recipient", payload.channel, payload.winners, payload.months))


# ---------------------------------------------------------------------------
# Purchases
# ---------------------------------------------------------------------------


@app.post("/purchase/stars", response_model=Any, tags=["Purchases"])
async def purchase_stars(
    payload: schemas.PurchaseStarsRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Purchase Telegram Stars.

    Works in full mode and no-cookie mode.  When an EVM payment method is
    selected the response is an EvmPaymentResult with an invoice address.
    """
    return _serialize(
        await _call(
            client,
            "purchase_stars",
            payload.username,
            payload.amount,
            payload.show_sender,
            payload.payment_method,
        )
    )


@app.post("/purchase/premium", response_model=Any, tags=["Purchases"])
async def purchase_premium(
    payload: schemas.PurchasePremiumRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Gift Telegram Premium."""
    return _serialize(
        await _call(
            client,
            "purchase_premium",
            payload.username,
            payload.months,
            payload.show_sender,
            payload.payment_method,
        )
    )


@app.post("/topup/ton", response_model=schemas.AdsTopupResult, tags=["Purchases"])
async def topup_ton(
    payload: schemas.TopupTonRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Top up TON to a Telegram Ads balance."""
    return _serialize(await _call(client, "topup_ton", payload.username, payload.amount, payload.show_sender))


@app.post("/purchase/batch", response_model=schemas.BatchResult, tags=["Purchases"])
async def batch_purchase(
    payload: schemas.BatchPurchaseRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Execute multiple Stars / Premium / TON purchases in batched transactions."""
    items = [item.model_dump() for item in payload.items]
    return _serialize(await _call(client, "batch_purchase", items, payload.payment_method))


# ---------------------------------------------------------------------------
# Giveaways
# ---------------------------------------------------------------------------


@app.post("/giveaway/stars", response_model=Any, tags=["Giveaways"])
async def giveaway_stars(
    payload: schemas.GiveawayStarsRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Start a Telegram Stars giveaway for a channel."""
    return _serialize(await _call(client, "giveaway_stars", payload.channel, payload.winners, payload.amount, payload.payment_method))


@app.post("/giveaway/premium", response_model=Any, tags=["Giveaways"])
async def giveaway_premium(
    payload: schemas.GiveawayPremiumRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Start a Telegram Premium giveaway for a channel."""
    return _serialize(await _call(client, "giveaway_premium", payload.channel, payload.winners, payload.months, payload.payment_method))


# ---------------------------------------------------------------------------
# Marketplace
# ---------------------------------------------------------------------------


@app.post("/bids", response_model=schemas.BidResult, tags=["Marketplace"])
async def place_bid(
    payload: schemas.PlaceBidRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Place a bid or instant-buy on a username, number or gift."""
    return _serialize(await _call(client, "place_bid", payload.item_type, payload.slug, payload.bid))


@app.get("/search/usernames", response_model=schemas.UsernamesResult, tags=["Marketplace"])
async def search_usernames(
    query: str = "",
    sort: str | None = None,
    filter: str | None = None,
    offset_id: str | None = None,
    client: FragmentClient = Depends(get_fragment_client),
):
    """Search Fragment username marketplace."""
    return _serialize(await _call(client, "search_usernames", query, sort, filter, offset_id))


@app.get("/search/numbers", response_model=schemas.NumbersResult, tags=["Marketplace"])
async def search_numbers(
    query: str = "",
    sort: str | None = None,
    filter: str | None = None,
    offset_id: str | None = None,
    client: FragmentClient = Depends(get_fragment_client),
):
    """Search Fragment anonymous-number marketplace."""
    return _serialize(await _call(client, "search_numbers", query, sort, filter, offset_id))


@app.get("/search/gifts", response_model=schemas.GiftsResult, tags=["Marketplace"])
async def search_gifts(
    query: str = "",
    collection: str | None = None,
    sort: str | None = None,
    filter: str | None = None,
    view: str | None = None,
    attr: str | None = None,
    offset: int | None = None,
    client: FragmentClient = Depends(get_fragment_client),
):
    """Search Fragment gifts marketplace.

    The `attr` query parameter is expected as a JSON string, e.g.
    `{"Background": ["Space"], "Model": ["Gold"]}`.
    """
    parsed_attr: dict[str, list[str]] | None = None
    if attr:
        parsed_attr = json.loads(attr)
    return _serialize(await _call(client, "search_gifts", query, collection, sort, filter, view, parsed_attr, offset))


@app.get("/usernames/{username}", response_model=schemas.UsernameInfo, tags=["Marketplace"])
async def get_username_info(
    username: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Get detailed information about a username."""
    return _serialize(await _call(client, "get_username_info", username))


@app.get("/numbers/{number}", response_model=schemas.NumberInfo, tags=["Marketplace"])
async def get_number_info(
    number: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Get detailed information about an anonymous number."""
    return _serialize(await _call(client, "get_number_info", number))


@app.get("/gifts/{slug}", response_model=schemas.GiftInfo, tags=["Marketplace"])
async def get_gift_info(
    slug: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Get detailed information about a gift."""
    return _serialize(await _call(client, "get_gift_info", slug))


# ---------------------------------------------------------------------------
# Prices
# ---------------------------------------------------------------------------


@app.get("/prices/stars", response_model=schemas.StarsPrices, tags=["Prices"])
async def get_stars_prices(
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Return all available Telegram Stars package prices."""
    return _serialize(await _call(client, "get_stars_prices"))


@app.get("/prices/stars/{quantity}", response_model=schemas.StarsPrice, tags=["Prices"])
async def get_stars_price(
    quantity: int,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Return price for a specific Stars quantity."""
    return _serialize(await _call(client, "get_stars_price", quantity))


@app.get("/prices/premium", response_model=schemas.PremiumPrices, tags=["Prices"])
async def get_premium_prices(
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Return Telegram Premium subscription prices."""
    return _serialize(await _call(client, "get_premium_prices"))


# ---------------------------------------------------------------------------
# Transaction history
# ---------------------------------------------------------------------------


@app.get("/history/stars", response_model=list[schemas.StarsTransaction], tags=["History"])
async def get_stars_history(
    sort: str = "desc",
    client: FragmentClient = Depends(get_fragment_client),
):
    """Return Telegram Stars transaction history."""
    return _serialize(await _call(client, "get_stars_history", sort))


@app.get("/history/premium", response_model=list[schemas.PremiumTransaction], tags=["History"])
async def get_premium_history(
    sort: str = "desc",
    client: FragmentClient = Depends(get_fragment_client),
):
    """Return Telegram Premium gift history."""
    return _serialize(await _call(client, "get_premium_history", sort))


@app.get("/history/topup", response_model=list[schemas.TopupTransaction], tags=["History"])
async def get_topup_history(
    sort: str = "asc",
    client: FragmentClient = Depends(get_fragment_client),
):
    """Return Telegram Ads TON top-up history."""
    return _serialize(await _call(client, "get_topup_history", sort))


# ---------------------------------------------------------------------------
# My assets
# ---------------------------------------------------------------------------


@app.get("/my/bids", response_model=schemas.MyBidsResult, tags=["My Assets"])
async def get_my_bids(
    item_type: str = "usernames",
    sort: str = "desc",
    client: FragmentClient = Depends(get_fragment_client),
):
    """Return My Bid History."""
    return _serialize(await _call(client, "get_my_bids", item_type, sort))


@app.get("/my/assets", response_model=schemas.MyAssetsResult, tags=["My Assets"])
async def get_my_assets(
    item_type: str = "usernames",
    client: FragmentClient = Depends(get_fragment_client),
):
    """Return My Assets (usernames, numbers or gifts)."""
    return _serialize(await _call(client, "get_my_assets", item_type))


# ---------------------------------------------------------------------------
# Assignment
# ---------------------------------------------------------------------------


@app.get("/assign/accounts", response_model=schemas.AssignAccountsResult, tags=["Assignment"])
async def get_assign_accounts(
    item_type: int,
    slug: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Return Telegram accounts available for assignment of a username or gift."""
    return _serialize(await _call(client, "get_assign_accounts", item_type, slug))


@app.post("/assign", response_model=schemas.AssignResult, tags=["Assignment"])
async def assign_to_telegram(
    payload: schemas.AssignRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Assign a username or gift to a Telegram account."""
    return _serialize(await _call(client, "assign_to_telegram", payload.item_type, payload.slug, payload.assign_to))


# ---------------------------------------------------------------------------
# Auctions / selling
# ---------------------------------------------------------------------------


@app.post("/auctions/start", response_model=schemas.StartAuctionResult, tags=["Auctions"])
async def start_auction(
    payload: schemas.StartAuctionRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Start an auction for a username or gift."""
    return _serialize(await _call(client, "start_auction", payload.item_type, payload.slug, payload.min_amount, payload.max_amount))


@app.post("/assets/sell", response_model=schemas.StartAuctionResult, tags=["Auctions"])
async def sell_asset(
    payload: schemas.SellAssetRequest,
    client:Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Sell a username or gift at a fixed price."""
    return _serialize(await _call(client, "sell_asset", payload.item_type, payload.slug, payload.price))


# ---------------------------------------------------------------------------
# NFT transfer
# ---------------------------------------------------------------------------


@app.post("/nft/transfer/recipient", response_model=schemas.NftTransferRecipient | None, tags=["NFT Transfer"])
async def search_nft_transfer_recipient(
    payload: schemas.NftTransferRecipientRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Search for a recipient to transfer a gift NFT to."""
    return _serialize(await _call(client, "search_nft_transfer_recipient", payload.query))


@app.post("/nft/transfer/init", response_model=schemas.NftTransferRequest, tags=["NFT Transfer"])
async def init_nft_transfer(
    payload: schemas.InitNftTransferRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Initialize a gift NFT transfer request."""
    return _serialize(await _call(client, "init_nft_transfer", payload.slug, payload.recipient))


@app.post("/nft/transfer", response_model=schemas.TransactionResult, tags=["NFT Transfer"])
async def transfer_nft(
    payload: schemas.TransferNftRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Execute a gift NFT transfer."""
    return _serialize(await _call(client, "transfer_nft", payload.req_id, payload.show_sender))


# ---------------------------------------------------------------------------
# Sessions
# ---------------------------------------------------------------------------


@app.get("/sessions", response_model=list[schemas.SessionInfo], tags=["Sessions"])
async def get_sessions(
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Return active Fragment sessions."""
    return _serialize(await _call(client, "get_sessions"))


@app.post("/sessions/{session_id}/terminate", response_model=bool, tags=["Sessions"])
async def terminate_session(
    session_id: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Terminate a Fragment session by ID."""
    return _serialize(await _call(client, "terminate_session", session_id))


# ---------------------------------------------------------------------------
# Item history pagination
# ---------------------------------------------------------------------------


@app.get("/items/{item_type}/{slug}/orders-history", response_model=dict[str, Any], tags=["Item History"])
async def get_orders_history(
    item_type: int,
    slug: str,
    offset_id: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Load more bid/order history for a username, number or gift."""
    return _serialize(await _call(client, "get_orders_history", item_type, slug, offset_id))


@app.get("/items/{item_type}/{slug}/owners-history", response_model=dict[str, Any], tags=["Item History"])
async def get_owners_history(
    item_type: int,
    slug: str,
    offset_id: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Load more ownership history for a username, number or gift."""
    return _serialize(await _call(client, "get_owners_history", item_type, slug, offset_id))


# ---------------------------------------------------------------------------
# Anonymous numbers
# ---------------------------------------------------------------------------


@app.get("/anonymous-numbers/{number}/login-code", response_model=schemas.LoginCodeResult, tags=["Anonymous Numbers"])
async def get_login_code(
    number: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Fetch the pending login code for an anonymous +888 number."""
    return _serialize(await _call(client, "get_login_code", number))


@app.post("/anonymous-numbers/{number}/toggle-login-codes", tags=["Anonymous Numbers"])
async def toggle_login_codes(
    number: str,
    can_receive: bool,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Enable or disable login code delivery for an anonymous number."""
    await _call(client, "toggle_login_codes", number, can_receive)
    return {"ok": True}


@app.post("/anonymous-numbers/{number}/terminate-sessions", response_model=schemas.TerminateSessionsResult, tags=["Anonymous Numbers"])
async def terminate_sessions(
    number: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Terminate all Telegram sessions for an anonymous number."""
    return _serialize(await _call(client, "terminate_sessions", number))


# ---------------------------------------------------------------------------
# Withdrawals
# ---------------------------------------------------------------------------


@app.get("/withdrawal/nft/state", response_model=dict[str, Any], tags=["Withdrawals"])
async def get_nft_withdrawal_state(
    transaction: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Get NFT withdrawal state from Fragment page."""
    return _serialize(await _call(client, "get_nft_withdrawal_state", transaction))


@app.post("/withdrawal/nft/init", response_model=schemas.NftWithdrawalInitResult, tags=["Withdrawals"])
async def init_nft_withdrawal(
    payload: schemas.InitNftWithdrawalRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Initialize NFT withdrawal to wallet."""
    return _serialize(await _call(client, "init_nft_withdrawal", payload.transaction, payload.keep_gift))


@app.post("/withdrawal/nft/confirm", response_model=schemas.NftWithdrawalConfirmResult, tags=["Withdrawals"])
async def confirm_nft_withdrawal(
    payload: schemas.ConfirmNftWithdrawalRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Confirm NFT withdrawal after user approval."""
    return _serialize(await _call(client, "confirm_nft_withdrawal", payload.transaction, payload.confirm_hash, payload.keep_gift))


@app.get("/withdrawal/stars/state", response_model=schemas.StarsWithdrawalState, tags=["Withdrawals"])
async def get_stars_withdrawal_state(
    transaction: str,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Get Stars withdrawal state from Fragment page."""
    return _serialize(await _call(client, "get_stars_withdrawal_state", transaction))


@app.post("/withdrawal/stars/init", response_model=schemas.StarsWithdrawalInitResult, tags=["Withdrawals"])
async def init_stars_withdrawal(
    payload: schemas.InitStarsWithdrawalRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Initialize Stars revenue withdrawal to wallet."""
    return _serialize(await _call(client, "init_stars_withdrawal", payload.transaction, payload.withdrawal_data))


@app.post("/withdrawal/stars/confirm", response_model=schemas.StarsWithdrawalConfirmResult, tags=["Withdrawals"])
async def confirm_stars_withdrawal(
    payload: schemas.ConfirmStarsWithdrawalRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Confirm Stars revenue withdrawal after user approval."""
    return _serialize(await _call(client, "confirm_stars_withdrawal", payload.transaction, payload.withdrawal_data, payload.confirm_hash))


# ---------------------------------------------------------------------------
# Confirm / raw API
# ---------------------------------------------------------------------------


@app.post("/confirm-request", response_model=dict[str, Any], tags=["Raw"])
async def confirm_request(
    payload: schemas.ConfirmRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Send confirmReq to Fragment after broadcasting a TON transaction."""
    return _serialize(await _call(client, "confirm_request", payload.req_id, payload.boc, payload.referer))


@app.post("/raw-call", response_model=dict[str, Any], tags=["Raw"])
async def raw_call(
    payload: schemas.RawCallRequest,
    client: Annotated[FragmentClient, Depends(get_fragment_client)],
):
    """Send an arbitrary method call to the Fragment API (full mode only)."""
    return _serialize(await _call(client, "call", payload.method, payload.data, page_url=payload.page_url))
