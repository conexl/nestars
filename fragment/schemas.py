"""Pydantic request/response models for the Fragment API wrapper."""

from __future__ import annotations

from typing import Any

from pydantic import BaseModel, Field


# ---------------------------------------------------------------------------
# Common / shared models
# ---------------------------------------------------------------------------


class FragmentHeaders(BaseModel):
    """Optional body-based credential payload for clients that prefer JSON.

    In normal operation credentials are supplied via HTTP headers.  This
    model is exposed so tools such as Swagger UI can easily populate them.
    """

    seed: str = Field(
        ...,
        description="12/18/24-word TON wallet mnemonic phrase.",
        examples=["word1 word2 ... word24"],
    )
    wallet_version: str = Field(
        "V5R1",
        description="TON wallet contract version.",
        examples=["V5R1", "V4R2"],
    )
    cookies: dict[str, str] | str | None = Field(
        None,
        description=(
            "Fragment session cookies for full mode. "
            "Pass as JSON object or semicolon-separated string."
        ),
    )
    api_key: str | None = Field(
        None,
        description="Optional TONAPI key (default key is used if omitted).",
    )
    timeout: float = Field(60.0, description="HTTP request timeout in seconds.")
    stats_enabled: bool = Field(
        True,
        description="Enable anonymous telemetry to the library author.",
    )


# ---------------------------------------------------------------------------
# Response models (mirroring FragmentAPI.types.results dataclasses)
# ---------------------------------------------------------------------------


class WalletInfo(BaseModel):
    address: str
    state: str
    balance_ton: float
    balance_usdt: float


class RecipientInfo(BaseModel):
    recipient: str
    name: str
    photo_url: str | None = None
    myself: bool = False


class ProfileInfo(BaseModel):
    name: str
    username: str
    photo_url: str | None = None
    identity_verified: bool
    wallet_address: str | None = None
    wallet_label: str | None = None
    wallet_verified: bool


class TransactionResult(BaseModel):
    tx_hash: str
    boc: str | None = None
    seqno_before: int | None = None
    seqno_after: int | None = None
    balance_before: float | None = None
    balance_after: float | None = None
    confirmed: bool = False


class StarsResult(BaseModel):
    transaction_id: str
    username: str
    amount: int
    payment_method: str = "ton"


class PremiumResult(BaseModel):
    transaction_id: str
    username: str
    amount: int
    payment_method: str = "ton"


class AdsTopupResult(BaseModel):
    transaction_id: str
    username: str
    amount: int


class GiveawayStarsResult(BaseModel):
    transaction_id: str
    channel: str
    winners: int
    amount: int
    payment_method: str = "ton"


class GiveawayPremiumResult(BaseModel):
    transaction_id: str
    channel: str
    winners: int
    amount: int
    payment_method: str = "ton"


class EvmInvoice(BaseModel):
    req_id: str
    invoice_address: str
    invoice_token: str
    invoice_chain_id: int
    invoice_chain_name: str
    invoice_amount_hex: str
    invoice_amount: float
    invoice_amount_raw: int
    token_symbol: str
    token_decimals: int
    expires_at: int
    payment_method: str
    api_hash: str
    page_url: str


class EvmPaymentResult(BaseModel):
    item_kind: str
    target: str
    amount: int
    payment_method: str
    invoice: EvmInvoice


class BidResult(BaseModel):
    transaction_id: str
    item_type: int
    slug: str
    bid: int
    confirm_method: str | None = None
    confirm_id: str | None = None


class AuctionInfo(BaseModel):
    highest_bid: str | None = None
    bid_step: str | None = None
    minimum_bid: str | None = None
    sell_price: str | None = None
    buy_now_price: str | None = None


class BidHistoryEntry(BaseModel):
    price: str | None = None
    date: str | None = None
    wallet: str | None = None


class OwnerHistoryEntry(BaseModel):
    price: str | None = None
    date: str | None = None
    wallet: str | None = None


class UsernameInfo(BaseModel):
    username: str
    status: str
    item_type: int
    ton_rate: float
    auction: AuctionInfo | None = None
    auction_end: str | None = None
    owner_wallet: str | None = None
    purchased_date: str | None = None
    bid_history: list[BidHistoryEntry] = []
    owner_history: list[OwnerHistoryEntry] = []
    bid_history_next_offset: str | None = None
    owner_history_next_offset: str | None = None


class NumberInfo(BaseModel):
    number: str
    display_number: str
    status: str
    item_type: int
    ton_rate: float
    restricted: bool = False
    auction: AuctionInfo | None = None
    auction_end: str | None = None
    owner_wallet: str | None = None
    purchased_date: str | None = None
    bid_history: list[BidHistoryEntry] = []
    owner_history: list[OwnerHistoryEntry] = []
    bid_history_next_offset: str | None = None
    owner_history_next_offset: str | None = None


class GiftAttribute(BaseModel):
    name: str
    value: str
    rarity: str | None = None


