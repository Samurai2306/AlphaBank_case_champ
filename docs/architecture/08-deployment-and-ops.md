# Deployment & Operations

## Demo

**Backend = Go** ([09-demo-backend-go.md](09-demo-backend-go.md)).  
**VPS runbook:** [../engineering/07-server-deployment.md](../engineering/07-server-deployment.md).

### Docker Compose

| File | Use |
|------|-----|
| `docker-compose.yml` | local: `web`, `api` (Go), `postgres`, `redis` |
| `deploy/docker-compose.prod.yml` | server + Nginx/Caddy TLS |

```text
make up          # local compose
make seed        # persona Маша (Go cmd/seed)
make test-api    # go test ./...
make test        # api + web
```

### Environments

| Env | LLM | Data | Edge |
|-----|-----|------|------|
| local | cloud key in `.env` | local postgres | optional |
| server / pitch | cloud + `DEMO_OFFLINE` fallback | docker volumes | Nginx/Caddy TLS + DNS |

### Observability demo

- `slog` JSON on Go (`request_id`, intent, latency).  
- `/api/v1/health` + `/api/v1/ready` (db+redis).  
- Optional `/metrics` Prometheus.

### Fallback без LLM

`DEMO_OFFLINE=1`: фикстуры SSE + SDUI из `apps/api/testdata/offline`.

---

## Production

### Topology

- K8s namespaces: `copilot-bff`, `copilot-ai`, `copilot-inference`, `copilot-data`.  
- GPU node pool для vLLM.  
- HPA по queue depth / tokens/s.  
- Multi-AZ; inference can be active-passive per bank DR.

### CI/CD

1. Lint + unit + contract tests (shared schemas).  
2. Eval suite (tax golden) gate.  
3. Security scan images.  
4. Deploy to NFT → UAT (bank) → prod with canary.  
5. Prompt/corpus changes требуют отдельного approve.

### SLOs

| SLO | Target |
|-----|--------|
| Availability chat API | 99.9% |
| p95 time-to-first-token | < 2s |
| Guardrail false-neg on jailbreak set | < 1% (internal suite) |
| Tax eval accuracy | ≥ 95% on golden |

### Runbooks

- **LLM down:** degrade to FAQ retrieval + human chat handoff.  
- **Vector DB down:** block TAX/LEGAL with apology; allow GENERAL bank FAQ cache.  
- **Hallucination incident:** kill switch prompts → safe mode; notify compliance.  
- **PII leak suspicion:** rotate keys, invalidate sessions, forensic on audit.

### Cost controls

- Token budgets per user/day.  
- Cache embeddings and frequent tariff answers.  
- Router на меньшей модели.

### Compliance ops

- Corpus version pinned; quarterly legal review НК.  
- Access to prod logs — break-glass only.  
- DPIA / банковские анкеты ИБ — до launch Phase 1.
