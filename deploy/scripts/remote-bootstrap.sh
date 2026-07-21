#!/usr/bin/env bash
# Idempotent VPS bootstrap stub — see docs/engineering/07-server-deployment.md
set -euo pipefail

echo "Install docker + compose plugin, create /opt/copilot, copy deploy/ files."
echo "Required inputs: SSH host, domain, DEMO_TOKEN, optional LLM_API_KEY."
echo "Then: ./remote-deploy.sh"