class GiftInfo(BaseModel):
    slug: str
    name: str
    status: str
    item_type: int
    ton_rate: float
    image_url: str | None = None
    sticker_url: str | None = None
    owner_wallet: str | None = None
    purchased_date: str | None = None
    auction: AuctionInfo | None = None
    auction_end: str | None = None
    attributes: list[GiftAttribute] = []
    issued: str | None = None
    bid_history: list[BidHistoryEntry] = []
    owner_history: list[OwnerHistoryEntry] = []
    bid_history_next_offset: str | None = None
    owner_history_next_offset: str | None = None


class StarsPrice(BaseModel):
    stars: int
    ton_price: str
    usd_price: str


class StarsPrices(BaseModel):
    packages: list[StarsPrice]
    ton_rate: float


class PremiumPriceOption(BaseModel):
    months: int
    label: str
    ton_price: str
    usd_price: str
    discount: str | None = None


class PremiumPrices(BaseModel):
    options: list[PremiumPriceOption]
    ton_rate: float


class StarsTransaction(BaseModel):
    recipient: str
    stars: int
    price_ton: str
    date: str


class PremiumTransaction(BaseModel):
    recipient: str
    duration: str
    price_ton: str
    date: str


class TopupTransaction(BaseModel):
    recipient: str
    amount: int
    date: str


class UsernamesResult(BaseModel):
    items: list[dict[str, Any]]
    next_offset_id: str | None = None


class NumbersResult(BaseModel):
    items: list[dict[str, Any]]
    next_offset_id: str | None = None


class GiftsResult(BaseModel):
    items: list[dict[str, Any]]
    next_offset: int | None = None


class MyBid(BaseModel):
    item_type: str
    slug: str
    name: str
    bid: float
    status: str
    date: str
    image_url: str | None = None
    description: str | None = None


class MyBidsResult(BaseModel):
    items: list[MyBid]
    ton_rate: float
    total_count: int


class MyAsset(BaseModel):
    item_type: str
    slug: str
    name: str
    description: str | None = None
    image_url: str | None = None
    assigned_to: str | None = None
    assigned_name: str | None = None


class MyAssetsResult(BaseModel):
    items: list[MyAsset]
    ton_rate: float
    total_count: int


class TelegramAccount(BaseModel):
    id: str
    name: str
    type: str
    photo_url: str | None = None


class AssignAccountsResult(BaseModel):
    accounts: list[TelegramAccount]
    can_disable: bool


class AssignResult(BaseModel):
    ok: bool
    message: str | None = None
    need_pay: bool = False
    req_id: str | None = None
    amount: str | None = None
    assign_name: str | None = None


class StartAuctionResult(BaseModel):
    ok: bool
    req_id: str | None = None


class NftTransferRecipient(BaseModel):
    myself: bool
    recipient: str
    name: str
    photo_url: str | None = None


class NftTransferRequest(BaseModel):
    req_id: str
    myself: bool
    item_title: str
    content: str
    button: str


class SessionInfo(BaseModel):
    session_id: str
    device: str
    location: str
    date: str | None = None
    is_current: bool


class LoginCodeResult(BaseModel):
    number: str
    code: str | None = None
    active_sessions: int


class TerminateSessionsResult(BaseModel):
    number: str
    message: str | None = None


class NftWithdrawalInitResult(BaseModel):
    ok: bool
    confirm_message: str | None = None
    confirm_button: str | None = None
    confirm_hash: str | None = None
    error: str | None = None


class NftWithdrawalConfirmResult(BaseModel):
    ok: bool
    need_update: bool
    mode: str
    html: str | None = None
    error: str | None = None


class StarsWithdrawalState(BaseModel):
    transaction: str
    withdrawal_data: str


class StarsWithdrawalInitResult(BaseModel):
    ok: bool
    confirm_message: str | None = None
    confirm_button: str | None = None
    confirm_hash: str | None = None
    error: str | None = None


class StarsWithdrawalConfirmResult(BaseModel):
    ok: bool
    need_update: bool
    mode: str
    html: str | None = None
    error: str | None = None


class BatchItemResult(BaseModel):
    type: str
    username: str
    amount: int
    ok: bool
    result: Any = None
    error: str | None = None
    chunk_index: int = 0


class BatchResult(BaseModel):
    total: int
    succeeded: int
    failed: int
    chunks_sent: int
    items: list[BatchItemResult] = []


# ---------------------------------------------------------------------------
# Request models
# ---------------------------------------------------------------------------


class AuthRequest(BaseModel):
    seed: str = Field(..., description="Wallet mnemonic phrase.")
    wallet_version: str = Field("V5R1", description="V4R2 or V5R1.")
    phone: str | None = Field(None, description="Phone number for SMS flow.")
    print_qr: bool = Field(True, description="Print QR code in server logs.")
    timeout: float = Field(60.0, description="Authentication timeout.")


class AuthResponse(BaseModel):
    cookies: dict[str, str]


