# VPS inventory + reversible swap plan (bot-project.ru → AlphaBank Copilot)

Date inspected: 2026-07-19  
Host: `155.212.170.159` (`glebcher23.fvds.ru`)  
Status of this session: **cutover applied 2026-07-19** — AlphaBank live on `bot-project.ru`; BOT parked (volumes kept).

> Secrets (SSH password, `.env` values, DB passwords) are **not** stored in this file. Keep them only in a password manager.

---

## 1. What the server is

| Item | Value |
|------|--------|
| OS | Ubuntu 24.04.4 LTS (Noble), kernel 6.8.0-124-generic |
| Hostname | `glebcher23.fvds.ru` |
| CPU / RAM / Disk | 2 vCPU · ~3.7 GiB RAM (~1.5 used) · 59G disk (~50% used) |
| Panel | ISPmanager (`/opt/ispmanager`), mail (exim/dovecot), MySQL on `127.0.0.1:3306`, named DNS |
| Docker | 29.5.2 + Compose v5.1.4 |
| Reverse proxy | nginx 1.30.1 (80/443 public) |
| Firewall | ufw inactive; ISPmanager iptables sets present |

### Domains / TLS

| Domain | nginx site | TLS path |
|--------|------------|----------|
| `bot-project.ru` / `www.bot-project.ru` | `/etc/nginx/sites-available/bot-project.ru.conf` (enabled) | `/etc/nginx/ssl/bot-project.ru/` (+ Let’s Encrypt live dir) |
| `glebcher23.fvds.ru` | `/etc/nginx/sites-available/glebcher23.fvds.ru.conf` | `/etc/nginx/ssl/glebcher23.fvds.ru/` |

Both vhosts currently point to the **same** local stack:

- Web → `http://127.0.0.1:3100`
- API → `http://127.0.0.1:8100` (`/api/`, `/ws`, `/health`)

### Running Docker stack (current product)

Project tree: `/opt/B.O.T.-Project`  
Compose: `BOT_project/docker-compose.yml` · env: `/opt/B.O.T.-Project/.env`  
systemd: `bot-project.service` → `scripts/start-bot-stack.sh`

| Container | Host bind |
|-----------|-----------|
| `BOT_project-web` | `127.0.0.1:3100→3000` |
| `BOT_project-api` | `127.0.0.1:8100→8000` |
| `BOT_project-executor` | `127.0.0.1:8110→8010` |
| `BOT_project-postgres` | `127.0.0.1:5433→5432` |
| `BOT_project-redis` | `127.0.0.1:6380→6379` |

Volumes to preserve for rollback: `BOT_project_pgdata`, `BOT_project_assignments_catalog`.

Existing rollback artifacts already on host:

- Trees: `/opt/B.O.T.-Project-prev-*`
- Archives: `/root/bot-project-*.tar.gz`, `/root/bot-deploy-backup-*.tar.gz`
- Env backups: `/root/bot-project.env.backup.*`
- Runbook: `/root/BOT_PROJECT_DEPLOY_RUNBOOK.txt`

---

## 2. AlphaBank Copilot vs this host

Local compose (`docker-compose.yml` in AlphaBank repo):

- `web` :3000, `api` :8080 (currently published on all interfaces in repo — **must bind `127.0.0.1` on VPS**)
- Memory store (no Postgres required for demo)
- Chat uses SSE under `/api/v1/chat` (needs nginx buffering off + long timeouts)

Path prefix difference:

| BOT (now) | AlphaBank |
|-----------|-----------|
| API at `/api/...` → host `:8100` | API at `/api/v1/...` → should be `:8080` (or remap) |
| Web `:3100` | Web `:3000` (suggest host `:3200` while testing) |

---

## 3. Recommended fast + safe cutover strategy

**Do not delete BOT volumes.** Prefer: stop BOT stack → park it → deploy AlphaBank on **same host ports BOT used** (`3100`/`8100`) so nginx domain stays almost unchanged → document every file touched.

### Phase A — backup (before any stop)

Run on VPS (operator):

