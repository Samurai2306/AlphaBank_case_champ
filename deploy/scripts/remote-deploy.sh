#!/usr/bin/env bash
set -euo pipefail
ROOT="${1:-/opt/copilot}"
cd "$ROOT"
git pull
docker compose -f deploy/docker-compose.prod.yml up -d --build
curl -fsS "http://127.0.0.1/api/v1/health" || true
echo "deploy done"
