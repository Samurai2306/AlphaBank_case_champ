# Server Deployment (VPS) — подготовка под ваш доступ

Когда дадите доступ к серверу, выкладка идёт **по этому runbook**. Документ написан так, чтобы агент мог выполнить шаги без догадок.

## Целевая топология на одном сервере

```text
Internet
   │
   ▼
┌─────────────────────────────┐
│  Nginx (or Caddy) :443/:80  │  TLS, HTTP→HTTPS, security headers
│  ├─ /           → web :3000 │  Next.js (standalone) or static+node
│  └─ /api/       → api :8080 │  Go copilot-api
└─────────────────────────────┘
         │              │
         ▼              ▼
   web container    api container
                         │
                         ▼
                   uploads volume
```

Опционально позже: `api.` subdomain → только Go; apex/www → web.

## Что нужно от вас перед деплоем

Передайте агенту (в чате или secrets store, **не коммитить**):

| Параметр | Пример | Зачем |
|----------|--------|-------|
| SSH host | `user@IP` или hostname | доступ |
| SSH port / key | `22` + private key path | вход |
| Домен(ы) | `copilot.example.ru` | DNS + TLS |
| DNS доступ | где A-записи (REG.RU / Cloudflare…) | зона |
| LLM key | OpenAI/OpenRouter | chat |
| Желаемый path | `/opt/copilot` | install dir |

Агент **не** хранит ключи в git. На сервере — `/opt/copilot/.env` chmod 600.

## DNS (зона)

Минимально:

| Type | Name | Value |
|------|------|-------|
| A | `@` или `copilot` | `SERVER_IP` |
| A | `www` (optional) | `SERVER_IP` |

Рекомендуемый вариант для чистоты:

| Host | Target |
|------|--------|
| `copilot.example.ru` | web+api через Nginx path `/api` |
| или `app.` + `api.` | split origins — тогда поправить `CORS_ORIGINS` и Next `NEXT_PUBLIC_API_URL` |

После делегирования: `dig +short copilot.example.ru` → IP сервера.

## Сервер: baseline (один раз)

OS: Ubuntu 22.04/24.04 LTS (предпочтительно).

```bash
# концептуальный чеклист (выполнит агент по SSH)
sudo apt update && sudo apt upgrade -y
# Docker Engine + Compose plugin
# UFW: allow OpenSSH, 80, 443 — deny остальное public
# create user deploy / use your user in docker group
# mkdir -p /opt/copilot && chown deploy:deploy
```

Пакеты: `docker.io` / official Docker CE, `docker compose`, `git`, `curl`, `fail2ban` (желательно).

## Артефакты в репо (будут созданы при коде)

```text
deploy/
  docker-compose.prod.yml
  nginx/copilot.conf
  caddy/Caddyfile          # альтернатива nginx+certbot
  scripts/remote-bootstrap.sh
  scripts/remote-deploy.sh
.env.example                 # без секретов
apps/api/Dockerfile
apps/web/Dockerfile
```

### `docker-compose.prod.yml` (сервисы)

| Service | Image/build | Notes |
|---------|-------------|-------|
| `api` | Go multi-stage | port 8080 internal, restart unless-stopped |
| `web` | Next standalone | port 3000 internal |
| `nginx` | `nginx:alpine` | 80/443 publish; mount certs |

Сеть: bridge `copilot_net`. Проверка готовности — `api` `/api/v1/ready`.

## TLS

**Вариант A — Caddy** (проще): автоматический Let's Encrypt.  
**Вариант B — Nginx + Certbot** — классика.

Redirect HTTP→HTTPS. HSTS после стабилизации.

## Env на сервере (`/opt/copilot/.env`)

```env
# public
PUBLIC_HOST=copilot.example.ru
CORS_ORIGINS=https://copilot.example.ru
NEXT_PUBLIC_API_URL=https://copilot.example.ru/api/v1

# api
HTTP_ADDR=:8080
DEMO_TOKEN=***long***
LLM_BASE_URL=https://api.openai.com/v1
LLM_API_KEY=***
LLM_MODEL=gpt-4o-mini
UPLOAD_DIR=/data/uploads
DEMO_OFFLINE=0
```

Сид Маши создаётся при старте API.

## Процедура деплоя (когда будет доступ)

1. Проверить SSH и `sudo`.  
2. Baseline Docker + UFW.  
3. Клонировать репо в `/opt/copilot` (или `git pull`).  
4. Создать `.env` из example.  
5. DNS A → IP; дождаться propagate.  
6. `docker compose -f deploy/docker-compose.prod.yml up -d --build`.  
7. Certbot/Caddy выпустить сертификат.  
8. Smoke: `curl -fsS https://$HOST/api/v1/health` и открыть Home в браузере.  
9. Прогнать demo script «Маша».

Обновления:

```bash
cd /opt/copilot && git pull
docker compose -f deploy/docker-compose.prod.yml up -d --build
```

Zero-downtime later: blue/green — не требуется для питча.

## Backup & safety

| Что | Как |
|-----|-----|
| Uploads | volume с PDF, если он включён |
| `.env` | копия в password manager, не в S3 публично |
| Rollback | `git checkout PREV && compose up -d --build` |

## Security checklist (сервер)

- [ ] 22/SSH key-only, optional non-default port  
- [ ] UFW: 22/80/443 only  
- [ ] `.env` chmod 600  
- [ ] Docker socket не доступен веб-контейнеру  
- [ ] Rate limit на `/api/v1/chat` (в Go)  
- [ ] Fail2ban на sshd  
- [ ] LLM key с лимитом spend  

## Локально vs сервер

| | Local | Server |
|--|-------|--------|
| Compose file | `docker-compose.yml` | `deploy/docker-compose.prod.yml` |
| TLS | no | yes |
| Domain | localhost | real DNS |
| Сид Маши | при старте API | при старте API |
| DEMO_OFFLINE | optional | pitch fallback if key fails |

## Что агент сделает после получения доступа

1. Зафиксирует в `docs/engineering/08-server-inventory.md` (создаст): IP, OS, domain, compose path — **без секретов**.  
2. Выполнит bootstrap + deploy.  
3. Вернёт URL и результат smoke.  
4. Не будет менять DNS у регистратора без вашего подтверждения, если нет API-токена зоны.

## Non-goals первого деплоя

- Kubernetes  
- Multi-region  
- On-prem GPU LLM  
- Внешний managed DB (можно позже)  