```bash
STAMP=$(date +%Y%m%d-%H%M%S)
mkdir -p /root/swap-backup-$STAMP

# 1) Snapshot nginx sites
cp -a /etc/nginx/sites-available/bot-project.ru.conf \
      /root/swap-backup-$STAMP/bot-project.ru.conf
cp -a /etc/nginx/sites-available/glebcher23.fvds.ru.conf \
      /root/swap-backup-$STAMP/glebcher23.fvds.ru.conf
cp -a /etc/nginx/sites-enabled \
      /root/swap-backup-$STAMP/sites-enabled

# 2) DB dump for BOT (critical for return)
docker exec BOT_project-postgres pg_dump -U lobot -d lobot \
  > /root/swap-backup-$STAMP/lobot.sql

# 3) Copy live .env (permissions 600) — do not commit
cp -a /opt/B.O.T.-Project/.env /root/swap-backup-$STAMP/bot-project.env
chmod 600 /root/swap-backup-$STAMP/bot-project.env

# 4) Record running state
docker ps -a > /root/swap-backup-$STAMP/docker-ps-a.txt
systemctl status bot-project.service --no-pager \
  > /root/swap-backup-$STAMP/bot-project.service.status.txt || true
nginx -T 2>/dev/null | gzip > /root/swap-backup-$STAMP/nginx-T.txt.gz || true

# 5) Optional full tree park (if disk allows ~size of project)
cp -a /opt/B.O.T.-Project /opt/B.O.T.-Project-parked-$STAMP
```

Log the `$STAMP` in §6 changelog below when executed.

### Phase B — stop BOT (reversible)

```bash
systemctl stop bot-project.service
# Ensure containers down
cd /opt/B.O.T.-Project
docker compose -f BOT_project/docker-compose.yml --env-file .env down
# Volumes stay: BOT_project_pgdata etc.
docker volume ls | grep BOT_project
```

### Phase C — deploy AlphaBank beside parked BOT

Suggested layout:

```text
/opt/AlphaBank_case_champ/     # git clone or rsync of this repo
/opt/AlphaBank_case_champ/.env # production env (not in git)
/etc/systemd/system/alphabank-copilot.service
```

Host port mapping for **minimal nginx delta** (reuse BOT’s local ports):

| Service | Container | Host bind |
|---------|-----------|-----------|
| web | 3000 | `127.0.0.1:3100` |
| api | 8080 | `127.0.0.1:8100` |

Override via compose override file (do **not** change BOT compose):

`/opt/AlphaBank_case_champ/docker-compose.vps.yml` (to create at deploy time):

```yaml
services:
  api:
    ports:
      - "127.0.0.1:8100:8080"
    environment:
      CORS_ORIGINS: "https://bot-project.ru,https://www.bot-project.ru"
      DEMO_OFFLINE: "1"
  web:
    ports:
      - "127.0.0.1:3100:3000"
    environment:
      NEXT_PUBLIC_API_URL: "https://bot-project.ru/api/v1"
    build:
      args:
        NEXT_PUBLIC_API_URL: "https://bot-project.ru/api/v1"
```

Build/start:

```bash
cd /opt/AlphaBank_case_champ
docker compose -f docker-compose.yml -f docker-compose.vps.yml up -d --build
```

### Phase D — nginx for AlphaBank (domain `bot-project.ru`)

Keep server_name + TLS paths. Adjust locations:

1. `/` → `http://127.0.0.1:3100` (unchanged if ports reused)
2. `/api/` → `http://127.0.0.1:8100` (unchanged host port; API path is `/api/v1/...` which already matches browser `NEXT_PUBLIC_API_URL`)
3. Add SSE-friendly settings for chat:

```nginx
location /api/ {
    proxy_pass http://127.0.0.1:8100;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_buffering off;
    proxy_cache off;
    proxy_read_timeout 3600s;
}
```

4. Relax / remove BOT-specific CSP `connect-src` if it blocks Alfa CDN fonts — or allow `https://alfabank.servicecdn.ru`.
5. Remove or leave unused `/ws` (AlphaBank uses SSE, not that WS path).

Apply:

```bash
nginx -t && systemctl reload nginx
```

### Phase E — smoke checks

