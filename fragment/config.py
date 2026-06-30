"""Application configuration."""

import os

from dotenv import load_dotenv

load_dotenv()


class Settings:
    """FastAPI service settings."""

    APP_NAME: str = os.getenv("APP_NAME", "Fragment API")
    APP_DESCRIPTION: str = os.getenv(
        "APP_DESCRIPTION",
        "RESTful wrapper around fragment-api-py with OpenAPI/Swagger docs",
    )
    APP_VERSION: str = os.getenv("APP_VERSION", "1.0.0")
    HOST: str = os.getenv("HOST", "0.0.0.0")
    PORT: int = int(os.getenv("PORT", "8000"))
    DOCS_URL: str = os.getenv("DOCS_URL", "/docs")
    REDOC_URL: str = os.getenv("REDOC_URL", "/redoc")
    OPENAPI_URL: str = os.getenv("OPENAPI_URL", "/openapi.json")


settings = Settings()
