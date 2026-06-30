"""Generate static openapi.json and swagger.yaml from the FastAPI app."""

import json

import yaml
from fastapi.openapi.utils import get_openapi

import main


openapi_schema = get_openapi(
    title=main.app.title,
    version=main.app.version,
    description=main.app.description,
    routes=main.app.routes,
)

with open("openapi.json", "w", encoding="utf-8") as f:
    json.dump(openapi_schema, f, indent=2, ensure_ascii=False)

with open("swagger.yaml", "w", encoding="utf-8") as f:
    yaml.safe_dump(openapi_schema, f, sort_keys=False, allow_unicode=True)

print("openapi.json and swagger.yaml generated")