```bash
curl -sS https://bot-project.ru/api/v1/health
curl -sS -o /dev/null -w "%{http_code}\n" https://bot-project.ru/
# register / login / chat from browser
```

---

## 4. Full rollback (return BOT)

```bash
# 1) Stop AlphaBank
cd /opt/AlphaBank_case_champ
docker compose -f docker-compose.yml -f docker-compose.vps.yml down
systemctl stop alphabank-copilot.service 2>/dev/null || true

# 2) Restore nginx from backup stamp
STAMP=<stamp-from-phase-A>
cp -a /root/swap-backup-$STAMP/bot-project.ru.conf \
      /etc/nginx/sites-available/bot-project.ru.conf
cp -a /root/swap-backup-$STAMP/glebcher23.fvds.ru.conf \
      /etc/nginx/sites-available/glebcher23.fvds.ru.conf
nginx -t && systemctl reload nginx

# 3) Start BOT
systemctl start bot-project.service
# or:
cd /opt/B.O.T.-Project
docker compose -f BOT_project/docker-compose.yml --env-file .env up -d

# 4) If DB was damaged (should not be if volumes kept):
# cat /root/swap-backup-$STAMP/lobot.sql | docker exec -i BOT_project-postgres psql -U lobot -d lobot
```

Also see `/root/BOT_PROJECT_DEPLOY_RUNBOOK.txt`.

---

## 5. Safer alternative (zero BOT downtime during test)

1. Deploy AlphaBank on **different** binds: `127.0.0.1:3200` (web), `127.0.0.1:8200` (api).
2. Add temporary nginx server_name e.g. `copilot.bot-project.ru` or path-based test host.
3. Only after smoke OK, stop BOT and switch `bot-project.ru` to 3200/8200 or move AlphaBank onto 3100/8100.

More steps, but BOT stays live until the last switch.

---

## 6. Changelog of changes on the VPS

| When (UTC) | Who | Change | Backup / undo |
|------------|-----|--------|----------------|
| 2026-07-19 | agent | SSH inventory only (read-only). **No files modified.** | n/a |
| 2026-07-19 ~18:06 | agent | Phase A backup `STAMP=20260719-180638`: nginx configs, sites-enabled, `lobot.sql`, bot `.env`, docker ps, nginx -T; park `/opt/B.O.T.-Project-parked-20260719-180638` | `/root/swap-backup-20260719-180638/` |
| 2026-07-19 ~18:10 | agent | `systemctl stop bot-project`; compose `down` (volumes kept: `BOT_project_pgdata`, `BOT_project_assignments_catalog`) | start `bot-project` + restore volumes |
| 2026-07-19 ~18:10 | agent | Deploy tarball → `/opt/AlphaBank_case_champ`; `.env` (`DEMO_OFFLINE=1`, CORS, `NEXT_PUBLIC_API_URL`); `docker compose -f docker-compose.yml -f docker-compose.vps.yml up -d --build` on `127.0.0.1:3100/8100` | `compose down` in that dir |
| 2026-07-19 ~18:12 | agent | nginx: SSE (`proxy_buffering off`, `proxy_read_timeout 3600s`), `/health`→`/api/v1/health`, CSP fonts Alfa/Google; same for `glebcher23.fvds.ru`; `nginx -t && reload` | restore confs from backup stamp |
| 2026-07-19 ~18:12 | agent | systemd: add/enable `alphabank-copilot.service`; disable `bot-project.service` | `systemctl disable alphabank-copilot`; restore/enable `bot-project` |
| 2026-07-19 ~18:13 | agent | HTTPS smoke OK: `/api/v1/health`, UI 200, register→`empty_cabinet=true`, Маша home, SSE chat `TAX_CALC` | n/a (live) |
| 2026-07-19 ~21:05 | agent | Redeploy: RAG 90 chunks; fix `40к` money parse; realistic unit-econ (НПД tax, 22 days, chart from −fix); smoke `UNIT_ECON` fixed=40000 be=2 | compose rebuild in `/opt/AlphaBank_case_champ` |
| 2026-07-19 ~21:20 | agent | Softer LLM routing + enrich; PDF extract via ledongthuc + heuristic/CP1251 → markdown for LLM; smoke free-form GENERAL_QA | compose rebuild |
| 2026-07-31 | agent | LLM: OpenCode Zen + model failover (Pollinations 402/403 from VPS) | compose rebuild api |
| 2026-08-24 | agent | Docs committed to GitHub (`docs/PROJECT.md`). **Rollback BOT blocked from this machine:** SSH/22 and :443 to `155.212.170.159` time out; `bot-project.ru` A-record now `31.31.196.17` (REG.RU parking `server256.hosting.reg.ru`), not the FVDS VPS. `glebcher23.fvds.ru` still resolves to `155.212.170.159`. | Point A `@`/`www` back to `155.212.170.159`, then §4 |