class RecipientUsernameRequest(BaseModel):
    username: str = Field(..., description="Telegram username (with or without @).")


class PremiumRecipientRequest(BaseModel):
    username: str
    months: int = Field(3, description="3, 6 or 12 months.")


class GiveawayStarsRecipientRequest(BaseModel):
    channel: str = Field(..., description="Telegram channel username.")
    winners: int = Field(1, ge=1)
    amount: int = Field(500, ge=50)


class GiveawayPremiumRecipientRequest(BaseModel):
    channel: str
    winners: int = Field(1, ge=1)
    months: int = Field(3, description="3, 6 or 12 months.")


class PurchaseStarsRequest(BaseModel):
    username: str
    amount: int = Field(..., ge=50, le=10_000_000)
    show_sender: bool = Field(True)
    payment_method: str = Field(
        "ton",
        description=(
            "ton, usdt_ton, usdt_eth, usdt_pol, usdc_eth, usdc_base, usdc_pol"
        ),
    )


class PurchasePremiumRequest(BaseModel):
    username: str
    months: int = Field(..., description="3, 6 or 12 months.")
    show_sender: bool = Field(True)
    payment_method: str = Field(
        "ton",
        description=(
            "ton, usdt_ton, usdt_eth, usdt_pol, usdc_eth, usdc_base, usdc_pol"
        ),
    )


class TopupTonRequest(BaseModel):
    username: str
    amount: int = Field(..., ge=1)
    show_sender: bool = Field(True)


class BatchItem(BaseModel):
    type: str = Field(
        ...,
        pattern="^(stars|premium|ton)$",
        description="Purchase kind.",
    )
    username: str
    amount: int | None = Field(None, description="For stars or ton items.")
    months: int | None = Field(None, description="For premium items.")
    show_sender: bool = Field(True)


class BatchPurchaseRequest(BaseModel):
    items: list[BatchItem]
    payment_method: str = Field(
        "ton",
        description="Only ton and usdt_ton are supported.",
    )


class GiveawayStarsRequest(BaseModel):
    channel: str
    winners: int = Field(..., ge=1, le=24_000)
    amount: int = Field(..., ge=50)
    payment_method: str = Field("ton")


class GiveawayPremiumRequest(BaseModel):
    channel: str
    winners: int = Field(..., ge=1, le=24_000)
    months: int = Field(..., description="3, 6 or 12 months.")
    payment_method: str = Field("ton")


class PlaceBidRequest(BaseModel):
    item_type: int = Field(..., description="1 = username, 3 = number, 5 = gift.")
    slug: str
    bid: int = Field(..., ge=1)


class SearchUsernamesRequest(BaseModel):
    query: str = ""
    sort: str | None = None
    filter: str | None = None
    offset_id: str | None = None


class SearchNumbersRequest(BaseModel):
    query: str = ""
    sort: str | None = None
    filter: str | None = None
    offset_id: str | None = None


class SearchGiftsRequest(BaseModel):
    query: str = ""
    collection: str | None = None
    sort: str | None = None
    filter: str | None = None
    view: str | None = None
    attr: dict[str, list[str]] | None = None
    offset: int | None = None


class AssignRequest(BaseModel):
    item_type: int = Field(..., description="1 = username, 5 = gift.")
    slug: str
    assign_to: str | None = Field(None, description="Telegram account ID to assign to.")


class StartAuctionRequest(BaseModel):
    item_type: int
    slug: str
    min_amount: int = Field(..., ge=1)
    max_amount: int = Field(0, description="0 starts an auction; equal to min_amount for fixed sale.")


class SellAssetRequest(BaseModel):
    item_type: int
    slug: str
    price: int = Field(..., ge=1)


class NftTransferRecipientRequest(BaseModel):
    query: str


class InitNftTransferRequest(BaseModel):
    slug: str
    recipient: str


class TransferNftRequest(BaseModel):
    req_id: str
    show_sender: bool = Field(True)


class ToggleLoginCodesRequest(BaseModel):
    number: str
    can_receive: bool


class NftWithdrawalStateRequest(BaseModel):
    transaction: str


class InitNftWithdrawalRequest(BaseModel):
    transaction: str
    keep_gift: bool = Field(False)


class ConfirmNftWithdrawalRequest(BaseModel):
    transaction: str
    confirm_hash: str
    keep_gift: bool = Field(False)


class StarsWithdrawalStateRequest(BaseModel):
    transaction: str


class InitStarsWithdrawalRequest(BaseModel):
    transaction: str
    withdrawal_data: str


class ConfirmStarsWithdrawalRequest(BaseModel):
    transaction: str
    withdrawal_data: str
    confirm_hash: str


class ConfirmRequest(BaseModel):
    req_id: str
    boc: str
    referer: str = Field("stars/buy", description="Fragment referer page path.")


class RawCallRequest(BaseModel):
    method: str
    data: dict[str, Any] | None = None
    page_url: str = Field("https://fragment.com", description="Base page URL for hash fetch.")