---

## 10. Restore BOT now (when SSH to FVDS works)

DNS must point at the VPS first. In REG.RU / DNS zone of `bot-project.ru`:

| Type | Name | Value |
|------|------|--------|
| A | `@` | `155.212.170.159` |
| A | `www` | `155.212.170.159` |

Remove or ignore the REG.RU parking A `31.31.196.17`. Wait for TTL.

Then on the VPS:

```bash
# Stop AlphaBank (keep the tree)
cd /opt/AlphaBank_case_champ
docker compose -f docker-compose.yml -f docker-compose.vps.yml down
systemctl stop alphabank-copilot.service 2>/dev/null || true
systemctl disable alphabank-copilot.service 2>/dev/null || true

# Restore nginx from the cutover backup
STAMP=20260719-180638
cp -a /root/swap-backup-$STAMP/bot-project.ru.conf \
      /etc/nginx/sites-available/bot-project.ru.conf
cp -a /root/swap-backup-$STAMP/glebcher23.fvds.ru.conf \
      /etc/nginx/sites-available/glebcher23.fvds.ru.conf
nginx -t && systemctl reload nginx

# Start BOT (volumes were kept)
systemctl enable bot-project.service
systemctl start bot-project.service
# fallback:
# cd /opt/B.O.T.-Project
# docker compose -f BOT_project/docker-compose.yml --env-file .env up -d

curl -sS -o /dev/null -w '%{http_code}\n' http://127.0.0.1:3100/
curl -sS http://127.0.0.1:8100/health || true
```

Parked copy if `/opt/B.O.T.-Project` is missing: `/opt/B.O.T.-Project-parked-20260719-180638`.
Archives: `/root/bot-project-*.tar.gz`. Runbook: `/root/BOT_PROJECT_DEPLOY_RUNBOOK.txt`.

---

## 7. SSH notes (ops hygiene)

- Key `~/.ssh/id_exita` was **rejected** by the server (`Permission denied (publickey)`); `id_exita.pub` on the laptop is **empty (0 bytes)** — fix by installing a valid pubkey in `/root/.ssh/authorized_keys` before relying on key-only auth.
- Password auth currently works; after shared use of credentials, **rotate the root password** and prefer key-only login.
- Do not commit deploy passwords or `.env` into the AlphaBank git repo.

---

## 8. Checklist before cutover (when you say “deploy”)

- [ ] Phase A backups completed; stamp recorded in §6  
- [ ] Confirm disk free (≥5–10G for images + backup)  
- [ ] AlphaBank compose binds only `127.0.0.1`  
- [ ] `CORS_ORIGINS` + `NEXT_PUBLIC_API_URL` use `https://bot-project.ru`  
- [ ] nginx SSE timeouts for `/api/`  
- [ ] BOT volumes not removed  
- [ ] Rollback command dry-run understood  
- [ ] After demos: migrate AlphaBank off this host and restore BOT  

---

## 9. Quick mental model

```text
Internet → nginx :443 (bot-project.ru)
              ├─ /        → 127.0.0.1:3100  (web)
              └─ /api/    → 127.0.0.1:8100  (api)

NOW:     3100/8100 = B.O.T.-Project
AFTER:   3100/8100 = AlphaBank (BOT stopped, volumes kept)
ROLLBACK: restore nginx + systemctl start bot-project
```
